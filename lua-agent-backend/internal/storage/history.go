package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const insertHistoryQuery = `
INSERT INTO histories (
    id,
    session_id,
    user_prompt,
    clarified_prompt,
    generated_code,
    validation_status,
    validation_errors,
    success,
    model_name,
    latency_ms,
    input_tokens,
    output_tokens,
    metadata
)
VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8, $9, $10, $11, $12, $13::jsonb)
RETURNING
    id,
    session_id,
    user_prompt,
    clarified_prompt,
    generated_code,
    validation_status,
    validation_errors,
    success,
    model_name,
    latency_ms,
    input_tokens,
    output_tokens,
    metadata,
    created_at,
    updated_at
`

const recentSuccessQuery = `
SELECT
    id,
    session_id,
    user_prompt,
    clarified_prompt,
    generated_code,
    validation_status,
    validation_errors,
    success,
    model_name,
    latency_ms,
    input_tokens,
    output_tokens,
    metadata,
    created_at,
    updated_at
FROM histories
WHERE success = TRUE
ORDER BY created_at DESC
LIMIT $1
`

func (r *PostgresRepository) Save(ctx context.Context, input SaveHistoryInput) (*DBHistory, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("storage repository is not initialized")
	}
	if strings.TrimSpace(input.UserPrompt) == "" {
		return nil, fmt.Errorf("user prompt is required")
	}

	sessionID := strings.TrimSpace(input.SessionID)
	if sessionID == "" {
		sessionID = uuid.NewString()
	}

	status := strings.TrimSpace(input.ValidationStatus)
	if status == "" {
		status = ValidationStatusUnknown
	}

	validationErrorsRaw, err := marshalStringSlice(input.ValidationErrors)
	if err != nil {
		return nil, fmt.Errorf("marshal validation errors: %w", err)
	}

	metadataRaw, err := marshalMetadata(input.Metadata)
	if err != nil {
		return nil, fmt.Errorf("marshal metadata: %w", err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err = ensureSession(ctx, tx, sessionID, input.UserPrompt); err != nil {
		return nil, err
	}

	row := tx.QueryRowContext(
		ctx,
		insertHistoryQuery,
		uuid.NewString(),
		sessionID,
		input.UserPrompt,
		input.ClarifiedPrompt,
		input.GeneratedCode,
		status,
		validationErrorsRaw,
		input.Success,
		input.ModelName,
		input.LatencyMS,
		input.InputTokens,
		input.OutputTokens,
		metadataRaw,
	)

	record, err := scanHistory(row)
	if err != nil {
		return nil, fmt.Errorf("insert history: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return record, nil
}

func (r *PostgresRepository) GetRecentSuccess(ctx context.Context, limit int) ([]DBHistory, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("storage repository is not initialized")
	}
	if limit <= 0 {
		limit = 10
	}

	rows, err := r.db.QueryContext(ctx, recentSuccessQuery, limit)
	if err != nil {
		return nil, fmt.Errorf("query recent success: %w", err)
	}
	defer rows.Close()

	result := make([]DBHistory, 0, limit)
	for rows.Next() {
		record, scanErr := scanHistory(rows)
		if scanErr != nil {
			return nil, fmt.Errorf("scan recent success row: %w", scanErr)
		}
		result = append(result, *record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate recent success rows: %w", err)
	}

	return result, nil
}

func ensureSession(ctx context.Context, tx *sql.Tx, sessionID string, userPrompt string) error {
	const query = `
INSERT INTO sessions (id, last_prompt)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE
SET last_prompt = EXCLUDED.last_prompt,
    updated_at = NOW()
`

	if _, err := tx.ExecContext(ctx, query, sessionID, truncateForSession(userPrompt)); err != nil {
		return fmt.Errorf("upsert session: %w", err)
	}

	return nil
}

func truncateForSession(prompt string) string {
	const maxLen = 512

	prompt = strings.TrimSpace(prompt)
	if len(prompt) <= maxLen {
		return prompt
	}

	return prompt[:maxLen]
}

type historyScanner interface {
	Scan(dest ...any) error
}

func scanHistory(scanner historyScanner) (*DBHistory, error) {
	var (
		record              DBHistory
		validationErrorsRaw []byte
		metadataRaw         []byte
	)

	if err := scanner.Scan(
		&record.ID,
		&record.SessionID,
		&record.UserPrompt,
		&record.ClarifiedPrompt,
		&record.GeneratedCode,
		&record.ValidationStatus,
		&validationErrorsRaw,
		&record.Success,
		&record.ModelName,
		&record.LatencyMS,
		&record.InputTokens,
		&record.OutputTokens,
		&metadataRaw,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(validationErrorsRaw, &record.ValidationErrors); err != nil {
		return nil, fmt.Errorf("unmarshal validation errors: %w", err)
	}
	if err := json.Unmarshal(metadataRaw, &record.Metadata); err != nil {
		return nil, fmt.Errorf("unmarshal metadata: %w", err)
	}
	if record.Metadata == nil {
		record.Metadata = map[string]any{}
	}

	return &record, nil
}

func marshalStringSlice(values []string) ([]byte, error) {
	if len(values) == 0 {
		values = []string{}
	}

	return json.Marshal(values)
}

func marshalMetadata(metadata map[string]any) ([]byte, error) {
	if metadata == nil {
		metadata = map[string]any{}
	}

	return json.Marshal(metadata)
}
