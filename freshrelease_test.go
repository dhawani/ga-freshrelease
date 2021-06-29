package main

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetIssue(t *testing.T) {
	sampleIssueJson := `{
    "issue": {
        "id": 797579,
        "key": "ED-50",
        "title": "Fix something awesome",
        "name": "Fix something awesome",
        "status_id": 18
    },
    "statuses": [
        {
            "label": "Open",
            "description": "This is the description for Open",
            "id": 18,
            "name": "open"
        },
        {
            "label": "In Progress",
            "description": "This is the description for Open",
            "id": 19,
            "name": "in_progress"
        }
    ]
}`
	sampleIssue := &IssueResponse{}
	err := json.Unmarshal([]byte(sampleIssueJson), sampleIssue)
	require.NoError(t, err)
	type args struct {
		token      string
		projectKey string
		issueKey   string
	}
	tests := []struct {
		name    string
		args    args
		want    *IssueResponse
		wantErr bool
		handler func(rw http.ResponseWriter, req *http.Request)
	}{
		{
			name: "invalid token",
			args: args{
				token:      "invalid-token",
				projectKey: "",
				issueKey:   "",
			},
			wantErr: true,
			handler: func(rw http.ResponseWriter, req *http.Request) {
				require.Equal(t, "Token invalid-token", req.Header.Get("Authorization"))
				rw.WriteHeader(http.StatusUnauthorized)
			},
		},
		{
			name: "valid request",
			args: args{
				token:      "token",
				projectKey: "VD",
				issueKey:   "VD-50",
			},
			wantErr: false,
			handler: func(rw http.ResponseWriter, req *http.Request) {
				require.Equal(t, "/VD/issues/VD-50", req.URL.String())
				require.Equal(t, "Token token", req.Header.Get("Authorization"))
				_, _ = rw.Write([]byte(sampleIssueJson))
			},
			want: sampleIssue,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(tt.handler))
			defer server.Close()
			got, err := GetIssue(server.URL, tt.args.token, tt.args.projectKey, tt.args.issueKey)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
