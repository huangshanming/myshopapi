package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "homestead:secret@tcp(127.0.0.1:3306)/mymall?parseTime=true&loc=Local")
	if err != nil { panic(err) }
	defer db.Close()
	checks := []string{
		`SELECT COUNT(*) FROM shops WHERE province IS NULL OR city IS NULL OR district IS NULL OR address IS NULL OR logo IS NULL`,
		`SELECT COUNT(*) FROM banners`,
		`SELECT COUNT(*) FROM homepage_theme_slots`,
		`SELECT COUNT(*) FROM coupons WHERE status='active' OR 1=1`,
	}
	for _, q := range checks {
		var n int
		err := db.QueryRow(q).Scan(&n)
		fmt.Println(n, err, q)
	}
	// sample banner columns
	rows, err := db.Query(`SHOW COLUMNS FROM banners`)
	if err != nil { fmt.Println("banners show:", err); return }
	for rows.Next() {
		var f, t, n, k, d, e sql.NullString
		_ = rows.Scan(&f, &t, &n, &k, &d, &e)
		fmt.Printf("banner col %s %s null=%s\n", f.String, t.String, n.String)
	}
	rows.Close()
}
