package amqptable

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestUnmarshalFromTable(t *testing.T) {
	table := map[string]any{
		"name": "John",
		"age":  30,
	}

	var testStruct TestStruct
	err := unmarshalFromTable(table, &testStruct)
	assert.NoError(t, err)
	assert.Equal(t, "John", testStruct.Name)
	assert.Equal(t, 30, testStruct.Age)
}

func TestMarshalToTable(t *testing.T) {
	var testStruct TestStruct
	testStruct.Name = "John"
	testStruct.Age = 30
	table, err := marshalToTable(testStruct)
	assert.NoError(t, err)
	assert.Equal(t, "John", table["name"])
	assert.Equal(t, 30, table["age"])
}
