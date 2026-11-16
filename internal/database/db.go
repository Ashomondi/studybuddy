package database

import (
	"database/sql"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"studybuddy/internal/migrations"

	_ "github.com/mattn/go-sqlite3"
)

var migrationFS = migrations.FS

func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path+"?_busy_timeout=5000&_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`); err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	entries, err := fs.ReadDir(migrationFS, ".")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	type mig struct {
		version int
		name    string
	}
	var migs []mig
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		v := versionFromName(e.Name())
		if v < 0 {
			continue
		}
		migs = append(migs, mig{version: v, name: e.Name()})
	}
	sort.Slice(migs, func(i, j int) bool { return migs[i].version < migs[j].version })

	for _, m := range migs {
		var exists int
		if err := db.QueryRow("SELECT 1 FROM schema_migrations WHERE version = ?", m.version).Scan(&exists); err == nil {
			continue
		} else if err != sql.ErrNoRows {
			return fmt.Errorf("check migration %d: %w", m.version, err)
		}
		sqlBytes, err := fs.ReadFile(migrationFS, m.name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", m.name, err)
		}
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", m.name, err)
		}
		if _, err := tx.Exec("INSERT INTO schema_migrations(version) VALUES (?)", m.version); err != nil {
			tx.Rollback()
			return fmt.Errorf("record migration %s: %w", m.name, err)
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func versionFromName(name string) int {
	// expects prefix like "001_"
	base := strings.TrimSuffix(name, ".sql")
	parts := strings.SplitN(base, "_", 2)
	if len(parts) < 2 {
		return -1
	}
	var v int
	if _, err := fmt.Sscanf(parts[0], "%d", &v); err != nil {
		return -1
	}
	return v
}
