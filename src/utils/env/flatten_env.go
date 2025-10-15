package env

import "fmt"

func FlattenEnv(maps ...map[string]string) []string {
	var result []string
	for _, m := range maps {
		for k, v := range m {
			result = append(result, fmt.Sprintf("%s=%s", k, v))
		}
	}
	return result
}
