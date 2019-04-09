package api

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

const (
	never   = "never"
	onError = "onError"
	all     = "all"
)

var (
	green = regexp.MustCompile(`(?m)\[0;32m(.*?)\[0m`)
	blue  = regexp.MustCompile(`(?m)\[0;34m(.*?)\[0m`)
)

func (a *App) sendEmailNotification(ctx context.Context, project string, output []byte, success bool) error {
	if a.notification == never || (success && a.notification == onError) {
		return nil
	}

	content := string(output)
	content = green.ReplaceAllString(content, "<span style=\"color: limegreen\">${1}</span>")
	content = blue.ReplaceAllString(content, "<span style=\"color: royalblue\">${1}</span>")
	finalOutput := strings.Split(content, "\n")

	notificationContent := map[string]interface{}{
		"success": success,
		"app":     project,
		"output":  finalOutput,
	}

	recipients := []string{a.notificationEmail}

	if err := a.mailerApp.SendEmail(ctx, "deploy", "deploy@vibioh.fr", "Deploy", fmt.Sprintf("[deploy] Deploy of %s", project), recipients, notificationContent); err != nil {
		return err
	}

	return nil
}
