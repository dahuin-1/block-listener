package main

import (
	"flag"
	"github.com/cc-ping-listener/env"
	"github.com/cc-ping-listener/unmarshalers"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"log"
)

var startBlockNum uint64
var endBlockNum uint64

//type User struct {
//	Cert       []byte
//	PrivateKey []byte
//}
//const (
//	channelID  = "kiesnet-dev"
//	configPath = "/Users/dhkim/Projects/cc-ping-listener/config/network.yaml"
//	credPath   = "/Users/dhkim/Projects/kiesnet-chaincode-dev-network/crypto-config/peerOrganizations/kiesnet.dev/users"
//)
//
//
//func getChannelProvider() (context.ChannelProvider, error) {
//	fabricUser, err := setUser()
//	if err != nil {
//		return nil, err
//	}
//	networkConfig := config.FromFile(configPath)
//	sdk, err := fabsdk.New(networkConfig)
//	if err != nil {
//		return nil, err
//	}
//	client, err := mspclient.New(sdk.Context()) //sdk 객체를 이용해서 channel client 생성
//	if err != nil {
//		return nil, err
//	}
//	signingIdentity, err := client.CreateSigningIdentity(mspctx.WithCert(fabricUser.Cert), mspctx.WithPrivateKey(fabricUser.PrivateKey))
//	if err != nil {
//		return nil, err
//	}
//	channelProvider := sdk.ChannelContext(channelID, fabsdk.WithIdentity(signingIdentity))
//	return channelProvider, nil
//}
//
//func setUser() (*User, error) {
//	mspPath := filepath.Join(credPath, "dhkim", "msp")
//	certPath := filepath.Join(mspPath, "signcerts", "cert.pem")
//	cert, err := os.ReadFile(certPath)
//	if err != nil {
//		return nil, err
//	}
//	keyStore := filepath.Join(mspPath, "keystore")
//	keys, err := os.ReadDir(keyStore)
//	if err != nil {
//		return nil, err
//	}
//	keyPath := filepath.Join(keyStore, keys[0].Name())
//	key, err := os.ReadFile(keyPath)
//	if err != nil {
//		return nil, err
//	}
//	return &User{Cert: cert, PrivateKey: key}, nil
//}

func getStartBlock(startBlock uint64) error {
	log.Printf("hello world, %d", startBlock)
	return nil
}

func init() {
	flag.Uint64Var(&startBlockNum, "startBlock", 0, "set start block number if needed")
	flag.Uint64Var(&endBlockNum, "endBlock", 0, "set end block number if needed")
	flag.Parse()
}

func main() {
	err := getStartBlock(startBlockNum)
	if err != nil {
		log.Fatalf("err")
	}
	channelProvider, err := env.GetChannelProvider()
	if err != nil {
		log.Fatalf("failed to get Channel Provider, err: %s", err)
	}
	client, err := event.New(channelProvider, event.WithBlockEvents())
	if err != nil {
		log.Fatalf("failed to return Client instance, err: %s", err)
	}
	registration, eventChannel, err := client.RegisterBlockEvent()
	if err != nil {
		log.Fatalf("failed to register Block Event, err: %s", err)
	}
	defer client.Unregister(registration)
	for {
		log.Printf("🎹👂🎹👂🎹👂🎹👂🎹👂🏻listen👂🏻🎹👂🎹👂🎹👂🎹👂🎹")
		e := <-eventChannel
		blockData := e.Block.Data.Data
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
		log.Printf("=================================== Received block number: %d ===================================", e.Block.Header.Number)
		if chaincodeEvent.EventName == "" {
			log.Println("event did not happen")
		} else {
			log.Printf("#################### Block event : %v ########### ", chaincodeEvent.EventName)
			log.Printf("#################### Block info - block %v ########### ", chaincodeAction.String())
		}
	}
}
