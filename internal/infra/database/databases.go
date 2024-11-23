package database

import (
	"fmt"
	"runtime"
	"sync"
	"teste/internal/infra/variables"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const BIND_PARAM string = "_bind_"

type Databases struct {
	Read  *gorm.DB
	Write *gorm.DB
	Redis *Redis
}

func NewDatabases(c *fiber.Ctx) *Databases {
	dbs := &Databases{}
	var waitGroup sync.WaitGroup
	defer waitGroup.Wait()

	// Inicializar as conexões com os bancos
	dbs.buildReadDatabase(&waitGroup)
	dbs.buildWriteDatabase(&waitGroup)
	dbs.buildRedisDatabase(c, &waitGroup)

	return dbs
}

func (d *Databases) Close() {
	// Fechar as conexões do GORM
	if d.Read != nil {
		d.Read.Session(&gorm.Session{DryRun: true})
	}
	if d.Write != nil {
		d.Write.Session(&gorm.Session{DryRun: true})
	}
}

func (d *Databases) buildReadDatabase(waitGroup *sync.WaitGroup) {
	lazyConnection := variables.DBLazyConnection()
	cfg := &SqlConfig{
		ConnectionName:        variables.ServiceName() + "-read",
		Host:                  variables.DBHost(),
		Port:                  variables.DBPort(),
		Database:              variables.DBName(),
		Username:              variables.DBUsername(),
		Password:              variables.DBPassword(),
		MinConnections:        variables.DBMinConnections(),
		MaxConnections:        variables.DBMaxConnections(),
		ConnectionMaxLifetime: variables.DBConnectionMaxLifeTime(),
		ConnectionMaxIdleTime: variables.DBConnectionMaxIdleTime(),
		LazyConnection:        lazyConnection,
	}

	// Usando GORM para criar a conexão
	if lazyConnection {
		d.Read = d.connectDatabase(cfg)
	} else {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			d.Read = d.connectDatabase(cfg)
		}()
	}
}

func (d *Databases) buildWriteDatabase(waitGroup *sync.WaitGroup) {
	lazyConnection := variables.DBLazyConnection()
	cfg := &SqlConfig{
		ConnectionName:        variables.ServiceName() + "-write",
		Host:                  variables.DBHost(),
		Port:                  variables.DBPort(),
		Database:              variables.DBName(),
		Username:              variables.DBUsername(),
		Password:              variables.DBPassword(),
		MinConnections:        variables.DBMinConnections(),
		MaxConnections:        variables.DBMaxConnections(),
		ConnectionMaxLifetime: variables.DBConnectionMaxLifeTime(),
		ConnectionMaxIdleTime: variables.DBConnectionMaxIdleTime(),
		LazyConnection:        lazyConnection,
	}

	// Usando GORM para criar a conexão
	if lazyConnection {
		d.Write = d.connectDatabase(cfg)
	} else {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			d.Write = d.connectDatabase(cfg)
		}()
	}
}

func (d *Databases) buildRedisDatabase(c *fiber.Ctx, waitGroup *sync.WaitGroup) {
	lazyConnection := variables.RedisLazyConnection()
	opt := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", variables.RedisHost(), variables.RedisPort()),
		Password:     variables.RedisPassword(),
		DB:           variables.RedisDB(),
		PoolSize:     10 * runtime.NumCPU(),
		MinIdleConns: 10,
	}

	// Configuração para o Redis
	if lazyConnection {
		d.Redis = NewRedis(c, opt, lazyConnection)
	} else {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			d.Redis = NewRedis(c, opt, lazyConnection)
		}()
	}
}

func (d *Databases) connectDatabase(cfg *SqlConfig) *gorm.DB {
	var dsn string
	var db *gorm.DB
	var err error

	// Verifica se o banco é Oracle ou Postgres
	if variables.IsOracle() {
		// Configuração para Oracle
		dsn = fmt.Sprintf("oracle://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		// Estabelecer conexão Oracle
		// (O código para Oracle pode ser ajustado aqui)
	} else {
		// Configuração para Postgres
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", cfg.Host, cfg.Port, cfg.Username, cfg.Database, cfg.Password)
		// Estabelecendo conexão GORM com o banco de dados
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database: " + err.Error())
		}

	}

	// Retorna a conexão com o banco
	return db
}
