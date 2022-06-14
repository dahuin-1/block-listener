package main

import (
	"errors"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"io/ioutil"
	"log"

	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"path/filepath"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type FabricUser struct {
	Cert       []byte
	Name       string
	PrivateKey []byte
}

const (
	cfgPath   = "/Users/dhkim/Projects/cc-ping-listener/config/network.yaml" //네트워크.야멜
	channelID = "kiesnet-dev"
	credPath  = "/Users/dhkim/Projects/kiesnet-chaincode-dev-network/crypto-config/peerOrganizations/kiesnet.dev/users"
	userName  = "dhkim"
)

func main() {
	channelProvider, err := getChannelProvider()
	if err != nil {
		log.Fatalf("failde to get channelProvider, err: %s", err)
	}
	client, err := event.New(channelProvider, event.WithBlockEvents())
	if err != nil {
		log.Fatalf("failed to get event.Client, err: %s", err)
	}
	registration, notifier, err := client.RegisterFilteredBlockEvent()
	if err != nil {
		log.Fatalf("failed to RegisterFilteredBlockEvent, err: %s", err)
	}
	defer client.Unregister(registration)

	for {
		log.Printf("♬♬♬♬♬♬♬♬♬♬♬♫♫♫♫♫♫♫♫♫♫♫♫♫♫♫♫🎹👂🏻listen👂🏻🎹♬♬♬♬♬♬♬♬♬♬♬♫♫♫♫♫♫♫♫♫♫♫♫♫♫♫♫")
		select {
		case e := <-notifier:
			log.Printf("#################### Block: %v ########################", e.FilteredBlock)
		case <-time.After(time.Second * 30):
			log.Printf("#################### NO event NO block ##################")
		}
	}
}

func getChannelProvider() (context.ChannelProvider, error) {
	fabricUser, err := NewFabricUser(userName)
	if err != nil {
		return nil, err
	}

	cert := fabricUser.Cert
	prvtKey := fabricUser.PrivateKey

	var configPath core.ConfigProvider
	configPath = config.FromFile(cfgPath)
	sdk, err := fabsdk.New(configPath) //sdk객체  //sdk, err := fabsdk.New(config.FromFile(configPath))
	if err != nil {
		return nil, err
	}
	mspClient, err := mspclient.New(sdk.Context())
	if err != nil {
		return nil, err
	}
	signingIdentity, err := mspClient.CreateSigningIdentity(mspctx.WithCert(cert), mspctx.WithPrivateKey(prvtKey))
	if err != nil {
		return nil, err
	}
	channelProvider := sdk.ChannelContext(channelID, fabsdk.WithIdentity(signingIdentity))
	return channelProvider, nil
}

func NewFabricUser(name string) (*FabricUser, error) {
	mspPath := filepath.Join(credPath, name, "msp")
	certPath := filepath.Join(mspPath, "signcerts", "cert.pem")
	cert, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	keyDir := filepath.Join(mspPath, "keystore")
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return nil, err
	}
	if len(files) != 1 {
		return nil, errors.New("keystore must have one value")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return &FabricUser{Name: name, Cert: cert, PrivateKey: key}, nil
}
