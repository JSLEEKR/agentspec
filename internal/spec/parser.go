package spec

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const maxSpecSize = 1 << 20 // 1MB

// ParseFile parses a single YAML spec file.
func ParseFile(path string) (*Spec, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open spec %s: %w", path, err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat spec %s: %w", path, err)
	}
	if info.Size() > maxSpecSize {
		return nil, fmt.Errorf("spec %s exceeds 1MB size limit (%d bytes)", path, info.Size())
	}

	return Parse(f)
}

// Parse parses a spec from a reader.
func Parse(r io.Reader) (*Spec, error) {
	data, err := io.ReadAll(io.LimitReader(r, maxSpecSize+1))
	if err != nil {
		return nil, fmt.Errorf("read spec: %w", err)
	}
	if int64(len(data)) > maxSpecSize {
		return nil, fmt.Errorf("spec data exceeds 1MB size limit")
	}

	var s Spec
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse YAML: %w", err)
	}

	return &s, nil
}

// ParseDir parses all YAML spec files in a directory (top-level only, non-recursive).
func ParseDir(dir string) ([]*Spec, []string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, fmt.Errorf("read dir %s: %w", dir, err)
	}

	var specs []*Spec
	var paths []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		s, err := ParseFile(path)
		if err != nil {
			return nil, nil, fmt.Errorf("parse %s: %w", path, err)
		}
		specs = append(specs, s)
		paths = append(paths, path)
	}

	return specs, paths, nil
}
