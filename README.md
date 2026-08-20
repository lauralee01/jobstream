# JobStream 🚀

JobStream is a full-stack job aggregation platform built with **Go (Golang)**, **Nuxt 4 (Vue 3)**, and **PostgreSQL**. It ingests, normalizes, and serves job postings from public APIs and company-specific ATS (Applicant Tracking System) platforms through a unified search and filtering experience.

---

## 🎨 Architecture & Core Features

### 1. Clean Backend Architecture (Go)
The backend follows a layered architecture that keeps business logic separate from HTTP and database concerns:
- **HTTP / Handler Layer**: Uses standard Go `net/http` (`ServeMux`) with CORS middleware and request handling.
- **Service Layer**: Coordinates business workflows, ingestion, category normalization, salary parsing, and job retrieval.
- **Domain Layer**: Defines core entities (`Job`, `Company`, `Source`) and repository interfaces independent of infrastructure concerns.
- **Repository Layer**: PostgreSQL data access powered by `pgxpool`, with dynamic filtering, pagination, sorting, full-text search, and batch UPSERTs.
- **Database Pooling & Shutdown**: Explicit `pgxpool` connection configuration and graceful `SIGINT` / `SIGTERM` handling allow active requests to finish and database resources to close cleanly.

### 2. Pluggable & Dynamic Multi-Source Aggregation
JobStream pulls listings concurrently from multiple external sources:
- **Public Aggregator APIs**: Integrations for **Adzuna**, **Remotive**, and **WeWorkRemotely**.
- **Company-Specific ATS Boards**: Integrations for **Greenhouse**, **Lever**, and **Ashby**.
- **Database-Driven Sources**: Greenhouse, Lever, and Ashby fetchers load company targets dynamically from active companies stored in PostgreSQL, making new sources easier to add without changing the core ingestion pipeline.

### 3. Concurrent Ingestion Pipeline
- Job sources are fetched concurrently using Go goroutines coordinated with `sync.WaitGroup` and error channels.
- Jobs are processed in **batches of 500** and written using multi-row PostgreSQL UPSERTs.
- `ON CONFLICT (source_id, platform) DO UPDATE` prevents duplicate listings while allowing existing jobs to be refreshed with updated information.
- A background **scheduler runs every 6 hours** to keep listings current automatically.
- Stale jobs can be marked inactive and older inactive listings removed separately.

### 4. Search, Normalization & Filtering
- **PostgreSQL Full-Text Search**: Uses a weighted GIN index with `to_tsvector` and `websearch_to_tsquery` for efficient multi-word keyword searches.
- **Relevance Ranking**: Job titles receive the highest search weight, with matching results ranked using `ts_rank` before recency.
- **Smart Category Normalization**: Raw job data is mapped using an ordered rule-based classifier into canonical categories including `Engineering`, `Data`, `Product`, `Design`, `Marketing`, `Sales`, `People`, `Finance`, `Security`, `Operations`, `Customer Success`, `Legal`, and `Other`.
- **Salary Parsing**: A regex-based parser handles salary formats such as `"$80,000"`, `"120k - 150k"`, and `"1.2m"` and extracts numeric `salary_min` and `salary_max` values for filtering.
- **Advanced Filters**: Supports filtering by location, category, minimum salary, remote status, platform, pagination, and sorting.
- **Optimized PostgreSQL Indexes**: Includes a GIN full-text search index and a partial composite index on active jobs ordered by `posted_at DESC`.

### 5. Metadata Caching & Testing
- An in-memory `MetadataCache` with configurable Time-To-Live (TTL) caches frequently requested platform and category metadata to reduce redundant PostgreSQL queries.
- Supports cache invalidation when underlying metadata changes.
- Automated tests cover core logic including **salary parsing, category classification, remote-role detection, and query parameter parsing**.

### 6. Interactive Nuxt 4 Frontend
- Built with **Nuxt 4**, **Vue 3**, **TypeScript**, and **Tailwind CSS** with server-side rendering (SSR).
- Supports **Dark Mode / Light Mode** using Nuxt Color Mode with a custom `ColorModeToggle.vue`.
- Includes reusable components for job cards, search, filtering, custom select inputs, pagination, and loading states.
- Search and filter state is integrated with Nuxt's reactive data-fetching flow.
- Manual **Sync Jobs** functionality allows users to trigger the ingestion pipeline with loading feedback and UI animations.

---

## 🏗️ Tech Stack

| Backend | Frontend | Database & Infra |
| :--- | :--- | :--- |
| Go (Golang) | Nuxt 4 / Vue 3 | PostgreSQL (Neon / local) |
| Standard `net/http` | TypeScript | GIN Full-Text Search |
| `pgx` / `pgxpool` | Tailwind CSS | Docker & Docker Compose |
| Goroutines | Nuxt Color Mode | Render (Backend) / Vercel (Frontend) |
| Go Testing | SSR | PostgreSQL Indexing |

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
Run the SQL scripts located in the `/migrations` folder to set up the schema, seed initial companies, and create the required indexes:
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
Create a `.env` file in the `backend` directory using `.env.example` as a guide:

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
- `POST /api/v1/jobs/sync` - Manually trigger the ingestion and normalization pipeline
- `GET /api/v1/jobs` - Paginated job listings with search, filtering, and sorting
- `GET /api/v1/jobs/categories` - Fetch normalized job categories
- `GET /api/v1/jobs/platforms` - Fetch active job platforms

`GET /api/v1/jobs` supports the following query parameters:

`keyword`, `location`, `category`, `remote`, `min_salary`, `platforms`, `page`, `limit`, `sort_by`, `sort_order`

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