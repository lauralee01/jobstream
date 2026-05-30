BEGIN;

-- Add new columns to jobs table
ALTER TABLE jobs
ADD COLUMN IF NOT EXISTS salary_min BIGINT,
ADD COLUMN IF NOT EXISTS salary_max BIGINT;

-- Index for filtering by minimum salary (most common)
CREATE INDEX IF NOT EXISTS idx_jobs_salary_min ON jobs(salary_min)
WHERE salary_min IS NOT NULL;

-- Index for filtering by maximum salary
CREATE INDEX IF NOT EXISTS idx_jobs_salary_max ON jobs(salary_max)
WHERE salary_max IS NOT NULL;

-- Composite index for range queries
CREATE INDEX IF NOT EXISTS idx_jobs_salary_range ON jobs(salary_min, salary_max)
WHERE salary_min IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_jobs_location ON jobs(location);
-- CREATE INDEX IF NOT EXISTS idx_jobs_title_tsvector ON jobs USING gin(to_tsvector('english', title));
-- CREATE INDEX IF NOT EXISTS idx_jobs_company_tsvector ON jobs USING gin(to_tsvector('english', company));

-- Update existing jobs by parsing their salary strings
-- This uses a simple numeric extraction regex for the initial backfill
UPDATE jobs
SET 
  salary_min = CASE
    WHEN salary ~ '^\d+' THEN
      (CAST(SUBSTRING(salary FROM '\d+') AS BIGINT)) *
      CASE
        WHEN salary ~* 'k' THEN 1000
        WHEN salary ~* 'm' THEN 1000000
        ELSE 1
      END
    ELSE NULL
  END,
  salary_max = CASE
    WHEN salary ~ '\d+' THEN
      (CAST(SUBSTRING(salary FROM '\d+') AS BIGINT)) *
      CASE
        WHEN salary ~* 'k' THEN 1000
        WHEN salary ~* 'm' THEN 1000000
        ELSE 1
      END
    ELSE NULL
  END
WHERE salary IS NOT NULL 
  AND salary != ''
  AND salary_min IS NULL;


-- Verify the migration worked
SELECT 
  COUNT(*) as total_jobs,
  COUNT(salary_min) as jobs_with_salary_parsed,
  COUNT(salary_max) as jobs_with_max_salary,
  ROUND(100.0 * COUNT(salary_min) / COUNT(*), 2) as parse_percentage
FROM jobs;

COMMIT;
