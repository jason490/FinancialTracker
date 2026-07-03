package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// TableExists reports whether a table with the given name exists. It is used at
// startup to distinguish a brand-new database (schema.sql creates everything)
// from an existing one that predates newer migrations.
func TableExists(db *sql.DB, name string) (bool, error) {
	var found string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?`, name).Scan(&found)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// RunMigrations applies ordered `*.sql` files from dir that have not yet been
// recorded in the schema_migrations table, executing each in its own
// transaction and recording its version on success.
//
// When freshInstall is true the base schema (schema.sql) already contains every
// change, so pending migrations are baselined (recorded as applied) without
// executing them. This keeps a single migration file usable both as the
// authoritative delta for existing databases and as a no-op for new ones.
//
// Adding a new migration is as simple as dropping another numbered file into the
// migrations directory (e.g. `0002_add_widget.sql`); files are applied in
// lexical order, so zero-pad the numeric prefix.
func RunMigrations(db *sql.DB, dir string, freshInstall bool) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at INTEGER NOT NULL DEFAULT (strftime('%s', 'now'))
	)`); err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read migrations dir %s: %w", dir, err)
	}

	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		version := strings.TrimSuffix(name, ".sql")

		applied, err := migrationApplied(db, version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		if freshInstall {
			if _, err := db.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, version); err != nil {
				return fmt.Errorf("baseline migration %s: %w", version, err)
			}
			continue
		}

		sqlBytes, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		if err := applyMigration(db, version, string(sqlBytes)); err != nil {
			return err
		}
	}

	return nil
}

// migrationApplied reports whether a migration version is already recorded.
func migrationApplied(db *sql.DB, version string) (bool, error) {
	var one int
	err := db.QueryRow(`SELECT 1 FROM schema_migrations WHERE version = ?`, version).Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check migration %s: %w", version, err)
	}
	return true, nil
}

// applyMigration runs a migration's SQL and records its version atomically.
func applyMigration(db *sql.DB, version, script string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin migration %s: %w", version, err)
	}
	if _, err := tx.Exec(script); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("apply migration %s: %w", version, err)
	}
	if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES (?)`, version); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("record migration %s: %w", version, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", version, err)
	}
	return nil
}
