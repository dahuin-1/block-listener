package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer"
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
	configPath = "/Users/dhkim/Projects/cc-ping-listener/config/network.yaml"
	credPath   = "/Users/dhkim/Projects/kiesnet-chaincode-dev-network/crypto-config/peerOrganizations/kiesnet.dev/users"
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

func setUser() (*User, error) {
	mspPath := filepath.Join(credPath, "dhkim", "msp")
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

//func getEnvelopeFromBlock(data []byte) (*common.Envelope, error) {
//	env := &common.Envelope{}
//	err := proto.Unmarshal(data, env)
//	if err != nil {
//		return nil, err
//	}
//	return env, nil
//}
//
//func getPayloadFromEnv(data []byte) (*common.Payload, error) {
//	payload := &common.Payload{}
//	err := proto.Unmarshal(data, payload)
//	if err != nil {
//		return nil, err
//	}
//	return payload, nil
//}

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
			//First Get the Envelope from the BlockData
			envelope, err := GetEnvelopeFromBlock(blockData[0])
			if err != nil {
				log.Fatalf("unmarshaling Envelope error: %s", err)
			}
			//Retrieve the Payload from the Envelope
			//payload := &common.Payload{}
			//err = proto.Unmarshal(envelope.Payload, payload)
			payload, err := getPayloadFromEnv(envelope.Payload)
			if err != nil {
				log.Fatalf("unmarshaling envelopePayload to payload error: %s", err)
			}
			//Read the Transaction from the Payload Data
			transaction := &peer.Transaction{}
			err = proto.Unmarshal(payload.Data, transaction)
			if err != nil {
				log.Fatalf("unmarshaling payloadData to transaction error: %s", err)
			}
			// Payload field is marshalled object of ChaincodeActionPayload
			chaincodeActionPayload := &peer.ChaincodeActionPayload{}
			err = proto.Unmarshal(transaction.Actions[0].Payload, chaincodeActionPayload)
			if err != nil {
				log.Fatalf("unmarshaling transactionActionPayload to chaincodeActionPayload error: %s", err)
			}
			// ProposalResponsePayload field is marshalled object of
			proposalResponsePayload := &peer.ProposalResponsePayload{}
			err = proto.Unmarshal(chaincodeActionPayload.Action.ProposalResponsePayload, proposalResponsePayload)
			if err != nil {
				log.Fatalf("unmarshaling chaincodeActionPayload Action ProposalResponsePayload to proposalResponsePayload error: %s", err)
			}
			chaincodeAction := &peer.ChaincodeAction{}
			err = proto.Unmarshal(proposalResponsePayload.Extension, chaincodeAction)
			if err != nil {
				log.Fatalf("unmarshaling proposalResponsePayload Extension to chaincodeAction error: %s", err)
			}
			chaincodeEvent := &peer.ChaincodeEvent{}
			err = proto.Unmarshal(chaincodeAction.Events, chaincodeEvent)
			if err != nil {
				log.Fatalf("unmarshaling chaincodeAction Events to chaincodeEvent error: %s", err)
			}
			var eventName string
			if chaincodeEvent.EventName != "" {
				eventName = chaincodeEvent.EventName
			} else {
				eventName = getEventByParsing(chaincodeActionPayload.ChaincodeProposalPayload)
			}
			log.Printf("#################### Block event : %v ########### ", eventName)
			log.Println("#############################################################")
			/////////////////////
			//log.Println(chaincodeEndorsedActionPayload.String())
			//log.Println(proposalResponsePayload.String()) //fruit/buy sell 등등 나옴
			//log.Println(chaincodeEvent.String())
		}
	}
}

func getEventByParsing(ChaincodeProposalPayload []byte) string {
	chaincodeAction := &peer.ChaincodeAction{}
	err := proto.Unmarshal(ChaincodeProposalPayload, chaincodeAction)
	if err != nil {
		log.Fatalf("unmarshaling Chaincode Action Payload error: %s", err)
	}
	chaincodeResults := &peer.ChaincodeAction{}
	err = proto.Unmarshal(chaincodeAction.Results, chaincodeResults)
	if err != nil {
		log.Fatalf("unmarshaling Chaincode Action Payload error: %s", err)
	}
	eventName := parseEvent(chaincodeResults.String())
	return eventName
}

func parseEvent(str string) string {
	firstIndex := strings.Index(str, "fruit")
	lastIndex := strings.LastIndex(str, `\n`)
	return str[firstIndex:lastIndex]
}
