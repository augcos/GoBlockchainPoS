package blockchainPoS

import (
	"os"
	"io"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
	"bufio"
	"strconv"
	"math/rand"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
)

// var blockTime indicates how many seconds between blocks
var blockTime = 30

// mutex for mutual exclusion (avoid data races between routines)
var mutex = &sync.Mutex{}

// channel declarations
var bcServer chan []Block
var CandidateBlocks = make(chan Block)
var tempBlocks []Block
var announcements = make(chan string)

// validator nodes list
var validators = make(map[string]int)
var proposedBlock = make(map[string]bool)

// RunTcp() starts the TCP server
func RunTcp() error {
	// starts the server on the chosen port
	server, err := net.Listen("tcp", ":" + os.Getenv("PORT"))
	log.Println("Listening on", os.Getenv("PORT"))
	if err != nil {
		return err
	}

	// closes down the server after exiting the function
	defer server.Close()

	// go routine to add candidate blocks to the temp block list
 	go func() {
		for candidate := range CandidateBlocks {
			mutex.Lock()
			tempBlocks = append(tempBlocks, candidate)
			mutex.Unlock()
		}
	}()

	// go routine to pick a winner every blockTime seconds
	go func() {
		for {
			pickWinner()
		}
	}()

	// go routine to handle any new incoming connections
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
	return nil
}


// handleConn() manages each new connection
func handleConn(conn net.Conn) {
	//we close the connection after exiting the function
	defer conn.Close()

	// go routine to write the announcements to the validator
	go func() {
		for {
			msg := <- announcements
			io.WriteString(conn, msg)
		}
	}()
	
	var address string
	t := time.Now()
	address = CalculateHash(t.String())
	io.WriteString(conn, fmt.Sprint("Your address is: ", address))

	// the user inputs its token balance and recieves an address
	io.WriteString(conn, "\nEnter token balance: ")
	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v not a number: %v - Please retry the connection", scanBalance.Text(), err)
			return
		}
		validators[address] = balance
		break
	}

	// the user inputs a string a new block is created
	io.WriteString(conn, "Enter a new string: ")
	scanner := bufio.NewScanner(conn)
	go func() {
		for scanner.Scan() {
			if proposedBlock[address]==false{ 
				newBlock, err := GenerateBlock(Blockchain[len(Blockchain)-1],scanner.Text(),address)
				if err != nil {
					log.Println(err)
					continue
				}
				if IsBlockValid(Blockchain[len(Blockchain)-1], newBlock) {
					CandidateBlocks <- newBlock
					proposedBlock[address] = true
				} 
				io.WriteString(conn, "A new block has been proposed.")
			} else {
				io.WriteString(conn, "You can not enter any new blocks until validation is decided.")
			}
		}	
	}()

	// the blockchain is printed to the nodes every minute
	for {
		time.Sleep(time.Minute)
		mutex.Lock()
		output, err := json.Marshal(Blockchain)
		mutex.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(conn, string(output) + "\n")
	}
}


// pickWinner() manages block validation and it addition to the chain
func pickWinner() {
	// only one block is validated every blockTime seconds
	time.Sleep(time.Duration(blockTime)*time.Second)

	// the tempBlocks are copied to the aux variable temp
	mutex.Lock()
	temp := tempBlocks
	mutex.Unlock()

	// the tempBlocks are copied to the aux variable temp
	lotteryPool := []string{}
	if len(temp)>0 {
	OUTER:
		for _, block := range temp {
			
	// the tempBlocks are copied to the aux variable temp
			for _, node := range lotteryPool {
				if block.Validator == node {
					continue OUTER
				}
			}

			// the validator variable is copied to the an aux variable
			mutex.Lock()
			setValidators := validators
			mutex.Unlock()

			// each validator gets proportional representation in the lotteryPool to its number of tokens
			tokens, ok := setValidators[block.Validator]
			if ok {
				for i:=0; i<tokens; i++ {
					lotteryPool = append(lotteryPool, block.Validator)
				}
			}
		}

		// randomly gets a winner for block validation
		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		lotteryWinner := lotteryPool[r.Intn(len(lotteryPool))]

		// goes through the temporary blokcs and appends the winner block to the blockchain
		for _, block := range temp {
			if block.Validator == lotteryWinner {
				mutex.Lock()
				Blockchain = append(Blockchain, block)
				mutex.Unlock()
				for addr := range validators {
					if addr==lotteryWinner {
						announcements <- "\nWinning validator: " + lotteryWinner + "\nEnter a new string: "
					} else {
						announcements <- "\nCongrats! You won the validation.\nEnter a new string: "
					}
					proposedBlock[addr] = false
				}
				break
			}
		}
	}

	// empties the tempBlocks slice
	mutex.Lock()
	tempBlocks = []Block{}
	mutex.Unlock()

	// print the blockchain to terminal
	spew.Dump(Blockchain)
}
