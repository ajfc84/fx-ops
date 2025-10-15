package utils

func GetRelease(envName string) string {
	switch envName {
	case "develop", "ops":
		return "alpha"
	case "qa":
		return "beta"
	case "uat", "nonprod":
		return "rc"
	default:
		return "main"
	}
}
