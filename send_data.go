package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nats-io/stan.go"
	"github.com/polonkoevv/wb-tech/internal/models"
)

func main() {
	data, err := os.ReadFile("model.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)
	v := models.Order{}
	err = json.Unmarshal(data, &v)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(v)

	clusterID := "test-cluster"
	clientID := "test-client2"
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("http://localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}

	err = sc.Publish("orders", data) // does not return until an ack has been received from NATS Streaming
	if err != nil {
		log.Fatal(err)
	}

	err = sc.Close()
	if err != nil {
		log.Fatal(err)
	}
}
