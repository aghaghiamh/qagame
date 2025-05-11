package broker

import "github.com/aghaghiamh/gocast/QAGame/entity"

type Broker interface {
	Publish(event entity.Event, payload string)
}
