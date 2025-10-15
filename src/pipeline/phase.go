package pipeline

import (
	"context"
	"fx-ops/myctx"
	"fx-ops/utils"

	"github.com/rs/zerolog/log"
)

func ExecutePhase(ctx context.Context, name, project string) (bool, error) {
	specsCfg := myctx.Get[*utils.SpecsData](ctx, myctx.Config)
	lg := log.With().Str("component", "phase").Str("phase", name).Str("project", project).Logger()

	lg.Info().Msg("Starting phase")

	for _, task := range specsCfg.Phases[name].Tasks {
		if err := ExecuteTask(ctx, task, project, name); err != nil {
			lg.Error().Err(err).Msg("Phase failed")
			return false, err
		}
	}

	if len(specsCfg.Phases[name].Tasks) == 0 {
		lg.Warn().Msg("No tasks found")
		return false, nil
	}

	lg.Info().Msg("Phase completed successfully")
	return true, nil
}
