package main

import (
	_ "aiyun_local_srv/boot"
	_ "aiyun_local_srv/router"
	"github.com/gogf/gf/frame/g"
	_ "github.com/lib/pq"
)

func main() {
	g.Server().Run()
}
