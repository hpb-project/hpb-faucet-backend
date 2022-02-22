package rpc

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

//普通返回值
type Reply struct {
	JsonRpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
	Result  interface{} `json:"result"`
}

type Reply1 struct {
	JsonRpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  string `json:"result"`
}

//签名信息返回值
type Reply2 struct {
	JsonRpc string   `json:"jsonrpc"`
	Id      int      `json:"id"`
	Result  *RawInfo `json:"result"`
}

//签名信息
type RawInfo struct {
	Raw string `json:"raw"`
}

//签名信息返回值
type ReplyReceipt struct {
	JsonRpc string   `json:"jsonrpc"`
	Id      int      `json:"id"`
	Result  *Receipt `json:"result"`
}

//交易收据
type Receipt struct {
	BlockHash         string `json:"blockHash"`
	BlockNumber       string `json:"blockNumber"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	From              string `json:"from"`
	GasUsed           string `json:"gasUsed"`
	Logs              string `json:"logs"`
	LogsBloom         string `json:"logs_bloom"`
	Status            string `json:"status"`
	To                string `json:"to"`
	TransactionHash   string `json:"transactionHash"`
	TransactionIndex  string `json:"transactionIndex"`
}
type Param struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Data  string `json:"data"`
	Value string `json:"value"`
}

/**
交易收据接口
*/
func GetTransactionReceipt(param1 interface{}) (interface{}, error) {
	url := beego.AppConfig.String("rpc::url")
	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0", "method": "eth_getTransactionReceipt", "id": 1, "params": []interface{}{param1}})
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result ReplyReceipt
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Result, nil
}

func EthGasPrice() (interface{}, error) {
	url := beego.AppConfig.String("rpc::url")
	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0", "method": "eth_gasPrice", "id": 67, "params": []interface{}{}})
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result Reply
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Result, nil
}

func EthEstimateGas(param1 interface{}) (interface{}, error) {
	url := beego.AppConfig.String("rpc::url")
	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0", "method": "eth_estimateGas", "id": 1, "params": []interface{}{param1}})
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result Reply
	err = json.Unmarshal(body, &result)
	if err != nil || result.Result == nil {
		logs.Error(err)
		return nil, fmt.Errorf("EthEstimateGas Gas estimation method is abnormal")
	}
	count := Hex2Dec(result.Result.(string))
	return count, nil
}

func EthGetTransactionCount(param, PARAM1 string) (interface{}, error) {
	url := beego.AppConfig.String("rpc::url")
	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0", "method": "eth_getTransactionCount", "id": 1, "params": []interface{}{param, PARAM1}})
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result Reply
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	count := Hex2Dec(result.Result.(string))
	return count, nil
}

func EthSendRawTransaction(param1 interface{}) (interface{}, error) {
	url := beego.AppConfig.String("rpc::url")
	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0", "method": "eth_sendRawTransaction", "id": 1, "params": []interface{}{param1}})
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result Reply
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result.Result, nil
}

func StringToPrivateKey(privateKeyStr string) (*ecdsa.PrivateKey, error) {
	privateKeyByte, err := hexutil.Decode(privateKeyStr)
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.ToECDSA(privateKeyByte)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func SignTransaction(chainID int64, tx *types.Transaction, privateKey *ecdsa.PrivateKey) (string, error) {
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
	if err != nil {
		return "", nil
	}

	b, err := rlp.EncodeToBytes(signTx)
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(b), nil
}

func OfflineTransfer(chainID int64, nonce uint64, toAddress string, value *big.Int, privk *ecdsa.PrivateKey) (string, error) {
	address := beego.AppConfig.String("token::address")
	param := Param{From: address, To: toAddress, Value: "0x56bc75e2d63100000"}
	gasUsed, err := EthEstimateGas(param)
	if err != nil {
		return "", err
	}
	gas := gasUsed.(uint64)
	fmt.Println("gas", gas)
	gasPrice, err := EthGasPrice()
	if err != nil {
		return "", err
	}
	gp := Hex2DecBig(gasPrice.(string))
	tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), value, gas, gp, nil)
	return SignTransaction(chainID, tx, privk)
}

func Hex2DecBig(val string) *big.Int {
	n := new(big.Int)
	na, _ := n.SetString(val, 0)
	return na
}

func CreateKey() (privs, addrs string) {
	//创建私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Println(err)
	}
	/*	//可通过此代码导入私钥
		privateKey,err=crypto.HexToECDSA("93d5d04256882aaad507ff09f510969f347758109793448aa79e1b4dbe5f6efa")
		if err != nil {
			log.Fatal(err)
		}
	*/
	privateKeyBytes := crypto.FromECDSA(privateKey)
	priv := hexutil.Encode(privateKeyBytes)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	return priv, address
}

func GetBalance(param string) (bool, interface{}, error) {
	//url := "https://mainnet.infura.io/v3/e15791a2ccc34c019f16d6aeeea732cf"
	url := "http://18.169.173.49:8580"
	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0", "method": "eth_getBalance", "id": 1, "params": []interface{}{param, "latest"}})
	if err != nil {
		return false, nil, err
	}
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return false, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, nil, err
	}
	var result Reply1
	err = json.Unmarshal(body, &result)
	if err != nil {
		return false, nil, err
	}
	return true, result.Result, nil
}

func Hex2Dec(val string) uint64 {
	val = val[2:]
	n, _ := strconv.ParseUint(val, 16, 32)
	return n
}
