package scenario

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Scenario describes a distributed system simulation loaded from YAML.
type Scenario struct {
	Simulation SimulationConfig         `yaml:"simulation"`
	Services   map[string]ServiceConfig `yaml:"services"`
}

// SimulationConfig contains global simulation settings.
type SimulationConfig struct {
	MaxDepth         int `yaml:"max_depth"`
	DefaultTimeoutMS int `yaml:"default_timeout_ms"`
}

// ServiceConfig describes one service and its outgoing calls.
type ServiceConfig struct {
	FailureRate float64       `yaml:"failure_rate"`
	Latency     LatencyConfig `yaml:"latency_ms"`
	Calls       []CallConfig  `yaml:"calls"`
}

// LatencyConfig describes inclusive latency bounds in milliseconds.
type LatencyConfig struct {
	Min int `yaml:"min"`
	Max int `yaml:"max"`
}

// CallConfig describes a call from one service to another.
type CallConfig struct {
	Service    string      `yaml:"service"`
	TimeoutMS  int         `yaml:"timeout_ms"`
	TimeoutSet bool        `yaml:"-"`
	Retry      RetryConfig `yaml:"retry"`
}

// UnmarshalYAML decodes a call while remembering whether timeout_ms was explicitly present.
func (c *CallConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("line %d: call must be a mapping", value.Line)
	}
	*c = CallConfig{}
	for i := 0; i < len(value.Content)-1; i += 2 {
		key := value.Content[i]
		val := value.Content[i+1]
		switch key.Value {
		case "service":
			if err := val.Decode(&c.Service); err != nil {
				return err
			}
		case "timeout_ms":
			if err := val.Decode(&c.TimeoutMS); err != nil {
				return err
			}
			c.TimeoutSet = true
		case "retry":
			if err := val.Decode(&c.Retry); err != nil {
				return err
			}
		default:
			return fmt.Errorf("line %d: field %s not found in type scenario.CallConfig", key.Line, key.Value)
		}
	}
	return nil
}

// RetryConfig describes retry behavior after failed or timed-out attempts.
type RetryConfig struct {
	Attempts  int `yaml:"attempts"`
	BackoffMS int `yaml:"backoff_ms"`
}

// UnmarshalYAML decodes retry configuration with strict field validation.
func (r *RetryConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("line %d: retry must be a mapping", value.Line)
	}
	*r = RetryConfig{}
	for i := 0; i < len(value.Content)-1; i += 2 {
		key := value.Content[i]
		val := value.Content[i+1]
		switch key.Value {
		case "attempts":
			if err := val.Decode(&r.Attempts); err != nil {
				return err
			}
		case "backoff_ms":
			if err := val.Decode(&r.BackoffMS); err != nil {
				return err
			}
		default:
			return fmt.Errorf("line %d: field %s not found in type scenario.RetryConfig", key.Line, key.Value)
		}
	}
	return nil
}
