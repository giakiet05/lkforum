# C4 Model Diagrams - LKForum

Kiến trúc hệ thống LKForum được mô tả theo mô hình C4 (Context, Container, Component, Code).

## Diagrams

### 1. System Context Diagram
**File:** `01-system-context.puml`

High-level overview of LKForum with actors and external systems:
- **Actors:** User, Admin
- **External Systems:** Google OAuth, Email Service
- **Relationships:** Interactions between actors and the LKForum system

### 2. Container Diagram
**File:** `02-container.puml`

Chi tiết các containers trong hệ thống:
- **User Web App** (Svelte 5, port 5173) - Giao diện người dùng
- **Admin Web App** (Svelte 5, port 5174) - Giao diện quản trị
- **Backend API** (Go + Gin, port 8080) - REST API & WebSocket
- **MongoDB Atlas** - Database chính
- **Redis** - Cache & session storage

### 3. Component Diagram - Backend API
**File:** `03-component-backend.puml`

Chi tiết các components trong Backend API:
- **HTTP Layer:** Router, Middlewares (Auth, CORS, Rate Limiter, Logger)
- **Handlers:** Auth, User, Community, Post, Comment, Message, Admin
- **Services:** Business logic layer
- **Repositories:** Data access layer
- **WebSocket Hub:** Real-time messaging

## Cách xem diagrams

### Option 1: PlantUML Online
Mở file `.puml` và copy content vào: https://www.plantuml.com/plantuml/uml/

### Option 2: VS Code Extension
1. Install extension: **PlantUML** by jebbs
2. Mở file `.puml`
3. Press `Alt+D` để preview

### Option 3: CLI (Local)
```bash
# Install PlantUML
brew install plantuml

# Generate PNG
plantuml 01-system-context.puml
plantuml 02-container.puml
plantuml 03-component-backend.puml

# Generate SVG (better quality)
plantuml -tsvg *.puml
```

### Option 4: Docker
```bash
docker run -it --rm \
  -v $(pwd):/data \
  plantuml/plantuml:latest \
  -tsvg /data/*.puml
```

## Technology Stack

### Frontend
- Svelte 5 (Runes API)
- TypeScript
- Vite

### Backend
- Go (Golang)
- Gin Framework
- Layered Architecture

### Database & Cache
- MongoDB Atlas (NoSQL)
- Redis (In-memory)

### External Services
- Google OAuth 2.0
- SMTP Email Service

## Architecture Highlights

✅ **Layered Architecture** - Clear separation: Handlers → Services → Repositories  
✅ **Real-time** - WebSocket Hub cho messaging  
✅ **Scalable** - Stateless API, Redis sessions  
✅ **Secure** - JWT auth, rate limiting, input validation  
✅ **Modular** - Dễ maintain và extend

## Notes

- Component diagram chỉ focus vào Backend API (frontend đơn giản)
- WebSocket Hub handle real-time messaging và online status
- Middleware chain: CORS → Logger → Rate Limiter → Auth
- Services chứa business logic, Repositories chỉ data access
- Redis dùng cho: sessions, cache, rate limiting, OTP storage
