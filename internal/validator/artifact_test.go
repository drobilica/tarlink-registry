package validator

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

const maximumArtifactBytes int64 = 8 << 30

// TestBlenderArtifact is opt-in because the reviewed upstream binary is
// hundreds of megabytes. The path-filtered artifact workflow enables it only
// when Blender's manifest or approved-source policy changes, then passes the
// verified bytes to TarLink's real safe-extraction acceptance test.
func TestBlenderArtifact(t *testing.T) {
	destination := os.Getenv("TARLINK_VERIFY_BLENDER_OUTPUT")
	if destination == "" {
		t.Skip("large upstream artifact verification is opt-in")
	}
	manifests, err := loadManifests(repositoryRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	item, exists := manifests["blender"]
	if !exists {
		t.Fatal("Blender manifest is missing")
	}
	parsed, err := url.Parse(item.Release.URL)
	if err != nil || parsed.Scheme != "https" {
		t.Fatalf("invalid Blender artifact URL %q", item.Release.URL)
	}
	policy, err := loadPolicy(repositoryRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	allowed := func(candidate *url.URL) bool {
		for _, prefix := range policy.Sources["blender"] {
			if policyAllows(prefix, candidate.String()) {
				return true
			}
		}
		return false
	}
	if !allowed(parsed) {
		t.Fatal("Blender artifact URL is outside approved sources")
	}

	dialer := &net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment, DialContext: dialer.DialContext,
			TLSClientConfig:     &tls.Config{MinVersion: tls.VersionTLS12},
			TLSHandshakeTimeout: 30 * time.Second, ResponseHeaderTimeout: 30 * time.Second,
		},
		Timeout: 2 * time.Hour,
		CheckRedirect: func(request *http.Request, via []*http.Request) error {
			if len(via) > 5 {
				return errors.New("redirect limit exceeded")
			}
			if request.URL.Scheme != "https" {
				return errors.New("HTTPS downgrade rejected")
			}
			if !allowed(request.URL) {
				return errors.New("redirect escaped approved source policy")
			}
			return nil
		},
	}
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, item.Release.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Accept-Encoding", "identity")
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("artifact HTTP status %d", response.StatusCode)
	}
	if response.ContentLength > maximumArtifactBytes {
		t.Fatal("artifact exceeds the 8 GiB client limit")
	}

	file, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	keep := false
	defer func() {
		_ = file.Close()
		if !keep {
			_ = os.Remove(destination)
		}
	}()
	hasher := sha256.New()
	written, err := io.Copy(io.MultiWriter(file, hasher), io.LimitReader(response.Body, maximumArtifactBytes+1))
	if err != nil {
		t.Fatal(err)
	}
	if written > maximumArtifactBytes {
		t.Fatal("artifact exceeds the 8 GiB client limit")
	}
	actual := hex.EncodeToString(hasher.Sum(nil))
	if actual != item.Release.SHA256 {
		t.Fatalf("artifact SHA-256 = %s, manifest declares %s", actual, item.Release.SHA256)
	}
	if err := file.Sync(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	keep = true
	t.Log(fmt.Sprintf("verified %d bytes at %s", written, destination))
}
