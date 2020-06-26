package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		ErrorIfNotNil(t, err, fmt.Sprintf("Error on marshalling sales order, %v", err))

		type response struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		}
		resp, err := http.Post("http://127.0.0.1:7070/salesorder", "application/json", bytes.NewBuffer(jsonData))
		ErrorIfNotNil(t, err, fmt.Sprintf("Error on posting sales order, %v", err))

		body, _ := ioutil.ReadAll(resp.Body)

		var respons response
		err = json.Unmarshal(body, &respons)
		ErrorIfNotNil(t, err, fmt.Sprintf("Error on unmarshalling post sales order response, %v", err))

		if respons.Message != "Sales order queued for processing." {
			LogAndFail(t, "Unable to process ssales order by sales order handler correctly, wrong response message.")
		}

		resp, err = http.Get("http://localhost:7070/salesorders")
		ErrorIfNotNil(t, err, fmt.Sprintf("Error while getting list of sales orders, %v", err))

		listBody, err := ioutil.ReadAll(resp.Body)
		ErrorIfNotNil(t, err, fmt.Sprintf("Error reading response body, %v", err))

		var list []SalesOrder
		err = json.Unmarshal(listBody, &list)
		ErrorIfNotNil(t, err, fmt.Sprintf("Error on unmarshalling list of sales orders, %v", err))

		for _, salesorder := range list {
			if salesorder.SalesOrderID == uuid3 {
				return
			}
		}
		// could not find the sales order created
		t.Fail()
	})

}
