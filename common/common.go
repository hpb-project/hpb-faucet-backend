package common

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/common/log"
	rds_conn "github.com/wuban/faucet/cacheTools"
	"github.com/wuban/faucet/common/rpc"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var count uint64

func SetCount() error {
	address := beego.AppConfig.String("token::address")
	n, err := rpc.EthGetTransactionCount(address, "latest")
	if err != nil {
		return err
	}
	count = n.(uint64)
	log.Info("账户nonce为：", count)
	return nil
}

func GetCount() uint64 {
	defer func() {
		count++
	}()
	return count
}

func TransferAmountFloatToInt(amount float64) *big.Int {
	bigVal := new(big.Float)
	bigVal.SetFloat64(amount)
	coin := new(big.Float)
	coin.SetInt(big.NewInt(1000000000000000000))
	bigVal.Mul(bigVal, coin)
	result := new(big.Int)
	f, _ := bigVal.Uint64()
	result.SetUint64(f)
	return result
}

func TransferAmount(amount *big.Int) *big.Int {
	n := new(big.Int)
	na, _ := n.SetString("1000000000000000000", 0)
	return na.Mul(na, amount)
}

func CheckAddress(address string) bool {
	if !strings.HasPrefix(address, "0x") {
		return false
	}
	return common.IsHexAddress(address)
}

func CheckIpAddress(ip string) (int, string, error) {
	var count string //计数
	var time int     //key 的过期时间
	c := "1"
	time, _ = beego.AppConfig.Int("token::time")
	ct := beego.AppConfig.String("token::count")

	if boo := rds_conn.SR.IsKeyExists(ip); boo {
		count = rds_conn.SR.Get(ip)
		time = rds_conn.SR.GetExp(ip)
		if count == ct {
			return 0, "", fmt.Errorf("each IP has only One hundred chance every 24 hours")
		}
		num, _ := strconv.Atoi(count)
		c = strconv.Itoa(num + 1)
	}
	return time, c, nil
}

func GetClientIP(ctx *context.Context) string {
	r := ctx.Request
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}

func SendMail(mailTo string, subject, body string, wg *sync.WaitGroup) error {
	return nil
	//defer wg.Done()
	//mailConn := map[string]string{
	//	"username": "tech@metarace.org",
	//	"authCode": "CGrB9KmsJ0Y7",
	//	"host":     "smtp.zoho.com",
	//	"port":     "587",
	//}
	//port, _ := strconv.Atoi(mailConn["port"])
	//m := gomail.NewMessage()
	//m.SetHeader("From", mime.QEncoding.Encode("UTF-8", "Support")+"<"+mailConn["username"]+">")
	//m.SetHeader("To", mailTo)
	//m.SetHeader("Subject", subject)
	//m.SetBody("text/html", body)
	//d := gomail.NewDialer(mailConn["host"], port, mailConn["username"], mailConn["authCode"])
	//err := d.DialAndSend(m)
	//if err != nil {
	//	log.Fatalln("To:", mailTo, "##", "Send Email Failed!Err:", err)
	//} else {
	//	log.Info("To:", mailTo, "##", "Send Email Successfully!")
	//}
	//return err
}

func WarnBalance() {

	time.Sleep(5 * time.Second)

	// Get account balance
	address := beego.AppConfig.String("token::address")
	ok, outAddressBalance, err := rpc.GetBalance(address)
	if err != nil || !ok {
		log.Errorf("Get balance Error!", err)

	}
	balanceByte, err := hex.DecodeString(outAddressBalance.(string)[2:])
	if err != nil {
		log.Errorf("DecodeString addressBalance Error!", err)
	}
	balanceInt := new(big.Int).SetBytes(balanceByte)

	// Get warning account value
	warningValue, _ := beego.AppConfig.Int64("token::warningValue")
	warningValueInt := TransferAmount(big.NewInt(warningValue))

	// Cmp=1 : balanceInt > warningValueInt
	if balanceInt.Cmp(warningValueInt) != 1 {
		var wg sync.WaitGroup
		mailTo := []string{
			"xueqian1991@163.com",
			"m13840625723@163.com",
		}
		subject := "水龙头账户余额不足提醒"
		body := "<h2>水龙头账户余额不足1000CMP，请相关工作人员进行处理!</h2>"
		for _, mail := range mailTo {
			wg.Add(1)
			go SendMail(mail, subject, body, &wg)
		}
		wg.Wait()
	}
}

//设置本地存储最后申请的accont地址
func SetLastAccount(account string) {
	var lastNum int //本地最多存储ip数量
	var key string  //存放ip的key
	key = "Last_Deposits"
	lastNum, _ = beego.AppConfig.Int("cache::lastNums")
	curNum, _ := rds_conn.SR.LLen(key)

	if curNum >= int64(lastNum) {
		rds_conn.SR.RPop(key)
	}
	err := rds_conn.SR.LpushByte(key, []byte(account))
	if err != nil {
		fmt.Println("Last Deposits save ips error ")
	}
}

//获取最后申请的account地址
func GetLastAccounts() ([]string, error) {
	key := "Last_Deposits"
	values, err := rds_conn.SR.LRange(key)
	if err != nil {
		fmt.Println("Last Deposits get values error", err)
		return nil, err
	}
	fmt.Println(values)
	return values, nil
}

//可以设置header得get 请求
func BasicGetHeader(url string, token string) (response string) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "token " + token)
	req.Header.Add("User-Agent","HPB fault")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	response = string(result)
	return

}

//发送GET请求
//url:请求地址
//response:请求返回的内容
func Get(url string) (response string) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, error := client.Get(url)
	defer resp.Body.Close()
	if error != nil {
		panic(error)
	}
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	response = result.String()
	return
}

//发送POST请求
//url:请求地址		data:POST请求提交的数据		contentType:请求体格式，如：application/json
//content:请求返回的内容
func Post(url string, data interface{}, contentType string) (content string) {
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", contentType)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	return
}
