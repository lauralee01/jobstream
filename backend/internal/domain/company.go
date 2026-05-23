package domain

import (
	"context"
	"time"
)

type Company struct {
	ID        string
	Name      string
	Slug      string
	Provider  string
	Enabled   bool
	CreatedAt time.Time
}

type CompanyRepository interface {
	GetEnabledByProvider(
		ctx context.Context,
		provider string,
	) ([]Company, error)
}
