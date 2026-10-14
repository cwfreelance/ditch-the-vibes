// Ditch the Vibes: a tiny step-by-step guide. One server binary for local dev
// (-tweaks adds a design A/B sidebar) and Vercel (Go framework preset, PORT env).
package main

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//go:embed page.html
var pageHTML string

var page = template.Must(template.New("page").Parse(pageHTML))

type Option struct{ Key, Letter, Name, Blurb, Icon string }

var options = []Option{
	{"a", "A", "Terminal on my Mac", "I type to an AI in the Terminal app.", "terminal"},
	{"b", "B", "Replit", "I work on it in my browser at replit.com.", "globe"},
	{"c", "C", "Another website", "Netlify, Bolt, Lovable, and the like.", "window"},
}

type Page struct {
	Title   string
	Lead    template.HTML
	Callout template.HTML
	Prompt  string
	Steps   []template.HTML
	Note    template.HTML
	Finish  bool
}

// Where Alexis sends things.
const (
	phone = "+13122851972"
	email = "colin@colinweaver.com"
)

const pasteNote = `Being asked to log in or click “Authorize” is normal. When it gives you a GitHub link, tap Next.`

var (
	pasteA = Page{
		Title:   "Paste this prompt into your AI",
		Lead:    `Start your AI tool in Terminal the way you normally do for this app, then paste everything in the box below and follow what it tells you.`,
		Callout: `If your Mac pops up a box about installing “command line developer tools,” click <b>Install</b>. If Terminal asks for your password, type your Mac password and press Enter. Nothing shows on screen while you type, and that’s normal.`,
		Prompt: `I'm not a programmer, so please explain everything in plain language and walk me through anything I need to do, one step at a time.

I want to share this app's code with my friend (GitHub username: cwfreelance) in a private GitHub repository on my account. If we're not already in this app's project folder, help me find it first.

1. Install the GitHub CLI (gh) if needed, and help me log in to my GitHub account with it.
2. Set up git for this project if it isn't already.
3. Make sure no secrets get uploaded: add a .gitignore covering .env files, node_modules, and build/cache folders; create a .env.example with variable names only; and check the code AND the git history for hardcoded API keys or passwords. Don't read or print any secret values. If you find a secret, stop and explain how I should replace it before continuing.
4. Create a private GitHub repository and push the code.
5. Add cwfreelance as a collaborator on the repository.
6. Give me the repository link.`,
		Note: pasteNote,
	}
	pasteB = Page{
		Title: "Paste this prompt into Replit",
		Lead:  `Open the app in Replit and paste everything in the box below into the <b>Agent</b> chat, then follow what it tells you.`,
		Prompt: `I'm not a programmer, so please explain in plain language and guide me through any clicks one step at a time.

I want to share this app's code with my friend in a private GitHub repository on my account.
1. Help me connect this project to GitHub using the Git tool (Tools → Git → Connect to GitHub) and create a new private repository.
2. Push the current code, keeping the Replit config files. Don't include secrets, .env files with real values, database dumps, or dependency/cache folders.
3. Give me the repository link.
Don't change how the app is published.`,
		Note: pasteNote,
	}
	pasteC = Page{
		Title: "Paste this prompt into the site’s AI",
		Lead:  `Open the app on the site where you built it and paste everything in the box below into its AI chat, then follow what it tells you.`,
		Prompt: `I'm not a programmer, so please explain in plain language and guide me through any clicks one step at a time.

I want to share this app's original source code with my friend in a private GitHub repository on my account.
1. Help me connect this project to GitHub and create a new private repository. (On Netlify this is under Project configuration → Build & deploy → Continuous deployment → Repository → Link to repository → Push to new repository.)
2. Check that the repository has the editable source code, not just the built website files. If it only has built files, help me find the original source.
3. Make sure no secrets, .env files with real values, or database dumps get uploaded.
4. Tell me in plain language whether anything about how my live site gets published will change.
5. Give me the repository link.`,
		Note: pasteNote,
	}
	invite = Page{
		Title: "Invite me to the repository",
		Lead:  `On your new repository’s page on github.com:`,
		Steps: []template.HTML{
			`Click <b>Settings</b> at the top of the page.`,
			`Click <b>Collaborators</b> in the left menu. GitHub may ask for your password.`,
			`Click <b>Add people</b>, type <code>cwfreelance</code>, and click <b>Add</b>.`,
		},
		Note: `Done? Tap Next.`,
	}
	textMe = Page{
		Title:  "Text me the link",
		Lead:   `Send me the GitHub link, and I’ll take it from here.`,
		Note:   `Have another app? Tap <svg class="inl"><use href="#i-chevrons-left"/></svg> to start over.`,
		Finish: true,
	}
)

// Pages after step 1, per option. Step n renders flows[opt][n-2].
var flows = map[string][]Page{
	"a": {pasteA, textMe},
	"b": {pasteB, invite, textMe},
	"c": {pasteC, invite, textMe},
}

