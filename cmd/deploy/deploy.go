package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/deploy/pkg/annotation"
	"github.com/ViBiOh/deploy/pkg/api"
	"github.com/ViBiOh/httputils/v4/pkg/alcotest"
	"github.com/ViBiOh/httputils/v4/pkg/flags"
	"github.com/ViBiOh/httputils/v4/pkg/health"
	"github.com/ViBiOh/httputils/v4/pkg/httputils"
	"github.com/ViBiOh/httputils/v4/pkg/logger"
	"github.com/ViBiOh/httputils/v4/pkg/owasp"
	"github.com/ViBiOh/httputils/v4/pkg/prometheus"
	"github.com/ViBiOh/httputils/v4/pkg/server"
	"github.com/ViBiOh/mailer/pkg/client"
	mailer "github.com/ViBiOh/mailer/pkg/client"
)

func main() {
	fs := flag.NewFlagSet("deploy", flag.ExitOnError)

	appServerConfig := server.Flags(fs, "", flags.NewOverride("WriteTimeout", "2m"))
	promServerConfig := server.Flags(fs, "prometheus", flags.NewOverride("Port", 9090), flags.NewOverride("IdleTimeout", "10s"), flags.NewOverride("ShutdownTimeout", "5s"))
	healthConfig := health.Flags(fs, "")

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

	appServer := server.New(appServerConfig)
	promServer := server.New(promServerConfig)
	prometheusApp := prometheus.New(prometheusConfig)
	healthApp := health.New(healthConfig)

	mailerApp, err := mailer.New(mailerConfig)
	logger.Fatal(err)
	defer mailerApp.Close()

	annotationApp := annotation.New(annotationConfig)
	apiApp := api.New(apiConfig, mailerApp, annotationApp)

	go apiApp.Start(healthApp.Done())

	go promServer.Start("prometheus", healthApp.End(), prometheusApp.Handler())
	go appServer.Start("http", healthApp.End(), httputils.Handler(apiApp.Handler(), healthApp, prometheusApp.Middleware, owasp.New(owaspConfig).Middleware))

	healthApp.WaitForTermination(appServer.Done())
	server.GracefulWait(appServer.Done(), promServer.Done())
}
