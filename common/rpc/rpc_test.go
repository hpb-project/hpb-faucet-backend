package rpc

import (
	"fmt"
	"math/big"
	"testing"
)

func Test1(t *testing.T) {
	p, _ := StringToPrivateKey("0xBC6E4646C083D15E918154E93B1315EB1B16A474C95B2047638E20B95DC5EDD8")
	nonce := uint64(1)
	value := big.NewInt(2)
	result, err := OfflineTransfer(269, nonce, "0x26D6d53F9eBfe7eC3151Dc9C418fc1685cA1423A", value, p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

func Test2(t *testing.T) {
	count := 1
	count++
	fmt.Println(count)
}

func TestCreateKey(t *testing.T) {

	private, address := CreateKey()
	fmt.Println("private", private)
	fmt.Println("address", address)

}

func TestGetTransactionReceipt(t *testing.T) {
	result, err := GetTransactionReceipt("0x6c5a429dd0805ab7290cc62a04f3e9cfa8ea249d9aab18bd7d706a90137554a0")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
}

func TestEthEstimateGas(t *testing.T) {
	param := Param{From: "0xaccf7aacd00bb765120f10685e58d916d3ec3057", To: "0x020D5741be5Af82aA9332C2d3B9cFCA3133035f5", Value: "0x64"}
	result, err := EthEstimateGas(param)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
