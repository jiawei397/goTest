package main

import (
	conf "hello/conf/loadConf.go"
)

var a = "G"

func main() {
	n()
	m()
	n()
	print(conf.HTTPPort)
}

func n() { print(a) }

func m() {
	// a := "O" //这是私有变量
	a = "O"
	print(a)
}
