package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/virtualeconomy/go-v-sdk/vsys"
	"strings"
	"time"
	"vsys-nft-bundle/config"
)

// RegisterNFTContract
// @Summary Register NFT contract
// @Description Before your register the contract, please make sure you have a maker wallet, and your wallet have 100 vsys at least.
// @Tags Contract
// @Accept  json
// @Produce json
// @Param description formData string false "NFT contract description"
// @Success 200 {object} Contract ""
// @Failure 400 {object} RespJsonObj ""
// @Router /api/contract/registerNFT [post]
func (s *Service) RegisterNFTContract(c *gin.Context) {
	desc, _ := c.GetPostForm("description")
	acc := vsys.InitAccount(config.Config.V.NetType)
	acc.BuildFromSeed(config.Config.V.MakerSeed, 0)
	paramData := vsys.DataEncoder{}
	paramData.EncodeArgAmount(0)
	trx := vsys.NewRegisterTransaction(CONTRACT_NFT, vsys.Base58Encode(paramData.Result()), desc)
	trx.SenderPublicKey = acc.PublicKey()
	trx.Signature = acc.SignData(trx.BuildTxData())
	data, err := vsys.GetVsysApi().Post(vsys.ApiContractBroadcastRegister, trx)
	if err != nil {
		respErr(c, err.Error())
		return
	}
	var contract Contract
	err = json.Unmarshal(data, &contract)
	_, err = retryToGetDetailWithTx(contract.TrxId)
	if err != nil {
		respErr(c, err.Error())
		return
	}
	respJson(c, contract)
}

// MintNFT
// @Summary Mint a NFT
// @Description This operation will coast 0.3vsys, please make sure your gas pool wallet has balance
// @Tags Contract
// @Accept  json
// @Produce json
// @Param description formData string false "NFT description"
// @Param attachment formData string false "Transaction attachment"
// @Success 200 {object} NFT ""
// @Failure 400 {object} RespJsonObj ""
// @Router /api/contract/mintNFT [post]
func (s *Service) MintNFT(c *gin.Context) {
	description, _ := c.GetPostForm("description")
	attachment, _ := c.GetPostForm("attachment")
	acc := vsys.InitAccount(config.Config.V.NetType)
	acc.BuildFromSeed(config.Config.V.MakerSeed, 0)
	data := vsys.DataEncoder{}
	data.EncodeArgAmount(1)
	data.Encode(description, vsys.DeTypeShortText)
	trx := acc.BuildExecuteContract(
		config.Config.V.NFTContractId,
		NFTFuncIdxIssue,
		data.Result(),
		attachment,
	)
	resp, err := vsys.SendExecuteContractTx(trx)
	if err != nil {
		respErr(c, err.Error())
		return
	}
	_, err = retryToGetDetailWithTx(resp.Id)
	if err != nil {
		respErr(c, err.Error())
		return
	}
	latestIndex, err := GetLatestNFTIndex()
	if err != nil {
		respErr(c, err.Error())
		return
	}
	// latestIndex is the next unmined NFT index, so we use latestIndex-1 represent for NFT minted just now.
	token, err := GetNFTByIndex(latestIndex-1)
	if err != nil {
		respErr(c, err.Error())
		return
	}

	respJson(c, NFT{
		Id:          token,
		TrxId:       resp.Id,
		Index:       latestIndex - 1,
		Description: description,
		Attachment:  attachment,
	})
}

