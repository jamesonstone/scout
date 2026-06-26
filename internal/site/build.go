package site

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jamesonstone/scout/internal/model"
)

func Build(cfg Config) (BuildResult, error) {
	cfg = normalizeConfig(cfg)
	data, err := loadSiteData(cfg)
	if err != nil {
		return BuildResult{}, err
	}
	if len(data.Daily) == 0 {
		return BuildResult{}, fmt.Errorf("no daily observations found under %s", filepath.Join(cfg.DataDir, "data", "daily"))
	}
	if err := resetOutputDir(cfg.OutDir); err != nil {
		return BuildResult{}, err
	}
	if err := writeStaticAssets(cfg.OutDir); err != nil {
		return BuildResult{}, err
	}
	if err := copyJSONData(cfg.DataDir, cfg.OutDir); err != nil {
		return BuildResult{}, err
	}
	if err := writeManifest(cfg.OutDir, data); err != nil {
		return BuildResult{}, err
	}
	if err := renderSite(cfg.OutDir, data); err != nil {
		return BuildResult{}, err
	}
	return BuildResult{OutDir: cfg.OutDir, DailyPages: len(data.Daily), MonthlyPages: len(data.Monthly), PaperPages: len(data.Papers)}, nil
}

func normalizeConfig(cfg Config) Config {
	if cfg.DataDir == "" {
		cfg.DataDir = "."
	}
	if cfg.OutDir == "" {
		cfg.OutDir = "public"
	}
	cfg.BasePath = normalizeBasePath(cfg.BasePath)
	return cfg
}

func normalizeBasePath(base string) string {
	base = strings.TrimSpace(base)
	if base == "" {
		return "/"
	}
	if !strings.HasPrefix(base, "/") {
		base = "/" + base
	}
	if !strings.HasSuffix(base, "/") {
		base += "/"
	}
	return base
}

func loadSiteData(cfg Config) (siteData, error) {
	records, err := loadPaperRecords(filepath.Join(cfg.DataDir, "data", "papers"))
	if err != nil {
		return siteData{}, err
	}
	daily, monthIDs, err := loadDailyPages(cfg, records)
	if err != nil {
		return siteData{}, err
	}
	monthly := make([]monthlyPage, 0, len(monthIDs))
	for month, ids := range monthIDs {
		page := monthlyFromIDs(cfg.BasePath, month, ids, records)
		monthly = append(monthly, page)
	}
	sort.Slice(monthly, func(i, j int) bool { return monthly[i].Month > monthly[j].Month })

	papers := make([]paperPage, 0, len(records))
	for _, record := range records {
		papers = append(papers, paperPageFromRecord(cfg.BasePath, record))
	}
	sort.Slice(papers, func(i, j int) bool {
		if papers[i].Score == papers[j].Score {
			return papers[i].Record.Title < papers[j].Record.Title
		}
		return papers[i].Score > papers[j].Score
	})
	return siteData{BasePath: cfg.BasePath, Daily: daily, Monthly: monthly, Papers: papers}, nil
}

func loadPaperRecords(root string) (map[string]model.PaperRecord, error) {
	files, err := filepath.Glob(filepath.Join(root, "*.json"))
	if err != nil {
		return nil, err
	}
	records := make(map[string]model.PaperRecord, len(files))
	for _, file := range files {
		var record model.PaperRecord
		if err := readJSON(file, &record); err != nil {
			return nil, fmt.Errorf("read paper record %s: %w", file, err)
		}
		if record.ID == "" {
			return nil, fmt.Errorf("paper record %s is missing id", file)
		}
		records[record.ID] = record
	}
	return records, nil
}

func loadDailyPages(cfg Config, records map[string]model.PaperRecord) ([]dailyPage, map[string][]string, error) {
	files, err := filepath.Glob(filepath.Join(cfg.DataDir, "data", "daily", "*", "*.json"))
	if err != nil {
		return nil, nil, err
	}
	sort.Strings(files)
	monthIDs := map[string][]string{}
	daily := make([]dailyPage, 0, len(files))
	for _, file := range files {
		var observation model.DailyObservation
		if err := readJSON(file, &observation); err != nil {
			return nil, nil, fmt.Errorf("read daily observation %s: %w", file, err)
		}
		papers := recordsForIDs(observation.PaperIDs, records)
		page := dailyFromRecords(cfg, observation.Date, papers)
		daily = append(daily, page)
		month := observation.Date[:7]
		monthIDs[month] = append(monthIDs[month], observation.PaperIDs...)
	}
	sort.Slice(daily, func(i, j int) bool { return daily[i].Date > daily[j].Date })
	return daily, monthIDs, nil
}

