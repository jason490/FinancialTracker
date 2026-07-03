package storage

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// writeMigration creates a migrations dir with a single ALTER migration.
func writeMigration(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	sqlText := "ALTER TABLE tags ADD COLUMN is_hidden INTEGER NOT NULL DEFAULT 0;\n"
	if err := os.WriteFile(filepath.Join(dir, "0001_tags_is_hidden.sql"), []byte(sqlText), 0o644); err != nil {
		t.Fatalf("write migration: %v", err)
	}
	return dir
}

// hasColumn reports whether a table has a column of the given name.
func hasColumn(t *testing.T, db *sql.DB, table, column string) bool {
	t.Helper()
	rows, err := db.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		t.Fatalf("pragma: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dflt sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			t.Fatalf("scan: %v", err)
		}
		if name == column {
			return true
		}
	}
	return false
}

func TestRunMigrations_ExistingDatabaseGetsColumn(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	// Simulate a pre-migration tags table (no is_hidden column).
	if _, err := db.Exec(`CREATE TABLE tags (id INTEGER PRIMARY KEY, name TEXT)`); err != nil {
		t.Fatalf("seed tags: %v", err)
	}

	dir := writeMigration(t)

	if err := RunMigrations(db, dir, false); err != nil {
		t.Fatalf("first migrate: %v", err)
	}
	if !hasColumn(t, db, "tags", "is_hidden") {
		t.Fatal("expected is_hidden column after migration")
	}

	// Idempotent: running again must not error or double-apply.
	if err := RunMigrations(db, dir, false); err != nil {
		t.Fatalf("second migrate: %v", err)
	}
}

func TestRunMigrations_FreshInstallBaselines(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	dir := writeMigration(t)

	// Fresh install: base schema already includes the change, so the migration
	// must be recorded without executing (there is no tags table here at all).
	if err := RunMigrations(db, dir, true); err != nil {
		t.Fatalf("baseline migrate: %v", err)
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = '0001_tags_is_hidden'`).Scan(&count); err != nil {
		t.Fatalf("query schema_migrations: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected migration to be baselined, got count=%d", count)
	}
}
