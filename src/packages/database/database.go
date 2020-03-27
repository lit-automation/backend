package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq" // Postgres driver
	log "github.com/sirupsen/logrus"
)

// database constants
const (
	dbConnection    = "dbname=%s user=%s password=%s port=%d host=%s sslmode=disable"
	dbConnectionLog = "Start connection to dbname=%s user=%s port=%d host=%s"
	createExtension = "CREATE EXTENSION IF NOT EXISTS \"%s\";"
	maxTimeout      = 120

	// Different test modes for db
	TestModeOff         TestMode = 1
	TestModeDropTables  TestMode = 2
	TestModeNoDropTable TestMode = 3
)

// TestMode Enum
type TestMode int64

// TestModeDesc description of different test modes
var TestModeDesc = map[TestMode]string{
	TestModeOff:         "Test mode off",
	TestModeDropTables:  "Test mode with dropping tables",
	TestModeNoDropTable: "Test mode without dropping tables",
}

// CreateDefaultExtensionConfig creates a default db extension config
func CreateDefaultExtensionConfig() []string {
	return []string{"uuid-ossp", "pg_trgm"}
}

// DBConfig config for database connection
type DBConfig struct {
	Host           string
	Port           int
	User           string
	Pass           string
	DBName         string
	RetryOnFailure bool
	LogMode        bool
	TestMode       TestMode
	TimeOut        int
	AutoMigrate    bool
	Extensions     []string
}

// CreateDefaultDBConfig creates a default postgres configuration
func CreateDefaultDBConfig() DBConfig {
	return DBConfig{
		RetryOnFailure: false,
		TestMode:       TestModeOff,
		Port:           5432,
		AutoMigrate:    true,
		Extensions:     CreateDefaultExtensionConfig(),
	}
}

// CreateTestDBConfig creates desired test DB configuration
func CreateTestDBConfig() DBConfig {
	return DBConfig{
		RetryOnFailure: true,
		TestMode:       TestModeNoDropTable,
		Port:           5432,
		TimeOut:        60,
		AutoMigrate:    true,
		Extensions:     CreateDefaultExtensionConfig(),
	}
}

// MustConnectDefault create must connection for default cluster
func MustConnectDefault(dbName string, models []interface{}) *gorm.DB {
	db, err := ConnectDefault(dbName, models)
	if err != nil {
		panic(err)
	}
	return db
}

// ConnectDefault create connection for default cluster
func ConnectDefault(dbName string, models []interface{}) (*gorm.DB, error) {
	config := CreateDefaultDBConfig()
	config.DBName = dbName
	client := NewClientFromConfig(config)
	return client.InitDB(models)
}

// MustConnectTest must connection for running tests
func MustConnectTest(dbName string, models []interface{}) *gorm.DB {
	db, err := ConnectTest(dbName, models)
	if err != nil {
		panic(err)
	}
	return db
}

// ConnectTest connection for running tests
func ConnectTest(dbName string, models []interface{}) (*gorm.DB, error) {
	config := CreateTestDBConfig()
	config.DBName = dbName
	client := NewClientFromConfig(config)
	return client.InitDB(models)
}

// MustConnectCustom must connect custom database connection with specified config
func MustConnectCustom(config DBConfig, models []interface{}) *gorm.DB {
	db, err := ConnectCustom(config, models)
	if err != nil {
		panic(err)
	}
	return db
}

// ConnectCustom create custom database connection with specfied config
func ConnectCustom(config DBConfig, models []interface{}) (*gorm.DB, error) {
	client := NewClientFromConfig(config)
	return client.InitDB(models)
}

// NewClientFromConfig create postgres client from config
func NewClientFromConfig(config DBConfig) PostgresClient {
	client := PostgresClient{
		DBConfig:          config,
		connectionRetries: 0,
	}
	// Config prio: config < env < flags
	client.ensureDefaults()
	client.evalEnvironment()
	return client
}

// PostgresClient client for connecting to postgres
type PostgresClient struct {
	DBConfig
	DB                *gorm.DB
	connectionRetries int
}

