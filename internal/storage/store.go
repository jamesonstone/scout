package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jamesonstone/scout/internal/model"
)

type Store struct {
	root string
}

func New(root string) Store {
	return Store{root: root}
}

func (s Store) SavePaper(record model.PaperRecord) error {
	path := s.paperPath(record.ID)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return writeJSON(path, record)
}

func (s Store) LoadPaper(id string) (model.PaperRecord, bool, error) {
	path := s.paperPath(id)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return model.PaperRecord{}, false, nil
	}
	if err != nil {
		return model.PaperRecord{}, false, err
	}
	var record model.PaperRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return model.PaperRecord{}, false, err
	}
	return record, true, nil
}

func (s Store) SaveObservation(date time.Time, ids []string) error {
	path := s.dailyObservationPath(date)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	unique := uniqueStrings(ids)
	sort.Strings(unique)
	return writeJSON(path, model.DailyObservation{Date: date.Format("2006-01-02"), PaperIDs: unique})
}

func (s Store) MonthRecords(month time.Time) ([]model.PaperRecord, error) {
	pattern := filepath.Join(s.root, "data", "daily", month.Format("2006-01"), "*.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	var records []model.PaperRecord
	for _, file := range files {
		var observation model.DailyObservation
		if err := readJSON(file, &observation); err != nil {
			return nil, err
		}
		for _, id := range observation.PaperIDs {
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			record, ok, err := s.LoadPaper(id)
			if err != nil {
				return nil, err
			}
			if ok {
				records = append(records, record)
			}
		}
	}
	return records, nil
}

func (s Store) SaveDailyReport(date time.Time, body string) (string, error) {
	path := filepath.Join(s.root, "reports", "daily", date.Format("2006-01"), date.Format("2006-01-02")+".md")
	return path, writeText(path, body)
}

func (s Store) SaveMonthlyReport(date time.Time, body string) (string, error) {
	path := filepath.Join(s.root, "reports", "monthly", date.Format("2006-01")+".md")
	return path, writeText(path, body)
}

func (s Store) paperPath(id string) string {
	return filepath.Join(s.root, "data", "papers", sanitize(id)+".json")
}

func (s Store) dailyObservationPath(date time.Time) string {
	return filepath.Join(s.root, "data", "daily", date.Format("2006-01"), date.Format("2006-01-02")+".json")
}

func sanitize(value string) string {
	replacer := strings.NewReplacer("/", "-", "\\", "-", ":", "-", " ", "-")
	return replacer.Replace(value)
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return writeFile(path, data)
}

func writeText(path, body string) error {
	return writeFile(path, []byte(body))
}

func writeFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create dir for %s: %w", path, err)
	}
	return os.WriteFile(path, data, 0o644)
}

func readJSON(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

func uniqueStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
