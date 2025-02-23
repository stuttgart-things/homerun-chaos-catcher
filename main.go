/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package main

import (
	"github.com/stuttgart-things/homerun-chaos-catcher/internal"
	"github.com/stuttgart-things/homerun-chaos-catcher/streams"
)

func main() {
	// PRINT BANNER + VERSION INFO
	internal.PrintBanner()

	// SUBSCRIBE TO REDIS STREAM
	streams.SubscribeToRedisStream()
}
