package core

import (
	"testing"
)

type mockGithubClient struct {
	getFunc    func(url string, response interface{}) error
	deleteFunc func(url string, response interface{}) error
}

func (m *mockGithubClient) Get(url string, response interface{}) error {
	return m.getFunc(url, response)
}

func (m *mockGithubClient) Delete(url string, response interface{}) error {
	return m.deleteFunc(url, response)
}

func TestGithubRepository_GetByTimePeriod(t *testing.T) {
	notifications := []Notification{
		{ID: "1", Subject: Subject{Title: "PR 1", Type: "PullRequest"}},
		{ID: "2", Subject: Subject{Title: "PR 2", Type: "PullRequest"}},
	}

	client := &mockGithubClient{
		getFunc: func(url string, response interface{}) error {
			resp := response.(*[]Notification)
			*resp = notifications
			return nil
		},
	}

	repo := NewGithubRepository(client)
	result, err := repo.GetByTimePeriod("7d")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(result))
	}
}

func TestGithubRepository_Delete(t *testing.T) {
	called := false
	client := &mockGithubClient{
		deleteFunc: func(url string, response interface{}) error {
			called = true
			if url != "notifications/threads/123" {
				t.Errorf("Expected url notifications/threads/123, got %s", url)
			}
			return nil
		},
	}

	repo := NewGithubRepository(client)
	err := repo.Delete("123")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !called {
		t.Error("Delete was not called")
	}
}
