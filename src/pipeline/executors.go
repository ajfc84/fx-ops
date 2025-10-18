package pipeline

import (
	"context"
	"fmt"
	"fx-ops/myctx"
	"fx-ops/utils"
	"fx-ops/utils/env"
	"os"
	"os/exec"
	"runtime"

	"github.com/rs/zerolog/log"
)

func ExecuteTask(ctx context.Context, task utils.TaskSpec, project, phase string) error {
	lg := log.With().Str("component", "task").Str("phase", phase).Str("type", task.Type).Str("exec", task.Exec).Logger()

	lg.Info().Msg("Starting task")

	var err error
	switch task.Type {
	case "go":
		handler, ok := utils.GoTaskHandlers[task.Exec]
		if !ok {
			err = fmt.Errorf("unknown Go function: %s", task.Exec)
			break
		}
		err = handler(ctx)

	case "shell":
		err = runShell(ctx, task.Exec)

	default:
		err = fmt.Errorf("invalid task type: %s", task.Type)
	}

	if err != nil {
		lg.Error().Err(err).Msg("Task failed")
		return err
	}

	lg.Info().Msg("Task completed successfully")
	return nil
}

func runShell(ctx context.Context, base string) error {
	envVars := myctx.Get[map[string]string](ctx, myctx.EnvVars)
	secrets := myctx.Get[map[string]string](ctx, myctx.Secrets)

	scriptPath := fmt.Sprintf("%s/%s", envVars["SUB_PROJECT_DIR"], base)

	var (
		cmd    *exec.Cmd
		script string
	)

	switch runtime.GOOS {
	case "windows":
		script = scriptPath + ".ps1"
	default:
		script = scriptPath + ".sh"
	}

	if _, err := os.Stat(script); err != nil {

		return fmt.Errorf("no shell script found for %s", script)
	}

	cmd = exec.Command(script)
	cmd.Dir = envVars["SUB_PROJECT_DIR"]
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env.FlattenEnv(env.Environ(), envVars, secrets)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run script:  %s", script)
	}

	return nil
}
