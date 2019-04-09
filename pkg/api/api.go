package api

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/ViBiOh/httputils/pkg/errors"
	"github.com/ViBiOh/httputils/pkg/httperror"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/ViBiOh/httputils/pkg/tools"
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

	mailerApp *client.App
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		tempFolder:        fs.String(tools.ToCamel(fmt.Sprintf("%sTempFolder", prefix)), "/tmp", "[api] Temp folder for uploading files"),
		notification:      fs.String(tools.ToCamel(fmt.Sprintf("%sNotification", prefix)), "all", "[api] Email notificiation when deploy ends (possibles values ares 'never', 'onError', 'all')"),
		notificationEmail: fs.String(tools.ToCamel(fmt.Sprintf("%sNotificationEmail", prefix)), "", "[api] Email address to notify"),
	}
}

// New creates new App from Config
func New(config Config, mailerApp *client.App) *App {
	return &App{
		tempFolder:        strings.TrimSpace(*config.tempFolder),
		notification:      strings.TrimSpace(*config.notification),
		notificationEmail: strings.TrimSpace(*config.notificationEmail),

		mailerApp: mailerApp,
	}
}

func validateRequest(r *http.Request) (string, string, error) {
	args := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/"), "/", 2)

	project := strings.TrimSpace(args[0])
	version := strings.TrimSpace(args[1])

	if project == "" {
		return "", "", errors.New("project name is required")
	}

	if version == "" {
		return "", "", errors.New("version sha is required")
	}

	return project, version, nil
}

// Handler for request. Should be use with net/http
func (a App) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		project, version, err := validateRequest(r)
		if err != nil {
			httperror.BadRequest(w, err)
			return
		}

		composeFilename := path.Join(a.tempFolder, fmt.Sprintf("docker-compose-%s-%s.yml", project, version))
		uploadFile, err := os.Create(composeFilename)
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}

		if _, err := io.Copy(uploadFile, r.Body); err != nil {
			httperror.InternalServerError(w, err)
			return
		}

		cmd := exec.Command("./deploy.sh", project, version, composeFilename)

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else if err := os.Remove(composeFilename); err != nil {
			logger.Error("%+s", errors.WithStack(err))
		}

		if _, err := w.Write(out.Bytes()); err != nil {
			httperror.InternalServerError(w, err)
		}
	})
}
