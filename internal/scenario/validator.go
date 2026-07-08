package scenario

import (
	"fmt"
	"sort"
	"strings"
)

// Validate checks scenario consistency and applies safe defaults in-place.
func Validate(s *Scenario) error {
	if s == nil {
		return fmt.Errorf("validate scenario: scenario is nil")
	}
	if len(s.Services) == 0 {
		return fmt.Errorf("validate scenario: at least one service is required")
	}
	if s.Simulation.MaxDepth < 0 {
		return fmt.Errorf("validate scenario: simulation.max_depth must be non-negative")
	}
	if s.Simulation.DefaultTimeoutMS < 0 {
		return fmt.Errorf("validate scenario: simulation.default_timeout_ms must be non-negative")
	}

	for name, svc := range s.Services {
		if strings.TrimSpace(name) == "" {
			return fmt.Errorf("validate scenario: service name cannot be empty")
		}
		if svc.FailureRate < 0 || svc.FailureRate > 1 {
			return fmt.Errorf("validate scenario: service %q failure_rate must be between 0 and 1", name)
		}
		if svc.Latency.Min < 0 || svc.Latency.Max < 0 {
			return fmt.Errorf("validate scenario: service %q latency must be non-negative", name)
		}
		if svc.Latency.Min > svc.Latency.Max {
			return fmt.Errorf("validate scenario: service %q latency min must be <= max", name)
		}
		for i := range svc.Calls {
			call := &svc.Calls[i]
			if strings.TrimSpace(call.Service) == "" {
				return fmt.Errorf("validate scenario: service %q call %d target cannot be empty", name, i)
			}
			if _, ok := s.Services[call.Service]; !ok {
				return fmt.Errorf("validate scenario: service %q calls unknown service %q", name, call.Service)
			}
			if call.TimeoutMS < 0 {
				return fmt.Errorf("validate scenario: service %q call %q timeout_ms must be non-negative", name, call.Service)
			}
			if !call.TimeoutSet && call.TimeoutMS == 0 {
				if s.Simulation.DefaultTimeoutMS <= 0 {
					return fmt.Errorf("validate scenario: service %q call %q requires timeout_ms or simulation.default_timeout_ms", name, call.Service)
				}
				call.TimeoutMS = s.Simulation.DefaultTimeoutMS
				call.TimeoutSet = true
			}
			if call.Retry.Attempts < 0 {
				return fmt.Errorf("validate scenario: service %q call %q retry.attempts must be non-negative", name, call.Service)
			}
			if call.Retry.BackoffMS < 0 {
				return fmt.Errorf("validate scenario: service %q call %q retry.backoff_ms must be non-negative", name, call.Service)
			}
		}
		s.Services[name] = svc
	}

	if unreferenced := unreferencedServices(s); len(unreferenced) > 1 {
		return fmt.Errorf("validate scenario: multiple unreferenced entry candidates: %s", strings.Join(unreferenced, ", "))
	}

	if hasCycle(s) && s.Simulation.MaxDepth == 0 {
		return fmt.Errorf("validate scenario: cycle detected; set simulation.max_depth to allow cycles")
	}
	return nil
}

// ValidateEntry checks that entry exists and all services are reachable from it.
func ValidateEntry(s *Scenario, entry string) error {
	if s == nil {
		return fmt.Errorf("validate scenario: scenario is nil")
	}
	if strings.TrimSpace(entry) == "" {
		return fmt.Errorf("validate scenario: entry service is required")
	}
	if _, ok := s.Services[entry]; !ok {
		return fmt.Errorf("validate scenario: entry service %q does not exist", entry)
	}

	reachable := map[string]bool{}
	var walk func(string)
	walk = func(name string) {
		if reachable[name] {
			return
		}
		reachable[name] = true
		for _, call := range s.Services[name].Calls {
			walk(call.Service)
		}
	}
	walk(entry)

	var missing []string
	for name := range s.Services {
		if !reachable[name] {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("validate scenario: services not reachable from entry %q: %s", entry, strings.Join(missing, ", "))
	}
	return nil
}

func unreferencedServices(s *Scenario) []string {
	incoming := map[string]int{}
	for name := range s.Services {
		incoming[name] = 0
	}
	for _, svc := range s.Services {
		for _, call := range svc.Calls {
			incoming[call.Service]++
		}
	}
	var out []string
	for name, count := range incoming {
		if count == 0 {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

func hasCycle(s *Scenario) bool {
	const (
		white = 0
		gray  = 1
		black = 2
	)
	color := map[string]int{}
	var visit func(string) bool
	visit = func(name string) bool {
		if color[name] == gray {
			return true
		}
		if color[name] == black {
			return false
		}
		color[name] = gray
		for _, call := range s.Services[name].Calls {
			if visit(call.Service) {
				return true
			}
		}
		color[name] = black
		return false
	}
	for name := range s.Services {
		if visit(name) {
			return true
		}
	}
	return false
}
