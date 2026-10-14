# Ditch the Vibes

A tiny step-by-step guide for a friend who vibe-coded some apps and needs to share the code on GitHub. One Go binary, zero dependencies, no database, no forms.

```sh
go run .            # http://localhost:3000
go run . -tweaks    # same, plus a design A/B sidebar (fonts, accent, theme, sizing)
go test ./...
```

Deploys to Vercel with the Go framework preset (`vercel.json`). The server reads `PORT`.
