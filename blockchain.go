package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	. "github.com/augcos/GoBlockchainPoS/blockchainPoS"
)



func main() {
	// loads the enviroment variables (.env file)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading the .env file")
	}



	// creates the genesis block and starts the blockchain
	go func() {
		var genesisBlock Block
		genesisBlock = Block{0, time.Now().String(), "genesis", "", "", ""}
		genesisBlock.Hash = CalculateBlockHash(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)
	}()
	// runs the server
	log.Fatal(RunTcp())
}
