package variables

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type variable struct {
	key          string
	defaultValue string
}

var (
	serviceName             = &variable{key: "SERVICE_NAME", defaultValue: "go-project-template"}
	serviceVersion          = &variable{key: "SERVICE_VERSION", defaultValue: "0.0.1"}
	environment             = &variable{key: "ENVIRONMENT", defaultValue: "local"}
	isLambda                = &variable{key: "LAMBDA", defaultValue: "false"}
	logLevel                = &variable{key: "LOG_LEVEL", defaultValue: "debug"}
	serverHost              = &variable{key: "SERVER_HOST", defaultValue: "0.0.0.0"}
	serverPort              = &variable{key: "SERVER_PORT", defaultValue: "5000"}
	serverTimeout           = &variable{key: "SERVER_TIMEOUT", defaultValue: "30"}
	databaseType            = &variable{key: "DATABASE_TYPE", defaultValue: "postgres"}
	dbHost                  = &variable{key: "DB_HOST", defaultValue: "localhost"}
	dbPort                  = &variable{key: "DB_PORT", defaultValue: "5432"}
	dbName                  = &variable{key: "DB_NAME", defaultValue: "go-project-template"}
	dbUsername              = &variable{key: "DB_USERNAME", defaultValue: "postgres"}
	dbPassword              = &variable{key: "DB_PASSWORD", defaultValue: "postgres123"}
	dbLazyConnection        = &variable{key: "DB_LAZY_CONNECTION", defaultValue: "true"}
	dbMinConnections        = &variable{key: "DB_MIN_CONNECTIONS", defaultValue: "2"}
	dbMaxConnections        = &variable{key: "DB_MAX_CONNECTIONS", defaultValue: "10"}
	dbConnectionMaxLifeTime = &variable{key: "DB_CONNECTION_MAX_LIFE_TIME", defaultValue: "900"}
	dbConnectionMaxIdleTime = &variable{key: "DB_CONNECTION_MAX_IDLE_TIME", defaultValue: "60"}
	redisHost               = &variable{key: "REDIS_HOST", defaultValue: "localhost"}
	redisPort               = &variable{key: "REDIS_PORT", defaultValue: "6379"}
	redisPassword           = &variable{key: "REDIS_PASSWORD", defaultValue: ""}
	redisDB                 = &variable{key: "REDIS_DB", defaultValue: "1"}
	redisLazyConnection     = &variable{key: "REDIS_LAZY_CONNECTION", defaultValue: "true"}
	prefixRoute             = &variable{key: "PREFIX_ROUTE", defaultValue: os.Getenv("PREFIX_ROUTE")}
	jwtTokenKey             = &variable{key: "JWT_TOKEN_KEY", defaultValue: os.Getenv("JWT_TOKEN_KEY")}
	dirLog                  = &variable{key: "DIR_LOG", defaultValue: "logs"}
)

func ServiceName() string {
	return get(serviceName)
}

func ServiceVersion() string {
	return get(serviceVersion)
}

func Environment() string {
	return get(environment)
}

func IsLambda() bool {
	return getBool(isLambda)
}

func LogLevel() string {
	return get(logLevel)
}

func ServerHost() string {
	return get(serverHost)
}

func ServerPort() int {
	return getInt(serverPort)
}

func ServerTimeout() int {
	return getInt(serverTimeout)
}

func DatabaseType() string {
	return get(databaseType)
}

func IsOracle() bool {
	return strings.ToUpper(DatabaseType()) == "ORACLE"
}

func DBHost() string {
	return get(dbHost)
}

func DBPort() string {
	return get(dbPort)
}

func DBName() string {
	return get(dbName)
}

func DBUsername() string {
	return get(dbUsername)
}

func DBPassword() string {
	return get(dbPassword)
}

func DBLazyConnection() bool {
	return getBool(dbLazyConnection)
}

func DBMinConnections() int {
	return getInt(dbMinConnections)
}

func DBMaxConnections() int {
	return getInt(dbMaxConnections)
}

func DBConnectionMaxLifeTime() time.Duration {
	return time.Second * time.Duration(getInt(dbConnectionMaxLifeTime))
}

func DBConnectionMaxIdleTime() time.Duration {
	return time.Second * time.Duration(getInt(dbConnectionMaxIdleTime))
}

func RedisHost() string {
	return get(redisHost)
}

func RedisPort() int {
	return getInt(redisPort)
}

func RedisPassword() string {
	return get(redisPassword)
}

func RedisDB() int {
	return getInt(redisDB)
}

func RedisLazyConnection() bool {
	return getBool(redisLazyConnection)
}

func PrefixRoute() string {
	return get(prefixRoute)
}

func JWTTokenKey() []byte {
	return []byte(get(jwtTokenKey))
}

func get(env *variable) string {
	value := os.Getenv(env.key)

	if len(value) == 0 {
		return env.defaultValue
	}

	return value
}

func getInt(env *variable) int {
	value := get(env)
	intValue, err := strconv.Atoi(value)

	if err != nil {
		logFatal(env, "int", value, err)
	}

	return intValue
}

func getBool(env *variable) bool {
	value := get(env)
	boolValue, err := strconv.ParseBool(value)

	if err != nil {
		logFatal(env, "bool", value, err)
	}

	return boolValue
}

func logFatal(env *variable, varType string, returnedValue string, err error) {

}

func DirLog() string {
	return get(dirLog)
}
