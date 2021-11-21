package serializer

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Convert protobuff message to JSON bytes
func ProtobufToJsonBytes(message proto.Message) ([]byte, error) {
	marshaler := protojson.MarshalOptions {
		UseEnumNumbers: false,
		EmitUnpopulated: true,
		Indent: "  ",
		UseProtoNames: true,
	}

	return marshaler.Marshal(message)
}