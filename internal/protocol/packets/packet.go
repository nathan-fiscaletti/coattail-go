package packets

import (
	"context"
	"reflect"
)

type Packet interface {
	Handle(ctx context.Context) (Packet, error)
}

// GetPacketData returns a map of the fields of a packet.
func GetPacketData(packet Packet) map[string]interface{} {
	// Retrieve the fields of the packet using their JSON tags.
	data := make(map[string]interface{})
	for i := 0; i < reflect.ValueOf(packet).NumField(); i++ {
		field := reflect.ValueOf(packet).Field(i)
		fieldName := reflect.TypeOf(packet).Field(i).Tag.Get("json")
		data[fieldName] = field.Interface()
	}

	return data
}
