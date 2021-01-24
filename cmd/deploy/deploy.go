package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/deploy/pkg/annotation"
	"github.com/ViBiOh/deploy/pkg/api"
	"github.com/ViBiOh/httputils/v3/pkg/alcotest"
	"github.com/ViBiOh/httputils/v3/pkg/flags"
	"github.com/ViBiOh/httputils/v3/pkg/httputils"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
	"github.com/ViBiOh/httputils/v3/pkg/owasp"
	"github.com/ViBiOh/httputils/v3/pkg/prometheus"
	"github.com/ViBiOh/mailer/pkg/client"
	mailer "github.com/ViBiOh/mailer/pkg/client"
)

func main() {
	fs := flag.NewFlagSet("deploy", flag.ExitOnError)

	serverConfig := httputils.Flags(fs, "", flags.NewOverride("WriteTimeout", "2m"))
	alcotestConfig := alcotest.Flags(fs, "")
	loggerConfig := logger.Flags(fs, "logger")
	prometheusConfig := prometheus.Flags(fs, "prometheus")
	owaspConfig := owasp.Flags(fs, "")

	apiConfig := api.Flags(fs, "api")
	mailerConfig := client.Flags(fs, "mailer")
	annotationConfig := annotation.Flags(fs, "annotation")

	logger.Fatal(fs.Parse(os.Args[1:]))

	alcotest.DoAndExit(alcotestConfig)
	logger.Global(logger.New(loggerConfig))
	defer logger.Close()

	mailerApp, err := mailer.New(mailerConfig)
	logger.Fatal(err)
	defer mailerApp.Close()

	annotationApp := annotation.New(annotationConfig)
	apiApp := api.New(apiConfig, mailerApp, annotationApp)

	server := httputils.New(serverConfig)
	go apiApp.Start(server.GetDone())

	server.ListenAndServe(apiApp.Handler(), nil, prometheus.New(prometheusConfig).Middleware, owasp.New(owaspConfig).Middleware)
}
