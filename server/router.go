package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "net/http/pprof"
	_ "vsys-nft-bundle/server/docs"
)

func (s *Service) initRouter(r gin.IRouter) {
	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.GET("/api/wallet/detail", s.GetWalletDetail)
	r.GET("/api/wallet/create", s.CreateWallet)
	r.POST("/api/wallet/import", s.ImportWallet)
	r.POST("/api/contract/registerNFT", s.RegisterNFTContract)
	r.POST("/api/contract/mintNFT", s.MintNFT)
	r.POST("/api/contract/sendNFT", s.SendNFT)
	r.POST("/api/contract/transferNFT", s.TransferNFT)
	r.POST("/api/transaction/applyGasFee", s.ApplyGasFee)
}
