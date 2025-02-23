/*
Copyright © 2024 PATRICK HERMANN patrick.hermann@sva.de
*/

package streams

import (
	"fmt"
	"os"
	"time"

	"github.com/stuttgart-things/homerun-chaos-catcher/internal"

	"github.com/redis/go-redis/v9"

	"github.com/stuttgart-things/redisqueue"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
)

var (
	redisServer        = os.Getenv("REDIS_SERVER")
	redisPort          = os.Getenv("REDIS_PORT")
	redisPassword      = os.Getenv("REDIS_PASSWORD")
	redisStream        = os.Getenv("REDIS_STREAM")
	redisConsumerGroup = os.Getenv("REDIS_CONSUMER_GROUP")
	log                = sthingsBase.StdOutFileLogger(logfilePath, "2006-01-02 15:04:05", 50, 3, 28)
	logfilePath        = "/tmp/homerun-chaos-catcher.log"
)

func SubscribeToRedisStream() {
	// PRINT BANNER + VERSION INFO
	log.Info("HOMERUN-CHAOS-CATCHER STARTED")
	log.Info("REDIS SERVER " + redisServer)
	log.Info("REDIS PORT " + redisPort)
	log.Info("REDIS STREAM " + redisStream)

	c, err := redisqueue.NewConsumerWithOptions(&redisqueue.ConsumerOptions{
		VisibilityTimeout: 60 * time.Second,
		BlockingTimeout:   5 * time.Second,
		GroupName:         redisConsumerGroup,
		ReclaimInterval:   1 * time.Second,
		BufferSize:        100,
		Concurrency:       1,
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     redisServer + ":" + redisPort,
			Password: redisPassword,
			DB:       0,
		}),
	})

	if err != nil {
		panic(err)
	}

	c.Register(redisStream, internal.ProcessStreams)

	go func() {
		for err := range c.Errors {
			fmt.Printf("err: %+v\n", err)
		}
	}()

	log.Info("START READING STREAM: ", redisStream+" ON "+redisServer+":"+redisPort)

	c.Run()

	log.Warn("READING STOPPED")
}
