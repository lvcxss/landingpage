# Terminus CSS Library Documentation

**Terminus CSS** is a lightweight, zero-dependency, dark-purple terminal-inspired CSS design system created to replace external frameworks like `terminal.css` with a modernized, customizable, self-contained stylesheet.

---

## 1. Design Tokens & Palette

Terminus CSS is built around the dark purple and dark-grey-close-to-black theme:

| Variable                | Hex Value | Purpose                                                                                            |
| :---------------------- | :-------- | :------------------------------------------------------------------------------------------------- |
| `--primary-color`       | `#bf00ff` | Signature electric / neon dark purple for primary accents, prompts, active borders, and highlights |
| `--primary-hover`       | `#d74bff` | Brighter purple for interactive hovers                                                             |
| `--primary-dim`         | `#9900d1` | Deep rich purple for secondary accents, error borders, and branch connectors                       |
| `--primary-dark`        | `#320046` | Deep purple wash for active backgrounds and subtle backdrops                                       |
| `--secondary-color`     | `#383a42` | Primary neutral dark grey for buttons, header strips, and card outlines                            |
| `--bg-canvas`           | `#101014` | Deep terminal screen background (close to black)                                                   |
| `--bg-surface`          | `#16161c` | Slightly elevated background for cards, windows, and panels                                        |
| `--bg-surface-elevated` | `#1e1e26` | Active tab panels, popovers, and elevated hover surfaces                                           |
| `--bg-code`             | `#121217` | Monospace code block background                                                                    |
| `--font-color`          | `#e8e9ed` | High-contrast readable terminal text                                                               |
| `--font-muted`          | `#8b8e9d` | Secondary logs, timestamps, comments                                                               |

---

## 2. Quickstart & Installation

To use Terminus CSS in any project:

```html
<!-- Stylesheet -->
<link rel="stylesheet" href="/static/terminus.css" />

<!-- Optional companion script for tabbed subtimeline switching -->
<script src="/static/terminus.js"></script>
```

---

## 3. Subtimeline Feature Guide

When an activity started and ended _during_ another broader milestone (for example, living in Bosnia between 2014–2022 and obtaining a Computer Science degree between 2016–2019), Terminus CSS provides the **Subtimeline Component**.

### 3.1 Tabbed Subtimeline (Recommended)

In this pattern, the parent timeline card contains an indented subtimeline panel with interactive tabs. Each tab represents a concurrent activity.

```html
<div class="terminal-timeline">
  <!-- Parent Milestone -->
  <div class="terminal-timeline-item">
    <div class="terminal-card">
      <header>
        <span>Bosnia & Herzegovina</span>
        <span class="timeline-period">2014 — 2022</span>
      </header>
      <div>
        <p>
          Living and working in Bosnia, building core technical competencies...
        </p>

        <!-- Nested Subtimeline -->
        <div class="terminal-subtimeline">
          <div class="subtimeline-meta">
            <span class="subtimeline-label">
              <span class="branch-symbol">↳</span> SUBTIMELINE // CONCURRENT
              ACTIVITIES
            </span>
            <span class="badge badge-purple">TABBED VIEW</span>
          </div>

          <!-- Tab Navigation -->
          <div class="subtimeline-tabs">
            <button
              class="subtimeline-tab-btn active"
              data-target="#subtab-degree"
            >
              <span class="tab-bracket">[</span> 2016 — 2019 : B.Sc. Computer
              Science <span class="tab-bracket">]</span>
            </button>
            <button class="subtimeline-tab-btn" data-target="#subtab-labs">
              <span class="tab-bracket">[</span> 2018 — 2020 : Network & Systems
              Lab <span class="tab-bracket">]</span>
            </button>
          </div>

          <!-- Tab 1: CS Degree -->
          <div id="subtab-degree" class="subtimeline-panel active">
            <div class="subtimeline-panel-header">
              <h4>Bachelor of Science in Computer Science</h4>
              <span class="subtimeline-timespan">2016 — 2019</span>
            </div>
            <div class="subtimeline-body">
              <p>
                Completed Bachelor's degree with focus on Operating Systems,
                Networking Architecture, and Security.
              </p>
              <div class="subtimeline-tags">
                <span class="tag">#c-cpp</span>
                <span class="tag">#linux</span>
                <span class="tag">#networks</span>
              </div>
            </div>
          </div>

          <!-- Tab 2: Concurrent Lab Experience -->
          <div id="subtab-labs" class="subtimeline-panel">
            <div class="subtimeline-panel-header">
              <h4>Network & Systems Laboratory</h4>
              <span class="subtimeline-timespan">2018 — 2020</span>
            </div>
            <div class="subtimeline-body">
              <p>
                Virtualization clusters, hypervisors, and firewall routing
                policies.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
```

### 3.2 Sequential Branch Track View

If you prefer a continuous vertical branch spine inside the parent event:

```html
<div class="terminal-subtimeline">
  <div class="subtimeline-meta">
    <span class="subtimeline-label"
      ><span class="branch-symbol">├──</span> CONCURRENT MILESTONES</span
    >
  </div>
  <div style="padding: 1rem 1.25rem;">
    <div class="subtimeline-branch-track">
      <div class="subtimeline-branch-node">
        <div class="subtimeline-branch-title">
          <span>2016 — 2019: B.Sc. in Computer Science</span>
          <span class="badge badge-purple">DEGREE</span>
        </div>
        <p>Graduated with honors during Bosnia residency.</p>
      </div>
      <div class="subtimeline-branch-node">
        <div class="subtimeline-branch-title">
          <span>2018 — 2020: Systems Administration</span>
          <span class="badge">LABS</span>
        </div>
      </div>
    </div>
  </div>
</div>
```

### 3.3 Collapsible CLI Accordion

For pure HTML5 collapsible subtimelines:

```html
<details class="subtimeline-details">
  <summary>
    <span>2016 — 2019: B.Sc. Computer Science (Click to expand)</span>
    <span class="badge badge-purple" style="margin-left: auto;">[+]</span>
  </summary>
  <div class="subtimeline-details-content">
    <p>Detailed breakdown of coursework and achievements.</p>
  </div>
</details>
```

---

## 4. Other Core Components

### 4.1 Terminal Windows & Cards

- `.terminal-card`: Clean card with top header strip, body, and hover glow.
- `.terminal-window`: Retro terminal window with macOS dots (`.terminal-dot.dot-red`, `dot-yellow`, `dot-green`) and titlebar.

### 4.2 Buttons

- `.btn.btn-default`: Secondary dark grey button
- `.btn.btn-primary`: Neon dark purple button with glow
- `.btn.btn-error`: Deep purple accent button
- `.btn.btn-ghost`: Transparent button with crisp 1px border

### 4.3 CLI File Rows

- `.cli-file-list`: Vertical interactive listing
- `.cli-file-row`: Sliding hover row with purple left-border accent
- `.file-type`, `.file-name`, `.file-date`

### 4.4 Terminal Prompts & Menus

- `.terminal-prompt`: Animated blinking block cursor `>_ ▋`
- `.terminal-menu`: Bracketed navigation items `[ item ]` with active purple borders

---

## 5. Live Showcase

A complete live demonstration and playground is available at:
`http://localhost:4000/static/showcase.html`
