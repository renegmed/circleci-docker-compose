package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
)

type SalesItem struct {
	ItemID    string  `json:"ID"`
	Name      string  `json:"name"`
	Qty       float32 `json:"qty"`
	UnitPrice float32 `json:"unitprice"`
}
type SalesOrder struct {
	SalesOrderID string `json:"ID"`
	Amount       float32
	Items        []SalesItem `json:"items"`
}

func TestOriginator(t *testing.T) {
	server := os.Getenv("SERVER")
	natsServerAddr := os.Getenv("NATS_SERVER_ADDR")
	clientID := os.Getenv("CLIENT_ID")

	t.Run("Sending Sales order to Originator", func(t *testing.T) {
		// if server != "7070" {
		// 	fmt.Println(server)
		// 	t.Fail()
		// 	return
		// }

		// if natsServerAddr != "nats://demo.nats.io" {
		// 	fmt.Println(natsServerAddr)
		// 	t.Fail()
		// 	return
		// }

		// if clientID != "tester" {
		// 	fmt.Println(clientID)
		// 	t.Fail()
		// 	return
		// }

		// resp, err := http.Get("http://localhost:7070/salesorders")
		// if err != nil {
		// 	log.Println("Error while getting list of sales orders,", err)
		// 	t.Fail()
		// }
		// body, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	log.Println("Error reading response body, ", err)
		// 	t.Fail()
		// }
		// t.Log(string(body))

		salesOrder := SalesOrder{}

		salesItem := SalesItem{"1001", "computer table", 2, 150.00}
		salesOrder.Items = append(salesOrder.Items, salesItem)
		salesItem = SalesItem{"1002", "computer chair", 2, 65.00}
		salesOrder.Items = append(salesOrder.Items, salesItem)
		salesOrder.Amount = 430.00
		salesOrder.SalesOrderID = "5051-1001"

		jsonData, err := json.Marshal(salesOrder)
		if err != nil {
			log.Println("Error on marshalling sales order, ", err)
			t.Fail()
		}
		type response struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		}
		resp, err := http.Post("http://127.0.0.1:7070/salesorder", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("Error on posting sales order, ", err)
			t.Fail()
		}

		body, _ := ioutil.ReadAll(resp.Body)

		var respons response
		err = json.Unmarshal(body, &respons)
		if err != nil {
			log.Println("Error on unmarshalling post sales order response, ", err)
			t.Fail()
		}
		t.Log(respons)
		if respons.ID != salesOrder.SalesOrderID {
			t.Log("Unable to process sales order by sales order handler correctly, wrong sales order ID.")
			t.Fail()
		}
		if respons.Message != "Sales order queued for processing." {
			t.Log("Unable to process ssales order by sales order handler correctly, wrong response message.")
			t.Fail()
		}

		resp, err = http.Get("http://localhost:7070/salesorders")
		if err != nil {
			log.Println("Error while getting list of sales orders,", err)
			t.Fail()
		}
		listBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading response body, ", err)
			t.Fail()
		}
		var list []SalesOrder
		err = json.Unmarshal(listBody, &list)
		if err != nil {
			log.Println("Error on unmarshalling list of sales orders ", err)
			t.Fail()
		}
		if len(list) != 1 {
			log.Println("Error on getting list of orders from order sales handler,", err)
			t.Fail()
		}
		t.Log(string(listBody))

		t.Log("Success")
	})

	t.Run("Sending Job order to Originator", func(t *testing.T) {
		if server != "7070" {
			fmt.Println(server)
			t.Fail()
			return
		}

		if natsServerAddr != "nats://demo.nats.io" {
			fmt.Println(natsServerAddr)
			t.Fail()
			return
		}

		if clientID != "tester" {
			fmt.Println(clientID)
			t.Fail()
			return
		}

		t.Log("Success")
	})

}

// func main() {
// 	fmt.Println("main ...")
// }
