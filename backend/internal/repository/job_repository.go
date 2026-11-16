package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/deadlium/eventbeep/backend/internal/models"
)

var (
	ErrJobNotFound = errors.New("job not found")
)

// JobRepository defines the data access contract for jobs.
type JobRepository interface {
	GetJobs(ctx context.Context, filter models.JobFilter) ([]models.Job, int64, error)
	GetJobByID(ctx context.Context, id uuid.UUID) (*models.Job, error)
	CreateJob(ctx context.Context, job *models.Job) error
}

// PostgresJobRepository implements JobRepository backed by pgxpool.
type PostgresJobRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresJobRepository creates a new repository using the provided pool.
func NewPostgresJobRepository(pool *pgxpool.Pool) *PostgresJobRepository {
	return &PostgresJobRepository{pool: pool}
}

// GetJobs fetches a paginated, filtered list of jobs and the matching total count.
func (r *PostgresJobRepository) GetJobs(ctx context.Context, filter models.JobFilter) ([]models.Job, int64, error) {
	var whereClauses []string
	var args []interface{}
	argIdx := 1

	if strings.TrimSpace(filter.Role) != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("role ILIKE $%d", argIdx))
		args = append(args, "%"+strings.TrimSpace(filter.Role)+"%")
		argIdx++
	}

	if strings.TrimSpace(filter.City) != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("city ILIKE $%d", argIdx))
		args = append(args, "%"+strings.TrimSpace(filter.City)+"%")
		argIdx++
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// 1. Count query
	countSQL := "SELECT COUNT(*) FROM jobs" + whereSQL
	var total int64
	countErr := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total)
	if countErr != nil {
		return nil, 0, fmt.Errorf("failed to count jobs: %w", countErr)
	}

	// 2. Fetch paginated records
	selectSQL := fmt.Sprintf(
		"SELECT id, hash_key, title, company, role, city, source_url, created_at FROM jobs%s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		whereSQL,
		argIdx,
		argIdx+1,
	)

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	queryArgs := append(args, limit, offset)

	rows, err := r.pool.Query(ctx, selectSQL, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute get jobs query: %w", err)
	}
	defer rows.Close()

	jobs := make([]models.Job, 0)
	for rows.Next() {
		var job models.Job
		var role sql.NullString
		var city sql.NullString

		if err := rows.Scan(
			&job.ID,
			&job.HashKey,
			&job.Title,
			&job.Company,
			&role,
			&city,
			&job.SourceURL,
			&job.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan job row: %w", err)
		}

		if role.Valid {
			job.Role = &role.String
		}
		if city.Valid {
			job.City = &city.String
		}

		jobs = append(jobs, job)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return jobs, total, nil
}

// GetJobByID retrieves a single job by its UUID.
func (r *PostgresJobRepository) GetJobByID(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	query := `
		SELECT id, hash_key, title, company, role, city, source_url, created_at
		FROM jobs
		WHERE id = $1
	`

	var job models.Job
	var role sql.NullString
	var city sql.NullString

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&job.ID,
		&job.HashKey,
		&job.Title,
		&job.Company,
		&role,
		&city,
		&job.SourceURL,
		&job.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrJobNotFound
		}
		return nil, fmt.Errorf("failed to fetch job by id: %w", err)
	}

	if role.Valid {
		job.Role = &role.String
	}
	if city.Valid {
		job.City = &city.String
	}

	return &job, nil
}

// CreateJob inserts a job record, using ON CONFLICT (hash_key) DO NOTHING for idempotency.
func (r *PostgresJobRepository) CreateJob(ctx context.Context, job *models.Job) error {
	_, err := r.CreateJobWithStatus(ctx, job)
	return err
}

// CreateJobWithStatus inserts a job record and returns whether a new row was inserted.
// Returns inserted = false when skipped due to ON CONFLICT (hash_key) DO NOTHING.
func (r *PostgresJobRepository) CreateJobWithStatus(ctx context.Context, job *models.Job) (bool, error) {
	if job.ID == uuid.Nil {
		job.ID = uuid.New()
	}
	if job.CreatedAt.IsZero() {
		job.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO jobs (id, hash_key, title, company, role, city, source_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (hash_key) DO NOTHING
	`

	tag, err := r.pool.Exec(ctx, query,
		job.ID,
		job.HashKey,
		job.Title,
		job.Company,
		job.Role,
		job.City,
		job.SourceURL,
		job.CreatedAt,
	)

	if err != nil {
		return false, fmt.Errorf("failed to insert job: %w", err)
	}

	return tag.RowsAffected() > 0, nil
}
