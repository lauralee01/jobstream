package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"jobstream/internal/db"
	"jobstream/internal/fetcher"
	"jobstream/internal/fetcher/adzuna"
	"jobstream/internal/fetcher/greenhouse"
	"jobstream/internal/fetcher/lever"
	"jobstream/internal/fetcher/remotive"
	"jobstream/internal/fetcher/remoteok"
	"jobstream/internal/fetcher/weworkremotely"
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

	// 3. Initialize Job Repository
	repo := db.NewPostgresJobRepository(pool)

	// 4. Initialize Company Repository
	companyRepo := db.NewPostgresCompanyRepository(pool)

	greenhouseCompanies, err := companyRepo.GetEnabledByProvider(
		context.Background(),
		"greenhouse",
	)
	if err != nil {
		log.Fatalf("failed to load greenhouse companies: %v", err)
	}

	leverCompanies, err := companyRepo.GetEnabledByProvider(
		context.Background(),
		"lever",
	)
	if err != nil {
		log.Fatalf("failed to load lever companies: %v", err)
	}

	// 5. Initialize Fetchers
	fetchers := []fetcher.Fetcher{
		remotive.NewClient(),
		adzuna.NewAPIClient(),
		weworkremotely.NewClient(),
		remoteok.NewClient(),
	}

	for _, company := range greenhouseCompanies {

		log.Printf(
			"Registering Greenhouse fetcher for: %s",
			company.Slug,
		)

		fetchers = append(
			fetchers,
			greenhouse.NewClient(company.Slug),
		)
	}

	for _, company := range leverCompanies {

		log.Printf(
			"Registering Lever fetcher for: %s",
			company.Slug,
		)

		fetchers = append(
			fetchers,
			lever.NewClient(company.Slug),
		)
	}

	// 5. Initialize Job Service
	jobService := jobs.NewJobService(repo, fetchers)

	// 6. Start Scheduler (runs in background)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	scheduler := scheduler.NewScheduler(jobService, 6*time.Hour)
	scheduler.Start(ctx)

	// 6. Initialize HTTP Router with job service
	router := apphttp.NewRouter(jobService)

	// 7. Start server
	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(server.ListenAndServe())
}
