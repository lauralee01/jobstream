CREATE INDEX IF NOT EXISTS idx_jobs_active_posted ON jobs (posted_at DESC) WHERE active = true;

-- Drop previous unweighted/category FTS index if exists
DROP INDEX IF EXISTS idx_jobs_fts_combined;

-- Weighted text search index across title (weight A) and company & location (weight B)
CREATE INDEX IF NOT EXISTS idx_jobs_fts_weighted ON jobs USING gin(
  (setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
   setweight(to_tsvector('english', coalesce(company, '') || ' ' || coalesce(location, '')), 'B'))
);
