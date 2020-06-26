package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	nats "github.com/nats-io/nats.go"
)

func TestNewSalesOrder(t *testing.T) {
	clientID := os.Getenv("CLIENT_ID")
	natsServers := os.Getenv("NATS_SERVER_ADDR")

	connector, err := Connect(clientID, natsServers)
	ErrorIfNotNil(t, err, fmt.Sprintf("Problem connecting to NATS server, %v", err))

	defer connector.Shutdown()

	t.Run("Add sales order and verify saved in sales order handler", func(t *testing.T) {

		// Post sales order to Sales Order Handler
		salesOrder := SalesOrder{}

		uuid1 := NewUuid()

		salesItem := SalesItem{uuid1, "computer table", 2, 150.00, "", ""}
		salesOrder.Items = append(salesOrder.Items, salesItem)

		uuid2 := NewUuid()

		salesItem = SalesItem{uuid2, "computer chair", 2, 65.00, "", ""}
		salesOrder.Items = append(salesOrder.Items, salesItem)
		salesOrder.Amount = 430.00

		uuid3 := NewUuid()

		salesOrder.SalesOrderID = uuid3

		response, err := postSalesOrder(connector.NATS(), salesOrder)
		ErrorIfNotNil(t, err, fmt.Sprintf("Error, sales order handler could not post sales order id %v", uuid3))

		if response.Message != "Received sales orders "+uuid3 {
			LogAndFail(t, fmt.Sprintf("Error, sales order %s could not be posted by sales order handler.", uuid3))
		}
		// get list of sales orders from sales orders handler, regardless of status

		salesOrders, err := listSalesOrders(connector.NATS())
		ErrorIfNotNil(t, err, fmt.Sprintf("%v", err))

		for _, salesOrder := range salesOrders {
			if salesOrder.SalesOrderID == uuid3 {
				return
			}
		}
		// Posted sales order not found in Sales Order Handler list
		t.Fail()

	})
}

func listSalesOrders(conn *nats.Conn) ([]SalesOrder, error) {
	resp, err := conn.Request("All.SalesOrder.List", nil, 500*time.Millisecond)
	if err != nil {
		return nil, fmt.Errorf("Error on request 'All.SalesOrder.List'")

	}

	if resp == nil {
		return nil, fmt.Errorf("Problem, has response but no message.")
	}

	var salesOrders []SalesOrder
	err = json.Unmarshal(resp.Data, &salesOrders)
	if err != nil {
		return nil, fmt.Errorf("Error on unmarshal salesOrders, %v", err)
	}
	return salesOrders, nil
}

func postSalesOrder(conn *nats.Conn, salesOrder SalesOrder) (Response, error) {

	body, err := json.Marshal(salesOrder)
	if err != nil {
		return Response{}, fmt.Errorf("Error on marshalling sales order, %v", err)
	}

	resp, err := conn.Request("Process.New.SalesOrder", body, 500*time.Millisecond)
	if err != nil {
		return Response{}, fmt.Errorf("Error on request 'Process.New.SalesOrder'")
	}
	var response Response
	err = json.Unmarshal(resp.Data, &response)
	if err != nil {
		return Response{}, fmt.Errorf("Error on marshaling response, %v", err)
	}

	//og.Println("+++ Reponse, \n\t ", response)
	return response, nil
}
