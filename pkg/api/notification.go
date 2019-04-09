package api

import (
	"context"
	"fmt"
	"strings"
)

const (
	never   = "never"
	onError = "onError"
	all     = "all"
)

func (a *App) sendEmailNotification(ctx context.Context, project string, output []byte, success bool) error {
	if a.notification == never || (success && a.notification == onError) {
		return nil
	}

	content := strings.Split(string(output), "\n")

	notificationContent := map[string]interface{}{
		"success": success,
		"app":     project,
		"output":  content,
	}

	recipients := []string{a.notificationEmail}

	if err := a.mailerApp.SendEmail(ctx, "deploy", "deploy@vibioh.fr", "Deploy", fmt.Sprintf("[deploy] Deploy of %s", project), recipients, notificationContent); err != nil {
		return err
	}

	return nil
}
