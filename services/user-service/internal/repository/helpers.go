package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func countQuery(ctx context.Context, session sqlx.Session, query string, args ...any) (int64, error) {
	var n int64
	err := session.QueryRowCtx(ctx, &n, query, args...)
	return n, err
}

func execRows(ctx context.Context, session sqlx.Session, query string, args ...any) (int64, error) {
	res, err := session.ExecCtx(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
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

func lastInsertID(res sql.Result) (uint64, error) {
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}
