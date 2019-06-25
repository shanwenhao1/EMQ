package model

type BaseClient struct {
	Broker       string // the broker URI. ex: tcp://10.10.1.1:1883
	Topic        string // The topic name to/from which to publish/subscribe
	ClientId     string
	User         string
	Password     string
	CleanSession bool   // Set Clean Session(default false)
	Qos          int    // The Quality of Service 0,1,2 (default 0)
	Num          int    // The number of messages to publish or subscribe (default 1)
	Payload      string // The message text to publish (default empty)
	Action       string // Action publish or subscribe (required)
	Store        string // The Store Directory (default use memory store)
}
