package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/deploy/pkg/api"
	httputils "github.com/ViBiOh/httputils/v2/pkg"
	"github.com/ViBiOh/httputils/v2/pkg/alcotest"
	"github.com/ViBiOh/httputils/v2/pkg/logger"
	"github.com/ViBiOh/httputils/v2/pkg/owasp"
	"github.com/ViBiOh/httputils/v2/pkg/prometheus"
	"github.com/ViBiOh/mailer/pkg/client"
)

func main() {
	fs := flag.NewFlagSet("deploy", flag.ExitOnError)

	serverConfig := httputils.Flags(fs, "")
	alcotestConfig := alcotest.Flags(fs, "")
	prometheusConfig := prometheus.Flags(fs, "prometheus")
	owaspConfig := owasp.Flags(fs, "")

	apiConfig := api.Flags(fs, "api")
	mailerConfig := client.Flags(fs, "mailer")

	logger.Fatal(fs.Parse(os.Args[1:]))

	alcotest.DoAndExit(alcotestConfig)

	prometheusApp := prometheus.New(prometheusConfig)
	owaspApp := owasp.New(owaspConfig)

	mailerApp := client.New(mailerConfig)
	apiApp := api.New(apiConfig, mailerApp)

	handler := httputils.ChainMiddlewares(apiApp.Handler(), prometheusApp, owaspApp)

	httputils.New(serverConfig).ListenAndServe(handler, httputils.HealthHandler(nil), nil)
}
