package server

type RespJsonObj struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Wallet struct {
	Seed string `json:"seed"`
	PrivateKey string `json:"private_key"`
	PublicKey string `json:"public_key"`
	Address string `json:"address"`
	NetType string `json:"net_type"`
}

type Contract struct {
	TrxId  		string `json:"trx_id"`
	ContractId  string `json:"contract_id"`
}
