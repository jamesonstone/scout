package site

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jamesonstone/scout/internal/artifact"
)

type ValidationResult struct {
	CheckedPages int
	CheckedLinks int
}

var hrefPattern = regexp.MustCompile(`href="([^"]+)"`)

func Validate(cfg Config) (ValidationResult, error) {
	cfg = normalizeConfig(cfg)
	required := []string{
		"index.html",
		"daily/index.html",
		"monthly/index.html",
		"data/index.json",
		"assets/styles.css",
		".nojekyll",
	}
	for _, rel := range required {
		if _, err := os.Stat(filepath.Join(cfg.OutDir, rel)); err != nil {
			return ValidationResult{}, fmt.Errorf("missing required output %s: %w", rel, err)
		}
	}

	var result ValidationResult
	pages, err := htmlFiles(cfg.OutDir)
	if err != nil {
		return ValidationResult{}, err
	}
	if len(pages) == 0 {
		return ValidationResult{}, fmt.Errorf("no HTML pages found under %s", cfg.OutDir)
	}
	for _, page := range pages {
		body, err := os.ReadFile(page)
		if err != nil {
			return ValidationResult{}, err
		}
		result.CheckedPages++
		links := hrefPattern.FindAllStringSubmatch(string(body), -1)
		for _, match := range links {
			result.CheckedLinks++
			if err := validateLink(cfg, match[1]); err != nil {
				return ValidationResult{}, fmt.Errorf("%s: %w", page, err)
			}
		}
	}
	if err := validateRequiredContent(cfg.OutDir); err != nil {
		return ValidationResult{}, err
	}
	if err := validateScoreJSON(filepath.Join(cfg.OutDir, "data", "papers")); err != nil {
		return ValidationResult{}, err
	}
	return result, nil
}

func htmlFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && filepath.Ext(path) == ".html" {
			files = append(files, path)
		}
		return nil
	})
	sortStrings(files)
	return files, err
}

func validateLink(cfg Config, href string) error {
	if href == "" || strings.HasPrefix(href, "#") || strings.HasPrefix(href, "data:") || strings.HasPrefix(href, "mailto:") {
		return nil
	}
	for _, prefix := range []string{"https://huggingface.co/", "https://arxiv.org/", "https://github.com/", "https://", "http://"} {
		if strings.HasPrefix(href, prefix) {
			return nil
		}
	}
	base := normalizeBasePath(cfg.BasePath)
	if !strings.HasPrefix(href, base) {
		return fmt.Errorf("internal link %q does not use base path %q", href, base)
	}
	rel := strings.TrimPrefix(href, base)
	if rel == "" {
		if _, err := os.Stat(filepath.Join(cfg.OutDir, "index.html")); err != nil {
			return fmt.Errorf("broken internal link %q: %w", href, err)
		}
		return nil
	}
	target := filepath.Join(cfg.OutDir, filepath.FromSlash(rel))
	if strings.HasSuffix(href, "/") {
		target = filepath.Join(target, "index.html")
	}
	if _, err := os.Stat(target); err != nil {
		return fmt.Errorf("broken internal link %q: %w", href, err)
	}
	return nil
}

func validateRequiredContent(outDir string) error {
	dailyPages, err := filepath.Glob(filepath.Join(outDir, "daily", "????-??-??", "index.html"))
	if err != nil {
		return err
	}
	if len(dailyPages) == 0 {
		return fmt.Errorf("no daily briefing pages generated")
	}
	monthlyPages, err := filepath.Glob(filepath.Join(outDir, "monthly", "????-??", "index.html"))
	if err != nil {
		return err
	}
	if len(monthlyPages) == 0 {
		return fmt.Errorf("no monthly ranking pages generated")
	}
	paperPages, err := filepath.Glob(filepath.Join(outDir, "papers", "*", "index.html"))
	if err != nil {
		return err
	}
	if len(paperPages) == 0 {
		return fmt.Errorf("no paper detail pages generated")
	}
	homeBody, err := os.ReadFile(filepath.Join(outDir, "index.html"))
	if err != nil {
		return err
	}
	home := string(homeBody)
	for _, token := range []string{"<h1>Scout</h1>", "Papers fetched by Scout", "Latest fetched date", "Published"} {
		if !strings.Contains(home, token) {
			return fmt.Errorf("home page missing %q", token)
		}
	}
	for _, token := range []string{"Install From Source", "Quick Start", "Storage Contract", "Repository Notes", "Maintainer"} {
		if strings.Contains(home, token) {
			return fmt.Errorf("home page contains README-style content %q", token)
		}
	}
	dailyBody, err := os.ReadFile(dailyPages[0])
	if err != nil {
		return err
	}
	for _, token := range []string{"Papers fetched on", "Executive Signal", "Top Papers", "Additional Papers", "Watchlist", "Archive", "Innovation Summary", "Executive Summary", "Why It Matters", "Implementation Angle", "Caveat", "Estimated Reading Priority", "Published"} {
		if !strings.Contains(string(dailyBody), token) {
			return fmt.Errorf("daily page missing %q", token)
		}
	}
	paperBody, err := os.ReadFile(paperPages[0])
	if err != nil {
		return err
	}
	for _, token := range []string{"Executive Summary", "Estimated Reading Priority", "Published", "First fetched", "Observation History", "Score Breakdown"} {
		if !strings.Contains(string(paperBody), token) {
			return fmt.Errorf("paper page missing %q", token)
		}
	}
	monthlyBody, err := os.ReadFile(monthlyPages[0])
	if err != nil {
		return err
	}
	for _, token := range []string{"Top 10 Papers", "Rising Papers", "Themes", "Complete Monthly Index"} {
		if !strings.Contains(string(monthlyBody), token) {
			return fmt.Errorf("monthly page missing %q", token)
		}
	}
	return nil
}

func validateScoreJSON(root string) error {
	files, err := filepath.Glob(filepath.Join(root, "*.json"))
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no public paper JSON records found")
	}
	for _, file := range files {
		var payload struct {
			Score struct {
				Novelty                 *int `json:"novelty"`
				PracticalImpact         *int `json:"practical_impact"`
				TechnicalDepth          *int `json:"technical_depth"`
				ImplementationPotential *int `json:"implementation_potential"`
				Relevance               *int `json:"relevance"`
				CommunitySignal         *int `json:"community_signal"`
				SummaryConfidence       *int `json:"summary_confidence"`
				Overall                 *int `json:"overall"`
			} `json:"score"`
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if err := artifact.ValidatePaperRecordJSON(file, data); err != nil {
			return err
		}
		if err := json.Unmarshal(data, &payload); err != nil {
			return fmt.Errorf("decode %s: %w", file, err)
		}
		if payload.Score.Novelty == nil ||
			payload.Score.PracticalImpact == nil ||
			payload.Score.TechnicalDepth == nil ||
			payload.Score.ImplementationPotential == nil ||
			payload.Score.Relevance == nil ||
			payload.Score.CommunitySignal == nil ||
			payload.Score.SummaryConfidence == nil ||
			payload.Score.Overall == nil {
			return fmt.Errorf("%s missing score breakdown fields", file)
		}
	}
	return nil
}

func sortStrings(values []string) {
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j] < values[j-1]; j-- {
			values[j], values[j-1] = values[j-1], values[j]
		}
	}
}
