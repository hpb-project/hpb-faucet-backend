package routers

import (
	"github.com/astaxie/beego"
	"github.com/wuban/faucet/controllers"
)

func init() {
	beego.Router("/api/faucet/v1/transfer", &controllers.FaucetController{}, "post:Transfer")
	beego.Router("/api/faucet/v1/getLastAccounts", &controllers.FaucetController{}, "get:GetLastAccounts")
	beego.Router("/api/faucet/v1/getToken", &controllers.FaucetController{}, "post:GetToken")
	//beego.Router("/api/faucet/v1/getUserInfo", &controllers.FaucetController{}, "post:GetUserInfo")
}
