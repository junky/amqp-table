package amqptable

import (
	"fmt"
	"reflect"
)

type InvalidUnmarshalError struct {
	Type reflect.Type
}

type InvalidMarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Pointer {
		return "json: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "json: Unmarshal(nil " + e.Type.String() + ")"
}

func (e *InvalidMarshalError) Error() string {
	if e.Type == nil {
		return "json: Marshal(nil)"
	}

	return "json: Marshal(nil " + e.Type.String() + ")"
}

func Unmarshal(data []byte, v any) error {
	table, err := readTable(data)
	if err != nil {
		return err
	}
	return unmarshalFromTable(table, v)
}

func Marshal(v any) ([]byte, error) {
	table, err := marshalToTable(v)
	if err != nil {
		return nil, err
	}

	return writeTable(table)
}

func unmarshalFromTable(table map[string]any, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	elem := rv.Elem()
	numFields := elem.NumField()
	for i := 0; i < numFields; i++ {
		field := elem.Field(i)
		fieldType := elem.Type().Field(i)
		jsonTag := fieldType.Tag.Get("json")

		if jsonTag == "" {
			continue // Skip fields without json tags
		}

		if value, ok := table[jsonTag]; ok {
			if err := setField(field, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func marshalToTable(v any) (map[string]any, error) {
	table := make(map[string]any)

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Marshal requires a struct, got %T", v)
	}

	numFields := rv.NumField()
	for i := 0; i < numFields; i++ {
		field := rv.Field(i)
		fieldType := rv.Type().Field(i)
		jsonTag := fieldType.Tag.Get("json")

		if jsonTag == "" {
			continue // Skip fields without json tags
		}

		table[jsonTag] = field.Interface()
	}

	return table, nil
}

func setField(field reflect.Value, value interface{}) error {
	switch field.Kind() {
	case reflect.Int8:
		if v, ok := value.(int8); ok {
			field.SetInt(int64(v))
		} else {
			return fmt.Errorf("expected int8, got %T", value)
		}
	case reflect.Int16:
		if v, ok := value.(int16); ok {
			field.SetInt(int64(v))
		} else {
			return fmt.Errorf("expected int16, got %T", value)
		}
	case reflect.Int32:
		if v, ok := value.(int32); ok {
			field.SetInt(int64(v))
		} else {
			return fmt.Errorf("expected int32, got %T", value)
		}
	case reflect.Int64:
		if v, ok := value.(int64); ok {
			field.SetInt(v)
		} else {
			return fmt.Errorf("expected int64, got %T", value)
		}
	case reflect.Int:
		if v, ok := value.(int); ok {
			field.SetInt(int64(v))
		} else {
			return fmt.Errorf("expected int, got %T", value)
		}
	case reflect.Uint8:
		if v, ok := value.(uint8); ok {
			field.SetUint(uint64(v))
		} else {
			return fmt.Errorf("expected int8, got %T", value)
		}
	case reflect.Uint16:
		if v, ok := value.(uint16); ok {
			field.SetUint(uint64(v))
		} else {
			return fmt.Errorf("expected int16, got %T", value)
		}
	case reflect.Uint32:
		if v, ok := value.(uint32); ok {
			field.SetUint(uint64(v))
		} else {
			return fmt.Errorf("expected int32, got %T", value)
		}
	case reflect.Uint64:
		if v, ok := value.(uint64); ok {
			field.SetUint(v)
		} else {
			return fmt.Errorf("expected uint64, got %T", value)
		}
	case reflect.Uint:
		if v, ok := value.(uint); ok {
			field.SetUint(uint64(v))
		} else {
			return fmt.Errorf("expected uint, got %T", value)
		}
	case reflect.String:
		if v, ok := value.(string); ok {
			field.SetString(v)
		} else {
			return fmt.Errorf("expected string, got %T", value)
		}
	case reflect.Bool:
		if v, ok := value.(bool); ok {
			field.SetBool(v)
		} else {
			return fmt.Errorf("expected bool, got %T", value)
		}
	case reflect.Float32:
		if v, ok := value.(float32); ok {
			field.SetFloat(float64(v))
		} else {
			return fmt.Errorf("expected float32, got %T", value)
		}
	case reflect.Float64:
		if v, ok := value.(float64); ok {
			field.SetFloat(v)
		} else {
			return fmt.Errorf("expected float64, got %T", value)
		}
	case reflect.Slice:
		if v, ok := value.([]byte); ok {
			field.SetBytes(v)
		} else {
			return fmt.Errorf("expected []byte for slice, got %T", value)
		}
	}
	return nil
}
