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

type SalesItem struct {
	ItemID    string  `json:"ID"`
	Name      string  `json:"name"`
	Qty       float32 `json:"qty"`
	UnitPrice float32 `json:"unitprice"`
	Status    string  `json:"status"`
	Location  string  `json:"location"`
}
type SalesOrder struct {
	SalesOrderID string      `json:"ID"`
	Amount       float32     `json:"amount"`
	Item         []SalesItem `json:"items"`
	Status       string      `json:"status"`
}

type SalesOrders struct {
	mutex       sync.Mutex
	salesorders []SalesOrder
}

func (s *SalesOrders) save(salesOrder SalesOrder) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.salesorders = append(s.salesorders, salesOrder)
}

func (s *SalesOrders) list() []SalesOrder {
	return s.salesorders
}

func (s *SalesOrders) ToString() string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var str strings.Builder
	for _, so := range s.salesorders {
		fmt.Fprintf(&str, "%s,", so.SalesOrderID)
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

func execute(salesOrders *SalesOrders) error {
	clientID := os.Getenv("CLIENT_ID")
	natsServers := os.Getenv("NATS_SERVER_ADDR")

	connector := NewConnector(clientID)

	log.Printf("clientID: %s, natsServers: %s\n", clientID, natsServers)

	err := connector.SetupConnectionToNATS(natsServers, nats.MaxReconnects(-1))
	if err != nil {
		return fmt.Errorf("Problem setting up connection to NATS servers, %v", err)
	}
	defer connector.Shutdown()

	go subscribeToSalesOrders(connector, salesOrders, clientID)

	select {}

	return nil
}

func subscribeToSalesOrders(conn *Connector, salesOrders *SalesOrders, clientID string) {

	conn.nc.Subscribe("Process.New.SalesOrder", func(m *nats.Msg) {
		log.Println("Subscribing to [Process.New.SalesOrder]")
		var salesOrder SalesOrder
		err := json.Unmarshal(m.Data, &salesOrder)
		if err != nil {
			log.Println("Error while unmarshaling [Process.New.SalesOrder] request")
		}

		salesOrders.save(salesOrder)

		resp := Response{clientID, "Received sales orders " + salesOrder.SalesOrderID}
		data, err := json.Marshal(resp)
		if err == nil {
			log.Println("Responding to [Process.New.SalesOrder]\n\t", resp)
			m.Respond(data)
		} else {
			log.Println("Error while marshalling [Process.New.SalesOrder] response", err)
		}

		// allocate inventory for the sales order
		// if response, no inventory,
		//     process new job order
		// if response, insufficient inventory
		//     process new job order for insufficient inventory

		// request to process new job order

	})

	conn.nc.Subscribe("All.SalesOrder.List", func(m *nats.Msg) {
		data, _ := json.Marshal(salesOrders.list())
		m.Respond(data)
	})

	conn.nc.Subscribe("SalesOrder.Details", func(m *nats.Msg) {
		// TODO
	})

	conn.nc.Subscribe("Process.New.JobOrder", func(m *nats.Msg) {
		// TODO
	})

	conn.nc.Subscribe("JobOrder.Details", func(m *nats.Msg) {
		// TODO
	})
}

func main() {
	salesOrders := SalesOrders{}

	err := execute(&salesOrders)
	if err != nil {
		log.Println(err)
	}
	runtime.Goexit()
}
