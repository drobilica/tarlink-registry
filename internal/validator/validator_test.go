package validator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"go.yaml.in/yaml/v3"
)

var idPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)
var sha256Pattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

const maxManifestBytes = 1 << 20

var discoveryCategories = map[string]bool{
	"game-development": true,
	"emulation":        true,
	"graphics":         true,
	"development":      true,
	"utilities":        true,
}

var desktopCategories = map[string]bool{
	"Development": true,
	"Emulator":    true,
	"Game":        true,
	"Graphics":    true,
	"Utility":     true,
}

type manifest struct {
	Schema      int         `yaml:"schema"`
	ID          string      `yaml:"id"`
	Name        string      `yaml:"name"`
	Summary     string      `yaml:"summary"`
	Homepage    string      `yaml:"homepage"`
	Categories  []string    `yaml:"categories"`
	Platform    platform    `yaml:"platform"`
	Release     release     `yaml:"release"`
	Application application `yaml:"application"`
	Desktop     desktop     `yaml:"desktop"`
}

type platform struct {
	OS   string `yaml:"os"`
	Arch string `yaml:"arch"`
}

type release struct {
	Version string `yaml:"version"`
	URL     string `yaml:"url"`
	SHA256  string `yaml:"sha256"`
	Archive string `yaml:"archive"`
}

type application struct {
	Executable string `yaml:"executable"`
}

type desktop struct {
	Enabled    *bool    `yaml:"enabled"`
	Categories []string `yaml:"categories"`
}

type sourcePolicy struct {
	Schema  int                 `yaml:"schema"`
	Sources map[string][]string `yaml:"sources"`
}

type indexFile struct {
	Schema int         `json:"schema"`
	Apps   []indexItem `json:"apps"`
}

type indexItem struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	Categories []string `json:"categories"`
}

func TestRegistry(t *testing.T) {
	root := repositoryRoot(t)
	manifests, err := loadManifests(root)
	if err != nil {
		t.Fatal(err)
	}
	policy, err := loadPolicy(root)
	if err != nil {
		t.Fatal(err)
	}
	validatePolicy(t, manifests, policy)
	validateIndex(t, root, manifests)
}

func repositoryRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate validator source")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
}

func loadManifests(root string) (map[string]manifest, error) {
	entries, err := os.ReadDir(filepath.Join(root, "apps"))
	if err != nil {
		return nil, fmt.Errorf("read applications: %w", err)
	}
	result := make(map[string]manifest)
	names := make(map[string]string)
	for _, entry := range entries {
		if entry.Type()&os.ModeSymlink != 0 || !entry.IsDir() || !idPattern.MatchString(entry.Name()) || len(entry.Name()) > 80 {
			return nil, fmt.Errorf("invalid application directory %q", entry.Name())
		}
		directory := filepath.Join(root, "apps", entry.Name())
		children, err := os.ReadDir(directory)
		if err != nil {
			return nil, fmt.Errorf("read application %s: %w", entry.Name(), err)
		}
		if len(children) != 1 || children[0].Name() != "manifest.yaml" || children[0].IsDir() || children[0].Type()&os.ModeSymlink != 0 {
			return nil, fmt.Errorf("application directory %q must contain only manifest.yaml", entry.Name())
		}
		data, err := os.ReadFile(filepath.Join(directory, "manifest.yaml"))
		if err != nil {
			return nil, fmt.Errorf("read %s manifest: %w", entry.Name(), err)
		}
		var item manifest
		if len(data) > maxManifestBytes {
			return nil, fmt.Errorf("manifest %s exceeds %d bytes", entry.Name(), maxManifestBytes)
		}
		if err := decodeYAML(data, &item); err != nil {
			return nil, fmt.Errorf("decode %s: %w", entry.Name(), err)
		}
		if item.ID != entry.Name() {
			return nil, fmt.Errorf("manifest ID %q does not match directory %q", item.ID, entry.Name())
		}
		if _, exists := result[item.ID]; exists {
			return nil, fmt.Errorf("duplicate manifest id %q", item.ID)
		}
		if err := validateManifest(item, entry.Name()); err != nil {
			return nil, fmt.Errorf("manifest %s: %w", entry.Name(), err)
		}
		foldedName := strings.ToLower(item.Name)
		if previous, duplicate := names[foldedName]; duplicate {
			return nil, fmt.Errorf("duplicate application name %q for %s and %s", item.Name, previous, item.ID)
		}
		names[foldedName] = item.ID
		result[item.ID] = item
	}
	if len(result) == 0 {
		return nil, errors.New("no YAML manifests found")
	}
	return result, nil
}

