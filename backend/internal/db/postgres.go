package db

import (
	"context"
	"jobstream/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresJobRepository implements the domain.JobRepository interface.
// This is called the 'Repository Pattern'. It keeps all database-specific code
type PostgresJobRepository struct {
	db *pgxpool.Pool
}

// NewPostgresJobRepository creates a new repository instance and injects the database connection pool.
func NewPostgresJobRepository(pool *pgxpool.Pool) *PostgresJobRepository {
	return &PostgresJobRepository{
		db: pool,
	}
}

// TODO: Implement the methods defined in the domain.JobRepository interface.
func (r *PostgresJobRepository) Save(ctx context.Context, job *domain.Job) error {
	_, err := r.db.Exec(ctx, "INSERT INTO jobs (id, source_id, platform, title, company, location, category, description, url, salary, posted_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT (id) DO NOTHING", job.ID, job.SourceID, job.Platform, job.Title, job.Company, job.Location, job.Category, job.Description, job.URL, job.Salary, job.PostedAt, job.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresJobRepository) FindAll(ctx context.Context) ([]domain.Job, error) {
	rows, err := r.db.Query(ctx, "SELECT id, source_id, platform, title, company, location, category, description, url, salary, posted_at, created_at FROM jobs")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []domain.Job
	for rows.Next() {
		var job domain.Job
		err := rows.Scan(&job.ID, &job.SourceID, &job.Platform, &job.Title, &job.Company, &job.Location, &job.Category, &job.Description, &job.URL, &job.Salary, &job.PostedAt, &job.CreatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return jobs, nil
}
