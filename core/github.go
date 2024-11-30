package core

import (
	"fmt"
	"strings"
	"time"
)

type GithubClient interface {
	Get(url string, response interface{}) error
	Delete(url string, response interface{}) error
}

type githubRepository struct {
	client GithubClient
}

func NewGithubRepository(client GithubClient) NotificationRepository {
	return &githubRepository{client: client}
}

func (r *githubRepository) Delete(id string) error {
	return r.client.Delete("notifications/threads/"+id, nil)
}

func (r *githubRepository) GetByTimePeriod(since string) ([]Notification, error) {
	sinceTime, err := parseDuration(since)
	if err != nil {
		return nil, err
	}

	sinceDate := time.Now().Add(-sinceTime).Format(time.RFC3339)
	url := fmt.Sprintf("notifications?all=true&since=%s", sinceDate)

	var notifications []Notification
	err = r.client.Get(url, &notifications)
	return notifications, err
}

func formatGithubURL(url string) string {
	apiURL := strings.Replace(url, "https://api.github.com/repos/", "", 1)
	return fmt.Sprintf("repos/%s", apiURL)
}
