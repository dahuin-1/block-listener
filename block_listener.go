package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
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

func getEnvelopeFromBlock(data []byte) (*common.Envelope, error) {
	var err error
	env := &common.Envelope{}
	if err = proto.Unmarshal(data, env); err != nil {
		return nil, err
	}
	return env, nil
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
			//log.Printf("#################### Block event : %v ########### ", e.Block.Data)
			blockData := e.Block.Data.Data
			//First Get the Envelope from the BlockData
			envelope, err := getEnvelopeFromBlock(blockData[0])
			if err != nil {
				log.Fatalf("unmarshaling Envelope error: %s", err)
			}
			//Retrieve the Payload from the Envelope
			payload := &common.Payload{}
			err = proto.Unmarshal(envelope.Payload, payload)
			if err != nil {
				log.Fatalf("unmarshaling Payload error: %s", err)
			}
			//Read the Transaction from the Payload Data
			transaction := &peer.Transaction{}
			err = proto.Unmarshal(payload.Data, transaction)
			if err != nil {
				log.Fatalf("unmarshaling Payload Transaction error: %s", err)
			}
			// Payload field is marshalled object of ChaincodeActionPayload
			chaincodeActionPayload := &peer.ChaincodeActionPayload{}
			err = proto.Unmarshal(transaction.Actions[0].Payload, chaincodeActionPayload)
			if err != nil {
				log.Fatalf("unmarshaling Chaincode Action Payload error: %s", err)
			}
			chaincodeEndorsedActionPayload := &peer.ProposalResponsePayload{}
			err = proto.Unmarshal(chaincodeActionPayload.Action.ProposalResponsePayload, chaincodeEndorsedActionPayload)
			if err != nil {
				log.Fatalf("unmarshaling chaincode Endorsed Action Payload error: %s", err)
			}
			proposalResponsePayload := &peer.ChaincodeAction{}
			err = proto.Unmarshal(chaincodeEndorsedActionPayload.Extension, proposalResponsePayload)
			if err != nil {
				log.Fatalf("unmarshaling Chaincode Action Payload error: %s", err)
			}
			chaincodeEvent := &peer.ChaincodeEvent{}
			err = proto.Unmarshal(proposalResponsePayload.Events, chaincodeEvent)
			//log.Println(proposalResponsePayload.String()) //fruit/buy sell 등등 나옴
			if err != nil {
				log.Fatalf("unmarshaling Chaincode Action Payload error: %s", err)
			}
			log.Printf("#################### Block event : %v ########### ", chaincodeEvent.EventName)
			log.Println("#############################################################")
		}
	}
}
