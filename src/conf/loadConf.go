package conf

import (
	"fmt"
	"log"

	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	HTTPPort     int
	ReadTimeout  string
	WriteTimeout string
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/config.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	LoadServer()
}

func LoadServer() {
	sec, err := Cfg.GetSection("self")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	HTTPPort = sec.Key("port").MustInt(3000)
	ReadTimeout = sec.Key("flag").MustString("300")
	WriteTimeout = sec.Key("tag").MustString("300")

	fmt.Println(HTTPPort)
}
