package enum

type blockchain string

const (
	GenesisFilepath    blockchain = "genesis.json"
	GenesisBlockNumber blockchain = "0"
	BlockchianDb       blockchain = "./blockcahinDb"
	BlockstateDb       blockchain = "./blockstateDb"
	LastHash           blockchain = "lh"
	CurrentBlockNumber blockchain = "blockNumber"
	ConfigFilePath     blockchain = "config.toml"
)
