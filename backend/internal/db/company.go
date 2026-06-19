package db

import (
	"context"
	"jobstream/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCompanyRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCompanyRepository(
	pool *pgxpool.Pool,
) *PostgresCompanyRepository {
	return &PostgresCompanyRepository{
		db: pool,
	}
}

func (r *PostgresCompanyRepository) GetEnabledByProvider(
	ctx context.Context,
	provider string,
) ([]domain.Company, error) {

	query := `
		SELECT
			id,
			name,
			slug,
			provider,
			enabled,
			created_at
		FROM companies
		WHERE provider = $1
		AND enabled = true
		ORDER BY name ASC;
	`

	rows, err := r.db.Query(
		ctx,
		query,
		provider,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	companies := []domain.Company{}

	for rows.Next() {

		var company domain.Company

		err := rows.Scan(
			&company.ID,
			&company.Name,
			&company.Slug,
			&company.Provider,
			&company.Enabled,
			&company.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return companies, nil
}
