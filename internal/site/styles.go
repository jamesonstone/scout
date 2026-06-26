package site

const stylesCSS = `:root {
  color-scheme: light;
  --bg: #f5f7f6;
  --surface: #ffffff;
  --surface-strong: #f9fbfa;
  --text: #1f2723;
  --muted: #65716b;
  --line: #dce4df;
  --green: #2f6f5e;
  --green-dark: #214f45;
  --gold: #a67316;
  --red: #a5463b;
  --shadow: 0 18px 55px rgba(31, 39, 35, 0.08);
}

* {
  box-sizing: border-box;
}

body {
  margin: 0;
  background: var(--bg);
  color: var(--text);
  font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  line-height: 1.6;
}

a {
  color: var(--green-dark);
  text-decoration-thickness: 0.08em;
  text-underline-offset: 0.18em;
}

.site-header,
.site-footer,
main {
  width: min(1120px, calc(100% - 40px));
  margin: 0 auto;
}

.site-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 28px 0;
}

.brand {
  font-weight: 800;
  font-size: 1.2rem;
  text-decoration: none;
  letter-spacing: 0;
}

nav {
  display: flex;
  gap: 22px;
}

nav a {
  color: var(--muted);
  text-decoration: none;
  font-weight: 650;
}

main {
  padding-bottom: 64px;
}

.hero {
  padding: 96px 0 80px;
  max-width: 880px;
}

.page-heading {
  padding: 72px 0 42px;
  max-width: 880px;
}

.eyebrow {
  margin: 0 0 14px;
  color: var(--gold);
  font-size: 0.78rem;
  font-weight: 800;
  letter-spacing: 0.08em;
  text-transform: uppercase;
}

h1,
h2,
h3,
p {
  margin-top: 0;
}

h1 {
  margin-bottom: 22px;
  font-size: clamp(2.4rem, 7vw, 5.4rem);
  line-height: 0.98;
  letter-spacing: 0;
}

h2 {
  margin-bottom: 22px;
  font-size: 1.45rem;
  line-height: 1.18;
}

h3 {
  font-size: 1.08rem;
  line-height: 1.3;
}

h4 {
  margin: 8px 0 0;
  color: var(--green-dark);
  font-size: 0.78rem;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.lede,
.page-heading p {
  color: var(--muted);
  font-size: 1.2rem;
  max-width: 760px;
}

.hero-actions,
.meta-links,
.card-meta,
.card-topline,
.score-line {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
}

.button {
  display: inline-flex;
  min-height: 44px;
  align-items: center;
  padding: 10px 18px;
  border-radius: 6px;
  background: var(--green-dark);
  color: white;
  text-decoration: none;
  font-weight: 750;
}

.button.secondary {
  background: transparent;
  color: var(--green-dark);
  border: 1px solid var(--line);
}

.section-grid,
.paper-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 24px;
}

.paper-grid {
  align-items: stretch;
}

.panel,
.paper-card,
.list-surface {
  background: var(--surface);
  border: 1px solid var(--line);
  border-radius: 8px;
  box-shadow: var(--shadow);
}

.panel,
.paper-card {
  padding: 28px;
}

.paper-section {
  margin: 52px 0;
}

.paper-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.paper-card p,
.panel p,
.paper-detail li {
  color: var(--muted);
}

.innovation {
  color: var(--text) !important;
  font-weight: 650;
}

.score-pill {
  display: inline-flex;
  align-items: center;
  min-height: 32px;
  padding: 4px 10px;
  border-radius: 999px;
  background: #e9efe9;
  color: var(--green-dark);
  font-size: 0.86rem;
  font-weight: 800;
}

.score-high {
  background: #dfeee6;
  color: var(--green-dark);
}

.score-mid {
  background: #f2e6c8;
  color: #6f4c12;
}

.score-low {
  background: #f2ddd9;
  color: var(--red);
}

.link-list,
.paper-rank,
.index-list,
.theme-list,
.link-stack {
  padding-left: 20px;
}

.link-stack.compact {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 18px;
  padding-left: 0;
  list-style: none;
}

.link-list li,
.paper-rank li,
.theme-list li {
  margin: 14px 0;
}

.paper-rank li,
.archive-row,
.theme-list li {
  display: flex;
  justify-content: space-between;
  gap: 18px;
}

.paper-rank li {
  align-items: center;
}

.paper-rank a {
  flex: 1 1 320px;
}

.score-bar {
  display: inline-flex;
  flex: 0 0 110px;
  height: 8px;
  overflow: hidden;
  border-radius: 999px;
  background: #e1e8e4;
}

.score-bar span {
  display: block;
  background: linear-gradient(90deg, var(--green), var(--gold));
}

.archive-row {
  padding: 20px 24px;
  border-bottom: 1px solid var(--line);
  text-decoration: none;
}

.archive-row:last-child {
  border-bottom: 0;
}

.index-list li {
  margin: 18px 0;
}

.index-list span,
.card-meta,
.score-line {
  color: var(--muted);
}

.paper-detail .panel {
  margin-bottom: 24px;
}

.score-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 14px;
  margin: 0;
}

.score-grid div {
  padding: 16px;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: var(--surface-strong);
}

.score-grid dt {
  color: var(--muted);
  font-size: 0.85rem;
}

.score-grid dd {
  margin: 4px 0 0;
  font-size: 1.45rem;
  font-weight: 800;
}

.empty-state,
.archive-note {
  color: var(--muted);
}

.site-footer {
  padding: 32px 0 56px;
  color: var(--muted);
  border-top: 1px solid var(--line);
}

@media (max-width: 760px) {
  .site-header {
    align-items: flex-start;
    flex-direction: column;
    gap: 16px;
  }

  nav {
    width: 100%;
    justify-content: space-between;
  }

  .section-grid,
  .paper-grid,
  .score-grid {
    grid-template-columns: 1fr;
  }

  .hero {
    padding-top: 64px;
  }

  .paper-rank li,
  .archive-row,
  .theme-list li {
    align-items: flex-start;
    flex-direction: column;
  }
}
`
