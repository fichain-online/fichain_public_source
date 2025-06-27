package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"FichainCore/client"
	"FichainCore/common"
	"FichainCore/config"
	"FichainCore/params"

	logger "github.com/HendrickPhan/golang-simple-logger"
)

const (
	REPLACE_ADDRESS = "1510151015101510151015101510151015101510"
)

func main() {
	configFile := flag.String("config", "", "Config file path")
	dataFile := flag.String("data", "", "data file path")
	flag.Parse()

	if *configFile == "" {
		logger.Error("You must provide a config file --config")
		panic("")
	}
	config.InitConfig(*configFile)

	if *dataFile == "" {
		logger.Error("You must provide a data file --data")
		panic("")
	}
	// create tx client
	txClient := client.NewClient(config.GetConfig())

	data, err := os.ReadFile(*dataFile)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	var txData []*TransactionData
	err = json.Unmarshal(data, &txData)
	if err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	// 5. Print the result to verify
	logger.Info("Successfully decoded JSON from file stream!, total data: ", len(txData))
	addressList := []string{}
	writeData := ""
	for _, v := range txData {
		amount := new(big.Int)
		amount.SetString(v.Amount, 10)

		for j := 0; j < len(v.ReplaceAddress); j++ {
			v.Data = strings.Replace(v.Data, REPLACE_ADDRESS, addressList[v.ReplaceAddress[j]], 1)
		}

		receipt, err := txClient.SendTransaction(
			common.HexToAddress(v.To),
			amount,
			common.FromHex(v.Data),
			v.Gas,
			params.TempGasPrice,
			v.Name,
		)
		if err != nil {
			logger.Error("error when send tx", err)
			panic(err)
		}
		logger.Info("receipt", receipt)
		addressList = append(addressList, hex.EncodeToString(receipt.ContractAddress.Bytes()))

		writeData = fmt.Sprintf("%v%v:%v\n", writeData, v.Name, receipt.String())
	}

	logger.Info("write data", writeData)
	os.WriteFile(
		fmt.Sprintf("result_%v.dat", time.Now().Unix()),
		[]byte(writeData),
		0644,
	)
}
