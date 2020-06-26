package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func TestOriginator(t *testing.T) {

	t.Run("Sending Sales order to Originator", func(t *testing.T) {

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

		jsonData, err := json.Marshal(salesOrder)
		if err != nil {
			fmt.Println("Error on marshalling sales order, ", err)
			t.Fail()
		}
		type response struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		}
		resp, err := http.Post("http://127.0.0.1:7070/salesorder", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Error on posting sales order, ", err)
			t.Fail()
		}

		body, _ := ioutil.ReadAll(resp.Body)

		var respons response
		err = json.Unmarshal(body, &respons)
		if err != nil {
			fmt.Println("Error on unmarshalling post sales order response, ", err)
			t.Fail()
		}
		t.Log("++++ RESPONS:\n\t", respons)

		if respons.Message != "Sales order queued for processing." {
			fmt.Println("Unable to process ssales order by sales order handler correctly, wrong response message.")
			t.Fail()
		}

		resp, err = http.Get("http://localhost:7070/salesorders")
		if err != nil {
			fmt.Println("Error while getting list of sales orders,", err)
			t.Fail()
		}
		listBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body, ", err)
			t.Fail()
		}

		var list []SalesOrder
		err = json.Unmarshal(listBody, &list)
		if err != nil {
			log.Println("Error on unmarshalling list of sales orders ", err)
			t.Fail()
		}

		for _, salesorder := range list {
			if salesorder.SalesOrderID == uuid3 {
				return
			}
		}
		// could not find the sales order created
		t.Fail()
	})

}
