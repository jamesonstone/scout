package storage

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/jamesonstone/scout/internal/model"
)

func TestStorePersistsPaperAndMonthlyLookup(t *testing.T) {
	root := t.TempDir()
	store := New(root)
	record := model.PaperRecord{ID: "2501.00001", Title: "Paper One", FirstSeen: "2026-01-02", ObservedDates: []string{"2026-01-02"}}
	if err := store.SavePaper(record); err != nil {
		t.Fatalf("save paper: %v", err)
	}
	day, _ := time.Parse("2006-01-02", "2026-01-02")
	if err := store.SaveObservation(day, []string{"2501.00001", "2501.00001"}); err != nil {
		t.Fatalf("save observation: %v", err)
	}
	loaded, ok, err := store.LoadPaper("2501.00001")
	if err != nil || !ok {
		t.Fatalf("load paper: ok=%v err=%v", ok, err)
	}
	if loaded.Title != record.Title {
		t.Fatalf("unexpected title: %s", loaded.Title)
	}
	monthRecords, err := store.MonthRecords(day)
	if err != nil {
		t.Fatalf("month records: %v", err)
	}
	if len(monthRecords) != 1 {
		t.Fatalf("expected 1 monthly record, got %d", len(monthRecords))
	}
	dailyPath := filepath.Join(root, "data", "daily", "2026-01", "2026-01-02.json")
	if _, err := filepath.Abs(dailyPath); err != nil {
		t.Fatalf("expected absolute path resolution: %v", err)
	}
}
