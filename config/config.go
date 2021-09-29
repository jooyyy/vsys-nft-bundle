package config

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/virtualeconomy/go-v-sdk/vsys"
	"io/ioutil"
)

var Config = struct {
	Port  uint `default:"7005" json:"port"`
	V struct {
		Env      string `json:"env"`
		NetType  vsys.NetType `json:"net_type"`
		Node 	 string `json:"node"`
		Maker 	 string `json:"maker"`
		MakerSeed     string `json:"maker_seed"`
		GasPool  string `json:"gas_pool"`
		NFTContractId string `json:"nft_contract_id"`
	} `json:"v"`
}{}

func init() {
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &Config)
	if err != nil {
		panic(err)
	}
	if Config.V.Env == "M" {
		Config.V.NetType = vsys.Mainnet
		vsys.InitApi (
			Config.V.Node,
			vsys.Mainnet,
		)
	}else {
		Config.V.NetType = vsys.Testnet
		vsys.InitApi (
			"http://veldidina.vos.systems:9928",
			vsys.Testnet,
		)
	}
	acc := vsys.InitAccount(Config.V.NetType)
	acc.BuildFromSeed(Config.V.MakerSeed, 0)
	Config.V.Maker = acc.Address()
	if Config.V.NFTContractId == "" {
		logrus.Warning("NFT contract id is empty.")
	}
}