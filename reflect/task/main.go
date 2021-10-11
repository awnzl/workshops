package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

type User struct {
	Name string
	Age  int64
}

type City struct {
	Name       string
	Population int64
	GDP        int64
	Mayor      string
}

func main() {
	var u User = User{"bob", 10}

	res, err := JSONEncode(u)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))

	c := City{"sf", 5000000, 567896, "mr jones"}
	res, err = JSONEncode(c)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))
}

func JSONEncode(v interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	buf.WriteString("{")

	// TODO: check if v is a struct else return error
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return nil, fmt.Errorf("v is not a struct")
	}

	wrapValue := func(name, val string, last bool) {
		buf.WriteString(`"` + name + `": `)
		buf.WriteString(`"` + val + `"`)

		if !last {
			buf.WriteString(`, `)
		}
	}

	// TODO: iterate over v`s reflect value using NumField()
	// use type switch to create string result of "{field}" + ": " + "{value}"
	// start with just 2 types - reflect.String and reflect.Int64
	value := reflect.ValueOf(v)
	indir := reflect.Indirect(value)

	for i := 0; i < value.NumField(); i++ {
		f := value.Field(i)

		last := i + 1 == value.NumField()
		name := indir.Type().Field(i).Name

		switch f.Interface().(type) {
		case int64:
			wrapValue(name, strconv.Itoa(int(f.Int())), last)
		case string:
			wrapValue(name, f.String(), last)
		}
	}

	buf.WriteString("}")
	return buf.Bytes(), nil
}
