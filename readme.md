# Crypto Wallet Demo
## Instructions To Run
[![Build Status](https://travis-ci.org/joemccann/dillinger.svg?branch=master)](https://travis-ci.org/joemccann/dillinger)

## Prerequisites

- Alchemy Account
- Alchemy Api key
- Metamask
- Go 

## After cloning
```sh
cp .env.example .env
```

In the env file add your alchemy api key

> ALCHEMY_URL=https://eth-sepolia.g.alchemy.com/YOUR_API_KEY


```sh
go mod tidy
go run main.go
```

The server will start on [localhost:8080](http://localhost:8080)

## Apis

| Function | Path | Response
| ------ | ------ |------ |
| Create Wallet | /create-wallet| address, private key |
| Get Balance | /get-balance/:address | balance of eth |
| Send Eth | /send-eth | transaction hash |


## Working
Create a wallet using the create wallet api, it will return the wallet address and the priavte key. Save both in a text file for now.

Setup metamask and save your key phrase, it can be used for wallet recovery. After the setup is complete you will have an address in meta mask. You will need test etherium to test the app. Copy the address that you got from the create wallet api and use this [Site](https://sepolia-faucet.pk910.de/) to transfer some eth to that address. 

After the test eth have been transferred you will recieve a transaction hash which you can use to verify the transaction on the [Sepolia test net](https://sepolia.etherscan.io/). 

Now use the get balance api to verify the balance in your wallet. Use the address that you saved from the create wallet api.

After balance is verified, use the send eth api to transfer some eth from your wallet to your metamask address. The request body is as follows:
```
{
    "private_key": "YOUR_PRIVATE_KEY",    //from create wallet api
    "from_address":"YOUR_WLLET_ADDRESS",  //from create wallet api
    "to_address": "METAMASK_ADDRESS",     //copy from metamask wallet
    "amount_in_eth": "0.01"               //should be less than balance
}
```
This api will return a transaction hash which you use to verify the transaction on the sepolia network. Your Alchemy account will also show this transaction in the mempool tab in your account. Also verify from metamask. It should reieve the eth