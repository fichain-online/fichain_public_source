package main

import (
	"flag"
	"fmt"
	"math/big"
	"strings"
	"time"

	"FichainCore/client"
	"FichainCore/common"
	core_config "FichainCore/config"
	"FichainCore/params"

	"github.com/ethereum/go-ethereum/accounts/abi"
	logger "github.com/hieuphanuit/golang-simple-logger"

	"FichainGoldInvestServer/config"
	"FichainGoldInvestServer/price_crawler"
)

const (
	FichainGasLimit = uint64(200000)
)

var goldInvestAbi abi.ABI

func initABI() error {
	var err error
	goldInvestAbi, err = abi.JSON(
		strings.NewReader(
			`[
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_buyPrice",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "_sellPrice",
				"type": "uint256"
			}
		],
		"name": "setPrices",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "getPrices",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`,
		),
	)
	if err != nil {
		return fmt.Errorf("failed to parse ABI: %w", err)
	}
	return nil
}

func createQueryPayload() ([]byte, error) {
	// Pack the arguments with the method name.
	return goldInvestAbi.Pack("getPrices")
}

func createUpdatePayload(buy, sell *big.Int) ([]byte, error) {
	// Pack the arguments with the method name.
	return goldInvestAbi.Pack("setPrices", buy, sell)
}

func main() {
	configFile := flag.String("config", "", "Config file path")
	flag.Parse()

	if *configFile == "" {
		logger.Error("You must provide a config file --config")
		panic("")
	}
	config.InitConfig(*configFile)
	core_config.SetConfig(config.GetConfig().Core)
	initABI()
	// create tx client
	txClient := client.NewClient(config.GetConfig().Core)
	update(txClient)

	for {
		<-time.After(5 * time.Minute)
		update(txClient)
	}
}

func update(txClient *client.Client) {
	// let call to get current price
	data, err := createQueryPayload()
	if err != nil {
		logger.Error("err", err)
		return
	}
	rs, err := txClient.CallSmartContract(
		common.HexToAddress(config.GetConfig().GoldInvestAddress),
		data,
	)
	if err != nil {
		logger.Error("err", err)
		return
	}

	unpacked, err := goldInvestAbi.Unpack("getPrices", rs.Data)
	unpackedBuy := big.NewInt(0)
	unpackedSell := big.NewInt(0)
	for i, v := range unpacked {
		if i == 0 {
			unpackedBuy = v.(*big.Int)
		} else {
			unpackedSell = v.(*big.Int)
		}
	}

	buy, sell, err := price_crawler.Crawl()
	if err != nil {
		logger.Error("err", err)
		return
	}
	logger.Info("Crawled Price", buy, sell)
	mul := big.NewInt(0)
	mul.SetString("100000000000000000000", 10)
	// format buy & sell crawled
	fmtBuy := big.NewInt(int64(buy))
	fmtBuy = fmtBuy.Mul(fmtBuy, mul)
	fmtSell := big.NewInt(int64(sell))
	fmtSell = fmtSell.Mul(fmtSell, mul)
	if fmtBuy.Cmp(unpackedBuy) != 0 || fmtSell.Cmp(unpackedSell) != 0 {
		logger.Info("price diff")
		payload, err := createUpdatePayload(fmtBuy, fmtSell)
		if err != nil {
			logger.Error("error when create update payload", err)
			return
		}
		receipt, err := txClient.SendTransaction(
			common.HexToAddress(config.GetConfig().GoldInvestAddress),
			big.NewInt(0),
			payload,
			FichainGasLimit,
			params.TempGasPrice,
			"update gold price",
		)
		if err != nil {
			logger.Error("error when send update transaction", err)
			return
		}
		logger.Info("update price success, receipt", receipt)
	} else {
		logger.Info("Price not change")
	}
}
