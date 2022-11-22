package indexer

import (
	"encoding/binary"
	"math"

	"github.com/qmuntal/gltf"
)

var (
	littleEndian = binary.LittleEndian
)

// Get the nth value in the buffer described by an accessor with accessorId
func getBufferForValueAt(gltf *gltf.Document, accesorId, n uint32) []byte {
	accessor := gltf.Accessors[accesorId]
	bufferView := gltf.BufferViews[*accessor.BufferView]
	bufferId := bufferView.Buffer
	buffer := gltf.Buffers[bufferId].Data
	buffer = buffer[bufferView.ByteOffset : bufferView.ByteOffset+bufferView.ByteLength]
	valueSize := accessor.ComponentType.ByteSize() * accessor.Type.Components()

	// if no byteStride specified, the buffer is tightly packed
	byteStride := bufferView.ByteStride
	if byteStride == 0 {
		byteStride = valueSize
	}
	// fmt.Println("bytestride: ", byteStride)
	pos := accessor.ByteOffset + n*byteStride
	valueBuffer := buffer[pos : pos+valueSize]

	return valueBuffer
}

func readComponent(buff []byte, componentType gltf.ComponentType, n uint32) interface{} {
	// fmt.Println("ComponentType: ", componentType)
	switch componentType {
	case (gltf.ComponentFloat):
		inte := littleEndian.Uint32(buff[n : n+componentType.ByteSize()])
		return math.Float32frombits(inte)
	case (gltf.ComponentByte):
		return buff[n]
	case (gltf.ComponentUbyte):
		return uint8(buff[n])
	case (gltf.ComponentShort):
		return int16(littleEndian.Uint16(buff[n : n+componentType.ByteSize()]))
	case (gltf.ComponentUshort):
		return littleEndian.Uint16(buff[n : n+componentType.ByteSize()])
	case (gltf.ComponentUint):
		return littleEndian.Uint32(buff[n:componentType.ByteSize()])
	}
	return nil
}

func readValueAt(gltf *gltf.Document, accesorId, n uint32) []interface{} {
	// fmt.Println("Reached HERE")
	buffer := getBufferForValueAt(gltf, accesorId, n)
	// fmt.Println("buffer: ", buffer)
	accessor := gltf.Accessors[accesorId]
	numOfComponents := accessor.Type.Components()
	var valueComponents []interface{}
	componentType := accessor.ComponentType

	// fmt.Println("numOfComponents: ", numOfComponents)

	for i := uint32(0); i < numOfComponents*componentType.ByteSize(); i += componentType.ByteSize() {
		valueComponents = append(valueComponents, readComponent(buffer, componentType, i))
	}

	return valueComponents
}
