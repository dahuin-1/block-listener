package main

import (
	"errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"io/ioutil"
	"log"
	"path/filepath"
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
		log.Fatalf("failde to getChannelProvider, err: %s", err)
	}
	client, err := event.New(channelProvider, event.WithBlockEvents())
	if err != nil {
		log.Fatalf("failed to return Client instance, err: %s", err)
	}
	registration, eventChannel, err := client.RegisterFilteredBlockEvent()
	if err != nil {
		log.Fatalf("failed to Register Filtered Block Event, err: %s", err)
	}
	defer client.Unregister(registration)

	for {
		log.Printf("♬♬♬♬♬♬♬♬♬♬♬♫♫♫♫♫♫♫♫♫♫♫♫♫♫♫♫🎹👂🏻listen👂🏻🎹♬♬♬♬♬♬♬♬♬♬♬♫♫♫♫♫♫♫♫♫♫♫♫♫♫♫♫")
		select {
		case e := <-eventChannel:
			log.Printf("#################### Block: %v ########################", e.FilteredBlock)
		}
	}
}

func getChannelProvider() (context.ChannelProvider, error) {
	fabricUser, err := setFabricUser(userName)
	if err != nil {
		return nil, err
	}

	cert := fabricUser.Cert
	privateKey := fabricUser.PrivateKey

	networkConfig := config.FromFile(configPath) //네트워크컨피그설정
	sdk, err := fabsdk.New(networkConfig)        //sdk객체를 얻음 //sdk, err := fabsdk.New(config.FromFile(configPath))
	if err != nil {
		return nil, err
	}
	client, err := mspclient.New(sdk.Context()) //sdk 객체를 이용해서 channel client 생성
	if err != nil {
		return nil, err
	}
	signingIdentity, err := client.CreateSigningIdentity(mspctx.WithCert(cert), mspctx.WithPrivateKey(privateKey))
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
	keys, err := ioutil.ReadDir(keyStore)
	if err != nil {
		return nil, err
	}
	if len(keys) != 1 {
		return nil, errors.New("keystore must have one value")
	}
	keyPath := filepath.Join(keyStore, keys[0].Name())
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return &FabricUser{Name: name, Cert: cert, PrivateKey: key}, nil
}
