package server

type RespJsonObj struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Wallet struct {
	Seed string `json:"seed,omitempty"`
	PrivateKey string `json:"private_key,omitempty"`
	PublicKey string `json:"public_key,omitempty"`
	Address string `json:"address,omitempty"`
	NetType string `json:"net_type,omitempty"`
	WalletDetail
}

type WalletDetail struct {
	Address 	string `json:"address,omitempty"`
	Regular  	int64 `json:"regular,omitempty"`
	MintingAverage int64 `json:"mintingAverage,omitempty"`
	Available 	int64 `json:"available,omitempty"`
	Effective 	int64 `json:"effective,omitempty"`
	Height 		int64 `json:"height,omitempty"`
}

type Contract struct {
	Type 		int `json:"type"`
	TrxId  		string `json:"id"`
	Fee 		int64 `json:"fee"`
	FeeScale 	int64 `json:"feeScale"`
	Timestamp 	int64 `json:"timestamp"`
	Proofs 		[]struct{
		ProofType 	string `json:"proofType"`
		PublicKey 	string `json:"publicKey"`
		Address 	string `json:"address"`
		Signature 	string `json:"signature"`
	} `json:"proofs"`
	ContractId  string `json:"contractId"`
	Contract 	struct {
		LanguageCode 	string `json:"languageCode"`
		LanguageVersion int `json:"languageVersion"`
		Triggers 		[]string `json:"triggers"`
		Descriptors  	[]string `json:"descriptors"`
		StateVariables  []string `json:"stateVariables"`

	} `json:"contract"`
	InitData 	string `json:"initData"`
	Description string `json:"description"`
}

type NFT struct {
	Id 		string `json:"id"`
	TrxId  	string `json:"trx_id"`
	Index 	int `json:"index"`
	Description string `json:"description"`
	Attachment string `json:"attachment"`
}
