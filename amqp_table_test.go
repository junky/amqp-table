package amqptable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadTable(t *testing.T) {
	data := []byte{
		0, 0, 0, 24, // table length
		4,                 // "name" length
		110, 97, 109, 101, // "name" in bytes
		83,         // "S" - Longstr type (ASCII code for 'S')
		0, 0, 0, 5, // "topic" length
		116, 111, 112, 105, 99, // "topic"
		3,            // "age" length
		97, 103, 101, // "age"
		73,          // "I" - int32 type (ASCII code for 'I')
		0, 0, 0, 30, // 30 as int32
	}

	table, err := ReadTable(data)
	if err != nil {
		t.Fatalf("Failed to read table: %v", err)
	}

	expectedTable := map[string]any{
		"name": "topic",
		"age":  int32(30),
	}

	assert.Equal(t, expectedTable, table)
}

func TestWriteTable(t *testing.T) {
	table := map[string]any{
		"name": "topic",
		"age":  int32(30),
	}

	data, err := WriteTable(table)
	if err != nil {
		t.Fatalf("Failed to write table: %v", err)
	}

	expectedData := []byte{
		0, 0, 0, 24, // table length
		4,                 // "name" length
		110, 97, 109, 101, // "name" in bytes
		83,         // "S" - Longstr type (ASCII code for 'S')
		0, 0, 0, 5, // "topic" length
		116, 111, 112, 105, 99, // "topic"
		3,            // "age" length
		97, 103, 101, // "age"
		73,          // "I" - int32 type (ASCII code for 'I')
		0, 0, 0, 30, // 30 as int32
	}

	assert.Equal(t, expectedData, data)
}
