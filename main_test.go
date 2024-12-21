package main

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	cbpb "google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
)

func TestBuildMessage(t *testing.T) {
	n := new(discordNotifier)
	b := &cbpb.Build{
		ProjectId: "my-project-id",
		Id:        "some-build-id",
		Status:    cbpb.Build_SUCCESS,
		LogUrl:    "https://some.example.com/log/url?foo=bar",
		Substitutions: map[string]string{
			"_APP_NAME": "my-app",
			"_URL":      "https://some.example.com",
		},
	}

	got, err := n.buildMessage(b)
	if err != nil {
		t.Fatalf("writeMessage failed: %v", err)
	}

	want, _ := json.Marshal(discordMessage{
		Embeds: []embed{
			{Title: "✅ SUCCESS",
				Color: 1127128,
				Description: `
				Build ID: ` + b.Id + `
				Service: ` + b.Substitutions["_APP_NAME"] + `
				Environment: ` + b.ProjectId + `
				Logs: ` + b.LogUrl + `
				Access: ` + b.Substitutions["_URL"] + `
			`,
			},
		},
	})

	gotJSON, _ := json.Marshal(got)

	if diff := cmp.Diff(gotJSON, want); diff != "" {
		t.Errorf("writeMessage got unexpected diff: %s", diff)
	}
}
