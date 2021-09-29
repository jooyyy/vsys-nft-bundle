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
// @Success 200 {object} Wallet ""
// @Failure 400 {object} RespJsonObj ""
// @Router /api/contract/registerNFT [post]
func (s *Service) RegisterNFTContract(c *gin.Context) {
	desc, _ := c.GetPostForm("description")
	acc := vsys.InitAccount(config.Config.V.NetType)
	acc.BuildFromSeed(config.Config.V.MakerSeed, 0)
	transaction := vsys.NewRegisterTransaction(CONTRACT_NFT, "", desc)
	transaction.SenderPublicKey = acc.PublicKey()
	transaction.Signature = acc.SignData(transaction.BuildTxData())
	data, err := vsys.GetVsysApi().Post(vsys.ApiContractBroadcastRegister, transaction)
	if err != nil {
		respErr(c, err.Error())
		return
	}
	logrus.Info("data", string(data))
	_, err = retryToGetDetailWithTx(transaction.TxId)
	if err != nil {
		respErr(c, err.Error())
		return
	}
	respJson(c, Contract{
		TrxId:   	 transaction.TxId,
		ContractId:  "",
	})
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
			logrus.Error(err)
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