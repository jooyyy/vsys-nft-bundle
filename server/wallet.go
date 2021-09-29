package server

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/virtualeconomy/go-v-sdk/vsys"
	"vsys-nft-bundle/config"
)

// CreateWallet
// @Summary Create a v chain wallet
// @Description
// @Tags Wallet
// @Accept  json
// @Produce json
//
// @Success 200 {object} Wallet ""
// @Failure 400 {object} RespJsonObj ""
// @Router /api/wallet/create [get]
func (s *Service) CreateWallet(c *gin.Context) {
	acc := vsys.InitAccount(config.Config.V.NetType)
	acc.BuildFromSeed(vsys.GenerateSeed(), 0)
	respJson(c, Wallet{
		Seed:       acc.AccountSeed(),
		PrivateKey: acc.PrivateKey(),
		PublicKey:  acc.PublicKey(),
		Address:    acc.Address(),
		NetType:    config.Config.V.Env,
	})
}

// ImportWallet
// @Summary Import a v chain wallet
// @Description
// @Tags Wallet
// @Accept  json
// @Produce json
// @Param seed formData string true "wallet mnemonic seed"
// @Success 200 {object} Wallet ""
// @Failure 400 {object} RespJsonObj ""
// @Router /api/wallet/import [post]
func (s *Service) ImportWallet(c *gin.Context) {
	seed, ok := c.GetPostForm("seed")
	if !ok {
		respErr(c, "seed is empty")
		return
	}
	acc := vsys.InitAccount(config.Config.V.NetType)
	acc.BuildFromSeed(seed, 0)
	respJson(c, Wallet{
		Seed:       acc.AccountSeed(),
		PrivateKey: acc.PrivateKey(),
		PublicKey:  acc.PublicKey(),
		Address:    acc.Address(),
		NetType:    config.Config.V.Env,
	})
}

// GetWalletDetail
// @Summary Get wallet detail
// @Description
// @Tags Wallet
// @Accept  json
// @Produce json
// @Param address query string true "Wallet address"
// @Success 200 {object} WalletDetail ""
// @Failure 400 {object} RespJsonObj ""
// @Router /api/wallet/detail [get]
func (s *Service) GetWalletDetail(c *gin.Context) {
	resp, err := vsys.GetVsysApi().Get(fmt.Sprintf(vsys.ApiGetAddressDetail, c.Query("address")))
	if err != nil {
		respErr(c, err.Error())
		return
	}
	var wallet WalletDetail
	err = json.Unmarshal(resp, &wallet)
	if err != nil {
		respErr(c, err.Error())
		return
	}
	respJson(c, wallet)
}
