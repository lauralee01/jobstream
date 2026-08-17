-- 012_optimize_search_and_filter_indexes.sql
-- Composite partial index for the primary active job listing feed
CREATE INDEX IF NOT EXISTS idx_jobs_active_posted ON jobs (posted_at DESC) WHERE active = true;

-- Combined text search index across title, company, location, and category for fast multi-field keyword search
CREATE INDEX IF NOT EXISTS idx_jobs_fts_combined ON jobs USING gin(
  to_tsvector('english', coalesce(title, '') || ' ' || coalesce(company, '') || ' ' || coalesce(location, '') || ' ' || coalesce(category, ''))
);
