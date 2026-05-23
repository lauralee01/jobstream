CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,

    provider TEXT NOT NULL,

    enabled BOOLEAN DEFAULT true,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);