-- EventBeep PostgreSQL Database Schema
-- Phase 1: Jobs Table Schema

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hash_key VARCHAR(64) UNIQUE NOT NULL, -- SHA-256 of lowercase(company_name + job_title)
    title VARCHAR(255) NOT NULL,
    company VARCHAR(255) NOT NULL,
    role VARCHAR(100),
    city VARCHAR(100),
    source_url VARCHAR(1024) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for optimized querying, filtering, and pagination
CREATE INDEX IF NOT EXISTS idx_jobs_hash_key ON jobs(hash_key);
CREATE INDEX IF NOT EXISTS idx_jobs_role ON jobs(role);
CREATE INDEX IF NOT EXISTS idx_jobs_city ON jobs(city);
CREATE INDEX IF NOT EXISTS idx_jobs_created_at ON jobs(created_at DESC);
