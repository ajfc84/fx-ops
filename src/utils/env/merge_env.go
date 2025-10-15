package env

func MergeEnv(base, layer map[string]string) map[string]string {
	out := make(map[string]string, len(base))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range layer {
		out[k] = v
	}
	return out
}
