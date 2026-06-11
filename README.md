# JobStream 🚀

JobStream is a premium, full-stack job aggregation platform built with **Go (Golang)**, **Nuxt 3 (Vue 3)**, and **PostgreSQL**. It is designed to ingest, normalize, and serve job postings from a variety of public APIs and company-specific ATS (Applicant Tracking System) platforms into a unified, high-performance search experience.

---

## 🎨 Architecture & Core Features

### 1. Clean DDD Backend Architecture (Go)
Adhering to Domain-Driven Design (DDD) and Clean Architecture principles, the backend provides an extremely scalable foundation:
- **HTTP / Handler Layer**: Custom router utilizing standard Go `net/http` (`ServeMux`) with pre-flight CORS middleware support.
- **Service Layer**: Coordinates business workflows, maps ingestion outputs, and processes jobs through category normalization and salary range parsing pipelines.
- **Domain Layer**: Clean, language-native entities (`Job`, `Company`, `Source`) and repository interfaces to keep the core business rules decoupled from framework dependencies.
- **Repository Layer**: PostgreSQL data access powered by `pgxpool` with dynamic query builders supporting page-based pagination, sorting, and advanced filters.

### 2. Pluggable & Dynamic Multi-Source Aggregation
JobStream pulls listings from multiple external environments concurrently:
- **Public Aggregator APIs**: Integrations for **Adzuna**, **Remotive**, and **WeWorkRemotely**.
- **Dynamic Company-Specific ATS Boards**: Native scrapers and clients for **Greenhouse**, **Lever**, and **Ashby** boards.
- **Database-Driven Sources**: Greenhouse, Lever, and Ashby fetchers load client targets dynamically from active companies stored in the PostgreSQL database, making the platform easily extensible.

### 3. Concurrent Pipelines & High-Performance Ingestion
- Ingestion runs in parallel using a Goroutine pool synchronized with `sync.WaitGroup` and monitored via dedicated error channels.
- Employs **batch processing (batch size of 500)** to execute high-performance bulk UPSERTs (`ON CONFLICT (source_id, platform) DO NOTHING`) in PostgreSQL, ensuring database writes remain rapid and deduplicated.
- Integrated with a background **scheduler** running every 6 hours to keep listings up-to-date automatically.

### 4. Advanced Search & Filtering Engines
- **Smart Category Normalization**: Raw categories are mapped using an ordered rule-based keyword classifier in `internal/category/classifier.go` into canonical classifications: `Engineering`, `Data`, `Product`, `Design`, `Marketing`, `Sales`, `People`, `Finance`, `Security`, `Operations`, `Customer Success`, `Legal`, or `Other`.
- **Robust Salary Parser**: A regex-based parser in `internal/salary/parser.go` that processes free-text salary strings (e.g. `"$80,000"`, `"120k - 150k"`, `"1.2m"`) and extracts numeric `salary_min` and `salary_max` values in USD.
- **Optimized PostgreSQL Indices**:
  - GIN indexes with `to_tsvector` for rapid, English full-text search on job `title`, `company`, `location`, and `category`.
  - B-tree and composite indexing on `posted_at DESC`, `platform`, `category`, and parsed `salary_min`/`salary_max` ranges for instantaneous filtering.

### 5. Thread-Safe Metadata Cache
- Employs an in-memory `MetadataCache` with a configurable Time-To-Live (TTL) to cache distinct platforms and categories, mitigating heavy redundant queries to PostgreSQL.
- Supports manual and scheduled cache invalidation utilities.

### 6. Interactive Nuxt 3 Frontend
- Built using **Nuxt 3**, **Vue 3**, and **Tailwind CSS** with full server-side rendering (SSR) support.
- Full **Dark Mode / Light Mode** support out-of-the-box (using Nuxt Color Mode with a customized `ColorModeToggle.vue` switcher).
- Custom components including a custom-styled `CustomSelect.vue` select-box dropdown, dynamic search/filter panels, and detailed cards.
- Manual **Sync Jobs** integration: A live user action triggering background synchronization pipelines with micro-animations, loading indicators, and skeleton loaders.

---

## 🏗️ Tech Stack

| Backend | Frontend | Database & Infra |
| :--- | :--- | :--- |
| Go (Golang) | Nuxt 3 / Vue 3 | PostgreSQL (Neon / local) |
| Standard `net/http` | Tailwind CSS | Docker & Docker Compose |
| `pgx` Connection Pool | TypeScript | GIN Full-Text Indexing |
| `godotenv` | Nuxt Color Mode | Render (Backend) / Vercel (Frontend) |

---

## 🚀 Running Locally

### 📋 Prerequisites
Ensure you have the following installed:
- Go 1.21+
- Node.js 18+ (with npm)
- Docker & Docker Compose

### 📦 Setup & Ingestion

#### 1. Clone the Repository
```bash
git clone https://github.com/lauralee01/jobstream.git
cd jobstream
```

#### 2. Start PostgreSQL
Start a local PostgreSQL instance via Docker Compose:
```bash
docker compose up -d
```

#### 3. Run Database Migrations
Run the SQL scripts located in the `/migrations` folder to set up schema, seed initial companies, and build optimal indices:
```bash
for file in migrations/*.sql; do
psql "postgres://username:password@localhost:5432/db_name?sslmode=disable" -f "$file"
done
```

---

## 🖥️ Backend Setup

```bash
cd backend
go mod tidy
```

### 🔑 Environment Variables
Create a `.env` file in the `backend` directory (using `.env.example` as a guide):
```env
DATABASE_URL=postgres://username:password@localhost:5432/db_name?sslmode=disable
PORT=8080

# External API Integrations (Optional)
ADZUNA_APP_ID=your_app_id
ADZUNA_APP_KEY=your_app_key
```

### ⚡ Run Server
```bash
go run cmd/api/main.go
```
The REST API will be available at `http://localhost:8080`.

### 🔌 API Routes
- `GET /health` - Service health status
- `POST /api/v1/jobs/sync` - Manually trigger ingestion and normalization pipelines
- `GET /api/v1/jobs` - Paginated job listings (supports query params: `keyword`, `location`, `category`, `remote`, `min_salary`, `platforms`, `page`, `limit`, `sort_by`, `sort_order`)
- `GET /api/v1/jobs/categories` - Fetch all normalized job categories
- `GET /api/v1/jobs/platforms` - Fetch all active job platforms

---

## 🌐 Frontend Setup

```bash
cd frontend
npm install
```

### 🔑 Environment Variables
Create a `.env` file in the `frontend` directory:
```env
NUXT_PUBLIC_API_BASE=http://localhost:8080/api/v1
```

### ⚡ Run Development Server
```bash
npm run dev
```
The application will be running at `http://localhost:3000`.
