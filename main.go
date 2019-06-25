package main

import (
	"EMQ/domain/factory"
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"os"
)

func main() {
	topic := "test"
	message := map[string]interface{}{
		"testMessage": "I want you",
	}
	payload, _ := json.Marshal(message)
	action := ""
	qos := 0
	num := 1
	newClient, err := factory.GenerateClient(topic, string(payload), action, qos, num)
	if err != nil {
		panic(err)
	}
	newClient.Action = "pub"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(newClient.Broker)
	opts.SetClientID(newClient.ClientId)
	opts.SetUsername(newClient.User)
	opts.SetPassword(newClient.Password)
	opts.SetCleanSession(newClient.CleanSession)

	if newClient.Store != ":memory:" {
		opts.SetStore(mqtt.NewFileStore(newClient.Store))
	}

	if newClient.Action == "pub" {
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		fmt.Println("Sample Publisher Started")
		for i := 0; i < newClient.Num; i++ {
			fmt.Println("---- doing publish ----")
			token := client.Publish(newClient.Topic, byte(newClient.Qos), false, newClient.Payload)
			token.Wait()
		}
		client.Disconnect(250)
		fmt.Println("Sample Publisher Disconnected")
	} else {
		receiveCount := 0
		choke := make(chan [2]string)
		opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
			choke <- [2]string{msg.Topic(), string(msg.Payload())}
		})
		client := mqtt.NewClient(opts)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
		if token := client.Subscribe(newClient.Topic, byte(newClient.Qos), nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
		for receiveCount < newClient.Num {
			incoming := <-choke
			fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
			receiveCount++
		}
		client.Disconnect(250)
		fmt.Println("Sample Subscriber Disconnected")
	}

}
