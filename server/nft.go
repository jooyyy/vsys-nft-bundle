package server

import (
	"encoding/json"
	"fmt"
	"github.com/virtualeconomy/go-v-sdk/vsys"
	"vsys-nft-bundle/config"
)

func GetLatestNFTIndex() (index int, err error) {
	url := vsys.GetVsysApi().NodeAddress + fmt.Sprintf(ApiNFTGetLatestIndex, config.Config.V.NFTContractId)
	data, err := vsys.UrlGetContent(url)
	if err != nil {
		return 0, err
	}
	type resp struct {
		ContractId 		string `json:"contractId"`
		LastTokenIndex 	int `json:"lastTokenIndex"`
	}
	var d resp
	err = json.Unmarshal(data, &d)
	if err != nil {
		return 0, err
	}
	return d.LastTokenIndex, nil
}

func GetNFTByIndex(index int) (token string, err error) {
	url := fmt.Sprintf("%s%s",
		vsys.GetVsysApi().NodeAddress,
		fmt.Sprintf(ApiNFTGetByIndex, config.Config.V.NFTContractId, index),
	)
	data, err := vsys.UrlGetContent(url)
	if err != nil {
		return "", err
	}
	type resp struct {
		TokenId 		string `json:"tokenId"`
	}
	var r resp
	err = json.Unmarshal(data, &r)
	if err != nil {
		return "", err
	}
	return r.TokenId, nil
}
