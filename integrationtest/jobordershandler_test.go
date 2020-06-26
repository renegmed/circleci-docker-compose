package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	nats "github.com/nats-io/nats.go"
)

func TestNewJobOrder(t *testing.T) {
	clientID := os.Getenv("CLIENT_ID")
	natsServers := os.Getenv("NATS_SERVER_ADDR")

	connector, err := Connect(clientID, natsServers)
	if err != nil {
		t.Errorf("Problem connecting to NATS server, %v", err)
		t.Fail()
	}
	defer connector.Shutdown()

	t.Run("Add job order and verify saved in job orders handler", func(t *testing.T) {

		// Post job order to Job Orders Handler
		jobOrder := JobOrder{}

		uuid1 := NewUuid()

		item := Item{uuid1, "computer table", 2, "", ""}
		jobOrder.Items = append(jobOrder.Items, item)

		uuid2 := NewUuid()

		item = Item{uuid2, "chair", 2, "", ""}
		jobOrder.Items = append(jobOrder.Items, item)

		uuid3 := NewUuid()

		jobOrder.JobOrderID = uuid3

		response, err := postJobOrder(connector.NATS(), jobOrder)
		if err != nil {
			t.Logf("Error, job order handler could not post job order id %v", uuid3)
			t.Fail()
		}

		if response.Message != "Received job order "+uuid3 {
			t.Logf("Error, job order handler could not post job order %s.", uuid3)
			t.Fail()
		}
		// get list of job orders from sales orders handler, regardless of status

		jobOrders, err := listJobOrders(connector.NATS())
		if err != nil {
			t.Errorf("%v", err)
			t.Fail()
		}
		for _, jobOrder := range jobOrders {
			if jobOrder.JobOrderID == uuid3 {
				return
			}
		}
		// Posted sales order not found in Job Order Handler list
		t.Fail()

	})
}

func listJobOrders(conn *nats.Conn) ([]JobOrder, error) {
	resp, err := conn.Request("All.JobOrder.List", nil, 500*time.Millisecond)
	if err != nil {
		return nil, fmt.Errorf("Error on request 'All.JobOrder.List'")

	}

	if resp == nil {
		return nil, fmt.Errorf("Problem, has response but no message.")
	}

	var jobOrders []JobOrder
	err = json.Unmarshal(resp.Data, &jobOrders)
	if err != nil {
		return nil, fmt.Errorf("Error on unmarshal jobOrders, %v", err)
	}
	return jobOrders, nil
}

func postJobOrder(conn *nats.Conn, jobOrder JobOrder) (Response, error) {

	body, err := json.Marshal(jobOrder)
	if err != nil {
		return Response{}, fmt.Errorf("Error on marshalling job order, %v", err)
	}

	resp, err := conn.Request("Process.New.JobOrder", body, 500*time.Millisecond)
	if err != nil {
		return Response{}, fmt.Errorf("Error on request 'Process.New.JobOrder'")
	}
	var response Response
	err = json.Unmarshal(resp.Data, &response)
	if err != nil {
		return Response{}, fmt.Errorf("Error on marshaling response, %v", err)
	}

	log.Println("+++ Reponse, \n\t ", response)
	return response, nil
}
