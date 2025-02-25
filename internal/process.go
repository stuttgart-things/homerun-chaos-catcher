/*
Copyright © 2024 PATRICK HERMANN patrick.hermann@sva.de
*/

package internal

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nitishm/go-rejson/v4"
	"github.com/pterm/pterm"
	goredis "github.com/redis/go-redis/v9"
	k8s "github.com/stuttgart-things/homerun-chaos-catcher/kubernetes"
	homerun "github.com/stuttgart-things/homerun-library"
	"github.com/stuttgart-things/redisqueue"
	sthingsBase "github.com/stuttgart-things/sthingsBase"
)

var (
	redisClient = goredis.NewClient(&goredis.Options{
		Addr:     redisServer + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})
	redisJSONHandler = rejson.NewReJSONHandler()
	redisServer      = os.Getenv("REDIS_SERVER")
	redisPort        = os.Getenv("REDIS_PORT")
	redisPassword    = os.Getenv("REDIS_PASSWORD")
	profilePath      = os.Getenv("PROFILE_PATH")
	pathToKubeconfig = os.Getenv("KUBECONFIG")
	logger           = pterm.DefaultLogger.WithLevel(pterm.LogLevelTrace)
)

func ProcessStreams(msg *redisqueue.Message) error {
	redisJSONHandler.SetGoRedisClientWithContext(context.Background(), redisClient)

	fmt.Println(msg.ID)
	fmt.Println(msg.Stream)

	messageID := fmt.Sprintf("%v", msg.Values["messageID"])

	eventMessage, err := homerun.GetMessageJSON(messageID, redisJSONHandler)
	if err != nil {
		logger.Error("ERROR", logger.Args("", err))
		return err
	}

	// READ CONFIGURATION FROM YAML FILE
	filePath := profilePath

	// Load the YAML configuration
	config, err := loadConfiguration(filePath)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
	}

	// OUTPUT THE MESSAGE
	var buf bytes.Buffer
	homerun.PrintTable(
		&buf,
		table.Row{"SEVERITY", "TITLE", "SYSTEM"},
		table.Row{eventMessage.Severity, eventMessage.Title, eventMessage.System},
		table.StyleColoredBlackOnBlueWhite,
	)

	fmt.Println(buf.String())

	// CHECK FOR TIMESTAMP
	ts := sthingsBase.ConvertStringToInteger(eventMessage.Timestamp)
	timestamp := int64(ts)
	time_diff, err := strconv.Atoi(os.Getenv("TIME_DIFFERENCE_MESSAGES"))

	// CREATE KUBERNETES CLIENT (BEFORE EVENT, IF CONNECTION BREAKS)
	k8sClient := k8s.CreateKubernetesClient(pathToKubeconfig)

	if messageTimeValid(timestamp, int64(time_diff)) {

		for name, chaosConfig := range config.ChaosEvents {

			// CHECK CONFIGURATION
			systemEnabled := false
			severityEnabled := false

			fmt.Printf("CHAOS: %s\n", name)

			// CHECK FOR SYSTEMS
			if chaosConfig.Systems[0] == "*" {
				fmt.Println("ALL SYSTEMS - CHECK")
				systemEnabled = true
			} else if sthingsBase.CheckForStringInSlice(chaosConfig.Systems, eventMessage.System) {
				fmt.Println("SYSTEMS - CHECK")
				systemEnabled = true
			}

			// CHECK FOR SEVERITY
			if sthingsBase.CheckForStringInSlice(chaosConfig.Severity, eventMessage.Severity) {
				severityEnabled = true
				fmt.Println("SEVERITY FOUND - CHECK")
			}

			// KUBERNETES ACTIONS
			if systemEnabled && severityEnabled {
				// SendToWLED(eventMessage.Severity, eventMessage.System)
				CreateChaos(chaosConfig.Resource, chaosConfig.Count, chaosConfig.Operation, k8sClient)
			} else {
				fmt.Println("SYSTEM OR SEVERITY DOES NOT MATCH")
			}
		}
	} else {
		fmt.Println("MESSAGE TOO OLD")
	}
	return nil
}

func messageTimeValid(timestamp, maxDiff int64) (triggerLight bool) {
	// GET THE CURRENT TIME AS UNIX TIMESTAMP (SECONDS SINCE EPOCH)
	currentTime := time.Now().Unix()

	fmt.Println("CURRENT TIME", currentTime)
	fmt.Println("TIMESTAMP EVENT", timestamp)

	// CALCULATE THE DIFFERENCE IN SECONDS
	diff := int64(math.Abs(float64(currentTime - timestamp)))

	// CHECK IF THE TIMESTAMP IS WITHIN GIVEN SECONDS OF THE CURRENT TIME
	if diff < maxDiff {
		fmt.Println("IN TIME - FIRE")
		triggerLight = true
		return triggerLight
	} else {
		fmt.Printf("The timestamp is %d seconds off from the current time.\n", diff)
		triggerLight = false
		return triggerLight
	}
}
