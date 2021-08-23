package main

import (
	"context"
	"flag"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var (
	logger       = log.New()
	address      *string
	redisAddress *string
	cache        *bool
	rdb          *redis.Client
	ctx          context.Context
)

func init() {
	// initialize logging
	logLevel := flag.Int("logLevel", 0, "defines the log level. 0=production builds. 1=dev builds.")
	address = flag.String("address", "localhost:8080", "address:port to bind http server to.")
	redisAddress = flag.String("redisAddress", "localhost:6379", "address:port for redis server.")
	cache = flag.Bool("cache", false, "enable cache for coingecko")
	flag.Parse()
	logger.Out = os.Stdout
	switch *logLevel {
	case 0:
		logger.SetLevel(log.InfoLevel)
	default:
		logger.SetLevel(log.DebugLevel)
	}
}

func main() {
	var wg sync.WaitGroup

	// Redis is used a an optional cache for coingecko data
	if *cache {
		rdb = redis.NewClient(&redis.Options{
			Addr:     *redisAddress,
			Password: "",
			DB:       0,
		})
		ctx = context.Background()
	}

	// Create the bot manager
	wg.Add(1)
	NewManager(*address, tickerCount, rdb, ctx)

	// wait forever
	wg.Wait()
}
