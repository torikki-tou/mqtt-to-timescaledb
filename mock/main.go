package main

import (
	"context"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

func main() {
	signalCtx, signalStop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer signalStop()
	
	mqttURL, err := url.Parse("mqtt://localhost:1883")
	if err != nil {
		panic(err)
	}
	
	mqttConn, err := mqtt.NewConnection(context.TODO(), mqtt.ClientConfig{
		ServerUrls:                    []*url.URL{mqttURL},
		CleanStartOnInitialConnection: true,
	})
	if err != nil {
		panic(err)
	}
	
	connectCtx, connectCancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer connectCancel()

	if err := mqttConn.AwaitConnection(connectCtx); err != nil {
		panic(err)
	}
		
	ticker := time.NewTicker(10*time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if _, err := mqttConn.Publish(context.TODO(), &paho.Publish{
				QoS:   1,
				Topic: "test",
				Payload: []byte(`{"field":"test"}`),
			}); err != nil {
				panic(err)
			}
			continue
		case <-signalCtx.Done():
		}
		break
	}
	
	disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer disconnectCancel()
	mqttConn.Disconnect(disconnectCtx)
}
