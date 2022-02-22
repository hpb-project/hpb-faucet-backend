package main

import (
	"github.com/astaxie/beego"
	"github.com/prometheus/common/log"
	_ "github.com/wuban/faucet/cacheTools"
	"github.com/wuban/faucet/common"
	_ "github.com/wuban/faucet/routers"
)

func main() {
	if err := common.SetCount(); err != nil {
		log.Errorf("账户nonce初始化失败")
		return
	}
	beego.Run()
}
