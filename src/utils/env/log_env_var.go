package env

import "github.com/rs/zerolog/log"

func LogEnvVars(environment map[string]string) {
	for key, value := range environment {
		log.Info().Str(key, value).Msg("environment variable loaded")
	}
}
