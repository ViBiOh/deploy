package api

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/ViBiOh/deploy/pkg/annotation"
	"github.com/ViBiOh/httputils/v4/pkg/cron"
	"github.com/ViBiOh/httputils/v4/pkg/flags"
	"github.com/ViBiOh/httputils/v4/pkg/httperror"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/mailer/pkg/client"
)

//go:embed scripts
var scripts embed.FS

// App of package
type App struct {
	mailerApp     client.App
	annotationApp annotation.App

	tempFolder        string
	notification      string
	notificationEmail string
}

// Config of package
type Config struct {
	tempFolder        *string
	notification      *string
	notificationEmail *string
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		tempFolder:        flags.New(prefix, "deploy", "TempFolder").Default("/tmp", nil).Label("Temp folder for uploading files").ToString(fs),
		notification:      flags.New(prefix, "deploy", "Notification").Default("onError", nil).Label("Email notificiation when deploy ends (possibles values ares 'never', 'onError', 'all')").ToString(fs),
		notificationEmail: flags.New(prefix, "deploy", "NotificationEmail").Default("", nil).Label("Email address to notify").ToString(fs),
	}
}

// New creates new App from Config
func New(config Config, mailerApp client.App, annotationApp annotation.App) App {
	return App{
		tempFolder:        strings.TrimSpace(*config.tempFolder),
		notification:      strings.TrimSpace(*config.notification),
		notificationEmail: strings.TrimSpace(*config.notificationEmail),

		mailerApp:     mailerApp,
		annotationApp: annotationApp,
	}
}

func validateRequest(r *http.Request) (project string, err error) {
	project = strings.TrimSpace(strings.Trim(r.URL.Path, "/"))
	if project == "" {
		err = errors.New("project name is required")
		return
	}

	return
}

func copyEmbedScript(source, name string) error {
	destinationFile, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0700)
	if err != nil {
		return fmt.Errorf("unable to open destination file: %s", err)
	}
	defer func() {
		if err := destinationFile.Close(); err != nil {
			logger.Error("error while closing destination file: %s", err)
		}
	}()

	sourceFile, err := scripts.Open(path.Join(source, name))
	if err != nil {
		return fmt.Errorf("unable to open source file: %s", err)
	}
	defer func() {
		if err := sourceFile.Close(); err != nil {
			logger.Error("error while closing source file: %s", err)
		}
	}()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

// Start worker
func (a App) Start(done <-chan struct{}) {
	for _, script := range []string{"clean", "deploy-compose"} {
		if err := copyEmbedScript("scripts", script); err != nil {
			logger.Error("unable to copy embed `%s` script: %s", script, err)
		}
	}

	cron.New().Days().At("06:00").In("Europe/Paris").OnError(func(err error) {
		logger.Error("%s", err)
	}).Start(func(_ context.Context) error {
		cmd := exec.Command("./clean")

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		err := cmd.Run()
		logger.Info("%s", out.Bytes())

		return err
	}, done)
}

// Handler for request. Should be use with net/http
func (a App) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		project, err := validateRequest(r)
		if err != nil {
			httperror.BadRequest(w, err)
			return
		}

		composeFilename := path.Join(a.tempFolder, fmt.Sprintf("docker-compose-%s.yaml", project))
		uploadFile, err := os.Create(composeFilename)
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}

		if _, err := io.Copy(uploadFile, r.Body); err != nil {
			httperror.InternalServerError(w, err)
			return
		}

		cmd := exec.Command("./deploy-compose", project, composeFilename)

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		err = cmd.Run()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else if removeErr := os.Remove(composeFilename); removeErr != nil {
			logger.Error("%s", removeErr)
		}

		output := out.Bytes()
		logger.Info("%s", output)

		a.notify(project, output, err)

		if _, err := w.Write(output); err != nil {
			httperror.InternalServerError(w, err)
		}
	})
}

func (a App) notify(project string, output []byte, err error) {
	success := err == nil

	if err := a.sendEmailNotification(context.Background(), project, output, success); err != nil {
		logger.Error("%s", err)
	}

	if err := a.sendAnnotation(context.Background(), project, success); err != nil {
		logger.Error("%s", err)
	}
}
