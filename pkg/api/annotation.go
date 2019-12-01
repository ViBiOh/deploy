package api

import (
	"context"
	"fmt"
)

func (a app) sendAnnotation(ctx context.Context, project string, success bool) error {
	text := fmt.Sprintf("Deploy of %s", project)
	if success {
		text += " successful"
	} else {
		text += " failed"
	}

	return a.annotationApp.Send(ctx, text, "deploy", project)
}
