/*
gRPC Client
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	xr "github.com/nleiva/xrgrpc"
)

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	log.Printf("This process took %s\n", elapsed)
}

func main() {
	// To time this process
	defer timeTrack(time.Now())

	// Encoding option; defaults to JSON
	enc := flag.String("enc", "json", "Encoding: 'json' or 'text'")
	// CLI to issue; defaults to "show grpc status"
	cli := flag.String("cli", "show grpc status", "Command to execute")
	// Config file; defaults to "config.json"
	cfg := flag.String("cfg", "../input/config.json", "Configuration file")

	flag.Parse()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := r.Int63n(10000)
	output := "Empty"

	// Define target parameters from the configuration file
	targets := xr.NewDevices()
	err := xr.DecodeJSONConfig(targets, *cfg)
	if err != nil {
		log.Fatalf("Could not read the config: %v\n", err)
	}

	// Setup a connection to the target
	conn, err := xr.Connect(targets.Routers[0])
	if err != nil {
		log.Fatalf("Could not setup a client connection to %s, %v", targets.Routers[0].Host, err)
	}
	defer conn.Close()

	// Return show command output based on encoding selected
	switch *enc {
	case "text":
		output, err = xr.ShowCmdTextOutput(conn, *cli, id)
	case "json":
		output, err = xr.ShowCmdJSONOutput(conn, *cli, id)
	default:
		log.Fatalf("Do NOT recognize encoding: %v\n", *enc)
	}
	if err != nil {
		log.Fatalf("Couldn't get the cli output: %v\n", err)
	}
	fmt.Printf("\nOutput from %s\n %s\n", targets.Routers[0].Host, output)
}
