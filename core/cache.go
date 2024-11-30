package core

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Cache struct {
	PRStatus       map[string]bool `json:"pr_status"`
	ThreadsDeleted map[string]bool `json:"threads_deleted"`
}

type CacheService interface {
	Load() (*Cache, error)
	Save(*Cache) error
	IsThreadDeleted(id string) bool
	SetThreadDeleted(id string)
	GetPRStatus(url string) (bool, bool)
	SetPRStatus(url string, merged bool)
}

type fileCacheService struct {
	cache     *Cache
	cacheDir  string
	cacheFile string
}

func NewFileCacheService(homeDir string) CacheService {
	return &fileCacheService{
		cacheDir:  filepath.Join(homeDir, ".dailyare"),
		cacheFile: "cache.json",
	}
}

func (s *fileCacheService) Load() (*Cache, error) {
	if err := os.MkdirAll(s.cacheDir, 0o755); err != nil {
		return nil, err
	}

	filePath := filepath.Join(s.cacheDir, s.cacheFile)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			s.cache = &Cache{
				PRStatus:       make(map[string]bool),
				ThreadsDeleted: make(map[string]bool),
			}
			return s.cache, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, &s.cache); err != nil {
		return nil, err
	}

	if s.cache.PRStatus == nil {
		s.cache.PRStatus = make(map[string]bool)
	}
	if s.cache.ThreadsDeleted == nil {
		s.cache.ThreadsDeleted = make(map[string]bool)
	}

	return s.cache, nil
}

func (s *fileCacheService) Save(cache *Cache) error {
	if err := os.MkdirAll(s.cacheDir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(s.cacheDir, s.cacheFile), data, 0o644)
}

func (s *fileCacheService) IsThreadDeleted(id string) bool {
	return s.cache.ThreadsDeleted[id]
}

func (s *fileCacheService) SetThreadDeleted(id string) {
	s.cache.ThreadsDeleted[id] = true
}

func (s *fileCacheService) GetPRStatus(url string) (bool, bool) {
	status, exists := s.cache.PRStatus[url]
	return status, exists
}

func (s *fileCacheService) SetPRStatus(url string, merged bool) {
	s.cache.PRStatus[url] = merged
}
