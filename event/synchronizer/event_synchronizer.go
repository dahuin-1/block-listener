package main

import (
	"flag"
	"github.com/cc-ping-listener/env"
	"github.com/cc-ping-listener/unmarshalers"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"log"
	"time"
)

var (
	startBlockNum uint64
	endBlockNum   uint64
	isUpdated     bool
)

func init() {
	flag.Uint64Var(&startBlockNum, "startBlock", 0, "set start block number if needed")
	flag.Uint64Var(&endBlockNum, "endBlock", 0, "set end block number if needed")
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
	currentBlockHeight := blockchainInfo.BCI.Height
	if startBlockNum == currentBlockHeight {
		log.Fatal("nothing to sync")
	}
	if endBlockNum > currentBlockHeight {
		log.Println("end block number should be smaller than current block height")
	}
	if startBlockNum >= endBlockNum {
		log.Println("end block number should be bigger than start block number")
	}
	if endBlockNum > currentBlockHeight || startBlockNum >= endBlockNum {
		log.Println("set end block number to current block number")
		endBlockNum = currentBlockHeight - 1
		log.Printf("current block number: %d", endBlockNum)
	}
	if startBlockNum < 1 {
		log.Println("start block number should be more than 1")
		log.Println("set start block number to 1")
		startBlockNum = 1
	}
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		synchronize(ledgerClient, currentBlockHeight)
	}
	defer func() {
		ticker.Stop()
	}()
}

func synchronize(ledgerClient *ledger.Client, currentBlockHeight uint64) {
	blockchainInfo, err := ledgerClient.QueryInfo()
	if err != nil {
		log.Fatalf("failed to get new blockchain information, err: %s", err)
	}
	newBlockHeight := blockchainInfo.BCI.Height
	if newBlockHeight > currentBlockHeight {
		startBlockNum = endBlockNum + 1
		endBlockNum = newBlockHeight - 1
		isUpdated = false
	}

	log.Printf("start block number: %d, end block number: %d", startBlockNum, endBlockNum)
	//i := startBlockNum
	if !isUpdated {
		for i := startBlockNum; i <= endBlockNum; i++ {
			//for i < blockchainInfo.BCI.Height {
			log.Printf("=================================== Sync on block number: %d ===================================", i)
			block, err := ledgerClient.QueryBlock(i)
			if err != nil {
				log.Fatalf("failed to query Block, err: %s", err)
			}
			if block == nil {
				log.Printf("failed to retrieve the block from blocknumber %d. The block is nil: ", i)
				panic(err)
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
		isUpdated = true
	}
}
