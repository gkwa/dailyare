package core

import (
	"testing"
)

func TestFileCacheService_LoadNoExistingCache(t *testing.T) {
	tmpDir := t.TempDir()
	svc := NewFileCacheService(tmpDir)

	cache, err := svc.Load()
	if err != nil {
		t.Fatalf("Failed to load cache: %v", err)
	}

	if len(cache.PRStatus) != 0 {
		t.Error("Expected empty PRStatus map")
	}
	if len(cache.ThreadsDeleted) != 0 {
		t.Error("Expected empty ThreadsDeleted map")
	}
}

func TestFileCacheService_SaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	svc := NewFileCacheService(tmpDir)

	cache := &Cache{
		PRStatus: map[string]bool{
			"pr1": true,
			"pr2": false,
		},
		ThreadsDeleted: map[string]bool{
			"thread1": true,
		},
	}

	if err := svc.Save(cache); err != nil {
		t.Fatalf("Failed to save cache: %v", err)
	}

	loaded, err := svc.Load()
	if err != nil {
		t.Fatalf("Failed to load cache: %v", err)
	}

	if len(loaded.PRStatus) != 2 {
		t.Errorf("Expected 2 PR statuses, got %d", len(loaded.PRStatus))
	}
	if !loaded.PRStatus["pr1"] {
		t.Error("Expected pr1 to be true")
	}
	if loaded.PRStatus["pr2"] {
		t.Error("Expected pr2 to be false")
	}
	if !loaded.ThreadsDeleted["thread1"] {
		t.Error("Expected thread1 to be true")
	}
}
