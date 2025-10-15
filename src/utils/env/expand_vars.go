package env

import (
	"os"
	"regexp"
)

func ExpandVars(vars map[string]string) map[string]string {
	resolved := make(map[string]string, len(vars))
	refPattern := regexp.MustCompile(`\$\{([^}]+)\}`)

	for key, val := range vars {
		resolved[key] = refPattern.ReplaceAllStringFunc(val, func(match string) string {
			name := refPattern.FindStringSubmatch(match)[1]

			// 1. try resolved vars (already expanded)
			if v, ok := resolved[name]; ok {
				return v
			}

			// 2. try vars map (raw values)
			if v, ok := vars[name]; ok {
				return v
			}

			// 3. fallback to OS environment
			if v := os.Getenv(name); v != "" {
				return v
			}

			// 4. leave placeholder unchanged if unknown
			return match
		})
	}

	return resolved
}
