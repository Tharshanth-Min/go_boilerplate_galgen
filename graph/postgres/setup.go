package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"example/graph/logengine"
	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool

func InitDbPool() {
	host := os.Getenv("EZTRADE_HOST")
	port := os.Getenv("EZTRADE_PORT")
	user := os.Getenv("EZTRADE_USER")
	password := os.Getenv("EZTRADE_PASSWORD")
	dbname := os.Getenv("EZTRADE_DB_NAME")
	sslmode := os.Getenv("EZTRADE_SSLMODE")
	databaseUrl := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=" + sslmode

	maxOpenConnection, err := strconv.Atoi(os.Getenv("EZTRADE_MAX_CONN"))
	if err != nil {
		panic(err)
	}
	maxIdleTime, err := strconv.Atoi(os.Getenv("EZTRADE_MAX_IDLE_TIME"))
	if err != nil {
		panic(err)
	}
	maxConnectionLifetime, err := strconv.Atoi(os.Getenv("EZTRADE_MAX_LIFETIME"))
	if err != nil {
		panic(err)
	}

	healthcheckperiod, err := strconv.Atoi(os.Getenv("EZTRADE_HEALTHCHECK_PREIOD"))
	if err != nil {
		panic(err)
	}

	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		log.Print(err)
		log.Print("Error in config.")
		//return &pgxpool.Pool{}
	}
	config.MaxConns = int32(maxOpenConnection)
	config.MaxConnLifetime = time.Duration(maxConnectionLifetime) * time.Minute
	config.HealthCheckPeriod = time.Duration(healthcheckperiod) * time.Minute
	config.MaxConnIdleTime = time.Duration(maxIdleTime) * time.Minute

	pool, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		logengine.GetTelemetryClient().TrackException(err.Error())
		log.Print(err)
		log.Print("Could not connect to Postgres.")
	}
	fmt.Println("Postgres connected!")
}


func GetPool() *pgxpool.Pool {
	if pool == nil {
		InitDbPool()
	}

	connectedPoolSize := pool.AcquireAllIdle(context.Background())
	for connectedPoolSize == nil {
		// log.Println("Pg Connection Lost")
		logengine.GetTelemetryClient().TrackException(errors.New("Pg Connection Lost"))
		pool.Close()
		time.Sleep(2 * time.Second)
		logengine.GetTelemetryClient().TrackException(errors.New("Reconnecting..."))
		// log.Print("Reconnecting...")
		InitDbPool()
		connectedPoolSize = pool.AcquireAllIdle(context.Background())
	}
	return pool
}