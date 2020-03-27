package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

const (
	CreateDbSql = `CREATE DATABASE %s OWNER postgres`
	RevokeConnectionsSql = `
DO
$$BEGIN
	IF EXISTS (
	 SELECT datname FROM pg_catalog.pg_database WHERE lower(datname) = lower('%s')
	) THEN
		REVOKE CONNECT ON DATABASE %s FROM public;
	END IF;
END$$;
`
	KillConnectionsSql = `
SELECT
  pg_terminate_backend (pg_stat_activity.pid)
FROM
  pg_stat_activity
WHERE
  pg_stat_activity.datname = '%s';
`
	DropDbSql = `DROP DATABASE IF EXISTS %s;`
)

type Repository interface {
	CreateDatabase(cs string, dbName string) error
	Execute(cs string, dbName string, script string) error
	DropDatabase(cs string, dbName string) error
}

func NewRepository(provider string) (Repository, error) {
	switch provider {
	case "postgres":
		return &PostgresRepository{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown provider: %s", provider))
	}
}

type PostgresRepository struct{}

func (repo *PostgresRepository) CreateDatabase(cs string, dbName string) error {
	_dbName := strings.ToLower(dbName)
	db, err := sql.Open("postgres", cs)
	if err != nil {
		return fmt.Errorf("unable to open db connection: %w", err)
	}
	defer db.Close()

	sql := fmt.Sprintf(CreateDbSql, _dbName)
	_, err = db.Exec(sql)
	if err != nil {
		return fmt.Errorf("unable to create database: %w", err)
	}

	return nil
}

func (repo *PostgresRepository) Execute(cs string, dbName string, script string) error {
	connString := fmt.Sprintf("%s dbname=%s", cs, strings.ToLower(dbName))
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("unable to open db connection: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(script)
	if err != nil {
		return fmt.Errorf("error running script: %w", err)
	}

	return nil
}

func (repo *PostgresRepository) DropDatabase(cs string, dbName string) error {
	_dbName := strings.ToLower(dbName)
	db, err := sql.Open("postgres", cs)
	if err != nil {
		return fmt.Errorf("unable to open db connection: %w", err)
	}
	defer db.Close()

	// revoke connections
	sql := fmt.Sprintf(RevokeConnectionsSql, _dbName, _dbName)
	_, err = db.Exec(sql)
	if err != nil {
		return fmt.Errorf("error revoking connections: %w", err)
	}

	// kill connections
	sql = fmt.Sprintf(KillConnectionsSql, _dbName)
	_, err = db.Exec(sql)
	if err != nil {
		return fmt.Errorf("error killing connections: %w", err)
	}

	// drop db
	sql = fmt.Sprintf(DropDbSql, _dbName)
	_, err = db.Exec(sql)
	if err != nil {
		return fmt.Errorf("error dropping db: %w", err)
	}

	return nil
}
