package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	parts := make([]string, n)
	for i := range parts {
		parts[i] = "?"
	}
	return strings.Join(parts, ",")
}

func inArgs(ids []uint64) []any {
	out := make([]any, len(ids))
	for i, id := range ids {
		out[i] = id
	}
	return out
}

func countCtx(ctx context.Context, conn sqlx.SqlConn, query string, args ...any) (int64, error) {
	var n int64
	err := conn.QueryRowCtx(ctx, &n, query, args...)
	return n, err
}

func execAffected(ctx context.Context, conn sqlx.SqlConn, query string, args ...any) (int64, error) {
	res, err := conn.ExecCtx(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func lastInsertID(ctx context.Context, session interface {
	ExecCtx(ctx context.Context, query string, args ...any) (sql.Result, error)
}, query string, args ...any) (uint64, error) {
	res, err := session.ExecCtx(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	if id <= 0 {
		return 0, fmt.Errorf("empty last insert id")
	}
	return uint64(id), nil
}

func buildUpdate(table string, updates map[string]interface{}, where string, whereArgs ...any) (string, []any, error) {
	if len(updates) == 0 {
		return "", nil, fmt.Errorf("empty updates")
	}
	sets := make([]string, 0, len(updates))
	args := make([]any, 0, len(updates)+len(whereArgs))
	for k, v := range updates {
		sets = append(sets, k+"=?")
		args = append(args, v)
	}
	args = append(args, whereArgs...)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, strings.Join(sets, ", "), where)
	return query, args, nil
}
