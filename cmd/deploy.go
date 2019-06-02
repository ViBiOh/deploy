package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/deploy/pkg/api"
	httputils "github.com/ViBiOh/httputils/pkg"
	"github.com/ViBiOh/httputils/pkg/alcotest"
	"github.com/ViBiOh/httputils/pkg/gzip"
	"github.com/ViBiOh/httputils/pkg/healthcheck"
	"github.com/ViBiOh/httputils/pkg/logger"
	"github.com/ViBiOh/httputils/pkg/opentracing"
	"github.com/ViBiOh/httputils/pkg/owasp"
	"github.com/ViBiOh/httputils/pkg/prometheus"
	"github.com/ViBiOh/httputils/pkg/server"
	"github.com/ViBiOh/mailer/pkg/client"
)

func main() {
	fs := flag.NewFlagSet("deploy", flag.ExitOnError)

	serverConfig := httputils.Flags(fs, "")
	alcotestConfig := alcotest.Flags(fs, "")
	prometheusConfig := prometheus.Flags(fs, "prometheus")
	opentracingConfig := opentracing.Flags(fs, "tracing")
	owaspConfig := owasp.Flags(fs, "")

	apiConfig := api.Flags(fs, "api")
	mailerConfig := client.Flags(fs, "mailer")

	if err := fs.Parse(os.Args[1:]); err != nil {
		logger.Fatal("%#v", err)
	}

	alcotest.DoAndExit(alcotestConfig)

	serverApp, err := httputils.New(serverConfig)
	if err != nil {
		logger.Fatal("%#v", err)
	}

	healthcheckApp := healthcheck.New()
	prometheusApp := prometheus.New(prometheusConfig)
	opentracingApp := opentracing.New(opentracingConfig)
	gzipApp := gzip.New()
	owaspApp := owasp.New(owaspConfig)

	mailerApp := client.New(mailerConfig)
	apiApp := api.New(apiConfig, mailerApp)

	handler := server.ChainMiddlewares(apiApp.Handler(), prometheusApp, opentracingApp, gzipApp, owaspApp)

	serverApp.ListenAndServe(handler, nil, healthcheckApp)
}
