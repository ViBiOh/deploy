package annotation

import (
	"context"
	"flag"
	"strings"

	"github.com/ViBiOh/httputils/v4/pkg/flags"
	"github.com/ViBiOh/httputils/v4/pkg/request"
)

// App of package
type App struct {
	req request.Request
}

// Config of package
type Config struct {
	url  *string
	user *string
	pass *string
}

type annotationPayload struct {
	Text string
	Tags []string
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		url:  flags.New(prefix, "annotation", "URL").Default("", nil).Label("URL of Annotation server (e.g. my.grafana.com/api/annotations)").ToString(fs),
		user: flags.New(prefix, "annotation", "User").Default("", nil).Label("User").ToString(fs),
		pass: flags.New(prefix, "annotation", "Pass").Default("", nil).Label("Pass").ToString(fs),
	}
}

// New creates new App from Config
func New(config Config) App {
	return App{
		req: request.New().Post(strings.TrimSpace(*config.url)).BasicAuth(strings.TrimSpace(*config.user), *config.pass),
	}
}

// Enabled check requirements are met
func (a App) Enabled() bool {
	return !a.req.IsZero()
}

// Send Grafana annotation
func (a App) Send(ctx context.Context, text string, tags ...string) error {
	if !a.Enabled() {
		return nil
	}

	_, err := a.req.JSON(ctx, annotationPayload{
		Text: text,
		Tags: tags,
	})

	return err
}
