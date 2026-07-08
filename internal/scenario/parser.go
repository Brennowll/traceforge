package scenario

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadFile reads, parses, validates, and applies defaults to a YAML scenario file.
func LoadFile(path string) (*Scenario, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read scenario file: %w", err)
	}
	return Parse(data)
}

// Parse parses YAML bytes into a Scenario, then validates and applies defaults.
func Parse(data []byte) (*Scenario, error) {
	if err := detectDuplicateServiceNames(data); err != nil {
		return nil, err
	}

	var s Scenario
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&s); err != nil {
		return nil, fmt.Errorf("parse scenario yaml: %w", err)
	}
	if err := Validate(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

func detectDuplicateServiceNames(data []byte) error {
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil // Parse will return the real YAML error later.
	}
	if len(root.Content) == 0 || root.Content[0].Kind != yaml.MappingNode {
		return nil
	}
	mapping := root.Content[0]
	for i := 0; i < len(mapping.Content)-1; i += 2 {
		key := mapping.Content[i]
		value := mapping.Content[i+1]
		if key.Value != "services" || value.Kind != yaml.MappingNode {
			continue
		}
		seen := map[string]struct{}{}
		for j := 0; j < len(value.Content)-1; j += 2 {
			name := value.Content[j].Value
			if _, ok := seen[name]; ok {
				return fmt.Errorf("validate scenario: duplicate service name %q", name)
			}
			seen[name] = struct{}{}
		}
	}
	return nil
}
