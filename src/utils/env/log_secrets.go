package env

import "github.com/rs/zerolog/log"

func LogSecrets(secrets map[string]string) {
	for key := range secrets {
		log.Info().Str(key, "*****").Msg("secret loaded")
	}
}
