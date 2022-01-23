# GoBlockchainPoS
## Introduction
This is my Go implementation of the proof-of-stake blockchain proposed by [Coral Health](https://github.com/nosequeldeebee/blockchain-tutorial). If you want to see the proof-of-work version of this blockchain, click [here](https://github.com/augcos/GoBlockchainPoW). This project was developed using Go v1.17.5 for Linux systems.

## How to install?
First, you will have to make sure to have preinstalled the required third-party packages. You can install them using the following commands:
```
go get github.com/joho/godotenv
```
You can clone this repository to your local system using the command:
```
git clone github.com/augcos/GoBlockchainPoS
```

## How to run?
### Using a TCP connection
GoBlockchainPoS runs using a TCP server. In order to launch both the server and the blockchain, run the main.go file:
```
go run main.go
```
Then, open a different terminal and connect to the TCP server:
```
nc localhost 8080
```
You will be prompted to input your number of tokens and a string for a new block. Every 30 seconds, a validator will be chosen, and its block will be appended to the blockchain.