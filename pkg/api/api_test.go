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
		wantVersion string
		wantErr     error
	}{
		{
			"get empty",
			httptest.NewRequest(http.MethodPost, "/", nil),
			"",
			"",
			errors.New("project name is required"),
		},
		{
			"missing version",
			httptest.NewRequest(http.MethodPost, "/deploy", nil),
			"deploy",
			"",
			errors.New("version sha is required"),
		},
		{
			"non-trailing slash",
			httptest.NewRequest(http.MethodPost, "/deploy/1234abcd", nil),
			"deploy",
			"1234abcd",
			nil,
		},
		{
			"trailing slash",
			httptest.NewRequest(http.MethodPost, "/deploy/1234abcd/", nil),
			"deploy",
			"1234abcd",
			nil,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.intention, func(t *testing.T) {
			project, version, err := validateRequest(testCase.input)

			failed := false

			if err == nil && testCase.wantErr != nil {
				failed = true
			} else if err != nil && testCase.wantErr == nil {
				failed = true
			} else if err != nil && err.Error() != testCase.wantErr.Error() {
				failed = true
			} else if project != testCase.wantProject {
				failed = true
			} else if version != testCase.wantVersion {
				failed = true
			}

			if failed {
				t.Errorf("ValidateRequest(%#v) = (`%s`, `%s`, %#v), want (`%s`, `%s`, %#v)", testCase.input, project, version, err, testCase.wantProject, testCase.wantVersion, testCase.wantErr)
			}
		})
	}
}
