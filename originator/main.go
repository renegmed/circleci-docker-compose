package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	nats "github.com/nats-io/nats.go"
)

type JobOrder struct {
	ID   string `json:"ID"`
	Name string `json:"name"`
}

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

type server struct {
	nc *nats.Conn
}

var natsServer server

func salesOrderHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		resp, err := natsServer.nc.Request("All.SalesOrder.List", nil, 500*time.Millisecond)
		if err != nil {
			responseError(w, http.StatusBadRequest, "Error on request 'All.SalesOrder.List'")
			return
		}

		if resp == nil {
			responseError(w, http.StatusBadRequest, "Problem, has response but no message.")
			return
		}

		//log.Println("Response to request All.NotStarted.Tasks", string(resp.Data))

		var salesOrders []SalesOrder
		err = json.Unmarshal(resp.Data, &salesOrders)
		if err != nil {
			log.Println("Error on unmarshal salesOrders", err)
		}

		log.Println(salesOrders)
		responseOk(w, salesOrders)

	case http.MethodPost:
		body, _ := ioutil.ReadAll(r.Body)

		var salesOrder SalesOrder
		err := json.Unmarshal(body, &salesOrder)
		if err != nil {
			responseError(w, http.StatusBadRequest, "Invalid data")
			return
		}

		resp, err := natsServer.nc.Request("Process.New.SalesOrder", body, 500*time.Millisecond)
		if err != nil {
			responseError(w, http.StatusBadRequest, "Error on request 'Process.New.SalesOrder'")
			return
		}

		log.Println("Response to request Process.New.SalesOrder'", string(resp.Data))

		type response struct {
			ID      string `json:"id"`
			Message string `json:"message"`
		}

		responseOk(w, response{ID: salesOrder.SalesOrderID, Message: "Sales order queued for processing."})
	}
}

func jobOrderHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func main() {

	serverPort := os.Getenv("SERVER_PORT")
	clientID := os.Getenv("CLIENT_ID")
	natsServers := os.Getenv("NATS_SERVER_ADDR")

	connector := NewConnector(clientID)

	log.Printf("serverPort: %s, clientID: %s, natsServers: %s\n", serverPort, clientID, natsServers)

	err := connector.SetupConnectionToNATS(natsServers, nats.MaxReconnects(-1))
	if err != nil {
		log.Printf("Problem setting up connection to NATS servers, %v", err)
		runtime.Goexit()
	}

	defer connector.Shutdown()

	natsServer = server{nc: connector.NATS()}

	http.HandleFunc("/salesorder", salesOrderHandler)
	http.HandleFunc("/salesorders", salesOrderHandler)
	http.HandleFunc("/joborder", jobOrderHandler)
	log.Printf("====== Generator server listening on port %s...", serverPort)
	if err := http.ListenAndServe(":"+serverPort, nil); err != nil {
		log.Fatal(err)
	}
}

func responseError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	body := map[string]string{
		"error": message,
	}
	json.NewEncoder(w).Encode(body)
}

func responseOk(w http.ResponseWriter, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(body)
}
