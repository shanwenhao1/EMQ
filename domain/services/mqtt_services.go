package services

import (
	"EMQ/domain/factory"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"os"
)

type MQttServices struct {
	Topic string // The topic name to/from which to publish/subscribe
	Qos   int    // The Quality of Service 0,1,2 (default 0)
	Store string // The Store Directory (default use memory store)
}

func (mQttSer MQttServices) generateOpt() (*mqtt.ClientOptions, error) {
	newClient, err := factory.GenerateClient(mQttSer.Topic, mQttSer.Qos)
	if err != nil {
		return nil, err
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker(newClient.Broker)
	opts.SetClientID(newClient.ClientId)
	opts.SetUsername(newClient.User)
	opts.SetPassword(newClient.Password)
	opts.SetCleanSession(newClient.CleanSession)

	if mQttSer.Store != "" {
		opts.SetStore(mqtt.NewFileStore(mQttSer.Store))
	}
	return opts, nil
}

func (mQttSer MQttServices) Publish(msg string) error {
	opts, err := mQttSer.generateOpt()
	if err != nil {
		return err
	}
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println(fmt.Sprintf("Sample Publisher Started: topic: %s, msg: %s", mQttSer.Topic, msg))
	token := client.Publish(mQttSer.Topic, byte(mQttSer.Qos), false, msg)
	token.Wait()
	//// push message 10 times
	//for i := 0; i < 10; i++ {
	//	fmt.Println("---- doing publish ----")
	//	token := client.Publish(topic, byte(qos), false, msg)
	//	token.Wait()
	//}
	client.Disconnect(250)
	fmt.Println("Sample Publisher Disconnected")
	return nil
}

func (mQttSer MQttServices) Subscribe(num int) error {
	receiveCount := 0
	choke := make(chan [2]string)
	opts, err := mQttSer.generateOpt()
	if err != nil {
		return err
	}

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		choke <- [2]string{msg.Topic(), string(msg.Payload())}
	})
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := client.Subscribe(mQttSer.Topic, byte(mQttSer.Qos), nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	for receiveCount < num {
		incoming := <-choke
		fmt.Printf("RECEIVED TOPIC: %s MESSAGE: %s\n", incoming[0], incoming[1])
		receiveCount++
	}
	client.Disconnect(250)
	fmt.Println("Sample Subscriber Disconnected")
	return nil
}
