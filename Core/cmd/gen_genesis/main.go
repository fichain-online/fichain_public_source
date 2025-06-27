package main

import (
	"encoding/json"
	"flag"
	"os"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/config"
	"FichainCore/database"
	"FichainCore/genesis"
)

func main() {
	// Parse command-line flags
	configFile := flag.String("config", "", "Config file path")
	genesisFile := flag.String("genesis", "", "Config file path")
	flag.Parse()

	if *configFile == "" {
		panic("You must provide a config file --config")
	}
	if *genesisFile == "" {
		panic("You must provide a config file --config")
	}
	config.InitConfig(*configFile)
	gns := &genesis.Genesis{}
	// load genesis
	bData, err := os.ReadFile(*genesisFile)
	if err != nil {
		panic("error when loading genesis file" + err.Error())
	}
	err = json.Unmarshal(bData, gns)
	if err != nil {
		panic("error when unmarshal genesis json data" + err.Error())
	}
	// commit genesis
	// load databases
	stateDB, err := database.NewBadgerDB(config.GetConfig().StatesDBPath)
	if err != nil {
		panic("error when load stateDB" + err.Error())
	}
	authorityValidatorDB, err := database.NewBadgerDB(config.GetConfig().AuthorityValidatorDBPath)
	if err != nil {
		panic("error when load stateDB" + err.Error())
	}
	authorityObserverDB, err := database.NewBadgerDB(config.GetConfig().AuthorityObserverDBPath)
	if err != nil {
		panic("error when load stateDB" + err.Error())
	}
	authorityFiatReseveDB, err := database.NewBadgerDB(
		config.GetConfig().AuthorityFiatReserveDBPath,
	)
	if err != nil {
		panic("error when load stateDB" + err.Error())
	}
	genesisBlock := gns.MustCommit(stateDB)
	logger.Info("Statedb commited")
	logger.Info("Genesis block", genesisBlock)
	err = gns.CommitAuthorities(
		authorityValidatorDB,
		authorityObserverDB,
		authorityFiatReseveDB,
	)
	if err != nil {
		panic("error when commit authorities" + err.Error())
	}
	logger.Info("Authoriries commited")
}
