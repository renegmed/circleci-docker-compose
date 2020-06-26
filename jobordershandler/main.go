package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"

	nats "github.com/nats-io/nats.go"
)

type Item struct {
	ItemID   string  `json:"ID"`
	Name     string  `json:"name"`
	Qty      float32 `json:"qty"`
	Status   string  `json:"status"`
	Location string  `json:"location"`
}

type JobOrder struct {
	JobOrderID string `json:"ID"`
	Name       string `json:"name"`
	Items      []Item `json:"item"`
	Status     string `json:"status"`
}

type JobOrders struct {
	jobOrders []JobOrder
	mutex     sync.Mutex
}

func (j *JobOrders) save(jobOrder JobOrder) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.jobOrders = append(j.jobOrders, jobOrder)
}

func (j *JobOrders) list() []JobOrder {
	return j.jobOrders
}

func (j *JobOrders) ToString() string {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	var str strings.Builder
	for _, jobOrder := range j.jobOrders {
		fmt.Fprintf(&str, "%s,", jobOrder.JobOrderID)
	}
	result := str.String()

	if len(result) > 0 {
		// remove last comma
		return strings.TrimSuffix(result, ",")
	}
	return result
}

type Response struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func execute() error {

	clientID := os.Getenv("CLIENT_ID")
	natsServers := os.Getenv("NATS_SERVER_ADDR")

	connector := NewConnector(clientID)

	log.Printf("clientID: %s, natsServers: %s\n", clientID, natsServers)

	err := connector.SetupConnectionToNATS(natsServers, nats.MaxReconnects(-1))
	if err != nil {
		return fmt.Errorf("Problem setting up connection to NATS servers, %v", err)
	}
	defer connector.Shutdown()

	jobOrders := JobOrders{}

	go subscribeToJobOrders(connector, &jobOrders, clientID)

	select {}

	return nil
}

func subscribeToJobOrders(conn *Connector, jobOrders *JobOrders, clientID string) {

	conn.nc.Subscribe("Process.New.JobOrder", func(m *nats.Msg) {
		log.Println("Subscribing to [Process.New.JobOrder]")
		var jobOrder JobOrder
		err := json.Unmarshal(m.Data, &jobOrder)
		if err != nil {
			log.Println("Error while unmarshaling [Process.New.JobOrder] request")
		}

		jobOrders.save(jobOrder)

		resp := Response{clientID, "Received job order " + jobOrder.JobOrderID}
		data, err := json.Marshal(resp)
		if err == nil {
			log.Println("Responding to [Process.New.JobOrder]\n\t", resp)
			m.Respond(data)
		} else {
			log.Println("Error while marshalling [Process.New.JobOrder] response", err)
		}
	})

	conn.nc.Subscribe("All.JobOrder.List", func(m *nats.Msg) {
		data, _ := json.Marshal(jobOrders.list())
		m.Respond(data)
	})
}

func main() {
	err := execute()
	if err != nil {
		log.Println(err)
	}
	runtime.Goexit()
}
