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
	url  string
	user string
	pass string
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
		url:  flags.New(prefix, "annotation").Name("URL").Default("").Label("URL of Annotation server (e.g. my.grafana.com/api/annotations)").ToString(fs),
		user: flags.New(prefix, "annotation").Name("User").Default("").Label("User").ToString(fs),
		pass: flags.New(prefix, "annotation").Name("Pass").Default("").Label("Pass").ToString(fs),
	}
}

// New creates new App from Config
func New(config Config) App {
	url := strings.TrimSpace(*config.url)

	if url == "" {
		return App{}
	}

	return App{
		url:  strings.TrimSpace(*config.url),
		user: strings.TrimSpace(*config.user),
		pass: strings.TrimSpace(*config.pass),
	}
}

// Enabled check requirements are met
func (a App) Enabled() bool {
	return a.url != ""
}

// Send Grafana annotation
func (a App) Send(ctx context.Context, text string, tags ...string) error {
	if !a.Enabled() {
		return nil
	}

	req := request.New().Post(a.url)
	if a.pass != "" {
		req = req.BasicAuth(a.user, a.pass)
	}

	_, err := req.JSON(ctx, annotationPayload{
		Text: text,
		Tags: tags,
	})

	return err
}
