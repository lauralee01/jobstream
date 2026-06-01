-- Add text search indexes for keyword filtering
CREATE INDEX IF NOT EXISTS idx_jobs_title_tsvector ON jobs USING gin(to_tsvector('english', title));
CREATE INDEX IF NOT EXISTS idx_jobs_company_tsvector ON jobs USING gin(to_tsvector('english', company));
CREATE INDEX IF NOT EXISTS idx_jobs_location_tsvector ON jobs USING gin(to_tsvector('english', location));
CREATE INDEX IF NOT EXISTS idx_jobs_category_tsvector ON jobs USING gin(to_tsvector('english', category));

-- Basic indexes for other filtering
CREATE INDEX IF NOT EXISTS idx_jobs_title ON jobs(title);
CREATE INDEX IF NOT EXISTS idx_jobs_company ON jobs(company);