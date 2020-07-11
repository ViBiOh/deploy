package annotation

import (
	"context"
	"flag"
	"strings"

	"github.com/ViBiOh/httputils/v3/pkg/flags"
	"github.com/ViBiOh/httputils/v3/pkg/request"
)

// App of package
type App interface {
	Send(context.Context, string, ...string) error
}

// Config of package
type Config struct {
	url  *string
	user *string
	pass *string
}

type app struct {
	url  string
	user string
	pass string
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
		return &app{}
	}

	return &app{
		url:  strings.TrimSpace(*config.url),
		user: strings.TrimSpace(*config.user),
		pass: strings.TrimSpace(*config.pass),
	}
}

func (a app) Enabled() bool {
	return a.url != ""
}

func (a app) Send(ctx context.Context, text string, tags ...string) error {
	if !a.Enabled() {
		return nil
	}

	req := request.New().Post(a.url)
	if a.pass != "" {
		req.BasicAuth(a.user, a.pass)
	}

	_, err := req.JSON(ctx, annotationPayload{
		Text: text,
		Tags: tags,
	})

	return err
}