// InitDB connect, migrate & create extensions
func (c *PostgresClient) InitDB(models []interface{}) (*gorm.DB, error) {
	c.LogConfig()
	err := c.Connect()
	if err != nil {
		return nil, err
	}
	c.DB.LogMode(c.LogMode)
	c.CreateDBExtensions()
	c.Migrate(models)
	return c.DB, nil
}

// Connect try to create a connection
func (c *PostgresClient) Connect() error {
	var err error
	url := fmt.Sprintf(dbConnection, c.DBName, c.User, c.Pass, c.Port, c.Host)
	logURL := fmt.Sprintf(dbConnectionLog, c.DBName, c.User, c.Port, c.Host)
	log.Infof(logURL)
	c.DB, err = gorm.Open("postgres", url)
	if c.RetryOnFailure {
		for err != nil && c.connectionRetries < c.TimeOut {
			c.connectionRetries++
			log.Warningf("unable to connect to %s, error: %s", logURL, err.Error())
			time.Sleep(500 * time.Millisecond)
			c.DB, err = gorm.Open("postgres", url)
		}
	}
	return err
}

// Migrate migrate database according to client config
func (c *PostgresClient) Migrate(models []interface{}) {
	if c.TestMode == TestModeDropTables {
		//If running tests always drop db up front.
		//Because in case tests fail, they might mess with data for other tests
		c.DB.DropTableIfExists(models...)
	}
	if c.AutoMigrate {
		c.DB.AutoMigrate(models...)
	}
}

// CreateDBExtensions create database extensions
func (c *PostgresClient) CreateDBExtensions() {
	for _, extension := range c.DBConfig.Extensions {
		c.DB.Exec(fmt.Sprintf(createExtension, extension))
	}
}

// ensureDefaults ensure clients defaults are set
func (c *PostgresClient) ensureDefaults() {
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.User == "" {
		c.User = "postgres"
	}
	if c.TestMode != TestModeOff {
		c.DBName += "_test"
	}
	if c.Port == 0 {
		c.Port = 5432
	}
	if c.RetryOnFailure {
		if c.TimeOut <= 0 {
			c.TimeOut = maxTimeout
		}
		if c.TimeOut > maxTimeout {
			c.TimeOut = maxTimeout
		}
	}
}

// LogConfig log current client config
func (c *PostgresClient) LogConfig() {
	log.Infof("Configuration:")
	log.Infof("DB HOST: %s", c.Host)
	log.Infof("DB Name: %s", c.DBName)
	log.Infof("DB User: %s", c.User)
	log.Infof("DB Port: %d", c.Port)
	log.Infof("TestMode: %s", TestModeDesc[c.TestMode])
	log.Infof("Retry on failure: %v", c.RetryOnFailure)
	if c.RetryOnFailure {
		log.Infof("Timeout: %d", c.TimeOut)
	}
}

func (c *PostgresClient) evalEnvironment() {
	host := os.Getenv("PG_HOST")
	user := os.Getenv("PG_USERNAME")
	pass := os.Getenv("PG_PASSWORD")
	port := os.Getenv("PG_PORT")
	dbName := os.Getenv("PG_DBNAME")

	if host != "" {
		c.Host = host
	}
	if user != "" {
		c.User = user
	}
	if pass != "" {
		c.Pass = pass
	}
	if dbName != "" {
		c.DBName = dbName
	}
	if port != "" {
		pgPort, err := strconv.Atoi(port)
		if err != nil {
			c.Port = 5432
		}
		c.Port = pgPort
	}
}

// Transact handles a postgres transaction
func Transact(db *gorm.DB, tf func(tx *gorm.DB) error) (err error) {
	if commonDB, ok := db.CommonDB().(sqlTx); ok && commonDB != nil {
		// If the db is already in a transaction, just execute tf
		// and let the outer transaction handle Rollback and Commit.
		return tf(db)
	}

	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("could not start transaction. %s", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit().Error
		}
	}()
	return tf(tx)
}

// sqlTx is a helper interface to check if a gorm.DB.CommonDB() is already in a transaction.
type sqlTx interface {
	Commit() error
	Rollback() error
}
