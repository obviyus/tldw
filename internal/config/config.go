package config

import (
	"sync"

	"github.com/klauspost/cpuid/v2"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gorm.io/gorm"

	"tldw-server/pkg/rnd"
	"tldw-server/pkg/txt"

	"tldw-server/internal/event"
)

var log = event.Log
var once sync.Once

type Config struct {
	db    *gorm.DB
	token string
}

func NewConfig(ctx *cli.Context) *Config {
	initLogger(ctx.Bool("debug"))

	return &Config{
		token: rnd.Token(8),
	}
}

func initLogger(debug bool) {
	once.Do(
		func() {
			log.SetFormatter(
				&logrus.TextFormatter{
					DisableColors: false,
					FullTimestamp: true,
				},
			)

			if debug {
				log.SetLevel(logrus.DebugLevel)
			} else {
				log.SetLevel(logrus.InfoLevel)
			}
		},
	)
}

func (c *Config) Init() error {
	if cpuName := cpuid.CPU.BrandName; cpuName != "" {
		log.Debugf("config: running on %s", txt.Quote(cpuid.CPU.BrandName))
	}

	return c.connectDb()
}
