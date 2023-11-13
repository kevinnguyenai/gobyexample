// GO provides a built -in JSON coding decoding (serialized derivative) support,
// Including the conversion between built -in and custom types and JSON data.

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Below we will use these two structures to demonstrate the coding and decoding of custom types.
type response1 struct {
	Page   int
	Fruits []string
}

// Only the field that can be exported will be encoded/decoded by JSON.The fields that must start with uppercase letters are exported.
type response2 struct {
	Page   int      `json:"page"`
	Fruits []string `json:"fruits"`
}

func main() {

// First of all, let's take a look at the encoding process of the basic data type to the JSON string.
// This is an example of some atomic values.
	bolB, _ := json.Marshal(true)
	fmt.Println(string(bolB))

	intB, _ := json.Marshal(1)
	fmt.Println(string(intB))

	fltB, _ := json.Marshal(2.34)
	fmt.Println(string(fltB))

	strB, _ := json.Marshal("gopher")
	fmt.Println(string(strB))

// This is an example of some slices and maps into JSON array and objects.
	slcD := []string{"apple", "peach", "pear"}
	slcB, _ := json.Marshal(slcD)
	fmt.Println(string(slcB))

	mapD := map[string]int{"apple": 5, "lettuce": 7}
	mapB, _ := json.Marshal(mapD)
	fmt.Println(string(mapB))

// The json package can automatically encode your custom type.
// Code output only contains exportable fields, and uses field names as key names of JSON data by default.
	res1D := &response1{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res1B, _ := json.Marshal(res1D)
	fmt.Println(string(res1B))

// You can declare the key name of the JSON data that defines the coding of the structural field.
// The definition of the above `response2` is an example of this label.
	res2D := &response2{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res2B, _ := json.Marshal(res2D)
	fmt.Println(string(res2B))

// Now let's take a look at the process of decoding JSON data to go value.
// This is an example of an ordinary data structure.
	byt := []byte(`{"num":6.13,"strs":["a","b"]}`)

// We need to provide a variable that can store decoding data.
// The `Map [String] interface {}` is a Map with a key that is string and a value of any value.
	var dat map[string]interface{}

// This is the actual decoding and related errors.
	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)

// In order to use the value in MAP, we need to transform them appropriately.
// For example, here we convert the value of `NUM` into types of` Float64.
	num := dat["num"].(float64)
	fmt.Println(num)

// The value of accessing nested requires a series of conversion.
	strs := dat["strs"].([]interface{})
	str1 := strs[0].(string)
	fmt.Println(str1)

// We can also decod in JSON into custom data types.
// The advantage of doing this is that it can add additional type of security to our procedures.
// Do not need type assertions when accessing the data after accessing the decoding.
	str := `{"page": 1, "fruits": ["apple", "peach"]}`
	res := response2{}
	json.Unmarshal([]byte(str), &res)
	fmt.Println(res)
	fmt.Println(res.Fruits[0])

// On the standard output of the above example,
// We always use Byte and String as an intermediary between data and JSON.
// Of course, we can also pass the JSON coding directly like `os.stdout`
// `os.writer` and even HTTP response.
	enc := json.NewEncoder(os.Stdout)
	d := map[string]int{"apple": 5, "lettuce": 7}
	enc.Encode(d)
}
