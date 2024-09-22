# amqp-table

`amqp-table` is a Go library that provides efficient encoding and decoding of AMQP 0-9-1 tables. This library is designed to work seamlessly with AMQP (Advanced Message Queuing Protocol) implementations in Go.

## Features

- Fast and efficient encoding of Go types to AMQP table format
- Decoding of AMQP tables to Go types
- Support for nested tables and arrays
- Compatibility with the AMQP 0-9-1 specification

## Installation

To install `amqp-table`, use `go get`:

```bash
go get github.com/junky/amqp-table
```

## Usage

To use `amqp-table` in your Go project, follow these steps:

1. Import the package:

```go
import amqptable "github.com/junky/amqp-table"
```

2. Use the provided functions to encode and decode AMQP tables:

```go
// Example of encoding a Go map to an AMQP table

table := map[string]any{
    "name": "John",
    "age":  30,
}

msgBody, err := amqptable.WriteTable(table)
if err != nil {
    log.Fatalf("Error encoding: %v", err)
}

// msgBody is a []byte that can be sent over AMQP

// Example of decoding an AMQP table to a Go map
var decoded map[string]any

err = amqptable.ReadTable(msgBody, &decoded)
if err != nil {
    log.Fatalf("Error decoding: %v", err)
}

fmt.Printf("Decoded map: %+v\n", decoded)
```

3. Marshal and Unmarshal structs to and from AMQP tables:

```go
// Example of encoding a Go struct to an AMQP table
type ExampleStruct struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

example := ExampleStruct{Name: "John", Age: 30}

msgBody, err := amqptable.Marshal(example)
if err != nil {
    log.Fatalf("Error encoding: %v", err)
}

// msgBody is a []byte that can be sent over AMQP

// Example of decoding an AMQP table to a Go struct
var decoded ExampleStruct

err = amqptable.Unmarshal(msgBody, &decoded)
if err != nil {
    log.Fatalf("Error decoding: %v", err)
}

fmt.Printf("Decoded struct: %+v\n", decoded)
```

## Documentation

For detailed documentation and API reference, please visit [GoDoc](https://pkg.go.dev/github.com/junky/amqp-table).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- This library is inspired by the AMQP 0-9-1 specification and aims to provide a robust implementation for Go developers.
