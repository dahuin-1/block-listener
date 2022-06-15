package main

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

type FabricUser struct {
	Cert       []byte
	Name       string
	PrivateKey []byte
}

const (
	channelID  = "kiesnet-dev"
	configPath = "/Users/dhkim/Projects/cc-ping-listener/config/network.yaml"
	credPath   = "/Users/dhkim/Projects/kiesnet-chaincode-dev-network/crypto-config/peerOrganizations/kiesnet.dev/users"
	userName   = "dhkim"
)

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
		log.Fatalf("failed to register Filtered Block Event, err: %s", err)
	}
	defer client.Unregister(registration)
	var blockNum uint64
	for {
		log.Printf("🎹👂🎹👂🎹👂🎹👂🎹👂🏻listen👂🏻🎹👂🎹👂🎹👂🎹👂🎹")
		select {
		case e := <-eventChannel:
			blockNum = e.Block.Header.Number
			log.Println("#########################################################")
			log.Println("###################### Received event ######################")
			log.Printf("################### BlockNum : %d ######################", blockNum)
			log.Printf("#################### Block info: %v ########################", e.Block)
			log.Println("#########################################################")
		case <-time.After(time.Second * 10):
			log.Println("#########################################################################")
			log.Printf("#################### Event did not happen this time #####################")
			if blockNum != 0 {
				log.Printf("################### Block number until now : %d #########################", blockNum)
			}
			log.Println("#########################################################################")
		}
	}
}

func getChannelProvider() (context.ChannelProvider, error) {
	fabricUser, err := setFabricUser(userName)
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

func setFabricUser(name string) (*FabricUser, error) {
	mspPath := filepath.Join(credPath, name, "msp")
	certPath := filepath.Join(mspPath, "signcerts", "cert.pem")

	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}

	keyStore := filepath.Join(mspPath, "keystore")
	keys, err := os.ReadDir(keyStore)
	if err != nil {
		return nil, err
	}

	keyPath := filepath.Join(keyStore, keys[0].Name())
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	return &FabricUser{Name: name, Cert: cert, PrivateKey: key}, nil
}
