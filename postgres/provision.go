package postgres

import (
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// _schemaVersion is the curretn version the schema is at. If any updates are
// done to the schema this number will be increased by one to automatically
// update the schema in the database accordingly.
const _schemaVersion = 1

const (
	_pgCodeRelationDoesNotExist = "42P01"
)

const (
	_advisoryLockProvision = iota + 1
)

type provisionStatement struct {
	version int
	sql     string
	params  []any
}

// acquireAdvisoryLock acquires an advisory lock at transaction level for a
// given reason.
func acquireAdvisoryLock(ctx context.Context, tx pgx.Tx, reason int) error {
	_, err := tx.Exec(ctx, "SELECT pg_advisory_lock($1)", reason)
	return err
}

// isSchemaAtCurrentVersion reports if the schema of the database is at the
// current exepected version defined by `_schemaVersion`.
func isSchemaAtCurrentVersion(ctx context.Context, conn *pgxpool.Conn) (bool, error) {
	var pgErr *pgconn.PgError
	var schemaVerion string
	rows, err := conn.Query(ctx, `SELECT value FROM omp_metadata WHERE key = 'schema_version'`)
	if err != nil && errors.As(err, &pgErr) {
		if pgErr.Code == _pgCodeRelationDoesNotExist {
			return false, nil
		}
	}
	defer rows.Close()
	if !rows.Next() {
		return true, errors.New(
			"now row found for key 'schema_version'. This error should never occur. If it does then the key might have changed or the name of the table",
		)
	}
	err = rows.Scan(&schemaVerion)
	if err != nil {
		return true, err
	}
	schemaVersionInt, err := strconv.Atoi(schemaVerion)
	return schemaVersionInt == _schemaVersion, err
}

// provision will provision the database with all required entities (tables,
// views, indeces etc.)
func provision(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()
	if isAlreadyProvisioned, err := isSchemaAtCurrentVersion(ctx, conn); err != nil || isAlreadyProvisioned {
		return err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	err = acquireAdvisoryLock(ctx, tx, _advisoryLockProvision)
	if err != nil {
		return err
	}
	pseudoTx, err := tx.Begin(ctx)
	if err != nil {
		return err
	}
	statements := []*provisionStatement{
		provisionStatementV1(),
	}
	for _, stmt := range statements {
		// skip older version and only apply the newest ones
		if stmt.version != _schemaVersion {
			continue
		}
		_, err = pseudoTx.Exec(ctx, stmt.sql, stmt.params...)
		if err != nil {
			return err
		}
	}
	if err := pseudoTx.Commit(ctx); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO omp_metadata VALUES('schema_version', $1);`, strconv.Itoa(_schemaVersion))
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
