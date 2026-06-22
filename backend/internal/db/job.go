package db

import (
	"context"
	"fmt"
	"jobstream/internal/domain"
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

func (r *PostgresJobRepository) Save(ctx context.Context, jobs []domain.Job) error {
	if len(jobs) == 0 {
		return nil
	}

	const cols = 17 // number of columns in INSERT
	valueStrings := make([]string, 0, len(jobs))
	valueArgs := make([]interface{}, 0, len(jobs)*cols)

	for i, job := range jobs {
		// 1-based parameter index
		base := i*cols + 1

		valueStrings = append(valueStrings,
			fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
				base+0, base+1, base+2, base+3, base+4, base+5, base+6,
				base+7, base+8, base+9, base+10, base+11, base+12, base+13, base+14, base+15, base+16,
			),
		)

		valueArgs = append(valueArgs,
			job.ID,
			job.SourceID,
			job.Platform,
			job.Title,
			job.Company,
			job.Location,
			job.Category,
			job.Description,
			job.URL,
			job.Salary,
			job.SalaryMin,
			job.SalaryMax,
			job.IsRemote,
			job.Active,
			job.LastSeenAt,
			job.PostedAt,
			job.CreatedAt,
		)
	}

	query := `
        INSERT INTO jobs (
            id, source_id, platform, title, company, location, category,
            description, url, salary, salary_min, salary_max, is_remote, active, last_seen_at, posted_at, created_at
        ) VALUES ` + strings.Join(valueStrings, ",") + `
        ON CONFLICT (source_id, platform) DO UPDATE SET
            title = EXCLUDED.title,
            company = EXCLUDED.company,
            location = EXCLUDED.location,
            category = EXCLUDED.category,
            description = EXCLUDED.description,
            url = EXCLUDED.url,
            salary = EXCLUDED.salary,
            salary_min = EXCLUDED.salary_min,
            salary_max = EXCLUDED.salary_max,
            is_remote = EXCLUDED.is_remote,
            active = EXCLUDED.active,
            last_seen_at = EXCLUDED.last_seen_at,
            posted_at = EXCLUDED.posted_at,
    `

	_, err := r.db.Exec(ctx, query, valueArgs...)
	return err
}

// FindAll retrieves jobs from the database based on the provided filter.
func (r *PostgresJobRepository) FindAll(
	ctx context.Context,
	filter domain.JobFilter,
) ([]domain.Job, int64, error) {

	// List view omits description — large TEXT column not needed in cards.
	baseQuery := `
		SELECT
			id,
			source_id,
			platform,
			title,
			company,
			location,
			COALESCE(category, ''),
			url,
			salary,
			salary_min,
			salary_max,
			is_remote,
			active,
			last_seen_at,
			posted_at,
			created_at
		FROM jobs
	`

	countQuery := `SELECT COUNT(*) FROM jobs`

	conditions := []string{"active = true"}
	args := []interface{}{}
	paramIdx := 1

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
	// Min Salary Filter (OPTIMIZED)
	// =========================
	// Now uses simple numeric comparison on pre-parsed salary_min column
	// instead of expensive regex parsing on every row

	if filter.MinSalary != nil {

		conditions = append(
			conditions,
			fmt.Sprintf("salary_min IS NOT NULL AND salary_min >= $%d", paramIdx),
		)

		args = append(args, *filter.MinSalary)
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
	// Remote Only Filter
	// =========================

	if filter.IsRemote != nil && *filter.IsRemote {
		conditions = append(conditions, "is_remote = true")
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

	sortBy := "posted_at"

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

	// FIX: Save the args count BEFORE adding pagination params
	// This ensures the count query gets the right parameter values
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)

	// Now append pagination args to the main query args
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
		countArgs...,
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
			&job.URL,
			&job.Salary,
			&job.SalaryMin,
			&job.SalaryMax,
			&job.IsRemote,
			&job.Active,
			&job.LastSeenAt,
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

	return jobs, total, nil
}

func (r *PostgresJobRepository) GetCategories(ctx context.Context) ([]string, error) {
	query := `
		SELECT DISTINCT category 
		FROM jobs 
		WHERE category IS NOT NULL AND category != '' 
		ORDER BY category ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *PostgresJobRepository) GetPlatforms(ctx context.Context) ([]string, error) {
	query := `
		SELECT DISTINCT platform 
		FROM jobs 
		WHERE platform IS NOT NULL AND platform != '' 
		ORDER BY platform ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var platforms []string
	for rows.Next() {
		var platform string
		if err := rows.Scan(&platform); err != nil {
			return nil, err
		}
		platforms = append(platforms, platform)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return platforms, nil
}

func (r *PostgresJobRepository) MarkStaleInactive(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `
		UPDATE jobs
		SET active = false
		WHERE last_seen_at < NOW() - INTERVAL '30 days'
	`)
	return err
}

func (r *PostgresJobRepository) DeleteOldInactive(ctx context.Context) error {
	_, err := r.db.Exec(ctx, `
		DELETE FROM jobs
		WHERE active = false
		AND last_seen_at < NOW() - INTERVAL '180 days'
	`)
	return err
}