func loadPolicy(root string) (sourcePolicy, error) {
	policyDir := filepath.Join(root, "policy")
	entries, err := os.ReadDir(policyDir)
	if err != nil {
		return sourcePolicy{}, fmt.Errorf("read source policy directory: %w", err)
	}
	if len(entries) != 1 || entries[0].Name() != "approved-sources.yaml" || entries[0].IsDir() || entries[0].Type()&os.ModeSymlink != 0 {
		return sourcePolicy{}, errors.New("policy directory must contain only approved-sources.yaml")
	}
	data, err := os.ReadFile(filepath.Join(policyDir, "approved-sources.yaml"))
	if err != nil {
		return sourcePolicy{}, fmt.Errorf("read source policy: %w", err)
	}
	if len(data) > maxManifestBytes {
		return sourcePolicy{}, fmt.Errorf("source policy exceeds %d bytes", maxManifestBytes)
	}
	var policy sourcePolicy
	if err := decodeYAML(data, &policy); err != nil {
		return sourcePolicy{}, fmt.Errorf("decode source policy: %w", err)
	}
	if policy.Schema != 1 {
		return sourcePolicy{}, fmt.Errorf("source policy schema must be 1, got %d", policy.Schema)
	}
	if len(policy.Sources) == 0 {
		return sourcePolicy{}, errors.New("source policy has no sources")
	}
	return policy, nil
}

func decodeYAML(data []byte, target any) error {
	var document yaml.Node
	if err := yaml.Unmarshal(data, &document); err != nil {
		return err
	}
	if err := validateYAMLNode(&document); err != nil {
		return err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("multiple YAML documents are not allowed")
		}
		return err
	}
	return nil
}

func validateYAMLNode(node *yaml.Node) error {
	if node == nil {
		return errors.New("invalid empty YAML node")
	}
	if node.Kind == yaml.AliasNode || node.Alias != nil || node.Anchor != "" {
		return errors.New("YAML aliases and anchors are not allowed")
	}
	allowedTags := map[string]bool{"": true, "!!map": true, "!!seq": true, "!!str": true, "!!int": true, "!!bool": true, "!!null": true}
	if !allowedTags[node.Tag] || node.Tag == "!!merge" || node.Value == "<<" {
		return fmt.Errorf("YAML tag %q or merge key is not allowed", node.Tag)
	}
	for _, child := range node.Content {
		if err := validateYAMLNode(child); err != nil {
			return err
		}
	}
	return nil
}

func validateManifest(item manifest, filenameID string) error {
	if item.Schema != 1 {
		return fmt.Errorf("schema must be 1, got %d", item.Schema)
	}
	if len(item.ID) > 80 || !idPattern.MatchString(item.ID) {
		return fmt.Errorf("id %q does not match lowercase kebab-case", item.ID)
	}
	if item.ID != filenameID {
		return fmt.Errorf("id %q does not match filename id %q", item.ID, filenameID)
	}
	if err := constrainedText(item.Name, 1, 120); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err := constrainedText(item.Summary, 1, 240); err != nil {
		return fmt.Errorf("summary: %w", err)
	}
	if err := validateHTTPS(item.Homepage); err != nil {
		return fmt.Errorf("homepage: %w", err)
	}
	if err := validateCategories(item.Categories, discoveryCategories, true); err != nil {
		return fmt.Errorf("categories: %w", err)
	}
	if item.Platform.OS != "linux" || item.Platform.Arch != "amd64" {
		return fmt.Errorf("platform must be linux/amd64, got %s/%s", item.Platform.OS, item.Platform.Arch)
	}
	if item.Release.Version == "" {
		return errors.New("release version is required")
	}
	if err := constrainedText(item.Release.Version, 1, 128); err != nil {
		return fmt.Errorf("release version: %w", err)
	}
	if strings.ContainsAny(item.Release.Version, `/\\`) || item.Release.Version == "." || item.Release.Version == ".." {
		return errors.New("release version is not filesystem-safe")
	}
	if err := validateHTTPS(item.Release.URL); err != nil {
		return fmt.Errorf("release URL: %w", err)
	}
	if !sha256Pattern.MatchString(item.Release.SHA256) {
		return errors.New("release sha256 must be 64 lowercase hexadecimal characters")
	}
	if item.Release.Archive != "tar.gz" && item.Release.Archive != "tar.xz" && item.Release.Archive != "zip" {
		return fmt.Errorf("unsupported archive %q", item.Release.Archive)
	}
	if err := validateExecutable(item.Application.Executable); err != nil {
		return fmt.Errorf("application executable: %w", err)
	}
	if item.Desktop.Enabled == nil {
		return errors.New("desktop.enabled is required")
	}
	if *item.Desktop.Enabled && len(item.Desktop.Categories) == 0 {
		return errors.New("desktop categories are required when desktop integration is enabled")
	}
	if err := validateCategories(item.Desktop.Categories, desktopCategories, false); err != nil {
		return fmt.Errorf("desktop categories: %w", err)
	}
	return nil
}

