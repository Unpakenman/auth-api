package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"strconv"
	"time"

	"github.com/XSAM/otelsql"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	_ "github.com/lib/pq"
)

type Client struct {
	*sqlx.DB
	parser Parser
}

type PostgreSQL struct {
	ConnString     string
	PathsToQueries []string

	LogLevel                        LogLevel
	MaxOpenConns                    int
	MaxIdleConns                    int
	MaxLifeTimeConns                time.Duration
	StatementTimeout                time.Duration
	IdleInTransactionSessionTimeout time.Duration
	LockTimeout                     time.Duration
}

// PGClient wrapped sqlx and pgx clients
type PGClient interface {
	NamedExec(
		ctx context.Context,
		queryName string,
		templateParams map[string]interface{},
		transaction Transaction,
		args interface{},
	) (sql.Result, error)
	Exec(
		ctx context.Context,
		queryName string,
		templateParams map[string]interface{},
		transaction Transaction,
		args ...any,
	) (sql.Result, error)
	NamedQuery(
		ctx context.Context,
		queryName string,
		templateParams map[string]interface{},
		transaction Transaction,
		args interface{},
	) (*sqlx.Rows, error)
	NamedQueryxContext(
		ctx context.Context,
		queryName string,
		templateParams map[string]interface{},
		transaction Transaction,
		args ...interface{},
	) (*sqlx.Rows, error)
	NamedGetContext(
		ctx context.Context,
		dest any,
		queryName string,
		templateParams map[string]interface{},
		transaction Transaction,
		args ...interface{},
	) error
	NamedSelectContext(
		ctx context.Context,
		dest any,
		queryName string,
		templateParams map[string]interface{},
		transaction Transaction,
		args ...interface{},
	) error

	BeginTransaction() (Transaction, error)
	CloseConnections() error
	GetQueryByName(
		name string,
		params map[string]interface{},
	) (string, error)
	WithTransaction(ctx context.Context, fn func(context.Context, Transaction) error) (err error)
}

// Transaction wrap sqlx.Tx
type Transaction interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Commit() error
	Rollback() error
}

// Struct for mapping namespace and files
type MappingNamespaceFiles struct {
	Namespace       string
	QueryFiles      fs.FS
	PathToDbQueries string
}

const (
	driverName     = "pgx"
	queryTraceName = "db.pg QUERY"
	queryTagName   = "query.name"

	SQLFileExt = "*.sql"
)

func NewClient(
	cfg PostgreSQL,
	queryFiles fs.FS,
	mappingNamespaceFiles []MappingNamespaceFiles,
) (PGClient, error) {
	if len(cfg.PathsToQueries) == 0 {
		return nil, fmt.Errorf("empty param PathsToQueries")
	}

	parser := NewParser()
	if len(mappingNamespaceFiles) > 0 {
		err := parser.AddFiles(mappingNamespaceFiles, SQLFileExt)
		if err != nil {
			return nil, err
		}
	} else if queryFiles == nil {
		for _, path := range cfg.PathsToQueries {
			err := parser.AddRoot(path, SQLFileExt)
			if err != nil {
				return nil, err
			}
		}
	} else {
		err := parser.AddFSRoot(cfg.PathsToQueries, queryFiles, SQLFileExt)
		if err != nil {
			return nil, err
		}
	}

	connConfig, err := pgx.ParseConfig(cfg.ConnString)
	if err != nil {
		return nil, err
	}

	if cfg.StatementTimeout != 0 {
		connConfig.Config.RuntimeParams["statement_timeout"] = strconv.Itoa(int(cfg.StatementTimeout.Milliseconds()))
	}

	if cfg.IdleInTransactionSessionTimeout != 0 {
		connConfig.Config.RuntimeParams["idle_in_transaction_session_timeout"] = strconv.Itoa(int(cfg.IdleInTransactionSessionTimeout.Milliseconds()))
	}

	if cfg.LockTimeout != 0 {
		connConfig.Config.RuntimeParams["lock_timeout"] = strconv.Itoa(int(cfg.LockTimeout.Milliseconds()))
	}

	connStr := stdlib.RegisterConnConfig(connConfig)

	db, err := otelsql.Open(driverName, connStr, otelsql.WithAttributes(
		semconv.DBSystemPostgreSQL,
	))
	if err != nil {
		return nil, err
	}

	// Register DB stats to meter

	dbx := sqlx.NewDb(db, driverName)

	if cfg.MaxLifeTimeConns > 0 {
		dbx.SetConnMaxLifetime(cfg.MaxLifeTimeConns)
	}
	if cfg.MaxOpenConns > 0 {
		dbx.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		dbx.SetMaxIdleConns(cfg.MaxIdleConns)
	}

	if err := dbx.Ping(); err != nil {
		return nil, err
	}

	return &Client{
		dbx,
		parser,
	}, nil
}

