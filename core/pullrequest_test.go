package core

import (
	"errors"
	"testing"
)

func TestGithubPRService_GetPRStatus(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		mockResp   PullRequest
		mockErr    error
		wantMerged bool
		wantErr    bool
	}{
		{
			name: "Merged PR",
			url:  "https://api.github.com/repos/owner/repo/pulls/1",
			mockResp: PullRequest{
				Merged: true,
				Title:  "Test PR",
			},
			wantMerged: true,
		},
		{
			name: "Unmerged PR",
			url:  "https://api.github.com/repos/owner/repo/pulls/2",
			mockResp: PullRequest{
				Merged: false,
				Title:  "Test PR",
			},
			wantMerged: false,
		},
		{
			name:    "API Error",
			url:     "https://api.github.com/repos/owner/repo/pulls/3",
			mockErr: errors.New("API error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &mockGithubClient{
				getFunc: func(url string, response interface{}) error {
					if tt.mockErr != nil {
						return tt.mockErr
					}
					pr := response.(*PullRequest)
					*pr = tt.mockResp
					return nil
				},
			}

			service := NewGithubPRService(client)
			merged, err := service.GetPRStatus(tt.url)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPRStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && merged != tt.wantMerged {
				t.Errorf("GetPRStatus() = %v, want %v", merged, tt.wantMerged)
			}
		})
	}
}

func TestGithubPRService_GetPRStatus_URLFormatting(t *testing.T) {
	client := &mockGithubClient{
		getFunc: func(url string, response interface{}) error {
			expected := "repos/owner/repo/pulls/1"
			if url != expected {
				t.Errorf("Expected URL %s, got %s", expected, url)
			}
			pr := response.(*PullRequest)
			*pr = PullRequest{Merged: true}
			return nil
		},
	}

	service := NewGithubPRService(client)
	_, err := service.GetPRStatus("https://api.github.com/repos/owner/repo/pulls/1")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
