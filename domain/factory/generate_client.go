package factory

import (
	"EMQ/domain/model"
	"EMQ/infra/utils"
	"errors"
)

func GenerateClient(topic, payload, action string, qos, num int) (model.BaseClient, error) {
	if qos != 0 && qos != 1 && qos != 2 {
		return model.BaseClient{}, errors.New("Parameter error")
	}
	uuid, _ := utils.NewUuid()
	newClient := model.BaseClient{
		Broker:       "tcp://192.168.80.134:1883",
		Topic:        topic,
		ClientId:     uuid,
		User:         "admin",
		Password:     "public",
		CleanSession: false,
		Qos:          qos,
		Num:          num,
		Payload:      payload,
		Action:       action,
		Store:        ":memory:",
	}
	return newClient, nil
}
