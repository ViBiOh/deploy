package api

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/ViBiOh/mailer/pkg/client"
)

type outputLine struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

const (
	never   = "never"
	onError = "onError"
	all     = "all"
)

var (
	green = regexp.MustCompile(`(?m)\[0;32m(.*?)\[0m`)
	red   = regexp.MustCompile(`(?m)\[0;33m(.*?)\[0m`)
	blue  = regexp.MustCompile(`(?m)\[0;34m(.*?)\[0m`)
)

func formatLines(output []byte) []outputLine {
	lines := strings.Split(string(output), "\n")
	outputLines := make([]outputLine, len(lines))

	for index, line := range lines {
		value := line
		color := ""

		if match := green.MatchString(value); match {
			value = green.ReplaceAllString(value, "${1}")
			color = "limegreen"
		}

		if match := red.MatchString(value); match {
			value = red.ReplaceAllString(value, "${1}")
			color = "crimson"
		}

		if match := blue.MatchString(value); match {
			value = blue.ReplaceAllString(value, "${1}")
			color = "royalblue"
		}

		outputLines[index] = outputLine{
			Value: value,
			Color: color,
		}
	}

	return outputLines
}

func (a *App) sendEmailNotification(ctx context.Context, project string, output []byte, success bool) error {
	if a.notification == never || (success && a.notification == onError) {
		return nil
	}

	notificationContent := map[string]interface{}{
		"success": success,
		"app":     project,
		"output":  formatLines(output),
	}

	recipients := []string{a.notificationEmail}

	return client.NewEmail(a.mailerApp).From("deploy@vibioh.fr").As("Deploy").WithSubject(fmt.Sprintf("[deploy] Deploy of %s", project)).Data(notificationContent).To(recipients...).Send(ctx)
}
