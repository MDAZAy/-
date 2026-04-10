package storage

import (
	"context"
	"fmt"
	"strings"
)

func (r *PostgresRepository) GetStats(ctx context.Context, filter StatsFilter) (*Stats, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("storage repository is not initialized")
	}

	whereClause, args := buildHistoryFilter(filter)
	stats := &Stats{
		ValidationStatuses: make(map[string]int64),
	}

	summaryQuery := `
SELECT
    COUNT(*) AS total_runs,
    COUNT(*) FILTER (WHERE success) AS successful_runs,
    COUNT(*) FILTER (WHERE NOT success) AS failed_runs,
    COALESCE(AVG(latency_ms), 0) AS average_latency_ms
FROM histories
` + whereClause

	if err := r.db.QueryRowContext(ctx, summaryQuery, args...).Scan(
		&stats.TotalRuns,
		&stats.SuccessfulRuns,
		&stats.FailedRuns,
		&stats.AverageLatencyMS,
	); err != nil {
		return nil, fmt.Errorf("query stats summary: %w", err)
	}

	if stats.TotalRuns > 0 {
		stats.SuccessRate = float64(stats.SuccessfulRuns) / float64(stats.TotalRuns)
	}

	groupQuery := `
SELECT validation_status, COUNT(*) AS total
FROM histories
` + whereClause + `
GROUP BY validation_status
ORDER BY validation_status
`

	rows, err := r.db.QueryContext(ctx, groupQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("query validation status groups: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			status string
			total  int64
		)

		if err := rows.Scan(&status, &total); err != nil {
			return nil, fmt.Errorf("scan validation status group: %w", err)
		}

		stats.ValidationStatuses[status] = total
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate validation status groups: %w", err)
	}

	return stats, nil
}

func buildHistoryFilter(filter StatsFilter) (string, []any) {
	conditions := make([]string, 0, 2)
	args := make([]any, 0, 2)

	if filter.From != nil {
		args = append(args, *filter.From)
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", len(args)))
	}

	if filter.To != nil {
		args = append(args, *filter.To)
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", len(args)))
	}

	if len(conditions) == 0 {
		return "", args
	}

	return "WHERE " + strings.Join(conditions, " AND "), args
}
