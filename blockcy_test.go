//Tests for the BlockCypher Go SDK. Test functions
//try to mirror file names where possible.
package blockcy

import (
	"log"
	"os"
	"testing"
)

var keys1, keys2 AddrKeychain
var txhash1, txhash2 string
var bcy API

func TestMain(m *testing.M) {
	//Set Coin/Chain to BlockCypher testnet
	bcy.Coin = "bcy"
	bcy.Chain = "test"
	//Set Token
	bcy.Token = "test-token"
	//Create/fund the test addresses
	var err error
	keys1, err = bcy.GenAddrKeychain()
	keys2, err = bcy.GenAddrKeychain()
	if err != nil {
		log.Fatal("Error generating test addresses: ", err)
	}
	txhash1, err = bcy.Faucet(keys1, 1e5)
	txhash2, err = bcy.Faucet(keys2, 2e5)
	if err != nil {
		log.Fatal("Error funding test addresses: ", err)
	}
	os.Exit(m.Run())
}

//TestsGetTXConf runs first, to test
//Confidence factor
func TestGetTXConf(t *testing.T) {
	conf, err := bcy.GetTXConf(txhash2)
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	t.Logf("%+v\n", conf)
	return
}

func TestBlockchain(t *testing.T) {
	ch, err := bcy.GetChain()
	if err != nil {
		t.Error("GetChain error encountered: ", err)
	}
	t.Logf("%+v\n", ch)
	_, err = bcy.GetBlock(187621, "")
	if err != nil {
		t.Error("GetBlock via height error encountered: ", err)
	}
	bl, err := bcy.GetBlock(0, "0000ffeb0031885f2292475eac7f9c6f7bf5057e3b0017a09cd1994e71b431a4")
	if err != nil {
		t.Error("GetBlock via hash error encountered: ", err)
	}
	t.Logf("%+v\n", bl)
	_, err = bcy.GetBlock(187621, "0000ffeb0031885f2292475eac7f9c6f7bf5057e3b0017a09cd1994e71b431a4")
	if err == nil {
		t.Error("Expected error when querying both height and hash in GetBlock, did not receive one")
	}
	err = nil
	bl, err = bcy.GetBlockPage(0, "0000cb69e3c85ec1a4a17d8a66634c1cf136acc9dca9a5a71664a593f92bc46e", 0, 1)
	if err != nil {
		t.Error("GetBlockPage error encountered: ", err)
	}
	t.Logf("%+v\n", bl)
	bl2, err := bcy.GetBlockNextTXs(bl)
	if err != nil {
		t.Error("GetBlockNextTXs error encountered: ", err)
	}
	t.Logf("%+v\n", bl2)
	return
}

func TestAddress(t *testing.T) {
	addr, err := bcy.GetAddrBal(keys1.Address)
	if err != nil {
		t.Error("GetAddrBal error encountered: ", err)
	}
	t.Logf("%+v\n", addr)
	addr, err = bcy.GetAddr(keys1.Address)
	if err != nil {
		t.Error("GetAddr error encountered: ", err)
	}
	t.Logf("%+v\n", addr)
	addr, err = bcy.GetAddrFull(keys2.Address)
	if err != nil {
		t.Error("GetAddrFull error encountered: ", err)
	}
	t.Logf("%+v\n", addr)
	return
}

func TestGenAddrMultisig(t *testing.T) {
	pubkeys := []string{
		"02c716d071a76cbf0d29c29cacfec76e0ef8116b37389fb7a3e76d6d32cf59f4d3",
		"033ef4d5165637d99b673bcdbb7ead359cee6afd7aaf78d3da9d2392ee4102c8ea",
		"022b8934cc41e76cb4286b9f3ed57e2d27798395b04dd23711981a77dc216df8ca",
	}
	response, err := bcy.GenAddrMultisig(AddrKeychain{PubKeys: pubkeys, ScriptType: "multisig-2-of-3"})
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	if response.Address != "De2gwq9GvNgvKgHCYRMKnPqss3pzWGSHiH" {
		t.Error("Response does not match expected address")
	}
	t.Logf("%+v\n", response)
}

func TestTX(t *testing.T) {
	txs, err := bcy.GetUnTX()
	if err != nil {
		t.Error("GetUnTX error encountered: ", err)
	}
	t.Logf("%+v\n", txs)
	tx, err := bcy.GetTX(txhash1)
	if err != nil {
		t.Error("GetTX error encountered: ", err)
	}
	t.Logf("%+v\n", tx)
	//Create and Send TXSkeleton
	temp := TempNewTX(keys2.Address, keys1.Address, 25000, false)
	skel, err := bcy.NewTX(temp)
	if err != nil {
		t.Error("NewTX error encountered: ", err)
	}
	t.Logf("%+v\n", skel)
	/*Will uncomment once client-side signing integrated into SDK
	skel, err = bcy.SendTX(skel)
	if err != nil {
		t.Error("SendTX error encountered: ", err)
	}
	t.Logf("%+v\n", skel)*/
	return
}

func TestHooks(t *testing.T) {
	hook, err := bcy.PostHook(Hook{Event: "new-block", URL: "https://my.domain.com/api/callbacks/doublespend?secret=justbetweenus"})
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	t.Logf("%+v\n", hook)
	if err = bcy.DeleteHook(hook); err != nil {
		t.Error("Error encountered: ", err)
	}
	hooks, err := bcy.ListHooks()
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	//Should be empty
	t.Logf("%+v\n", hooks)
	return
}

func TestPayment(t *testing.T) {
	pay, err := bcy.PostPayment(PaymentFwd{Destination: keys1.Address})
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	t.Logf("%+v\n", pay)
	if err = bcy.DeletePayment(pay); err != nil {
		t.Error("Error encountered: ", err)
	}
	pays, err := bcy.ListPayments()
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	//Should be empty
	t.Logf("%+v\n", pays)
	return
}

func TestMicro(t *testing.T) {
	mic := MicroTX{Priv: keys2.Private, ToAddr: keys1.Address, Value: 25000}
	result, err := bcy.SendMicro(mic)
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	t.Logf("%+v\n", result)
	txmic, err := bcy.GetTX(result.Hash)
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	t.Logf("%+v\n", txmic)
	return
}
