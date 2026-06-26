package site

import "github.com/jamesonstone/scout/internal/model"

type Config struct {
	DataDir  string
	OutDir   string
	BasePath string
}

type BuildResult struct {
	OutDir       string
	DailyPages   int
	MonthlyPages int
	PaperPages   int
}

type siteData struct {
	BasePath string
	Daily    []dailyPage
	Monthly  []monthlyPage
	Papers   []paperPage
}

type dailyPage struct {
	Date            string
	URL             string
	JSONURL         string
	MonthlyURL      string
	ExecutiveSignal string
	Top             []paperPage
	Additional      []paperPage
	Watchlist       []paperPage
	ArchiveCount    int
}

type monthlyPage struct {
	Month     string
	URL       string
	Top       []paperPage
	Rising    []paperPage
	Themes    []theme
	AllPapers []paperPage
}

type paperPage struct {
	Record     model.PaperRecord
	URL        string
	JSONURL    string
	Categories string
	FirstSeen  string
	Observed   string
	Score      int
	ScoreClass string
	Links      []siteLink
}

type siteLink struct {
	Label string
	URL   string
}

type theme struct {
	Name  string
	Count int
}

type pageView struct {
	PageTitle string
	Site      siteData
	Day       dailyPage
	Month     monthlyPage
	Paper     paperPage
	Sections  []paperSection
}

type paperSection struct {
	Title  string
	Papers []paperPage
}

type manifest struct {
	Daily   []manifestEntry `json:"daily"`
	Monthly []manifestEntry `json:"monthly"`
	Papers  []manifestEntry `json:"papers"`
}

type manifestEntry struct {
	ID    string `json:"id,omitempty"`
	Date  string `json:"date,omitempty"`
	Month string `json:"month,omitempty"`
	Title string `json:"title,omitempty"`
	URL   string `json:"url"`
	JSON  string `json:"json,omitempty"`
	Score int    `json:"score,omitempty"`
}
