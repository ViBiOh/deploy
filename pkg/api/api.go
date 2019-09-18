package api

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/ViBiOh/httputils/v2/pkg/errors"
	"github.com/ViBiOh/httputils/v2/pkg/httperror"
	"github.com/ViBiOh/httputils/v2/pkg/logger"
	"github.com/ViBiOh/httputils/v2/pkg/tools"
	"github.com/ViBiOh/mailer/pkg/client"
)

// Config of package
type Config struct {
	tempFolder        *string
	notification      *string
	notificationEmail *string
}

// App of package
type App struct {
	tempFolder        string
	notification      string
	notificationEmail string

	mailerApp client.App
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		tempFolder:        tools.NewFlag(prefix, "deploy").Name("TempFolder").Default("/tmp").Label("Temp folder for uploading files").ToString(fs),
		notification:      tools.NewFlag(prefix, "deploy").Name("Notification").Default("onError").Label("Email notificiation when deploy ends (possibles values ares 'never', 'onError', 'all')").ToString(fs),
		notificationEmail: tools.NewFlag(prefix, "deploy").Name("NotificationEmail").Default("").Label("Email address to notify").ToString(fs),
	}
}

// New creates new App from Config
func New(config Config, mailerApp client.App) *App {
	return &App{
		tempFolder:        strings.TrimSpace(*config.tempFolder),
		notification:      strings.TrimSpace(*config.notification),
		notificationEmail: strings.TrimSpace(*config.notificationEmail),

		mailerApp: mailerApp,
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

		composeFilename := path.Join(a.tempFolder, fmt.Sprintf("docker-compose-%s.yml", project))
		uploadFile, err := os.Create(composeFilename)
		if err != nil {
			httperror.InternalServerError(w, errors.WithStack(err))
			return
		}

		if _, err := io.Copy(uploadFile, r.Body); err != nil {
			httperror.InternalServerError(w, errors.WithStack(err))
			return
		}

		cmd := exec.Command("./deploy.sh", project, composeFilename)

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		err = cmd.Run()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else if removeErr := os.Remove(composeFilename); removeErr != nil {
			logger.Error("%+s", errors.WithStack(removeErr))
		}

		output := out.Bytes()
		logger.Info("%s", output)

		if err := a.sendEmailNotification(context.Background(), project, output, err == nil); err != nil {
			logger.Error("%+s", err)
		}

		if _, err := w.Write(output); err != nil {
			httperror.InternalServerError(w, errors.WithStack(err))
		}

	})
}
