package blockchainPoS

import (
	"os"
	"io"
	"log"
	"net"
	"time"
	"bufio"
	"strconv"
	"math/rand"
	"encoding/json"
)

var bcServer chan []Block
var CandidateBlocks = make(chan Block)
var tempBlocks []Block

var announcements = make(chan string)
var validators = make(map[string]int)

// RunTcp() starts the TCP server
func RunTcp() error {
	bcServer = make(chan []Block)
	server, err := net.Listen("tcp", ":" + os.Getenv("PORT"))
	log.Println("Listening on", os.Getenv("PORT"))
	if err != nil {
		return err
	}
	defer server.Close()

 	go func() {
		for candidate := range CandidateBlocks {
			mutex.Lock()
			tempBlocks = append(tempBlocks, candidate)
			mutex.Unlock()
		}
	}()

	go func() {
		for {
			pickWinner()
		}
	}()

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

	// go routine to write the announcements
	go func() {
		for {
			msg := <- announcements
			io.WriteString(conn, msg)
		}
	}()
	
	// the user inputs its token balance and it's given an address
	var address string
	io.WriteString(conn, "Enter token balance:")
	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v not a number: %v", scanBalance.Text(), err)
			return
		}
		t := time.Now()
		address = CalculateHash(t.String())
		validators[address] = balance
		break
	}

	// the user inputs a string a new block is created
	io.WriteString(conn, "\nEnter a new string: ")
	scanner := bufio.NewScanner(conn)
	go func() {
		for {
			for scanner.Scan() {
				newBlock, err := GenerateBlock(Blockchain[len(Blockchain)-1],scanner.Text(),address)
				if err != nil {
					log.Println(err)
					continue
				}
				log.Println(IsBlockValid(Blockchain[len(Blockchain)-1], newBlock))
				if IsBlockValid(Blockchain[len(Blockchain)-1], newBlock) {
					CandidateBlocks <- newBlock
				}
				io.WriteString(conn, "\nEnter a new string: ")
			}
		}
	}()

	// the blockchain is printed every minute
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


func pickWinner() {
	time.Sleep(30*time.Second)

	mutex.Lock()
	temp := tempBlocks
	mutex.Unlock()

	lotteryPool := []string{}
	if len(temp)>0 {
	OUTER:
		for _, block := range temp {
			for _, node := range lotteryPool {
				if block.Validator == node {
					continue OUTER
				}
			}

			mutex.Lock()
			setValidators := validators
			mutex.Unlock()

			k, ok := setValidators[block.Validator]
			if ok {
				for i:=0; i<k; i++ {
					lotteryPool = append(lotteryPool, block.Validator)
				}
			}
		}

		s := rand.NewSource(time.Now().Unix())
		r := rand.New(s)
		lotteryWinner := lotteryPool[r.Intn(len(lotteryPool))]

		for _, block := range temp {
			if block.Validator == lotteryWinner {
				mutex.Lock()
				Blockchain = append(Blockchain, block)
				mutex.Unlock()
				for _ = range validators {
					announcements <- "\nWinning validator: " + lotteryWinner
				}
				break
			}
		}
	}
	
	mutex.Lock()
	tempBlocks = []Block{}
	mutex.Unlock()
}
