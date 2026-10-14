# Ditch the Vibes

A tiny step-by-step guide for a friend who vibe-coded some apps and needs to share the code on GitHub. One Go binary, zero dependencies, no database.

```sh
vercel env pull .env.local            # once: gets BLOB_READ_WRITE_TOKEN for screenshot uploads
set -a; . ./.env.local; set +a
go run .                              # http://localhost:3000
go run . -tweaks                      # same, plus a design A/B sidebar
go test ./...
```

Screenshots are uploaded to Vercel Blob (`POST /upload`, images under 10MB) and the resulting URL is dropped into a pre-filled text or email. Without the token the button reports that uploads aren't set up.

Deploys to Vercel with the Go framework preset (`vercel.json`). The server reads `PORT`.
