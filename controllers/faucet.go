package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/prometheus/common/log"
	rds_conn "github.com/wuban/faucet/cacheTools"
	"github.com/wuban/faucet/common"
	"github.com/wuban/faucet/common/rpc"
	"github.com/wuban/faucet/models"
	"strconv"
)

type FaucetController struct {
	beego.Controller
}

func (f *FaucetController) Transfer() {

	//ip := common.GetClientIP(f.Ctx)
	//t, c, err := common.CheckIpAddress(ip)
	//if err != nil {
	//	f.ResponseInfo(500, err.Error(), "")
	//	return
	//}


	var param models.Param
	data := f.Ctx.Input.RequestBody
	json.Unmarshal(data, &param)

	//判断传入的token是否存在
	if param.Token == "" {
		f.ResponseInfo(500, "", "please login github.")
		return
	}
	//获取github userInfo Url
 	userInfoUrl	:= beego.AppConfig.String("github::userInfoUrl")
	//获取用户信息
	userInfoStr := common.BasicGetHeader(userInfoUrl,param.Token)
	var userInfo models.UserInfo
	json.Unmarshal([]byte(userInfoStr), &userInfo)
	//判断获取github账户信息情况。错误消息是否存在
	if userInfo.Message !=""{
		f.ResponseInfo(500, "", userInfo.Message)
		return
	}

	githubId :=  strconv.FormatInt(userInfo.Id, 10)
	//验证地址是否正确
	if boo := common.CheckAddress(param.To); !boo {
		f.ResponseInfo(500, "Request address format exception, please re-enter.", "")
		return
	}

	//判断地址是否存在了
	//if boo := rds_conn.SR.IsKeyExists(param.To); boo {
	//	f.ResponseInfo(500, "", "Exceeding the daily limit.")
	//	return
	//}

	//限制每个github账号一天领取一次
	if boo := rds_conn.SR.IsKeyExists(githubId); boo {
		//判断申请次数 errCount 不为nil，则value存放的是githubId,可以放行 认为至少申请了一次
		applicationCount, errCount :=strconv.Atoi(rds_conn.SR.Get(githubId))
		if errCount == nil && applicationCount >=3 {
			//超过3次就提示超过次数限制
			f.ResponseInfo(500, "", "user " + userInfo.Login +  " Exceeding the daily limit.")
			return
		}

		//f.ResponseInfo(500, "", "user " + userInfo.Login +  " Exceeding the daily limit.")
		//return
	}


	chainID, _ := beego.AppConfig.Int("rpc::chainID" 	)
	myKey := beego.AppConfig.String("token::MYKEY")
	amount, _ := beego.AppConfig.Float("token::amount")
	tie, _ := beego.AppConfig.Int("token::time")
	p, err := rpc.StringToPrivateKey(myKey)

	if err != nil {
		log.Errorf("私钥解析失败")
		f.ResponseInfo(500, "Server exception", "")
		return
	}
	//value := big.NewFloat(amount)
	//v := common.TransferAmount(value)

	v := common.TransferAmountFloatToInt(amount)

	result, err := rpc.OfflineTransfer(int64(chainID), common.GetCount(), param.To, v, p)

	if err != nil {
		log.Errorf("签名失败", err)
		f.ResponseInfo(500, "Signature failure", "")
		return
	}

	hash, err := rpc.EthSendRawTransaction(result)

	if err != nil {
		f.ResponseInfo(500, "Radio failure", "false")
		return
	}

	//rds_conn.SR.SetKvAndExp(ip, c, t)



	//After the account address is successfully collected,
	//it is saved in redis to limit the collection frequency of users
	if hash != nil && len(hash.(string)) > 0 {

		//rds_conn.SR.SetKvAndExp(githubId, githubId, tie)
		//判断是否已经申请过
		if boo := rds_conn.SR.IsKeyExists(githubId); boo {
			//已经存在这个key,获取这个key申请的次数。 转化次数的时候可能会出去，之前申请的里面value存在的账户
			applicationCount, errCount :=strconv.Atoi(rds_conn.SR.Get(githubId))
			if errCount != nil {
				//获取不到次数，至少申请了一次，次数从2开始
				rds_conn.SR.SetKvAndExp(githubId, "2", tie)
			}else{
				applicationCount ++
				rds_conn.SR.SetKvAndExp(githubId, strconv.Itoa(applicationCount), tie)
			}
		}else{
			//第一次领取，插入次数1
			rds_conn.SR.SetKvAndExp(githubId, "1", tie)
		}


		f.ResponseInfo(200, "", hash)

		//存储最后几个account地址
		common.SetLastAccount(param.To)

		// Reminder when the account balance is insufficient
		go common.WarnBalance()
		return
	}

	f.ResponseInfo(500, "Radio failure", "false")
	return
}

func (f *FaucetController) GetLastAccounts() {
	arrIps, err := common.GetLastAccounts()
	if err != nil {
		f.ResponseInfo(500, "get errors ", "false")
		return
	}
	f.ResponseInfo(200, "", arrIps)
}


//获取github Token
func (f *FaucetController) GetToken() {
	//获取Code
	var oauth models.OauthParam
	data := f.Ctx.Input.RequestBody
	json.Unmarshal(data, &oauth)
	if oauth.Code == "" {
		f.ResponseInfo(500, "can't found code params", "false")
		return
	}
	//获取github的tokenUrl
	clientId := beego.AppConfig.String("github::clientId")
	ClientSecrets := beego.AppConfig.String("github::ClientSecrets")
	tokenUrl := beego.AppConfig.String("github::token_url") +
		"?code=" + oauth.Code + "&client_id=" + clientId + "&client_secret=" + ClientSecrets

	responseStr := common.Get(tokenUrl)
	f.ResponseInfo(200, "", responseStr)

}

//处理返回值信息
func (e *FaucetController) ResponseInfo(code int, err_msg string, result interface{}) {
	switch code {
	case 500:
		e.Data["json"] = map[string]interface{}{"error": "500", "err_msg": err_msg, "data": result}
	case 200:
		e.Data["json"] = map[string]interface{}{"error": "200", "err_msg": err_msg, "data": result}
	}
	e.ServeJSON()
}
