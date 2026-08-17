// Package indexgen deterministically derives the registry discovery index.
package indexgen

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"go.yaml.in/yaml/v3"
)

var idPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*$`)

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

type index struct {
	Schema int         `json:"schema"`
	Apps   []indexItem `json:"apps"`
}

type indexItem struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	Categories []string `json:"categories"`
}

// Generate returns the canonical index bytes for root's strict manifests.
func Generate(root string) ([]byte, error) {
	appsRoot := filepath.Join(root, "apps")
	entries, err := os.ReadDir(appsRoot)
	if err != nil {
		return nil, err
	}
	items := make([]indexItem, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || !idPattern.MatchString(entry.Name()) || len(entry.Name()) > 80 {
			return nil, fmt.Errorf("invalid application directory %q", entry.Name())
		}
		children, err := os.ReadDir(filepath.Join(appsRoot, entry.Name()))
		if err != nil {
			return nil, err
		}
		if len(children) != 1 || children[0].Name() != "manifest.yaml" || !children[0].Type().IsRegular() {
			return nil, fmt.Errorf("application directory %q must contain only manifest.yaml", entry.Name())
		}
		file, err := os.Open(filepath.Join(appsRoot, entry.Name(), "manifest.yaml"))
		if err != nil {
			return nil, err
		}
		data, readErr := io.ReadAll(io.LimitReader(file, 1<<20+1))
		closeErr := file.Close()
		if readErr != nil {
			return nil, readErr
		}
		if closeErr != nil {
			return nil, closeErr
		}
		if len(data) > 1<<20 {
			return nil, fmt.Errorf("manifest %s exceeds 1 MiB", entry.Name())
		}
		decoder := yaml.NewDecoder(bytes.NewReader(data))
		decoder.KnownFields(true)
		var item manifest
		decodeErr := decoder.Decode(&item)
		var trailing any
		trailingErr := decoder.Decode(&trailing)
		if decodeErr != nil {
			return nil, fmt.Errorf("decode %s: %w", entry.Name(), decodeErr)
		}
		if !errors.Is(trailingErr, io.EOF) {
			return nil, fmt.Errorf("manifest %s must contain one document", entry.Name())
		}
		if item.Schema != 1 || item.ID != entry.Name() || item.Name == "" || item.Release.Version == "" || len(item.Categories) == 0 {
			return nil, fmt.Errorf("manifest %s lacks required index metadata", entry.Name())
		}
		items = append(items, indexItem{
			ID: item.ID, Name: item.Name, Version: item.Release.Version,
			Categories: append([]string(nil), item.Categories...),
		})
	}
	if len(items) == 0 {
		return nil, errors.New("registry contains no manifests")
	}
	sort.Slice(items, func(left, right int) bool { return items[left].ID < items[right].ID })
	data, err := json.MarshalIndent(index{Schema: 1, Apps: items}, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func Check(root string) error {
	expected, err := Generate(root)
	if err != nil {
		return err
	}
	actual, err := os.ReadFile(filepath.Join(root, "index", "index.json"))
	if err != nil {
		return err
	}
	if !bytes.Equal(actual, expected) {
		return errors.New("index/index.json is stale; run `go run ./cmd/generate-index`")
	}
	return nil
}

func Write(root string) error {
	data, err := Generate(root)
	if err != nil {
		return err
	}
	directory := filepath.Join(root, "index")
	temporary, err := os.CreateTemp(directory, ".index-*.json")
	if err != nil {
		return err
	}
	name := temporary.Name()
	keep := false
	defer func() {
		_ = temporary.Close()
		if !keep {
			_ = os.Remove(name)
		}
	}()
	if err := temporary.Chmod(0o644); err != nil {
		return err
	}
	if _, err := temporary.Write(data); err != nil {
		return err
	}
	if err := temporary.Sync(); err != nil {
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Rename(name, filepath.Join(directory, "index.json")); err != nil {
		return err
	}
	opened, err := os.Open(directory)
	if err != nil {
		return err
	}
	syncErr := opened.Sync()
	closeErr := opened.Close()
	if syncErr != nil {
		return syncErr
	}
	if closeErr != nil {
		return closeErr
	}
	keep = true
	return nil
}
