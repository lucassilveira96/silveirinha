package database

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	oracleDriver = "godror"
	// oracleConnectionString = `user="%s" password="%s" connectString="%s:%s/%s?connect_timeout=3" poolSessionTimeout=900s heterogeneousPool=false standaloneConnection=false`
	oracleConnectionString = `user="%s" password="%s" connectString="%s:%s/%s" poolSessionTimeout=900s heterogeneousPool=false standaloneConnection=false`
)

func NewOracle(c *fiber.Ctx, cfg *SqlConfig) *Database {
	if len(cfg.Driver) == 0 {
		cfg.Driver = oracleDriver
	}

	return NewDatabase(c, cfg, oracleConnectionStringBuilder)
}

func oracleConnectionStringBuilder(cfg *SqlConfig) string {
	return fmt.Sprintf(oracleConnectionString, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}
