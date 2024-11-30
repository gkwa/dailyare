package core

import (
	"testing"

	"github.com/go-logr/logr/testr"
)

type mockNotificationRepo struct {
	deleteFunc          func(id string) error
	getByTimePeriodFunc func(since string) ([]Notification, error)
}

func (m *mockNotificationRepo) Delete(id string) error {
	return m.deleteFunc(id)
}

func (m *mockNotificationRepo) GetByTimePeriod(since string) ([]Notification, error) {
	return m.getByTimePeriodFunc(since)
}

type mockPRService struct {
	getPRStatusFunc func(url string) (bool, error)
}

func (m *mockPRService) GetPRStatus(url string) (bool, error) {
	return m.getPRStatusFunc(url)
}

type mockCacheService struct {
	cache *Cache
}

func newMockCacheService() *mockCacheService {
	return &mockCacheService{
		cache: &Cache{
			PRStatus:       make(map[string]bool),
			ThreadsDeleted: make(map[string]bool),
		},
	}
}

func (m *mockCacheService) Load() (*Cache, error)          { return m.cache, nil }
func (m *mockCacheService) Save(*Cache) error              { return nil }
func (m *mockCacheService) IsThreadDeleted(id string) bool { return m.cache.ThreadsDeleted[id] }
func (m *mockCacheService) SetThreadDeleted(id string)     { m.cache.ThreadsDeleted[id] = true }
func (m *mockCacheService) GetPRStatus(url string) (bool, bool) {
	v, ok := m.cache.PRStatus[url]
	return v, ok
}
func (m *mockCacheService) SetPRStatus(url string, merged bool) { m.cache.PRStatus[url] = merged }

func TestNotificationService_FetchNotifications(t *testing.T) {
	notificationRepo := &mockNotificationRepo{
		getByTimePeriodFunc: func(since string) ([]Notification, error) {
			return []Notification{
				{ID: "1", Subject: Subject{Type: "PullRequest", URL: "pr1"}},
				{ID: "2", Subject: Subject{Type: "PullRequest", URL: "pr2"}},
			}, nil
		},
		deleteFunc: func(id string) error {
			return nil
		},
	}

	prService := &mockPRService{
		getPRStatusFunc: func(url string) (bool, error) {
			return url == "pr1", nil
		},
	}

	cacheService := newMockCacheService()
	service := NewNotificationService(notificationRepo, prService, cacheService)
	logger := testr.New(t)

	err := service.FetchNotifications(logger, "7d", false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !cacheService.IsThreadDeleted("1") {
		t.Error("Expected thread 1 to be marked as deleted")
	}
	if cacheService.IsThreadDeleted("2") {
		t.Error("Expected thread 2 to not be marked as deleted")
	}
}

func TestNotificationService_FetchNotifications_NoCache(t *testing.T) {
	notificationRepo := &mockNotificationRepo{
		getByTimePeriodFunc: func(since string) ([]Notification, error) {
			return []Notification{
				{ID: "1", Subject: Subject{Type: "PullRequest", URL: "pr1"}},
			}, nil
		},
		deleteFunc: func(id string) error {
			return nil
		},
	}

	prService := &mockPRService{
		getPRStatusFunc: func(url string) (bool, error) {
			return true, nil
		},
	}

	cacheService := newMockCacheService()
	cacheService.SetThreadDeleted("1") // Pre-mark as deleted

	service := NewNotificationService(notificationRepo, prService, cacheService)
	logger := testr.New(t)

	err := service.FetchNotifications(logger, "7d", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
