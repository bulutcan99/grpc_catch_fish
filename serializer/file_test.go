package serializer

import (
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"pb/generator"
	"pb/proto"
	"testing"
)

func TestWriteProtobufToBinaryFile(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	newLaptop := generator.NewLaptop()
	err := WriteProtobufToBinaryFile(newLaptop, binaryFile)
	require.NoError(t, err)
}

func TestReadProtobufFromBinaryFile(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	newLaptop := generator.NewLaptop()
	err := WriteProtobufToBinaryFile(newLaptop, binaryFile)
	require.NoError(t, err)

	laptop := &pb.Laptop{}
	err = ReadProtobufFromBinaryFile(laptop, binaryFile)
	require.NoError(t, err)
	require.True(t, proto.Equal(newLaptop, laptop))
}

func TestWriteProtobufToJSONFile(t *testing.T) {
	t.Parallel()
	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"
	laptopSample := generator.NewLaptop()
	err := WriteProtobufToBinaryFile(laptopSample, binaryFile)
	require.NoError(t, err)

	laptopSample2 := &pb.Laptop{}
	err = ReadProtobufFromBinaryFile(laptopSample2, binaryFile)
	require.NoError(t, err)
	require.True(t, proto.Equal(laptopSample, laptopSample2))

	err = WriteProtobufToJSONFile(laptopSample2, jsonFile)
	require.NoError(t, err)

}
