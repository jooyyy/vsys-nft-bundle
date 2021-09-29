package server

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/virtualeconomy/go-v-sdk/vsys"
	"strconv"
	"vsys-nft-bundle/config"
)

// ApplyGasFee
// @Summary Apply gas fee from pool wallet
// @Description Every contract function execute operation coast 0.3 vsys
// @Tags Transaction
// @Accept  json
// @Produce json
// @Param recipient formData string true "NFT contract description"
// @Param amount formData int64 false "Amount in the smallest unit, 0.3 vsys for 30000000 gas by default"
// @Success 200 {object} string "Get transaction id when apply gas fee success."
// @Failure 400 {object} RespJsonObj ""
// @Router /api/transaction/applyGasFee [post]
func (s *Service) ApplyGasFee(c *gin.Context) {
	recipient, ok := c.GetPostForm("recipient")
	if !ok {
		respErr(c, "need recipient")
		return
	}
	amountStr, _ := c.GetPostForm("amount")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil || amount == 0 {
		amount = 30000000
	}

	trxId, err := applyGasFee(recipient, amount)
	if err != nil {
		respErr(c, err.Error())
		return
	}

	respJson(c, trxId)
}

func applyGasFee(recipient string, amount int64) (string, error) {
	if vsys.IsValidateAddress(recipient, config.Config.V.NetType) == false {
		return "", errors.New("invalid recipient")
	}
	acc := vsys.InitAccount(config.Config.V.NetType)
	acc.BuildFromSeed(config.Config.V.GasPool, 0)
	trx := acc.BuildPayment(recipient, amount, fmt.Sprintf("Gas fee from pool."))
	resp, err := vsys.SendPaymentTx(trx)
	if err != nil {
		return "", err
	}
	_, err = retryToGetDetailWithTx(resp.Id)
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}
