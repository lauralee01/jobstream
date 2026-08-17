package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"jobstream/internal/db"
	"jobstream/internal/fetcher"
	"jobstream/internal/fetcher/adzuna"
	"jobstream/internal/fetcher/ashby"
	"jobstream/internal/fetcher/greenhouse"
	"jobstream/internal/fetcher/lever"
	"jobstream/internal/fetcher/remotive"
	"jobstream/internal/fetcher/weworkremotely"
	"jobstream/internal/fetcher/workable"
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

	// 2. Parse and configure Postgres connection pool
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v", err)
	}
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.MaxConnLifetime = 1 * time.Hour

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	// 3. Initialize Job Repository
	jobRepo := db.NewPostgresJobRepository(pool)

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

	ashbyCompanies, err := companyRepo.GetEnabledByProvider(
		context.Background(),
		"ashby",
	)
	if err != nil {
		log.Fatalf("failed to load ashby companies: %v", err)
	}

	workableCompanies, err := companyRepo.GetEnabledByProvider(
		context.Background(),
		"workable",
	)
	if err != nil {
		log.Fatalf("failed to load workable companies: %v", err)
	}

	// 5. Initialize Fetchers
	fetchers := []fetcher.Fetcher{
		remotive.NewClient(),
		adzuna.NewAPIClient(),
		weworkremotely.NewClient(),
	}

	for _, company := range greenhouseCompanies {
		fetchers = append(
			fetchers,
			greenhouse.NewClient(company.Slug),
		)
	}

	for _, company := range leverCompanies {
		fetchers = append(
			fetchers,
			lever.NewClient(company.Slug),
		)
	}

	for _, company := range ashbyCompanies {
		fetchers = append(
			fetchers,
			ashby.NewClient(company.Slug),
		)
	}

	for _, company := range workableCompanies {
		fetchers = append(
			fetchers,
			workable.NewClient(company.Slug),
		)
	}

	// 5. Initialize Job Service
	jobService := jobs.NewJobService(jobRepo, fetchers)

	// 6. Start Scheduler (runs in background)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sched := scheduler.NewScheduler(jobService, 6*time.Hour)
	sched.Start(ctx)

	// 6. Initialize HTTP Router with job service
	router := apphttp.NewRouter(jobService)

	// 7. Start server
	server := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Println("🚀 Server running on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting cleanly")
}
