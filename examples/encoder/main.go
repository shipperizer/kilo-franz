package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/shipperizer/kilo-franz/encoder"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	enc := encoder.NewJSONEncoder()

	payload, err := enc.Encode(person{Name: "Alice", Age: 30})
	if err != nil {
		panic(err)
	}

	var pretty bytes.Buffer
	if err := json.Indent(&pretty, payload, "", "  "); err != nil {
		panic(err)
	}

	fmt.Println(pretty.String())
}
