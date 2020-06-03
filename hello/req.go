package main

import (
	"fmt"

	"github.com/kirinlabs/HttpRequest"
)

func main() {
	req := HttpRequest.NewRequest()
	res, err := req.Get("https://api.github.com/events")
	fmt.Println(res, err)
}
