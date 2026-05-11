CREATE TABLE jobs (
    id TEXT PRIMARY KEY,
    source_id TEXT,
    platform TEXT,
    title TEXT,
    company TEXT,
    location TEXT,
    description TEXT,
    url TEXT,
    salary TEXT,
    posted_at TIMESTAMP,
    created_at TIMESTAMP
);
