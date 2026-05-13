package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"jobstream/internal/db"
	"jobstream/internal/fetcher"
	"jobstream/internal/fetcher/linkedin"
	"jobstream/internal/jobs"
	"jobstream/internal/scheduler"

	apphttp "jobstream/internal/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("godotenv: %v (using environment variables only)", err)
	}
	// Load database URL from environment variables
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// 2. Create Postgres connection pool
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	// 3. Initialize Repository
	repo := db.NewPostgresJobRepository(pool)

	// 4. Register fetchers
	fetchers := []fetcher.Fetcher{
		&fetcher.MockFetcher{},
		&linkedin.Fetcher{},
	}

	// 5. Initialize Job Service
	jobService := jobs.NewJobService(repo, fetchers)

	// 6. Start Scheduler (runs in background)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler := scheduler.NewScheduler(jobService, 10*time.Hour)
	scheduler.Start(ctx)

	// 6. Initialize HTTP Router with job service
	router := apphttp.NewRouter(jobService)

	// 7. Start server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
