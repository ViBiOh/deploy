package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateRequest(t *testing.T) {
	var cases = []struct {
		intention   string
		input       *http.Request
		wantProject string
		wantErr     error
	}{
		{
			"get empty",
			httptest.NewRequest(http.MethodPost, "/", nil),
			"",
			errors.New("project name is required"),
		},
		{
			"missing version",
			httptest.NewRequest(http.MethodPost, "/deploy", nil),
			"deploy",
			nil,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			project, err := validateRequest(testCase.input)

			failed := false

			if err == nil && testCase.wantErr != nil {
				failed = true
			} else if err != nil && testCase.wantErr == nil {
				failed = true
			} else if err != nil && err.Error() != testCase.wantErr.Error() {
				failed = true
			} else if project != testCase.wantProject {
				failed = true
			}

			if failed {
				t.Errorf("ValidateRequest(%#v) = (`%s`, %#v), want (`%s`, %#v)", testCase.input, project, err, testCase.wantProject, testCase.wantErr)
			}
		})
	}
}
