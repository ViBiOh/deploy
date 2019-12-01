package main

import (
	"flag"
	"os"

	"github.com/ViBiOh/deploy/pkg/annotation"
	"github.com/ViBiOh/deploy/pkg/api"
	"github.com/ViBiOh/httputils/v3/pkg/alcotest"
	"github.com/ViBiOh/httputils/v3/pkg/httputils"
	"github.com/ViBiOh/httputils/v3/pkg/logger"
	"github.com/ViBiOh/httputils/v3/pkg/owasp"
	"github.com/ViBiOh/httputils/v3/pkg/prometheus"
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
	annotationConfig := annotation.Flags(fs, "annotation")

	logger.Fatal(fs.Parse(os.Args[1:]))

	alcotest.DoAndExit(alcotestConfig)

	mailerApp := client.New(mailerConfig)
	annotationApp := annotation.New(annotationConfig)
	apiApp := api.New(apiConfig, mailerApp, annotationApp)

	server := httputils.New(serverConfig)
	server.Middleware(prometheus.New(prometheusConfig))
	server.Middleware(owasp.New(owaspConfig))
	server.ListenServeWait(apiApp.Handler())
}