func validatePolicy(t *testing.T, manifests map[string]manifest, policy sourcePolicy) {
	t.Helper()
	if len(policy.Sources) != len(manifests) {
		t.Fatalf("source policy has %d IDs but there are %d manifests", len(policy.Sources), len(manifests))
	}
	for id, prefixes := range policy.Sources {
		if len(id) > 80 || !idPattern.MatchString(id) {
			t.Errorf("source policy ID %q is not lowercase kebab-case", id)
		}
		item, ok := manifests[id]
		if !ok {
			t.Errorf("source policy contains unknown ID %q", id)
			continue
		}
		if len(prefixes) == 0 {
			t.Errorf("source policy %q has no prefixes", id)
			continue
		}
		matched := false
		for _, prefix := range prefixes {
			_, err := validatePrefix(prefix)
			if err != nil {
				t.Errorf("source policy %q: %v", id, err)
				continue
			}
			if policyAllows(prefix, item.Release.URL) {
				matched = true
			}
		}
		if !matched {
			t.Errorf("release URL for %q is not covered by an approved source prefix", id)
		}
	}
}

func validatePrefix(prefix string) (*url.URL, error) {
	u, err := url.Parse(prefix)
	if err != nil {
		return nil, fmt.Errorf("invalid source prefix %q: %w", prefix, err)
	}
	if u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return nil, fmt.Errorf("source prefix %q must be a query-free HTTPS URL", prefix)
	}
	if !strings.HasSuffix(u.Path, "/") || u.Path == "/" {
		return nil, fmt.Errorf("source prefix %q must be a narrow non-root path ending in /", prefix)
	}
	if trimmed := strings.TrimSuffix(u.Path, "/"); path.Clean(trimmed) != trimmed {
		return nil, fmt.Errorf("source prefix %q path must be canonical", prefix)
	}
	return u, nil
}

func policyAllows(prefix, candidate string) bool {
	prefixURL, prefixErr := validatePrefix(prefix)
	candidateURL, candidateErr := url.Parse(candidate)
	if prefixErr != nil || candidateErr != nil || candidateURL.Scheme != "https" || candidateURL.User != nil || path.Clean(candidateURL.Path) != candidateURL.Path {
		return false
	}
	return strings.EqualFold(candidateURL.Host, prefixURL.Host) && strings.HasPrefix(candidateURL.EscapedPath(), prefixURL.EscapedPath())
}

func validateHTTPS(raw string) error {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.Fragment != "" {
		return fmt.Errorf("must be an HTTPS URL without user information or fragment: %q", raw)
	}
	if u.Path != "" && path.Clean(u.Path) != u.Path {
		return fmt.Errorf("URL path must be canonical: %q", raw)
	}
	return nil
}