// NamedExec uses sqlx.NamedExecContext
func (c *Client) NamedExec(
	ctx context.Context,
	queryName string,
	templateParams map[string]interface{},
	transaction Transaction,
	args interface{},
) (sql.Result, error) {
	query, err := c.parser.Exec(queryName, templateParams)
	if err != nil {
		return nil, err
	}
	if query == "" {
		return nil, fmt.Errorf("not found query by name: %v", queryName)
	}

	if transaction != nil {
		return transaction.NamedExecContext(ctx, query, args)
	}
	return c.NamedExecContext(ctx, query, args)
}

// Exec uses sqlx.ExecContext
func (c *Client) Exec(
	ctx context.Context,
	queryName string,
	templateParams map[string]interface{},
	transaction Transaction,
	args ...any,
) (sql.Result, error) {
	query, err := c.parser.Exec(queryName, templateParams)
	if err != nil {
		return nil, err
	}
	if query == "" {
		return nil, fmt.Errorf("not found query by name: %v", queryName)
	}

	if transaction != nil {
		return transaction.ExecContext(ctx, query, args...)
	}

	return c.ExecContext(ctx, query, args...)
}

// NamedQuery uses sqlx.NamedQueryContext and NamedQuery for transactions
func (c *Client) NamedQuery(
	ctx context.Context,
	queryName string,
	templateParams map[string]interface{},
	transaction Transaction,
	args interface{},
) (*sqlx.Rows, error) {
	query, err := c.parser.Exec(queryName, templateParams)
	if err != nil {
		return nil, err
	}
	if query == "" {
		return nil, fmt.Errorf("not found query by name: %v", queryName)
	}

	if transaction != nil {
		return transaction.NamedQuery(query, args)
	}
	return c.NamedQueryContext(ctx, query, args)
}

// NamedQueryxContext uses sqlx.QueryxContext
func (c *Client) NamedQueryxContext(
	ctx context.Context,
	queryName string,
	templateParams map[string]interface{},
	transaction Transaction,
	args ...interface{},
) (*sqlx.Rows, error) {
	query, err := c.parser.Exec(queryName, templateParams)
	if err != nil {
		return nil, err
	}
	if query == "" {
		return nil, fmt.Errorf("not found query by name: %v", queryName)
	}

	if transaction != nil {
		return transaction.QueryxContext(ctx, query, args...)
	}

	return c.QueryxContext(ctx, query, args...)
}

// NamedGetContext uses sqlx.GetContext
func (c *Client) NamedGetContext(
	ctx context.Context,
	dest any,
	queryName string,
	templateParams map[string]interface{},
	transaction Transaction,
	args ...interface{},
) error {
	query, err := c.parser.Exec(queryName, templateParams)
	if err != nil {
		return err
	}
	if query == "" {
		return fmt.Errorf("not found query by name: %v", queryName)
	}

	if transaction != nil {
		return transaction.GetContext(ctx, dest, query, args...)
	}

	return c.GetContext(ctx, dest, query, args...)
}

// NamedSelectContext uses sqlx.SelectContext
func (c *Client) NamedSelectContext(
	ctx context.Context,
	dest any,
	queryName string,
	templateParams map[string]interface{},
	transaction Transaction,
	args ...interface{},
) error {
	query, err := c.parser.Exec(queryName, templateParams)
	if err != nil {
		return err
	}
	if query == "" {
		return fmt.Errorf("not found query by name: %v", queryName)
	}

	if transaction != nil {
		return transaction.SelectContext(ctx, dest, query, args...)
	}

	return c.SelectContext(ctx, dest, query, args...)
}

// BeginTransaction create transaction *sqlx.TX
func (c *Client) BeginTransaction() (Transaction, error) {
	return c.Beginx()
}

// CloseConnections run (*sql.DB).Close()
func (c *Client) CloseConnections() error {
	return c.Close()
}

// GetQueryByName return parsed query text from templates
func (c *Client) GetQueryByName(
	name string,
	params map[string]interface{},
) (string, error) {
	return c.parser.Exec(name, params)
}

func (c *Client) WithTransaction(ctx context.Context, fn func(context.Context, Transaction) error) (err error) {
	tx, err := c.BeginTransaction()
	if err != nil {
		return fmt.Errorf("begin transaction %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = errors.Join(err, fmt.Errorf("rollback failed %w", rbErr))
			}
		}
	}()

	err = fn(ctx, tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("commit failed: %w", err)
	}

	return
}
