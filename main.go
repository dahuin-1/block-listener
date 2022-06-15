package main

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
	"os"
	"path/filepath"
	// "time"
)

/*
1. 네트워크 컨피그 설정
2. sdk 객체
3. channel client
4. block listening
*/

type User struct {
	Cert []byte
	//Name       string
	PrivateKey []byte
}

const (
	chainCodeID = "ping"
	channelID   = "kiesnet-dev"
	configPath  = "/Users/dhkim/Projects/cc-ping-listener/config/network.yaml"
	credPath    = "/Users/dhkim/Projects/kiesnet-chaincode-dev-network/crypto-config/peerOrganizations/kiesnet.dev/users"
)

func getChannelProvider() (context.ChannelProvider, error) {
	fabricUser, err := setUser()
	if err != nil {
		return nil, err
	}

	networkConfig := config.FromFile(configPath) //네트워크컨피그설정

	sdk, err := fabsdk.New(networkConfig) //sdk객체를 얻음 //sdk, err := fabsdk.New(config.FromFile(configPath))
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

func main() {
	channelProvider, err := getChannelProvider()
	if err != nil {
		log.Fatalf("failed to get Channel Provider, err: %s", err)
	}
	client, err := event.New(channelProvider, event.WithBlockEvents())
	if err != nil {
		log.Fatalf("failed to return Client instance, err: %s", err)
	}
	registration, eventChannel, err := client.RegisterChaincodeEvent(chainCodeID, "fruit/soldout|fruit/restock")
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
			log.Printf("################### BlockNum : %d ##########################", e.BlockNumber)
			log.Printf("#################### Block event : %v ########### ", e.EventName)
			log.Println("#############################################################")
		}
	}
}