func recordsForIDs(ids []string, records map[string]model.PaperRecord) []model.PaperRecord {
	out := make([]model.PaperRecord, 0, len(ids))
	for _, id := range ids {
		if record, ok := records[id]; ok {
			out = append(out, record)
		}
	}
	sortRecords(out)
	return out
}

func dailyFromRecords(cfg Config, date string, records []model.PaperRecord) dailyPage {
	page := dailyPage{
		Date:            date,
		URL:             siteURL(cfg.BasePath, "daily/"+date+"/"),
		JSONURL:         siteURL(cfg.BasePath, "data/daily/"+date[:7]+"/"+date+".json"),
		MonthlyURL:      siteURL(cfg.BasePath, "monthly/"+date[:7]+"/"),
		ExecutiveSignal: extractExecutiveSignal(reportPath(cfg.DataDir, date)),
		ArchiveCount:    len(records),
	}
	if page.ExecutiveSignal == "" {
		page.ExecutiveSignal = "Scout found papers for " + date + " and ranked them by implementation signal."
	}
	for _, record := range records {
		paper := paperPageFromRecord(cfg.BasePath, record)
		switch record.Recommendation {
		case model.RecommendationRead:
			page.Top = append(page.Top, paper)
		case model.RecommendationWorthWatching:
			page.Additional = append(page.Additional, paper)
		default:
			page.Watchlist = append(page.Watchlist, paper)
		}
	}
	return page
}

func monthlyFromIDs(basePath, month string, ids []string, records map[string]model.PaperRecord) monthlyPage {
	unique := uniqueStrings(ids)
	monthRecords := recordsForIDs(unique, records)
	papers := make([]paperPage, 0, len(monthRecords))
	for _, record := range monthRecords {
		papers = append(papers, paperPageFromRecord(basePath, record))
	}
	page := monthlyPage{Month: month, URL: siteURL(basePath, "monthly/"+month+"/"), AllPapers: papers, Themes: themesForRecords(monthRecords)}
	page.Top = limitPapers(papers, 10)
	for _, paper := range papers {
		if paper.Score >= 70 {
			page.Rising = append(page.Rising, paper)
		}
	}
	page.Rising = limitPapers(page.Rising, 5)
	return page
}

func paperPageFromRecord(basePath string, record model.PaperRecord) paperPage {
	return paperPage{
		Record:     record,
		URL:        siteURL(basePath, "papers/"+slug(record.ID)+"/"),
		JSONURL:    siteURL(basePath, "data/papers/"+sanitizeFile(record.ID)+".json"),
		Categories: categoriesLabel(record.Categories),
		FirstSeen:  record.FirstSeen,
		Observed:   strings.Join(record.ObservedDates, ", "),
		Score:      record.Score.Overall,
		ScoreClass: scoreClass(record.Score.Overall),
		Links:      linksForRecord(record.Links),
	}
}

func linksForRecord(links model.Links) []siteLink {
	var out []siteLink
	add := func(label, url string) {
		url = strings.TrimSpace(url)
		if url != "" {
			out = append(out, siteLink{Label: label, URL: url})
		}
	}
	add("Hugging Face", links.HuggingFace)
	add("arXiv", links.Arxiv)
	add("Paper", links.Paper)
	add("PDF", links.PDF)
	for _, url := range links.GitHub {
		add("GitHub", url)
	}
	for _, url := range links.Project {
		add("Project", url)
	}
	return out
}

func resetOutputDir(outDir string) error {
	if err := os.RemoveAll(outDir); err != nil {
		return err
	}
	return os.MkdirAll(outDir, 0o755)
}

func copyJSONData(dataDir, outDir string) error {
	source := filepath.Join(dataDir, "data")
	target := filepath.Join(outDir, "data")
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		dest := filepath.Join(target, rel)
		if entry.IsDir() {
			return os.MkdirAll(dest, 0o755)
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}
		return copyFile(path, dest)
	})
}

func copyFile(source, target string) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func writeManifest(outDir string, data siteData) error {
	m := manifest{}
	for _, day := range data.Daily {
		m.Daily = append(m.Daily, manifestEntry{Date: day.Date, URL: day.URL, JSON: day.JSONURL})
	}
	for _, month := range data.Monthly {
		m.Monthly = append(m.Monthly, manifestEntry{Month: month.Month, URL: month.URL})
	}
	for _, paper := range data.Papers {
		m.Papers = append(m.Papers, manifestEntry{ID: paper.Record.ID, Title: paper.Record.Title, URL: paper.URL, JSON: paper.JSONURL, Score: paper.Score})
	}
	return writeJSON(filepath.Join(outDir, "data", "index.json"), m)
}

func readJSON(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return writeFile(path, data)
}
