ALTER TABLE jobs
ADD CONSTRAINT unique_source_platform
UNIQUE(source_id, platform);