// Fonts offered in the Tweaks sidebar. Production loads only prodFonts.
var (
	bodyFonts   = []string{"Inter", "Geist", "Manrope", "DM Sans", "Plus Jakarta Sans", "Figtree", "Outfit", "Lexend", "Space Grotesk"}
	monoFonts   = []string{"JetBrains Mono", "IBM Plex Mono", "Geist Mono"}
	prodFonts   = []string{"Geist", "Manrope", "JetBrains Mono"} // keep in sync with :root in page.html
	googleFonts = map[string]string{
		"Inter":             "Inter:wght@400;500;600;700;800",
		"Geist":             "Geist:wght@400;500;600;700;800",
		"Manrope":           "Manrope:wght@400;500;600;700;800",
		"DM Sans":           "DM+Sans:wght@400;500;600;700;800",
		"Plus Jakarta Sans": "Plus+Jakarta+Sans:wght@400;500;600;700;800",
		"Figtree":           "Figtree:wght@400;500;600;700;800",
		"Outfit":            "Outfit:wght@400;500;600;700;800",
		"Lexend":            "Lexend:wght@400;500;600;700;800",
		"Space Grotesk":     "Space+Grotesk:wght@400;500;600;700",
		"JetBrains Mono":    "JetBrains+Mono:wght@400;500",
		"IBM Plex Mono":     "IBM+Plex+Mono:wght@400;500",
		"Geist Mono":        "Geist+Mono:wght@400;500",
	}
)

func fontsHref(names ...string) string {
	q := make([]string, len(names))
	for i, n := range names {
		q[i] = "family=" + googleFonts[n]
	}
	return "https://fonts.googleapis.com/css2?" + strings.Join(q, "&") + "&display=swap"
}

type view struct {
	Tweaks               bool
	Phone, Email, Ctx    string
	FontsHref            string
	BodyFonts, MonoFonts []string
	Options              []Option
	Option               *Option
	Page                 *Page
	Title                string
	Step, Total, Pct     int
	Prev, Next           string
}

func newMux(tweaks bool) http.Handler {
	base := view{Tweaks: tweaks, Phone: phone, Email: email, Options: options, BodyFonts: bodyFonts, MonoFonts: monoFonts, FontsHref: fontsHref(prodFonts...)}
	if tweaks {
		base.FontsHref = fontsHref(append(append([]string{}, bodyFonts...), monoFonts...)...)
	}
	render := func(w http.ResponseWriter, v view) {
		var b bytes.Buffer
		if err := page.Execute(&b, v); err != nil {
			log.Println(err)
			http.Error(w, "template error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if !tweaks {
			w.Header().Set("Cache-Control", "public, max-age=300, s-maxage=31536000")
		}
		w.Write(b.Bytes())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		v := base
		v.Title = "Step 1"
		render(w, v)
	})
	mux.HandleFunc("GET /{opt}/{n}", func(w http.ResponseWriter, r *http.Request) {
		opt := r.PathValue("opt")
		flow, ok := flows[opt]
		n, err := strconv.Atoi(r.PathValue("n"))
		if !ok || err != nil || n < 2 || n > len(flow)+1 {
			http.NotFound(w, r)
			return
		}
		v := base
		for i := range options {
			if options[i].Key == opt {
				v.Option = &options[i]
			}
		}
		v.Page = &flow[n-2]
		v.Step, v.Total = n, len(flow)+1
		v.Pct = n * 100 / v.Total
		v.Title = fmt.Sprintf("Step %d of %d", n, v.Total)
		v.Ctx = fmt.Sprintf("Step %d of %d, Option %s", n, v.Total, strings.ToUpper(opt))
		v.Prev = "/"
		if n > 2 {
			v.Prev = fmt.Sprintf("/%s/%d", opt, n-1)
		}
		if n < v.Total {
			v.Next = fmt.Sprintf("/%s/%d", opt, n+1)
		}
		render(w, v)
	})
	mux.HandleFunc("POST /upload", upload)
	return mux
}

// upload stores a screenshot in Vercel Blob and answers {"url": ...}.
// ponytail: unauthenticated and unthrottled; add a rate limit if the URL ever leaks.
func upload(w http.ResponseWriter, r *http.Request) {
	token := os.Getenv("BLOB_READ_WRITE_TOKEN")
	if token == "" {
		http.Error(w, "uploads aren't set up on this server", http.StatusServiceUnavailable)
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 10<<20))
	if err != nil {
		http.Error(w, "that file is over 10MB", http.StatusRequestEntityTooLarge)
		return
	}
	ct := http.DetectContentType(body)
	ext, ok := map[string]string{"image/png": ".png", "image/jpeg": ".jpg", "image/gif": ".gif", "image/webp": ".webp"}[ct]
	if !ok {
		http.Error(w, "that doesn't look like a screenshot (PNG or JPG)", http.StatusUnsupportedMediaType)
		return
	}
	var id [12]byte
	rand.Read(id[:])
	req, _ := http.NewRequestWithContext(r.Context(), "PUT", "https://vercel.com/api/blob?pathname=shots/"+hex.EncodeToString(id[:])+ext, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("x-api-version", "13")
	req.Header.Set("x-vercel-blob-access", "public")
	req.Header.Set("x-content-type", ct)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("blob:", err)
		http.Error(w, "couldn't reach storage", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	var out struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil || resp.StatusCode != http.StatusOK || out.URL == "" {
		log.Println("blob:", resp.Status, err)
		http.Error(w, "storage rejected the upload", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func main() {
	tweaks := flag.Bool("tweaks", false, "show the Tweaks sidebar for local design A/B")
	flag.Parse()
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("listening on http://localhost:%s (tweaks=%v)", port, *tweaks)
	log.Fatal(http.ListenAndServe(":"+port, newMux(*tweaks)))
}
