-- Speed up common list filters and default sort (posted_at DESC).
CREATE INDEX IF NOT EXISTS idx_jobs_posted_at ON jobs (posted_at DESC);
CREATE INDEX IF NOT EXISTS idx_jobs_platform ON jobs (platform);
CREATE INDEX IF NOT EXISTS idx_jobs_category ON jobs (category);
CREATE INDEX IF NOT EXISTS idx_jobs_created_at ON jobs (created_at DESC);
