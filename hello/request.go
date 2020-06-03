package main

import (
	"fmt"

	"github.com/winterssy/sreq"
)

func main() {

	resp, err := sreq.Get("https://www.google.com/").Raw()
	println(resp, err)
	if err != nil {
		panic(err)
	}
	// println(resp.Text())
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header["Content-Type"])
}
