package serializer_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arcbjorn/store-management-system/sample"
	"github.com/arcbjorn/store-management-system/serializer"
)

func TestFileSerializer(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"

	laptop1 := sample.NewLaptop()
	err := serializer.WriteProtobufToBinaryFile(laptop1, binaryFile)
	require.NoError(t, err)
}