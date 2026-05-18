package db

import (
	"context"
	"fmt"
	"jobstream/internal/domain"
	"log"
	"strings"

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

func (r *PostgresJobRepository) Save(ctx context.Context, job *domain.Job) error {
	_, err := r.db.Exec(ctx, "INSERT INTO jobs (id, source_id, platform, title, company, location, category, description, url, salary, posted_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) ON CONFLICT (id) DO NOTHING", job.ID, job.SourceID, job.Platform, job.Title, job.Company, job.Location, job.Category, job.Description, job.URL, job.Salary, job.PostedAt, job.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// FindAll retrieves jobs from the database based on the provided filter.
// It returns a slice of jobs, the total count of matching jobs, and an error if any occurs.
// The method supports filtering by keyword, category, location, and platform.
// It also supports sorting and pagination.
func (r *PostgresJobRepository) FindAll(
	ctx context.Context,
	filter domain.JobFilter,
) ([]domain.Job, int64, error) {

	baseQuery := `
		SELECT
			id,
			source_id,
			platform,
			title,
			company,
			location,
			COALESCE(category, ''),
			description,
			url,
			salary,
			posted_at,
			created_at
		FROM jobs
	`

	countQuery := `SELECT COUNT(*) FROM jobs`

	conditions := []string{}
	args := []interface{}{}
	paramIdx := 1

	log.Println("filter", filter)
	

	// =========================
	// Keyword Search
	// =========================

	if filter.Keyword != "" {

		conditions = append(
			conditions,
			fmt.Sprintf(`
				(
					title ILIKE $%d OR
					company ILIKE $%d OR
					location ILIKE $%d OR
					category ILIKE $%d
				)
			`,
				paramIdx,
				paramIdx,
				paramIdx,
				paramIdx,
			),
		)

		args = append(args, "%"+filter.Keyword+"%")
		paramIdx++
	}

	// =========================
	// Category Filter
	// =========================

	if filter.Category != "" {

		conditions = append(
			conditions,
			fmt.Sprintf("category ILIKE $%d", paramIdx),
		)

		args = append(args, "%"+filter.Category+"%")
		paramIdx++
	}

	// =========================
	// Location Filter
	// =========================

	if filter.Location != "" {

		conditions = append(
			conditions,
			fmt.Sprintf("location ILIKE $%d", paramIdx),
		)

		args = append(args, "%"+filter.Location+"%")
		paramIdx++
	}

	// =========================
	// Platform Filter
	// =========================

	if len(filter.Platforms) > 0 {

		platformConditions := []string{}

		for _, platform := range filter.Platforms {

			if platform == "" {
				continue
			}

			platformConditions = append(
				platformConditions,
				fmt.Sprintf("platform = $%d", paramIdx),
			)

			args = append(args, platform)

			paramIdx++
		}

		if len(platformConditions) > 0 {
			conditions = append(
				conditions,
				"("+strings.Join(platformConditions, " OR ")+")",
			)
		}
	}

	// =========================
	// Build WHERE clause
	// =========================

	query := baseQuery
	countSQL := countQuery

	if len(conditions) > 0 {

		whereClause := " WHERE " + strings.Join(conditions, " AND ")

		query += whereClause
		countSQL += whereClause
	}

	// =========================
	// Sorting
	// =========================

	allowedSortColumns := map[string]bool{
		"created_at": true,
		"posted_at":  true,
		"title":      true,
		"company":    true,
	}

	sortBy := "created_at"

	if allowedSortColumns[filter.SortBy] {
		sortBy = filter.SortBy
	}

	sortOrder := "DESC"

	if strings.ToUpper(filter.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}

	query += fmt.Sprintf(
		" ORDER BY %s %s",
		sortBy,
		sortOrder,
	)

	// =========================
	// Pagination
	// =========================

	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	offset := (filter.Page - 1) * filter.Limit

	query += fmt.Sprintf(
		" LIMIT $%d OFFSET $%d",
		paramIdx,
		paramIdx+1,
	)

	args = append(args, filter.Limit, offset)

	// =========================
	// Total Count
	// =========================

	var total int64

	err := r.db.QueryRow(
		ctx,
		countSQL,
		args[:len(args)-2]...,
	).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	// =========================
	// Execute Query
	// =========================

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	jobs := []domain.Job{}

	for rows.Next() {

		var job domain.Job

		err := rows.Scan(
			&job.ID,
			&job.SourceID,
			&job.Platform,
			&job.Title,
			&job.Company,
			&job.Location,
			&job.Category,
			&job.Description,
			&job.URL,
			&job.Salary,
			&job.PostedAt,
			&job.CreatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	log.Println("Found", total, "jobs", jobs)

	return jobs, total, nil
}