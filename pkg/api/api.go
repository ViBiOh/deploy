package api

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/ViBiOh/httputils/pkg/httperror"
	"github.com/ViBiOh/httputils/pkg/tools"
)

// Config of package
type Config struct {
	tempFolder *string
}

// App of package
type App struct {
	tempFolder string
}

// Flags adds flags for configuring package
func Flags(fs *flag.FlagSet, prefix string) Config {
	return Config{
		tempFolder: fs.String(tools.ToCamel(fmt.Sprintf("%sTempFolder", prefix)), "/tmp", "[api] Temp folder for uploading files"),
	}
}

// New creates new App from Config
func New(config Config) *App {
	return &App{
		tempFolder: strings.TrimSpace(*config.tempFolder),
	}
}

// Handler for request. Should be use with net/http
func (a App) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		projectName := r.URL.Query().Get("project")
		if projectName == "" {
			httperror.BadRequest(w, errors.New("project name is missing"))
			return
		}

		versionSha1 := r.URL.Query().Get("version")
		if versionSha1 == "" {
			httperror.BadRequest(w, errors.New("version sha1 is missing"))
			return
		}

		composeFilename := path.Join(a.tempFolder, fmt.Sprintf("docker-compose-%s.yml", versionSha1))
		uploadFile, err := os.Create(composeFilename)
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}

		if _, err := io.Copy(uploadFile, r.Body); err != nil {
			httperror.InternalServerError(w, err)
			return
		}

		cmd := exec.Command("./deploy.sh", projectName, versionSha1, composeFilename)

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write(out.Bytes())
	})
}
