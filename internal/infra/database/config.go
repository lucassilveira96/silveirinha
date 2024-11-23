package database

import "time"

type SqlConfig struct {
	ConnectionName        string
	Driver                string
	Host                  string
	Port                  string
	Database              string
	Username              string
	Password              string
	MinConnections        int
	MaxConnections        int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
	LazyConnection        bool
}
