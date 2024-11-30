package core

import (
	"github.com/go-logr/logr"
)

type Notification struct {
	ID      string  `json:"id"`
	Subject Subject `json:"subject"`
}

type Subject struct {
	Title string `json:"title"`
	Type  string `json:"type"`
	URL   string `json:"url"`
}

type NotificationService interface {
	FetchNotifications(logger logr.Logger, since string, noCache bool) error
}

type NotificationRepository interface {
	Delete(id string) error
	GetByTimePeriod(since string) ([]Notification, error)
}

type notificationService struct {
	notificationRepo NotificationRepository
	prService        PRService
	cacheService     CacheService
}

func NewNotificationService(repo NotificationRepository, pr PRService, cache CacheService) NotificationService {
	return &notificationService{
		notificationRepo: repo,
		prService:        pr,
		cacheService:     cache,
	}
}

func (s *notificationService) FetchNotifications(logger logr.Logger, since string, noCache bool) error {
	notifications, err := s.notificationRepo.GetByTimePeriod(since)
	if err != nil {
		return err
	}

	logger.V(1).Info("Fetched notifications", "count", len(notifications))

	cache, err := s.cacheService.Load()
	if err != nil {
		return err
	}

	for _, notification := range notifications {
		if notification.Subject.Type == "PullRequest" {
			if !noCache && s.cacheService.IsThreadDeleted(notification.ID) {
				logger.V(1).Info("Skipping already deleted thread",
					"title", notification.Subject.Title,
					"id", notification.ID)
				continue
			}

			logger.V(1).Info("Checking PR status",
				"title", notification.Subject.Title,
				"id", notification.ID)

			merged, err := s.handlePullRequest(notification.Subject.URL, noCache)
			if err != nil {
				logger.Error(err, "Failed to handle pull request")
				continue
			}

			if !merged {
				logger.V(1).Info("PR not merged, skipping deletion",
					"title", notification.Subject.Title,
					"id", notification.ID)
				continue
			}

			logger.V(1).Info("Deleting notification for merged PR",
				"title", notification.Subject.Title,
				"id", notification.ID)

			err = s.notificationRepo.Delete(notification.ID)
			if err != nil {
				logger.Error(err, "Failed to delete notification")
				continue
			}
			s.cacheService.SetThreadDeleted(notification.ID)
			logger.V(1).Info("Successfully deleted notification",
				"title", notification.Subject.Title,
				"id", notification.ID)
		}
	}

	return s.cacheService.Save(cache)
}

func (s *notificationService) handlePullRequest(url string, noCache bool) (bool, error) {
	if !noCache {
		if merged, exists := s.cacheService.GetPRStatus(url); exists {
			return merged, nil
		}
	}

	merged, err := s.prService.GetPRStatus(url)
	if err != nil {
		return false, err
	}

	if !noCache {
		s.cacheService.SetPRStatus(url, merged)
	}
	return merged, nil
}
