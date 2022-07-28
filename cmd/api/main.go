package main

import (
	"context"
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	doljara "github.com/stripe-island/dol-jara-server"
)

func main() {
	ctx := context.Background()
	if err := funcframework.RegisterHTTPFunctionContext(ctx, "/DoljaraRooms", doljara.DoljaraRooms); err != nil {
		log.Fatalf("Failed to register function:: %v", err)
	}

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("Failed to start functions framework: %v", err)
	}
}
