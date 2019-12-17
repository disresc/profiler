package main

import (
	"os"

	drReceiver "github.com/disresc/lib/receiver"
	drTransmitter "github.com/disresc/lib/transmitter"
)

var transmitter *drTransmitter.Service
var name string

func main() {
	var found bool
	name, found = os.LookupEnv("name")
	if !found {
		name = "profiler"
	}

	go startTransmitter(name)

	startReceiver(name)
}

func startReceiver(name string) {
	receiver := drReceiver.NewService(name)
	receiver.RegisterData("ves", "kvmtop-cpu", 10)
	receiver.RegisterData("ves", "kvmtop-net", 10)
	receiver.RegisterData("ves", "kvmtop-disk", 10)
	receiver.Start()
	for {
		event := <-receiver.EventChannel()
		handle(event)
	}
}

func startTransmitter(name string) {
	transmitter = drTransmitter.NewService(name)
	transmitter.Start()
}
