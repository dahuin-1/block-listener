package main

import (
	"github.com/cc-ping-listener/unmarshalers"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type User struct {
	Cert       []byte
	PrivateKey []byte
}

const (
	channelID  = "kiesnet-dev" 
	configPath = "/config/network.yaml"
	credPath   = ""
)

func getChannelProvider() (context.ChannelProvider, error) {
	fabricUser, err := setUser()
	if err != nil {
		return nil, err
	}
	networkConfig := config.FromFile(configPath) //네트워크컨피그설정
	sdk, err := fabsdk.New(networkConfig)        //sdk객체를 얻음
	if err != nil {
		return nil, err
	}
	client, err := mspclient.New(sdk.Context()) //sdk 객체를 이용해서 channel client 생성
	if err != nil {
		return nil, err
	}
	signingIdentity, err := client.CreateSigningIdentity(mspctx.WithCert(fabricUser.Cert), mspctx.WithPrivateKey(fabricUser.PrivateKey))
	if err != nil {
		return nil, err
	}
	channelProvider := sdk.ChannelContext(channelID, fabsdk.WithIdentity(signingIdentity))
	return channelProvider, nil
}

func getFunctionName(ChaincodeProposalPayload []byte) string {
	chaincodeAction, err := unmarshalers.GetChaincodeAction(ChaincodeProposalPayload)
	if err != nil {
		log.Fatalf("unmarshaling Chaincode Action Payload error: %s", err)
	}
	funcName := parse(chaincodeAction.String())
	return funcName
}

func parse(str string) string {
	firstIndex := strings.Index(str, "fruit")
	lastIndex := strings.LastIndex(str, `\n`)
	return str[firstIndex:lastIndex]
}

func setUser() (*User, error) {
	mspPath := filepath.Join(credPath, "username", "msp")
	certPath := filepath.Join(mspPath, "signcerts", "cert.pem")
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	keyStore := filepath.Join(mspPath, "keystore")
	keys, err := os.ReadDir(keyStore)
	if err != nil {
		return nil, err
	}
	keyPath := filepath.Join(keyStore, keys[0].Name())
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return &User{Cert: cert, PrivateKey: key}, nil
}

func main() {
	channelProvider, err := getChannelProvider()
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
		select {
		case e := <-eventChannel:
			log.Println("############################################################")
			log.Println("###################### Received event ######################")
			log.Printf("################### BlockNum : %d ##########################", e.Block.Header.Number)
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
			//log.Println("~~~~~~~~~")
			//log.Println(proposalResponsePayload.String())
			//log.Println("~~~~~~~~~")
			chaincodeAction, err := unmarshalers.GetChaincodeAction(proposalResponsePayload.Extension)
			if err != nil {
				log.Fatalf("unmarshaling proposalResponsePayload Extension to chaincodeAction error: %s", err)
			}
			chaincodeEvent, err := unmarshalers.GetChaincodeEvent(chaincodeAction.Events)
			if err != nil {
				log.Fatalf("unmarshaling chaincodeAction.Events to chaincodeEvent error: %s", err)
			}
			if chaincodeEvent.EventName != "" {
				log.Printf("#################### Block event : %v ########### ", chaincodeEvent.EventName)
			} else {
				log.Printf("#################### Function name : %v ########### ", getFunctionName(chaincodeActionPayload.ChaincodeProposalPayload))
			}
			log.Printf("#################### Block info : %v ########### ", chaincodeAction.String())
			log.Println("#############################################################")
		}
	}
}
