package myctx

import "context"

type ContextKey string

const (
	DockerClient ContextKey = "dockerClient"
	Args         ContextKey = "args"
	EnvVars      ContextKey = "envVars"
	Secrets      ContextKey = "secrets"
	Config       ContextKey = "config"
	SrcDir       ContextKey = "srcDir"
	Notes        ContextKey = "notes"
)

func Set(ctx context.Context, key ContextKey, val any) context.Context {
	return context.WithValue(ctx, key, val)
}

func Get[T any](ctx context.Context, key ContextKey) T {
	if v, ok := ctx.Value(key).(T); ok {
		return v
	}
	var zero T
	return zero
}