func constrainedText(value string, minRunes, maxRunes int) error {
	if !utf8.ValidString(value) || strings.TrimSpace(value) != value {
		return errors.New("must be valid, trimmed UTF-8")
	}
	length := utf8.RuneCountInString(value)
	if length < minRunes || length > maxRunes {
		return errors.New("has an invalid length or contains control characters")
	}
	for _, character := range value {
		if unicode.IsControl(character) {
			return errors.New("has an invalid length or contains control characters")
		}
	}
	return nil
}

func validateCategories(categories []string, allowed map[string]bool, required bool) error {
	if required && len(categories) == 0 {
		return errors.New("at least one category is required")
	}
	if len(categories) == 0 {
		return nil
	}
	seen := make(map[string]bool, len(categories))
	for _, category := range categories {
		if !allowed[category] {
			return fmt.Errorf("unsupported category %q", category)
		}
		if seen[category] {
			return fmt.Errorf("duplicate category %q", category)
		}
		seen[category] = true
	}
	return nil
}

func validateExecutable(executable string) error {
	if executable == "" || !utf8.ValidString(executable) || strings.ContainsRune(executable, 0) {
		return errors.New("path is empty or invalid UTF-8")
	}
	if len(executable) > 4096 {
		return errors.New("path is too long")
	}
	if strings.Contains(executable, `\`) {
		return errors.New("backslashes are not allowed")
	}
	if strings.ContainsAny(executable, "$%") {
		return errors.New("path interpolation syntax is not allowed")
	}
	if strings.HasPrefix(executable, "/") || path.IsAbs(executable) {
		return errors.New("absolute paths are not allowed")
	}
	if len(executable) >= 2 && ((executable[0] >= 'A' && executable[0] <= 'Z') || (executable[0] >= 'a' && executable[0] <= 'z')) && executable[1] == ':' {
		return errors.New("Windows drive paths are not allowed")
	}
	if path.Clean(executable) != executable {
		return errors.New("must be canonical (no . or .. path segments)")
	}
	if strings.Count(executable, "/")+1 > 64 {
		return errors.New("path is too deep")
	}
	for _, segment := range strings.Split(executable, "/") {
		if segment == "" || segment == "." || segment == ".." {
			return errors.New("must not contain empty, . or .. path segments")
		}
	}
	for _, character := range executable {
		if unicode.IsControl(character) {
			return errors.New("path contains a control character")
		}
	}
	return nil
}

func validateIndex(t *testing.T, root string, manifests map[string]manifest) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "index", "index.json"))
	if err != nil {
		t.Fatal(fmt.Errorf("read index: %w", err))
	}
	var got indexFile
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&got); err != nil {
		t.Fatal(fmt.Errorf("decode index: %w", err))
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		t.Fatal("index must contain exactly one JSON value")
	}
	if got.Schema != 1 {
		t.Fatalf("index schema must be 1, got %d", got.Schema)
	}
	ids := make([]string, 0, len(manifests))
	for id := range manifests {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	if len(got.Apps) != len(ids) {
		t.Fatalf("index has %d apps, expected %d", len(got.Apps), len(ids))
	}
	for i, id := range ids {
		item := manifests[id]
		want := indexItem{ID: id, Name: item.Name, Version: item.Release.Version, Categories: item.Categories}
		if got.Apps[i].ID != id {
			t.Errorf("index app %d has ID %q, expected sorted ID %q", i, got.Apps[i].ID, id)
		}
		if !equalIndexItem(got.Apps[i], want) {
			t.Errorf("index app %q does not match its manifest", id)
		}
	}
	want := indexFile{Schema: 1, Apps: make([]indexItem, 0, len(ids))}
	for _, id := range ids {
		item := manifests[id]
		want.Apps = append(want.Apps, indexItem{ID: id, Name: item.Name, Version: item.Release.Version, Categories: item.Categories})
	}
	canonical, err := json.MarshalIndent(want, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	canonical = append(canonical, '\n')
	if !bytes.Equal(data, canonical) {
		t.Fatal("index/index.json is not the deterministic generated index")
	}
}

func equalIndexItem(a, b indexItem) bool {
	if a.ID != b.ID || a.Name != b.Name || a.Version != b.Version || len(a.Categories) != len(b.Categories) {
		return false
	}
	for i := range a.Categories {
		if a.Categories[i] != b.Categories[i] {
			return false
		}
	}
	return true
}
