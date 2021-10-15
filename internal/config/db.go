package config

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"tldw-server/internal/models"
	"tldw-server/internal/mutex"
)

func (c *Config) Db() *gorm.DB {
	if c.db == nil {
		log.Fatalln("config: database not connected")
	}

	return c.db
}

// CloseDb closes the db connection (if any).
func (c *Config) CloseDb() error {
	if c.db != nil {
		psg, _ := c.db.DB()
		if err := psg.Close(); err == nil {
			c.db = nil
		} else {
			return err
		}
	}

	return nil
}

func (c *Config) InitDb() {
	models.MigrateDb()
}

func (c *Config) connectDb() error {
	mutex.Db.Lock()
	defer mutex.Db.Unlock()

	dbDsn := models.GetDbDsn()
	db, err := gorm.Open(postgres.Open(dbDsn))
	if err != nil || db != nil {
		for i := 1; i <= 12; i++ {
			db, err = gorm.Open(postgres.Open(dbDsn))

			if db != nil && err == nil {
				break
			}

			time.Sleep(5 * time.Second)
		}

		if err != nil || db == nil {
			log.Fatal(err)
		}
	}

	psg, _ := db.DB()
	psg.SetMaxIdleConns(4)
	psg.SetMaxOpenConns(256)
	psg.SetConnMaxLifetime(10 * time.Minute)

	c.db = db

	return err
}