// TransferNFT
// @Summary Transfer NFT
// @Description This operation will coast 0.3vsys, please make sure your gas pool wallet has balance
// @Tags Contract
// @Accept  json
// @Produce json
// @Param owner_seed formData string true "Seed for NFT owner"
// @Param nft_index formData int true "NFT index"
// @Param recipient formData string true "NFT receiver"
// @Param attachment formData string false "Transaction attachment"
// @Success 200 {object} string "transfer NFT transaction id"
// @Failure 400 {object} RespJsonObj ""
// @Router /api/contract/transferNFT [post]
func (s *Service) TransferNFT(c *gin.Context) {
	ownerSeed, ok := c.GetPostForm("owner_seed")
	if !ok {
		respErr(c, "need owner_seed")
		return
	}
	nftIndex, ok := c.GetPostForm("nft_index")
	if !ok {
		respErr(c, "need nft index")
		return
	}
	recipient, ok := c.GetPostForm("recipient")
	if !ok {
		respErr(c, "need recipient")
		return
	}
	attachment, _ := c.GetPostForm("attachment")
	if vsys.IsValidateAddress(recipient, config.Config.V.NetType) == false {
		respErr(c, "invalid recipient.")
		return
	}
	acc := vsys.InitAccount(config.Config.V.NetType)
	acc.BuildFromSeed(ownerSeed, 0)
	data := vsys.DataEncoder{}
	data.EncodeArgAmount(2)
	data.Encode(recipient, vsys.DeTypeAddress)
	data.Encode(nftIndex, vsys.DeTypeInt32)
	resp, err := vsys.SendExecuteContractTx(acc.BuildExecuteContract(
		config.Config.V.NFTContractId,
		NFTFuncIdxSend,
		data.Result(),
		attachment,
	))
	if err != nil {
		respErr(c, err.Error())
		return
	}
	_, err = retryToGetDetailWithTx(resp.Id)
	if err != nil {
		respErr(c, err.Error())
		return
	}

	respJson(c, resp.Id)
}

// SendNFT
// @Summary Send NFT to recipient wallet
// @Description NFT will be sent from maker wallet to recipient, this operation will coast 0.3vsys
// @Tags Contract
// @Accept  json
// @Produce json
// @Param nft_index formData int true "NFT index"
// @Param recipient formData string true "NFT receiver"
// @Param attachment formData string false "Transaction attachment"
// @Success 200 {object} string "transfer NFT transaction id"
// @Failure 400 {object} RespJsonObj ""
// @Router /api/contract/sendNFT [post]
func (s *Service) SendNFT(c *gin.Context) {
	nftIndex, ok := c.GetPostForm("nft_index")
	if !ok {
		respErr(c, "need nft index")
		return
	}
	recipient, ok := c.GetPostForm("recipient")
	if !ok {
		respErr(c, "need recipient")
		return
	}
	attachment, _ := c.GetPostForm("attachment")
	if vsys.IsValidateAddress(recipient, config.Config.V.NetType) == false {
		respErr(c, "invalid recipient.")
		return
	}
	acc := vsys.InitAccount(config.Config.V.NetType)
	acc.BuildFromSeed(config.Config.V.MakerSeed, 0)
	data := vsys.DataEncoder{}
	data.EncodeArgAmount(2)
	data.Encode(recipient, vsys.DeTypeAddress)
	data.Encode(nftIndex, vsys.DeTypeInt32)
	resp, err := vsys.SendExecuteContractTx(acc.BuildExecuteContract(
		config.Config.V.NFTContractId,
		NFTFuncIdxSend,
		data.Result(),
		attachment,
	))
	if err != nil {
		respErr(c, err.Error())
		return
	}
	_, err = retryToGetDetailWithTx(resp.Id)
	if err != nil {
		respErr(c, err.Error())
		return
	}

	respJson(c, resp.Id)
}

func retryToGetDetailWithTx(trxId string) (TrxStatus, error) {
	var retry int
	for retry < 3 {
		retry += 1
		time.Sleep(4 * time.Second)
		status, err := checkTransactionStatus(trxId)
		switch status {
		case TrxStatusPending:
			continue
		case TrxStatusSuccess:
			return status, nil
		case TrxStatusFailed:
			logrus.Errorln(err)
			continue
		}
	}
	logrus.Info("retry", retry)
	return TrxStatusFailed, errors.New("GetTransactionDetail timeout")
}

type TrxStatus int
const (
	TrxStatusSuccess TrxStatus = iota
	TrxStatusPending
	TrxStatusFailed
)

func checkTransactionStatus(txId string) (status TrxStatus, err error) {
	path := fmt.Sprintf("/transactions/info/%s", txId)
	content, err := vsys.GetVsysApi().Get(path)
	if err != nil {
		return TrxStatusFailed, err
	}

	type TransactionDetail struct {
		vsys.TransactionResponse
		Status 			string
	}
	var trx TransactionDetail
	err = json.Unmarshal(content, &trx)
	if err != nil {
		return TrxStatusFailed, err
	}

	if strings.ToLower(trx.Status) == "success" {
		return TrxStatusSuccess, nil
	}else if trx.Status == "" {
		return TrxStatusPending, nil
	}else {
		return TrxStatusFailed, errors.New(trx.Status)
	}
}