package main

import (
	"log"
	"testing"

	nats "github.com/nats-io/nats.go"
	"github.com/satori/uuid"
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
	Items        []SalesItem `json:"items"`
	Status       string      `json:"status"`
}

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

type Response struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func NewUuid() string {
	newUUID := uuid.NewV4()
	// if err != nil {
	// 	return "", fmt.Errorf("Error on generating new uuid, %v", err)
	// }
	return newUUID.String() //, nil
}

func Connect(clientID, natsServers string) (*Connector, error) {

	connector := NewConnector(clientID)

	//log.Printf("clientID: %s, natsServers: %s\n", clientID, natsServers)

	err := connector.SetupConnectionToNATS(natsServers, nats.MaxReconnects(-1))
	if err != nil {
		log.Printf("Problem setting up connection to NATS servers, %v", err)
		return nil, err
	}

	return connector, nil
}

func ErrorIfNotNil(t *testing.T, err error, s string) {
	if err != nil {
		LogAndFail(t, s)
	}
}

func LogAndFail(t *testing.T, s string) {
	t.Log(s)
	t.Fail()
}
