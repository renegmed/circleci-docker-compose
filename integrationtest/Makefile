init-project:
	go mod init github.com/renegmed/learn-circleci/integrationtest
.PHONY: init-project

test-originator:
	SERVER=7070 NATS_SERVER_ADDR=nats://demo.nats.io CLIENT_ID=originator-tester go test originator_test.go common.go connector.go -v

test-salesorders:
	SERVER=7070 NATS_SERVER_ADDR=nats://demo.nats.io CLIENT_ID=salesorders-tester go test salesordershandler_test.go common.go connector.go 

test-joborders:
	SERVER=7070 NATS_SERVER_ADDR=nats://demo.nats.io CLIENT_ID=joborders-tester go test jobordershandler_test.go common.go connector.go 