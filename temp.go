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
