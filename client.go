package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type Client struct {
	Name string
	id   uuid.UUID
}

type ClientMap map[uuid.UUID]Client

func (cm ClientMap) AddClient(name string) (uuid.UUID, error) {
	newUuid, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("Error while generating new UUID")
		return uuid.Nil, err
	}

	cm[newUuid] = Client{id: newUuid, Name: name}
	return newUuid, nil
}

func (cm ClientMap) GetClientByID(id string) (*Client, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		fmt.Printf("Could not parse string '%v' into an UUID", id)
		return nil, err
	}

	client, ok := cm[parsedId]
	if !ok {
		fmt.Printf("Did not find any client for UUID %v\n", parsedId)
		return nil, errors.New("no client found for the provided UUID")
	}

	return &client, nil
}
