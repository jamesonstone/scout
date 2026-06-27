package site

const layoutTemplate = `{{define "pageStart"}}<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.PageTitle}}</title>
  <link rel="icon" href="data:,">
  <link rel="stylesheet" href="{{asset .Site.BasePath "assets/styles.css"}}">
</head>
<body>
  <header class="site-header">
    <a class="brand" href="{{base .Site.BasePath}}">Scout</a>
    <nav aria-label="Primary">
      <a href="{{asset .Site.BasePath "daily/"}}">Daily</a>
      <a href="{{asset .Site.BasePath "monthly/"}}">Monthly</a>
      <a href="{{asset .Site.BasePath "data/index.json"}}">Data</a>
    </nav>
  </header>
  <main>
{{end}}
{{define "pageEnd"}}  </main>
  <footer class="site-footer">
    <p>Read-only intelligence briefings generated from Hugging Face Daily Papers. Source artifacts remain in Scout.</p>
  </footer>
</body>
</html>{{end}}
`

const pageTemplates = `
{{define "home"}}{{template "pageStart" .}}
<section class="hero">
  <p class="eyebrow">Papers fetched by Scout</p>
  <h1>Scout</h1>
  <p class="lede">Daily AI research signal from Hugging Face Daily Papers, distilled for builders who need to decide what deserves deeper reading.</p>
  <div class="hero-actions">
    {{with .Site.FeaturedDay}}<a class="button" href="{{.URL}}">Latest papers from {{.Date}}</a>{{end}}
    {{with index .Site.Monthly 0}}<a class="button secondary" href="{{.URL}}">Current monthly ranking</a>{{end}}
  </div>
</section>
{{with .Site.FeaturedDay}}
<section class="paper-section home-latest">
  <div class="section-heading">
    <div>
      <p class="eyebrow">Latest fetched date</p>
      <h2>Papers fetched on {{.Date}}</h2>
    </div>
    <a href="{{.URL}}">{{.ArchiveCount}} paper{{if ne .ArchiveCount 1}}s{{end}}</a>
  </div>
  {{if $.Site.FeaturedPapers}}<div class="paper-grid">{{range $.Site.FeaturedPapers}}{{template "paperCard" .}}{{end}}</div>{{else}}<p class="empty-state">No papers were fetched for this date yet. The previous populated briefings remain available below.</p>{{end}}
</section>
{{end}}
<section class="section-grid">
  <div class="panel">
    <h2>Daily Dates</h2>
    <ol class="link-list">{{range .Site.Daily}}<li><a href="{{.URL}}">{{.Date}}</a><span>{{.ArchiveCount}} paper{{if ne .ArchiveCount 1}}s{{end}} fetched</span></li>{{end}}</ol>
  </div>
  <div class="panel">
    <h2>Highest Signal</h2>
    <ol class="paper-rank">{{range .Site.Papers | limitHome}}<li><a href="{{.URL}}">{{.Record.Title}}</a><span>Published {{.PublishedDate}}</span><span class="score-bar" aria-label="Score {{.Score}} out of 100"><span style="width: {{.Score}}%"></span></span><strong>{{.Score}}</strong></li>{{end}}</ol>
  </div>
</section>
{{template "pageEnd" .}}{{end}}

{{define "dailyArchive"}}{{template "pageStart" .}}
<section class="page-heading">
  <p class="eyebrow">Archive</p>
  <h1>Daily Briefings</h1>
  <p>Papers are organized by the date Scout fetched them. Each paper page also shows the paper publish date when Hugging Face provides it.</p>
</section>
<section class="list-surface">{{range .Site.Daily}}
  <a class="archive-row" href="{{.URL}}"><span>Fetched {{.Date}}</span><strong>{{.ArchiveCount}} paper{{if ne .ArchiveCount 1}}s{{end}}</strong></a>
{{end}}</section>
{{template "pageEnd" .}}{{end}}

{{define "monthlyArchive"}}{{template "pageStart" .}}
<section class="page-heading">
  <p class="eyebrow">Rankings</p>
  <h1>Monthly Research Signal</h1>
  <p>Monthly pages are recomputed from durable paper records so rankings reflect cumulative signal rather than daily order.</p>
</section>
<section class="list-surface">{{range .Site.Monthly}}
  <a class="archive-row" href="{{.URL}}"><span>{{.Month}}</span><strong>{{len .AllPapers}} papers</strong></a>
{{end}}</section>
{{template "pageEnd" .}}{{end}}

{{define "daily"}}{{template "pageStart" .}}
<section class="page-heading">
  <p class="eyebrow">Daily briefing</p>
  <h1>Papers fetched on {{.Day.Date}}</h1>
  <h2>Executive Signal</h2>
  <p>{{.Day.ExecutiveSignal}}</p>
  <div class="meta-links"><a href="{{.Day.JSONURL}}">Daily JSON</a><a href="{{.Day.MonthlyURL}}">Monthly ranking</a></div>
</section>
{{range .Sections}}{{template "paperSection" .}}{{end}}
<section class="panel archive-note">
  <h2>Archive</h2>
  <p>Daily record count: {{.Day.ArchiveCount}}. Persistent paper JSON lives under <a href="{{asset .Site.BasePath "data/index.json"}}">public data</a>.</p>
  <ol class="index-list compact">{{range .Day.Archive}}<li><a href="{{.URL}}">{{.Record.Title}}</a><span>Published {{.PublishedDate}} · {{.Score}}/100 · {{.Record.Recommendation}}</span></li>{{end}}</ol>
</section>
{{template "pageEnd" .}}{{end}}

{{define "monthly"}}{{template "pageStart" .}}
<section class="page-heading">
  <p class="eyebrow">Monthly ranking</p>
  <h1>{{.Month.Month}}</h1>
  <p>Top 10 papers, rising papers, dominant themes, and the full monthly index.</p>
</section>
<section class="panel">
  <h2>Top 10 Papers</h2>
  <ol class="paper-rank">{{range .Month.Top}}<li><a href="{{.URL}}">{{.Record.Title}}</a><span class="score-bar" aria-label="Score {{.Score}} out of 100"><span style="width: {{.Score}}%"></span></span><strong>{{.Score}}</strong></li>{{end}}</ol>
</section>
<section class="section-grid">
  <div class="panel">
    <h2>Rising Papers</h2>
    <ol class="link-list">{{range .Month.Rising}}<li><a href="{{.URL}}">{{.Record.Title}}</a><span>{{.Score}}/100</span></li>{{end}}</ol>
  </div>
  <div class="panel">
    <h2>Themes</h2>
    <ul class="theme-list">{{range .Month.Themes}}<li><span>{{.Name}}</span><strong>{{.Count}}</strong></li>{{end}}</ul>
  </div>
</section>
<section class="panel">
  <h2>Complete Monthly Index</h2>
  <ol class="index-list">{{range .Month.AllPapers}}<li><a href="{{.URL}}">{{.Record.Title}}</a><span>{{.Categories}}</span></li>{{end}}</ol>
</section>
{{template "pageEnd" .}}{{end}}

{{define "paper"}}{{template "pageStart" .}}
<article class="paper-detail">
  <section class="page-heading">
    <p class="eyebrow">Paper detail</p>
    <h1>{{.Paper.Record.Title}}</h1>
    <div class="score-line"><span class="score-pill {{.Paper.ScoreClass}}">{{.Paper.Score}}/100</span><span>{{.Paper.Record.Recommendation}}</span><span>Published {{.Paper.PublishedDate}}</span><span>Fetched {{.Paper.FirstSeen}}</span><span>{{.Paper.Categories}}</span></div>
  </section>
  <section class="panel">
    <h2>Innovation Summary</h2>
    <p>{{.Paper.Record.InnovationSummary}}</p>
    <h2>Executive Summary</h2>
    <p>{{.Paper.ExecutiveSummary}}</p>
    <h2>Why It Matters</h2>
    <ul>{{range .Paper.Record.WhyItMatters}}<li>{{.}}</li>{{end}}</ul>
    <h2>Implementation Angle</h2>
    <ul>{{range .Paper.Record.ImplementationAngle}}<li>{{.}}</li>{{end}}</ul>
    <h2>Caveat</h2>
    <p>{{.Paper.Record.Caveat}}</p>
    <h2>Estimated Reading Priority</h2>
    <p>{{.Paper.EstimatedReadingPriority}}</p>
  </section>
  <section class="section-grid">
    <div class="panel">
      <h2>Observation History</h2>
      <p>Published {{.Paper.PublishedDate}}. First fetched {{.Paper.FirstSeen}}. Observed {{.Paper.Observed}}.</p>
      <p><a href="{{.Paper.JSONURL}}">Paper JSON record</a></p>
    </div>
    <div class="panel">
      <h2>Links</h2>
      <ul class="link-stack">
        {{range .Paper.Links}}<li><a href="{{.URL}}">{{.Label}}</a></li>{{end}}
      </ul>
    </div>
  </section>
  <section class="panel">
    <h2>Score Breakdown</h2>
    <dl class="score-grid">
      <div><dt>Novelty</dt><dd>{{.Paper.Record.Score.Novelty}}</dd></div>
      <div><dt>Practical Impact</dt><dd>{{.Paper.Record.Score.PracticalImpact}}</dd></div>
      <div><dt>Technical Depth</dt><dd>{{.Paper.Record.Score.TechnicalDepth}}</dd></div>
      <div><dt>Implementation</dt><dd>{{.Paper.Record.Score.ImplementationPotential}}</dd></div>
      <div><dt>Relevance</dt><dd>{{.Paper.Record.Score.Relevance}}</dd></div>
      <div><dt>Community</dt><dd>{{.Paper.Record.Score.CommunitySignal}}</dd></div>
      <div><dt>Confidence</dt><dd>{{.Paper.Record.Score.SummaryConfidence}}</dd></div>
    </dl>
  </section>
</article>
{{template "pageEnd" .}}{{end}}

{{define "paperSection"}}
<section class="paper-section">
  <h2>{{.Title}}</h2>
  {{if .Papers}}{{if .Compact}}<div class="list-surface compact-list">{{range .Papers}}{{template "paperRow" .}}{{end}}</div>{{else}}<div class="paper-grid">{{range .Papers}}{{template "paperCard" .}}{{end}}</div>{{end}}{{else}}<p class="empty-state">No papers in this section.</p>{{end}}
</section>
{{end}}

{{define "paperRow"}}
<article class="compact-row">
  <div>
    <h3><a href="{{.URL}}">{{.Record.Title}}</a></h3>
    <p class="paper-date">Published {{.PublishedDate}} · Fetched {{.FirstSeen}}</p>
    <p>{{.Record.InnovationSummary}}</p>
  </div>
  <div class="row-meta"><span>{{.Score}}/100</span><span>{{.Record.Recommendation}}</span></div>
</article>
{{end}}

{{define "paperCard"}}
<article class="paper-card">
  <div class="card-topline"><span class="score-pill {{.ScoreClass}}">{{.Score}}/100</span><span>{{.Record.Recommendation}}</span></div>
  <h3><a href="{{.URL}}">{{.Record.Title}}</a></h3>
  <p class="paper-date">Published {{.PublishedDate}} · Fetched {{.FirstSeen}}</p>
  <h4>Innovation Summary</h4>
  <p class="innovation">{{.Record.InnovationSummary}}</p>
  <h4>Executive Summary</h4>
  <p>{{.ExecutiveSummary}}</p>
  <h4>Why It Matters</h4>
  <ul>{{range .Record.WhyItMatters}}<li>{{.}}</li>{{end}}</ul>
  <h4>Implementation Angle</h4>
  <ul>{{range .Record.ImplementationAngle}}<li>{{.}}</li>{{end}}</ul>
  <h4>Caveat</h4>
  <p>{{.Record.Caveat}}</p>
  <h4>Estimated Reading Priority</h4>
  <p>{{.EstimatedReadingPriority}}</p>
  <h4>Links</h4>
  <ul class="link-stack compact">
    {{range .Links}}<li><a href="{{.URL}}">{{.Label}}</a></li>{{end}}
  </ul>
  <div class="card-meta"><span>{{.Categories}}</span><a href="{{.JSONURL}}">JSON</a></div>
</article>
{{end}}
`
