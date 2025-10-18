package pipeline

import (
	"context"
	"fx-ops/myctx"
	"fx-ops/utils"

	"github.com/rs/zerolog/log"
)

func ExecuteStage(ctx context.Context, name, project string) (bool, error) {
	specsCfg := myctx.Get[*utils.SpecsData](ctx, myctx.Config)
	lg := log.With().Str("component", "stage").Str("stage", name).Str("project", project).Logger()

	lg.Info().Msg("Starting stage")

	for _, task := range specsCfg.Stages[name].Tasks {
		if err := ExecuteTask(ctx, task, project, name); err != nil {
			lg.Error().Err(err).Msg("stage failed")
			return false, err
		}
	}

	if len(specsCfg.Stages[name].Tasks) == 0 {
		lg.Warn().Msg("No tasks found")
		return false, nil
	}

	lg.Info().Msg("Stage completed successfully")
	return true, nil
}
