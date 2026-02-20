package handlers

import (
	"html/template"
	"net/http"

	"typeten/internal/domain"
	"typeten/internal/usecases"
)

var (
	indexTpl   = template.Must(template.New("index").Parse(indexHTML))
	textTpl    = template.Must(template.New("text").Parse(textHTML))
	sessionTpl = template.Must(template.New("session").Parse(sessionHTML))
)

type indexViewModel struct {
	Texts []*domain.TextInfo
}

type textViewModel struct {
	Text *domain.TextInfo
}

type sessionViewModel struct {
	Session *domain.Session
}

// IndexPage renders the main page with list of texts and a form to add a new one.
func (h *Handlers) IndexPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	out, err := h.listTextsUseCase.Execute(r.Context(), usecases.ListTextsInput{
		UserID: h.currentUserID,
	})
	if err != nil {
		http.Error(w, "Failed to load texts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := indexTpl.Execute(w, indexViewModel{Texts: out.Texts}); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// CreateTextHTML handles form submission for creating a text and redirects back to index.
func (h *Handlers) CreateTextHTML(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	_, err := h.createTextUseCase.Execute(r.Context(), usecases.CreateTextInput{
		UserID:  h.currentUserID,
		Title:   title,
		Content: content,
	})
	if err != nil {
		http.Error(w, "Failed to create text: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// TextDetailPage renders a simple page for a single text with ability to start a session.
func (h *Handlers) TextDetailPage(w http.ResponseWriter, r *http.Request, textID string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// For MVP, reuse ListTexts use case and filter by ID.
	out, err := h.listTextsUseCase.Execute(r.Context(), usecases.ListTextsInput{
		UserID: h.currentUserID,
	})
	if err != nil {
		http.Error(w, "Failed to load texts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var found *domain.TextInfo
	for _, t := range out.Texts {
		if t.ID == domain.TextID(textID) {
			found = t
			break
		}
	}
	if found == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := textTpl.Execute(w, textViewModel{Text: found}); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// CreateSessionHTML handles form submission for creating a session and redirects to session page.
func (h *Handlers) CreateSessionHTML(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	textID := r.FormValue("text_id")

	out, err := h.createSessionUseCase.Execute(r.Context(), usecases.CreateSessionInput{
		UserID: h.currentUserID,
		TextID: domain.TextID(textID),
	})
	if err != nil {
		http.Error(w, "Failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/sessions/"+string(out.Session.ID), http.StatusSeeOther)
}

// SessionPage renders the session page with embedded JavaScript for typing UI.
func (h *Handlers) SessionPage(w http.ResponseWriter, r *http.Request, sessionID string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	out, err := h.getSessionUseCase.Execute(r.Context(), usecases.GetSessionInput{
		SessionID: domain.SessionID(sessionID),
	})
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := sessionTpl.Execute(w, sessionViewModel{Session: out.Session}); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

const indexHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>TypeTen - Practice Typing</title>
  <style>
    body {
      font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      margin: 0;
      padding: 0;
      background: #0f172a;
      color: #e5e7eb;
    }
    header {
      padding: 1.5rem 2rem;
      background: #020617;
      border-bottom: 1px solid #1f2937;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    header h1 {
      margin: 0;
      font-size: 1.4rem;
    }
    main {
      max-width: 960px;
      margin: 2rem auto;
      padding: 0 1.5rem 3rem;
      display: grid;
      grid-template-columns: minmax(0, 2fr) minmax(0, 3fr);
      gap: 2rem;
      align-items: flex-start;
    }
    .card {
      background: #020617;
      border-radius: 0.75rem;
      border: 1px solid #1f2937;
      padding: 1.25rem 1.5rem;
      box-shadow: 0 18px 40px rgba(15, 23, 42, 0.6);
    }
    h2 {
      margin-top: 0;
      font-size: 1.1rem;
      margin-bottom: 0.75rem;
    }
    .texts-list {
      list-style: none;
      padding: 0;
      margin: 0.5rem 0 0;
    }
    .texts-list li {
      padding: 0.6rem 0.4rem;
      border-bottom: 1px solid #111827;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    .texts-list li:last-child {
      border-bottom: none;
    }
    .texts-list a {
      color: #e5e7eb;
      text-decoration: none;
    }
    .texts-list a:hover {
      color: #a5b4fc;
    }
    .badge {
      font-size: 0.75rem;
      padding: 0.1rem 0.5rem;
      border-radius: 999px;
      background: #111827;
      color: #9ca3af;
    }
    label {
      display: block;
      font-size: 0.85rem;
      color: #9ca3af;
      margin-bottom: 0.25rem;
    }
    input[type="text"], textarea {
      width: 100%;
      border-radius: 0.5rem;
      border: 1px solid #1f2937;
      background: #020617;
      color: #e5e7eb;
      padding: 0.6rem 0.75rem;
      font-size: 0.9rem;
      resize: vertical;
    }
    textarea {
      min-height: 180px;
      line-height: 1.4;
      font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
    }
    input:focus, textarea:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 1px #4f46e5;
    }
    button {
      border: none;
      border-radius: 999px;
      padding: 0.5rem 1.1rem;
      background: linear-gradient(135deg, #4f46e5, #7c3aed);
      color: white;
      font-size: 0.9rem;
      font-weight: 500;
      cursor: pointer;
      margin-top: 0.75rem;
      display: inline-flex;
      align-items: center;
      gap: 0.4rem;
    }
    button:hover {
      filter: brightness(1.1);
    }
    button:active {
      transform: translateY(1px);
    }
    .empty {
      font-size: 0.85rem;
      color: #6b7280;
      padding: 0.5rem 0.25rem;
    }
    .subtitle {
      font-size: 0.8rem;
      color: #6b7280;
    }
  </style>
</head>
<body>
  <header>
    <h1>TypeTen</h1>
    <span class="subtitle">Practice touch typing with your own texts</span>
  </header>
  <main>
    <section class="card">
      <h2>Your texts</h2>
      {{if .Texts}}
      <ul class="texts-list">
        {{range .Texts}}
        <li>
          <div>
            <a href="/texts/{{.ID}}">{{.Title}}</a>
            <div class="subtitle">{{.TotalLines}} lines · {{.FragmentCount}} fragments</div>
          </div>
          <span class="badge">ID: {{.ID}}</span>
        </li>
        {{end}}
      </ul>
      {{else}}
      <p class="empty">You don't have any texts yet. Add one on the right to start practicing.</p>
      {{end}}
    </section>
    <section class="card">
      <h2>Add a new text</h2>
      <form method="post" action="/texts">
        <div style="margin-bottom:0.75rem;">
          <label for="title">Title</label>
          <input id="title" name="title" type="text" required placeholder="e.g. The quick brown fox">
        </div>
        <div>
          <label for="content">Text content (one paragraph per line)</label>
          <textarea id="content" name="content" required placeholder="Paste or type your text here...&#10;Each line will be used as a typing unit."></textarea>
        </div>
        <button type="submit">Save text</button>
      </form>
    </section>
  </main>
</body>
</html>`

const textHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>{{.Text.Title}} · TypeTen</title>
  <style>
    body {
      font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      margin: 0;
      padding: 0;
      background: #0f172a;
      color: #e5e7eb;
    }
    header {
      padding: 1.25rem 2rem;
      background: #020617;
      border-bottom: 1px solid #1f2937;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    header a {
      color: #9ca3af;
      text-decoration: none;
      font-size: 0.85rem;
    }
    header a:hover {
      color: #e5e7eb;
    }
    main {
      max-width: 800px;
      margin: 2rem auto;
      padding: 0 1.5rem 3rem;
    }
    .card {
      background: #020617;
      border-radius: 0.75rem;
      border: 1px solid #1f2937;
      padding: 1.5rem 1.75rem;
      box-shadow: 0 18px 40px rgba(15, 23, 42, 0.6);
    }
    h1 {
      margin: 0;
      font-size: 1.3rem;
    }
    .meta {
      margin-top: 0.4rem;
      font-size: 0.85rem;
      color: #9ca3af;
    }
    button {
      border: none;
      border-radius: 999px;
      padding: 0.55rem 1.2rem;
      background: linear-gradient(135deg, #4f46e5, #7c3aed);
      color: white;
      font-size: 0.9rem;
      font-weight: 500;
      cursor: pointer;
      margin-top: 1.2rem;
      display: inline-flex;
      align-items: center;
      gap: 0.4rem;
    }
    button:hover {
      filter: brightness(1.1);
    }
    button:active {
      transform: translateY(1px);
    }
  </style>
</head>
<body>
  <header>
    <a href="/">&larr; Back to texts</a>
    <span></span>
  </header>
  <main>
    <section class="card">
      <h1>{{.Text.Title}}</h1>
      <p class="meta">
        {{.Text.TotalLines}} lines · {{.Text.FragmentCount}} fragments · ID: {{.Text.ID}}
      </p>
      <form method="post" action="/sessions">
        <input type="hidden" name="text_id" value="{{.Text.ID}}">
        <button type="submit">Start practice session</button>
      </form>
    </section>
  </main>
</body>
</html>`

const sessionHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Session {{.Session.ID}} · TypeTen</title>
  <style>
    body {
      font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      margin: 0;
      padding: 0;
      background: #020617;
      color: #e5e7eb;
    }
    header {
      padding: 1.25rem 2rem;
      background: #020617;
      border-bottom: 1px solid #1f2937;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    header a {
      color: #9ca3af;
      text-decoration: none;
      font-size: 0.85rem;
    }
    header a:hover {
      color: #e5e7eb;
    }
    main {
      max-width: 960px;
      margin: 1.5rem auto 3rem;
      padding: 0 1.5rem;
      display: grid;
      grid-template-columns: minmax(0, 3fr) minmax(0, 2fr);
      gap: 1.75rem;
    }
    .card {
      background: #020617;
      border-radius: 0.75rem;
      border: 1px solid #1f2937;
      padding: 1.25rem 1.5rem;
      box-shadow: 0 18px 40px rgba(15, 23, 42, 0.6);
    }
    #current-line {
      font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
      background: #030712;
      border-radius: 0.5rem;
      padding: 0.8rem 0.9rem;
      margin-bottom: 0.75rem;
      border: 1px solid #111827;
      font-size: 0.95rem;
      min-height: 2.4em;
      display: flex;
      align-items: center;
    }
    textarea {
      width: 100%;
      border-radius: 0.5rem;
      border: 1px solid #1f2937;
      background: #020617;
      color: #e5e7eb;
      padding: 0.65rem 0.75rem;
      font-size: 0.95rem;
      resize: vertical;
      min-height: 120px;
      font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
    }
    textarea:focus {
      outline: none;
      border-color: #4f46e5;
      box-shadow: 0 0 0 1px #4f46e5;
    }
    button {
      border: none;
      border-radius: 999px;
      padding: 0.55rem 1.2rem;
      background: linear-gradient(135deg, #4f46e5, #7c3aed);
      color: white;
      font-size: 0.9rem;
      font-weight: 500;
      cursor: pointer;
      margin-top: 0.75rem;
      display: inline-flex;
      align-items: center;
      gap: 0.4rem;
    }
    button:hover {
      filter: brightness(1.1);
    }
    button:active {
      transform: translateY(1px);
    }
    .badge {
      font-size: 0.75rem;
      padding: 0.1rem 0.5rem;
      border-radius: 999px;
      background: #111827;
      color: #9ca3af;
    }
    .stat-label {
      font-size: 0.75rem;
      color: #9ca3af;
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }
    .stat-value {
      font-size: 1.3rem;
      font-weight: 600;
    }
    .stats-grid {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 0.75rem 1.5rem;
      margin-top: 0.75rem;
    }
    .pill {
      border-radius: 999px;
      padding: 0.25rem 0.6rem;
      font-size: 0.75rem;
      background: #0b1120;
      border: 1px solid #1f2937;
      display: inline-flex;
      align-items: center;
      gap: 0.4rem;
    }
    .pill-dot {
      width: 6px;
      height: 6px;
      border-radius: 999px;
      background: #22c55e;
    }
    .pill-dot.inactive {
      background: #6b7280;
    }
  </style>
</head>
<body>
  <header>
    <a href="/texts/{{.Session.TextID}}">&larr; Back to text</a>
    <span class="badge">Session ID: {{.Session.ID}}</span>
  </header>
  <main>
    <section class="card">
      <div style="margin-bottom:0.5rem;">
        <div class="pill" id="status-pill">
          <span class="pill-dot" id="status-dot"></span>
          <span id="status-text">Warm up</span>
        </div>
      </div>
      <div id="current-line">Loading text…</div>
      <textarea id="input" placeholder="Type the line above here…"></textarea>
      <button id="complete-line-btn" type="button">Complete line</button>
    </section>
    <section class="card">
      <div class="stat-label">Session progress</div>
      <div class="stats-grid">
        <div>
          <div class="stat-label">Completed lines</div>
          <div class="stat-value" id="stat-lines">0</div>
        </div>
        <div>
          <div class="stat-label">Average WPM</div>
          <div class="stat-value" id="stat-wpm">0</div>
        </div>
        <div>
          <div class="stat-label">Accuracy</div>
          <div class="stat-value" id="stat-accuracy">0%</div>
        </div>
        <div>
          <div class="stat-label">Elapsed</div>
          <div class="stat-value" id="stat-time">0s</div>
        </div>
      </div>
    </section>
  </main>
  <script>
    (function() {
      const sessionId = "{{.Session.ID}}";
      const textId = "{{.Session.TextID}}";

      const currentLineEl = document.getElementById("current-line");
      const inputEl = document.getElementById("input");
      const completeBtn = document.getElementById("complete-line-btn");
      const statLinesEl = document.getElementById("stat-lines");
      const statWpmEl = document.getElementById("stat-wpm");
      const statAccuracyEl = document.getElementById("stat-accuracy");
      const statTimeEl = document.getElementById("stat-time");
      const statusPillEl = document.getElementById("status-pill");
      const statusDotEl = document.getElementById("status-dot");
      const statusTextEl = document.getElementById("status-text");

      let lines = [];
      let currentIndex = 0;
      let sessionStart = Date.now();
      let lineStart = Date.now();
      let timerId = null;

      function updateTimer() {
        const elapsedSec = Math.floor((Date.now() - sessionStart) / 1000);
        const mins = Math.floor(elapsedSec / 60);
        const secs = elapsedSec % 60;
        statTimeEl.textContent = (mins > 0 ? mins + "m " : "") + secs + "s";
      }

      function startTimer() {
        if (timerId) return;
        timerId = setInterval(updateTimer, 1000);
      }

      function setStatus(text, active) {
        statusTextEl.textContent = text;
        if (active) {
          statusDotEl.classList.remove("inactive");
        } else {
          statusDotEl.classList.add("inactive");
        }
      }

      function loadLines() {
        fetch("/api/texts/" + encodeURIComponent(textId) + "/fragments")
          .then(function(res) {
            if (!res.ok) {
              throw new Error("Failed to load fragments");
            }
            return res.json();
          })
          .then(function(data) {
            const frags = data.fragments || [];
            frags.forEach(function(f) {
              (f.lines || []).forEach(function(line) {
                lines.push(line);
              });
            });
            if (lines.length === 0) {
              currentLineEl.textContent = "No lines to practice in this text.";
              completeBtn.disabled = true;
              inputEl.disabled = true;
              setStatus("Idle", false);
              return;
            }
            currentIndex = 0;
            currentLineEl.textContent = lines[currentIndex];
            inputEl.value = "";
            inputEl.focus();
            sessionStart = Date.now();
            lineStart = Date.now();
            startTimer();
            setStatus("Typing", true);
          })
          .catch(function(err) {
            console.error(err);
            currentLineEl.textContent = "Failed to load text.";
            setStatus("Error loading text", false);
          });
      }

      function computeAccuracy(expected, typed) {
        if (!expected && !typed) return 100.0;
        const maxLen = Math.max(expected.length, typed.length);
        if (maxLen === 0) return 100.0;
        let correct = 0;
        for (let i = 0; i < maxLen; i++) {
          if (expected[i] === typed[i]) {
            correct++;
          }
        }
        return (correct / maxLen) * 100.0;
      }

      function computeWPM(typed, millis) {
        const minutes = millis / 60000.0;
        if (minutes <= 0) return 0.0;
        const words = typed.length / 5.0;
        return words / minutes;
      }

      function sendProgress(accuracy, wpm) {
        return fetch("/api/sessions/" + encodeURIComponent(sessionId) + "/progress", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            accuracy_percent: accuracy,
            wpm: wpm
          })
        }).then(function(res) {
          if (!res.ok) {
            throw new Error("Failed to record progress");
          }
          return res.json();
        }).then(function(data) {
          statLinesEl.textContent = data.completed_lines;
          statWpmEl.textContent = data.average_wpm.toFixed(1);
          statAccuracyEl.textContent = data.total_accuracy_percent.toFixed(1) + "%";
          if (data.is_completed) {
            setStatus("Completed", false);
          }
        }).catch(function(err) {
          console.error(err);
          setStatus("Error sending progress", false);
        });
      }

      completeBtn.addEventListener("click", function() {
        if (currentIndex >= lines.length) {
          setStatus("Completed", false);
          return;
        }
        const expected = lines[currentIndex] || "";
        const typed = inputEl.value || "";
        const now = Date.now();
        const lineMillis = now - lineStart;
        const accuracy = computeAccuracy(expected, typed);
        const wpm = computeWPM(typed, lineMillis);

        sendProgress(accuracy, wpm).then(function() {
          currentIndex++;
          if (currentIndex >= lines.length) {
            currentLineEl.textContent = "All lines completed. Great job!";
            inputEl.disabled = true;
            completeBtn.disabled = true;
            setStatus("Completed", false);
          } else {
            currentLineEl.textContent = lines[currentIndex];
            inputEl.value = "";
            inputEl.focus();
            lineStart = Date.now();
            setStatus("Typing", true);
          }
        });
      });

      inputEl.addEventListener("keydown", function(e) {
        if (e.key === "Enter" && (e.ctrlKey || e.metaKey)) {
          e.preventDefault();
          completeBtn.click();
        }
      });

      loadLines();
    })();
  </script>
</body>
</html>`

