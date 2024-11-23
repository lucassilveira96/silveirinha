package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"teste/internal/infra/logger"
	"teste/internal/infra/logger/attributes"
	"time"

	"github.com/gofiber/fiber/v2"
)

type (
	ConnectionStringBuilder func(cfg *SqlConfig) string

	Database struct {
		db                      *sql.DB
		config                  *SqlConfig
		connectionStringBuilder ConnectionStringBuilder
		retries                 []time.Duration
		locker                  sync.Mutex
	}
)

func NewDatabase(c *fiber.Ctx, cfg *SqlConfig, connectionStringBuilder ConnectionStringBuilder) *Database {
	database := &Database{
		config:                  cfg,
		connectionStringBuilder: connectionStringBuilder,
		retries: []time.Duration{
			250 * time.Millisecond,
			500 * time.Millisecond,
			1000 * time.Millisecond,
			2500 * time.Millisecond,
			5000 * time.Millisecond,
		},
	}

	if !cfg.LazyConnection {
		database.initializeAndGetDB(c)
	}

	return database
}

func (d *Database) Connection(c *fiber.Ctx) *sql.DB {
	return d.initializeAndGetDB(c)
}

func (d *Database) Close(c *fiber.Ctx) {
	d.locker.Lock()
	defer d.locker.Unlock()

	if d.db == nil {
		return
	}

	if err := d.db.Close(); err != nil {
		logger.Error(
			c,
			fmt.Sprintf("Failed to close database [%s] with connection [%s]", d.config.Database, d.config.ConnectionName),
			d.configToAttribute().WithError(err),
		)
	}

	d.db = nil
}

func (d *Database) initializeAndGetDB(c *fiber.Ctx) *sql.DB {
	db := d.db
	if db != nil {
		return db
	}

	d.locker.Lock()
	defer d.locker.Unlock()

	// double-checked locking
	if db = d.db; db != nil {
		return db
	}

	// start := time.Now()
	// logger.Info(context.Background(), fmt.Sprintf("Initializing database [%s] with connection [%s]", d.config.Database, d.config.ConnectionName), d.configToAttribute())
	connected := false
	var err error

	for retry, duration := range d.retries {
		db, err = sql.Open(d.config.Driver, d.connectionStringBuilder(d.config))

		if err != nil {
			logger.Fatal(
				c,
				fmt.Sprintf("Failed to initialize the database [%s] with connection [%s]", d.config.Database, d.config.ConnectionName),
				d.configToAttribute().WithError(err), // Os atributos do erro
			)
		}

		if err = d.checkConnection(db); err != nil {
			logger.Warn(
				c,
				fmt.Sprintf("Connection retry [%d]: Database [%s] with connection [%s]", retry+1, d.config.Database, d.config.ConnectionName),
				d.configToAttribute().WithError(err),
			)
			time.Sleep(duration)
		} else {
			connected = true
			err = nil
			break
		}
	}

	if !connected {
		if err := d.checkConnection(db); err != nil {
			logger.Fatal(
				c,
				fmt.Sprintf("Failed to connect to the database [%s] with connection [%s]", d.config.Database, d.config.ConnectionName),
				d.configToAttribute().WithError(err),
			)
		}
	}

	db.SetMaxIdleConns(d.config.MinConnections)
	db.SetMaxOpenConns(d.config.MaxConnections)
	db.SetConnMaxLifetime(d.config.ConnectionMaxLifetime)
	db.SetConnMaxIdleTime(d.config.ConnectionMaxIdleTime)

	d.db = db

	// elapsed := time.Since(start)
	// logger.Info(context.Background(), fmt.Sprintf("Database [%s] initialized with connection [%s] in [%v]", d.config.Database, d.config.ConnectionName, elapsed), d.configToAttribute())

	return db
}

func (d *Database) checkConnection(db *sql.DB) error {
	timeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return db.PingContext(timeout)
}

func (d *Database) configToAttribute() attributes.Attributes {
	config := d.config
	return attributes.Attributes{
		"database.connection_name":          config.ConnectionName,
		"database.driver":                   config.Driver,
		"database.host":                     config.Host,
		"database.port":                     config.Port,
		"database.database":                 config.Database,
		"database.username":                 config.Username,
		"database.password":                 "[Masked]",
		"database.min_connections":          config.MinConnections,
		"database.max_connections":          config.MaxConnections,
		"database.connection_max_lifetime":  config.ConnectionMaxLifetime,
		"database.connection_max_idle_time": config.ConnectionMaxIdleTime,
		"database.lazy_connection":          config.LazyConnection,
	}
}
