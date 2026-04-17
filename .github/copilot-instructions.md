# Copilot Instructions

## Build, Test, Lint
- Backend (Go): `cd backend && go run main.go` to serve; `go test ./...` for all tests; single test `go test ./<pkg> -run TestName`. No repo-defined linter—use `gofmt`/`go vet` as needed.
- Admin web (Svelte): `cd frontend/admin-web && npm install && npm run dev`; `npm run build`; type/lint check `npm run check`. Same scripts available with `pnpm` (lockfile present).
- User web (Svelte): `cd frontend/user-web && npm install && npm run dev`; `npm run build`; `npm run check`.
- Docker: `docker-compose up --build` brings up redis, backend (8081→8080), user web (3001), admin web (3004) using `backend/.env`.
- Utility scripts: from `backend` run `go run cmd/seed_admin/main.go` (create admin), `go run cmd/migrate/main.go` (migrate moderator avatars), `go run cmd/clear_notifications/main.go` (clean legacy notices).

## Architecture Overview
- Backend: Go Gin API wired in `internal/bootstrap/init.go` (config → Mongo client/Redis/SMTP/Cloudinary/Google OAuth/Gemini). Repos (`internal/repo`) → services (`internal/service`) → controllers (`internal/controller`) → routes (`internal/route/*`), mounted under `/api` (auth, users, communities, posts, votes, comments, notifications, drafts, reports, channels/messages, websocket, admin suites). Event bus + WebSocket hub (`internal/platform/ws`) plus background services for reputation, notifications, messaging, channels, communities, and AI moderation (Gemini). Config via `.env` (see `internal/config/config.go`); Mongo collections enumerated in `internal/config/collections.go`.
- Frontend: Two Vite/Svelte 5 SPAs. `frontend/user-web` uses `src/routes.ts` (home, popular, explore, auth, profile, communities, posts, mod tools, messages, Google OAuth flows). `frontend/admin-web` routes in `src/routes.ts` (login, dashboard, users, communities, reports). Both are SPA-router based (`svelte-spa-router`), assets/components under `src`.
- Domain docs: requirements/features in `CHUONG1_TONG_QUAN_DE_TAI.md`; DB schema in `DATABASE_DESIGN.md`; UI flows in `UI_DESIGN.md`; user stories in `User Story.md`; admin setup in `backend/docs/CREATE_ADMIN.md`.

## Key Conventions
- Always load env via `backend/.env`; CORS origins driven by `ALLOWED_ORIGINS`; JWT, Redis, SMTP, Cloudinary, Google OAuth, and Gemini keys pulled from env. Default ports: backend 8080 (mapped 8081 in compose), user web 5173 (served at 3001 in compose), admin web 5174 (served at 3004).
- API additions should follow repo → service → controller → route pattern and be registered in `internal/bootstrap/init.go`; add collection constants in `internal/config/collections.go` when new Mongo collections appear.
- Background services (`Start()` on reputation/notification/message/channel/community/moderation) are kicked off in `bootstrap.Init`—keep them aligned with event bus expectations when adding features.
- WebSocket/event features rely on the shared bus (`internal/platform/bus`) and hub in `internal/platform/ws`; notifications and messages assume those hooks.
- Frontend navigation is centralized in each `src/routes.ts`; new pages should be registered there and wired to auth flows in `src/auth`/`stores`.
