package main

import (
	"flag"
	"github.com/cc-ping-listener/env"
	"github.com/cc-ping-listener/unmarshalers"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"log"
)

var startBlockNum uint64

func init() {
	flag.Uint64Var(&startBlockNum, "startBlock", 1, "set start block number if needed")
	flag.Parse()
}

func main() {
	channelProvider, err := env.GetChannelProvider()
	if err != nil {
		log.Fatalf("failed to get Channel Provider, err: %s", err)
	}
	eventClient, err := event.New(channelProvider, event.WithBlockEvents(), event.WithSeekType("from"), event.WithBlockNum(startBlockNum))
	if err != nil {
		log.Fatalf("failed to return Client instance, err: %s", err)
	}
	registration, eventChannel, err := eventClient.RegisterBlockEvent()
	if err != nil {
		log.Fatalf("failed to register Block Event, err: %s", err)
	}
	defer eventClient.Unregister(registration)
	for e := range eventChannel {
		log.Printf("🎹👂🎹👂🎹👂🎹👂🎹👂🏻listen👂🏻🎹👂🎹👂🎹👂🎹👂🎹")
		log.Printf("=================================== Received block number: %d ===================================", e.Block.Header.Number)
		blockData := e.Block.Data.Data
		envelope, err := unmarshalers.GetEnvelopeFromBlock(blockData[0])
		if err != nil {
			log.Fatalf("unmarshaling Envelope error: %s", err)
		}
		payload, err := unmarshalers.GetPayloadFromEnv(envelope.Payload)
		if err != nil {
			log.Fatalf("unmarshaling envelope Payload to payload error: %s", err)
		}
		transaction, err := unmarshalers.GetTransaction(payload.Data)
		if err != nil {
			log.Fatalf("unmarshaling payload Data to transaction error: %s", err)
		}
		chaincodeActionPayload, err := unmarshalers.GetChaincodeActionPayload(transaction.Actions[0].Payload)
		if err != nil {
			log.Fatalf("unmarshaling transaction Action Payload to chaincodeActionPayload error: %s", err)
		}
		proposalResponsePayload, err := unmarshalers.GetProposalResponsePayload(chaincodeActionPayload.Action.ProposalResponsePayload)
		if err != nil {
			log.Fatalf("unmarshaling ProposalResponsePayload to proposalResponsePayload error: %s", err)
		}
		chaincodeAction, err := unmarshalers.GetChaincodeAction(proposalResponsePayload.Extension)
		if err != nil {
			log.Fatalf("unmarshaling proposalResponsePayload Extension to chaincodeAction error: %s", err)
		}
		chaincodeEvent, err := unmarshalers.GetChaincodeEvent(chaincodeAction.Events)
		if err != nil {
			log.Fatalf("unmarshaling chaincodeAction Events to chaincodeEvent error: %s", err)
		}
		if chaincodeEvent.EventName == "" {
			log.Println("event did not happen")
		} else {
			log.Printf("#################### Block event : %v ########### ", chaincodeEvent.EventName)
			log.Printf("#################### Block info - block %v ########### ", chaincodeAction.String())
		}
	}
}
