package main

//
//import (
//	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
//	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
//	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
//	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
//	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
//	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
//	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
//	fabcfg "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
//	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
//	"github.com/kataras/golog"
//)
//
////type Channel struct {
//	sdk      *fabsdk.FabricSDK
//	identity mspctx.SigningIdentity
//	ctx      context.ChannelProvider
//}
//
//func main() {
//	//channel client
//	var cfg core.ConfigProvider
//	cfg = fabcfg.FromFile("./kiesnet-chaincode-dev-network/configtx.yaml")💃🏻🎹
//	sdk, err := fabsdk.New(cfg) //sdk객체
//	if err != nil {
//		return
//	}
//	client, err := mspclient.New(sdk.Context()) //getSigningIdentity
//	if err != nil {
//		return
//	}
//	//msp
//	si, err := client.CreateSigningIdentity(mspctx.WithCert([]byte(cert)), mspctx.WithPrivateKey([]byte(priKey)))
//	if err != nil {
//		return
//	}
//	// channel provider
//	//ctxOpts := append([]fabsdk.ContextOption{fabsdk.WithIdentity(si)}, ctxOpts...)
//	ctx := sdk.ChannelContext(channelID, []fabsdk.ContextOption{fabsdk.WithIdentity(si)}...)
//	//return
//	channel := &Channel{
//		sdk:      sdk,
//		identity: si,
//		ctx:      ctx,
//	}
//	e := &event.Client{}
//	err = listenBlockEvent(e)
//	if err != nil {
//		return
//	}
//	//client를 이용해서 block listening
//
//	//channel에 event가 도착하면
//	//for loop~~~
//	//event 발견하면 print
//	/*
//		network config 설정 -> 얠 이용해서 sdk 객체를 얻는다.
//	*/
//}
//
//func listenBlockEvent(client *event.Client) error {
//	registration, notifier, err := client.RegisterBlockEvent()
//	if err != nil {
//		return err
//	}
//	defer client.Unregister(registration)
//
//	for {
//		select {
//		case e := <-notifier:
//			blc := e.Block
//			blockNum := blc.Header.Number
//			if err != nil {
//				return err
//			}
//			golog.Infof("event %d listen", blockNum)
//		}
//	}
//}
/*
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
	user       = "dhkim"
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
			//case <-time.After(time.Second * 10):
			//	log.Printf("#################### NO event NO block ##################")
		}
	}
}

func getChannelProvider() (context.ChannelProvider, error) {
	fabricUser, err := setFabricUser(user)
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
	client, err := mspclient.New(sdk.Context())
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
*/
