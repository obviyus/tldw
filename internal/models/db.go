package models

import (
	"fmt"
	"os"
	"sync"
	"time"

	"tldw-server/internal/event"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var g Gorm

type Gorm struct {
	once sync.Once
	db   *gorm.DB
}

var log = event.Log

// Db returns the gorm db connection
func (g *Gorm) Db() *gorm.DB {
	g.once.Do(g.Connect)

	if g.db == nil {
		log.Fatalln("models: database not connected")
	}

	return g.db
}

// GetDbDsn checks env vars for database credentials and generates a config
func GetDbDsn() string {
	dbUser := os.Getenv("DBUSER")
	dbName := os.Getenv("DBNAME")
	dbPass := os.Getenv("DBPASS")
	dbHost := os.Getenv("DBHOST")
	if dbUser == "" || dbName == "" || dbPass == "" || dbHost == "" {
		panic("db: required env var not set")
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPass, dbName)
}

// UnscopedDb returns an unscoped database connection.
func UnscopedDb() *gorm.DB {
	return g.Db().Unscoped()
}

// Connect creates a new gorm db connector
func (g *Gorm) Connect() {
	dbDsn := GetDbDsn()
	db, err := gorm.Open(postgres.Open(dbDsn))

	if err != nil || db == nil {
		for i := 1; i <= 12; i++ {
			log.Infof("gorm.Open(%s, %s) %d\n", "postgres", dbDsn, i)
			db, err := gorm.Open(postgres.Open(dbDsn), &gorm.Config{})

			if db != nil && err == nil {
				break
			} else {
				time.Sleep(5 * time.Second)
			}
		}

		if err != nil || db == nil {
			log.Fatalln(err)
		}
	}

	pgs, _ := db.DB()
	pgs.SetMaxIdleConns(4)
	pgs.SetMaxOpenConns(256)

	g.db = db
}

// Close closes the gorm db connection.
func (g *Gorm) Close() {
	if g.db != nil {
		pgs, _ := g.db.DB()
		if err := pgs.Close(); err != nil {
			log.Fatal(err)
		}

		g.db = nil
	}
}
