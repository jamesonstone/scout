package site

import (
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jamesonstone/scout/internal/model"
)

func renderSite(outDir string, data siteData) error {
	pages := []struct {
		path string
		name string
		view pageView
	}{
		{"index.html", "home", pageView{PageTitle: "Scout Research Intelligence", Site: data}},
		{"daily/index.html", "dailyArchive", pageView{PageTitle: "Scout Daily Archive", Site: data}},
		{"monthly/index.html", "monthlyArchive", pageView{PageTitle: "Scout Monthly Rankings", Site: data}},
	}
	for _, page := range pages {
		if err := renderPage(filepath.Join(outDir, page.path), page.name, page.view); err != nil {
			return err
		}
	}
	for _, day := range data.Daily {
		view := pageView{
			PageTitle: "Scout Daily " + day.Date,
			Site:      data,
			Day:       day,
			Sections: []paperSection{
				{Title: "Top Papers", Papers: day.Top},
				{Title: "Additional Papers", Papers: day.Additional},
				{Title: "Watchlist", Papers: day.Watchlist},
			},
		}
		if err := renderPage(filepath.Join(outDir, "daily", day.Date, "index.html"), "daily", view); err != nil {
			return err
		}
	}
	for _, month := range data.Monthly {
		view := pageView{PageTitle: "Scout Monthly " + month.Month, Site: data, Month: month}
		if err := renderPage(filepath.Join(outDir, "monthly", month.Month, "index.html"), "monthly", view); err != nil {
			return err
		}
	}
	for _, paper := range data.Papers {
		view := pageView{PageTitle: paper.Record.Title, Site: data, Paper: paper}
		if err := renderPage(filepath.Join(outDir, "papers", slug(paper.Record.ID), "index.html"), "paper", view); err != nil {
			return err
		}
	}
	return nil
}

func renderPage(path, name string, data any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	tmpl, err := template.New("page").Funcs(template.FuncMap{
		"asset":     assetURL,
		"base":      baseURL,
		"limitHome": limitHome,
	}).Parse(layoutTemplate + pageTemplates)
	if err != nil {
		return err
	}
	return tmpl.ExecuteTemplate(file, name, data)
}

func writeStaticAssets(outDir string) error {
	if err := writeFile(filepath.Join(outDir, ".nojekyll"), []byte("")); err != nil {
		return err
	}
	return writeFile(filepath.Join(outDir, "assets", "styles.css"), []byte(stylesCSS))
}

func writeFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func siteURL(basePath, route string) string {
	route = strings.TrimLeft(route, "/")
	if route == "" {
		return normalizeBasePath(basePath)
	}
	return normalizeBasePath(basePath) + route
}

func assetURL(basePath, route string) string {
	return siteURL(basePath, route)
}

func baseURL(basePath string) string {
	return normalizeBasePath(basePath)
}

func reportPath(dataDir, date string) string {
	return filepath.Join(dataDir, "reports", "daily", date[:7], date+".md")
}

func extractExecutiveSignal(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	body := string(data)
	start := strings.Index(body, "## Executive Signal")
	if start < 0 {
		return ""
	}
	rest := strings.TrimSpace(body[start+len("## Executive Signal"):])
	if next := strings.Index(rest, "\n## "); next >= 0 {
		rest = rest[:next]
	}
	return strings.TrimSpace(rest)
}

func sortRecords(records []model.PaperRecord) {
	sort.Slice(records, func(i, j int) bool {
		if records[i].Score.Overall == records[j].Score.Overall {
			return records[i].Title < records[j].Title
		}
		return records[i].Score.Overall > records[j].Score.Overall
	})
}

func limitPapers(papers []paperPage, limit int) []paperPage {
	if len(papers) <= limit {
		return papers
	}
	return papers[:limit]
}

func limitHome(papers []paperPage) []paperPage {
	return limitPapers(papers, 6)
}

func themesForRecords(records []model.PaperRecord) []theme {
	counts := map[string]int{}
	for _, record := range records {
		for _, category := range record.Categories {
			counts[category]++
		}
	}
	themes := make([]theme, 0, len(counts))
	for name, count := range counts {
		themes = append(themes, theme{Name: name, Count: count})
	}
	sort.Slice(themes, func(i, j int) bool {
		if themes[i].Count == themes[j].Count {
			return themes[i].Name < themes[j].Name
		}
		return themes[i].Count > themes[j].Count
	})
	if len(themes) > 6 {
		themes = themes[:6]
	}
	return themes
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
	sort.Strings(out)
	return out
}

func categoriesLabel(categories []string) string {
	if len(categories) == 0 {
		return "N/A"
	}
	return strings.Join(categories, ", ")
}

func scoreClass(score int) string {
	switch {
	case score >= 80:
		return "score-high"
	case score >= 60:
		return "score-mid"
	default:
		return "score-low"
	}
}

func slug(value string) string {
	return sanitizeFile(strings.ToLower(value))
}

func sanitizeFile(value string) string {
	replacer := strings.NewReplacer("/", "-", "\\", "-", ":", "-", " ", "-")
	return replacer.Replace(value)
}
