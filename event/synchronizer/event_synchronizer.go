package main

import (
	"flag"
	"github.com/cc-ping-listener/env"
	"github.com/cc-ping-listener/unmarshalers"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"log"
	"time"
)

var startBlockNum uint64

func init() {
	flag.Uint64Var(&startBlockNum, "startBlock", 0, "set start block number if needed")
	flag.Parse()
}

func main() {
	channelProvider, err := env.GetChannelProvider()
	if err != nil {
		log.Fatalf("failed to get channel provider, err: %s", err)
	}
	ledgerClient, err := ledger.New(channelProvider)
	if err != nil {
		log.Fatalf("failed to return ledger client instance, err: %s", err)
	}
	blockchainInfo, err := ledgerClient.QueryInfo()
	if err != nil {
		log.Fatalf("failed to get blockchain information, err: %s", err)
	}
	if startBlockNum == blockchainInfo.BCI.Height {
		log.Fatalf("nothing to sync, now block number is %d", blockchainInfo.BCI.Height-1)
	}
	if startBlockNum < 1 {
		log.Println("start block number should be more than 1")
		log.Println("set start block number to 1")
		startBlockNum = 1
	}
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for ; true; <-ticker.C {
		synchronize(ledgerClient)
	}
}

func synchronize(ledgerClient *ledger.Client) {
	blockchainInfo, err := ledgerClient.QueryInfo()
	if err != nil {
		log.Fatalf("failed to get new blockchain information, err: %s", err)
	}
	log.Printf("set start block number: %d waiting for sync...", startBlockNum)
	i := startBlockNum
	for i < blockchainInfo.BCI.Height {
		log.Printf("=================================== Sync on block number: %d ===================================", i)
		block, err := ledgerClient.QueryBlock(i)
		if err != nil {
			log.Fatalf("failed to query Block, err: %s", err)
		}
		if block == nil {
			log.Fatalf("failed to retrieve the block from blocknumber %d. The block is nil: ", i)
		}
		blockData := block.Data.Data
		envelope, err := unmarshalers.GetEnvelopeFromBlock(blockData[0])
		if err != nil {
			log.Fatalf("unmarshaling Envelope error: %s", err)
		}
		payload, err := unmarshalers.GetPayloadFromEnv(envelope.Payload)
		if err != nil {
			log.Fatalf("unmarshaling envelopePayload to payload error: %s", err)
		}
		transaction, err := unmarshalers.GetTransaction(payload.Data)
		if err != nil {
			log.Fatalf("unmarshaling payloadData to transaction error: %s", err)
		}
		chaincodeActionPayload, err := unmarshalers.GetChaincodeActionPayload(transaction.Actions[0].Payload)
		if err != nil {
			log.Fatalf("unmarshaling transactionActionPayload to chaincodeActionPayload error: %s", err)
		}
		proposalResponsePayload, err := unmarshalers.GetProposalResponsePayload(chaincodeActionPayload.Action.ProposalResponsePayload)
		if err != nil {
			log.Fatalf("unmarshaling chaincodeActionPayload.Action ProposalResponsePayload to proposalResponsePayload error: %s", err)
		}
		chaincodeAction, err := unmarshalers.GetChaincodeAction(proposalResponsePayload.Extension)
		if err != nil {
			log.Fatalf("unmarshaling proposalResponsePayload Extension to chaincodeAction error: %s", err)
		}
		chaincodeEvent, err := unmarshalers.GetChaincodeEvent(chaincodeAction.Events)
		if err != nil {
			log.Fatalf("unmarshaling chaincodeAction.Events to chaincodeEvent error: %s", err)
		}
		if chaincodeEvent.EventName == "" {
			log.Println("event did not happen")
		} else {
			log.Printf("#################### Block event : %v ########### ", chaincodeEvent.EventName)
			log.Printf("#################### Block info - block %v ########### ", chaincodeAction.String())
		}
		i++
	}
	startBlockNum = blockchainInfo.BCI.Height
}
