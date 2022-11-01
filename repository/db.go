package repository

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var db *sqlx.DB

func initDB() {
	db = connectDB()
	db.SetMaxIdleConns(4)
	db.SetMaxOpenConns(4)
	db.SetConnMaxLifetime(time.Second * 15)
}

func DBGet() *sqlx.DB {
	if db == nil {
		initDB()
	}
	return db
}

func connectDB() *sqlx.DB {
	username, password, databasename, databaseHost := getDBConfig()
	//Define DB connection string
	dbURI := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", username, password, databaseHost, databasename)
	log.Println("connection string: ", dbURI)
	db, err := sqlx.Open("postgres", dbURI)
	if err != nil {
		log.Fatal(fmt.Errorf("error opening database: %w", err).Error())
	}
	if err := db.Ping(); err != nil {
		log.Fatal(fmt.Errorf("error connecting to database: %w", err).Error())
	}
	log.Println("Successfully connected to db!")
	return db
}

func getDBConfig() (username, password, databasename, databaseHost string) {
	dir, _ := os.Getwd()
	dir += "/configs"

	viper.SetConfigName("app")
	// Set the path to look for the configurations file
	viper.AddConfigPath(dir)
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	databasename = viper.GetString("DB_NAME")
	databaseHost = viper.GetString("DB_HOST")
	username = viper.GetString("DB_USERNAME")
	password = viper.GetString("DB_PASSWORD")

	return
}
