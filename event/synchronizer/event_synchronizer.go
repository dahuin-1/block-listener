package main

import (
	"flag"
	"github.com/cc-ping-listener/unmarshalers"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	mspctx "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

func main() {
	channelProvider, err := getChannelProvider()
	if err != nil {
		log.Fatalf("failed to get Channel Provider, err: %s", err)
	}
	ledgerClient, err := ledger.New(channelProvider)
	if err != nil {
		log.Fatalf("failed to return Client instance, err: %s", err)
	}
	blockchainInfo, err := ledgerClient.QueryInfo()
	if err != nil {
		log.Fatalf("unable to get fs.LedgerQueryInfo, err: %s", err)
		return
	}
	currentBlockHeight := blockchainInfo.BCI.Height
	startBlock := flag.String("startBlock", "", "set start block number if needed")
	endBlock := flag.String("endBlock", "", "set end block number if needed")
	flag.Parse()
	var startBlockNum, endBlockNum uint64
	if *startBlock == "" {
		startBlockNum = 1
	} else {
		startBlockNum, err = strconv.ParseUint(*startBlock, 10, 64)
	}
	if err != nil {
		log.Fatalf("failed to convert string startBlock to int, err: %s", err)
	}
	if *endBlock == "" {
		endBlockNum = currentBlockHeight
	} else {
		endBlockNum, err = strconv.ParseUint(*endBlock, 10, 64)
	}
	if err != nil {
		log.Fatalf("failed to convert string endBlock to int, err: %s", err)
	}
	if startBlockNum == currentBlockHeight {
		log.Fatal("nothing to sync on this batch")
	}
	if endBlockNum > currentBlockHeight {
		log.Fatal("endBlockNum should be smaller than currentBlockHeight")
	}
	if startBlockNum >= endBlockNum {
		log.Fatal("endBlockNum be bigger than startBlockNum")
	}
	////////////////
	//block, err := ledgerClient.QueryBlock(1)
	//if err != nil {
	//	log.Fatalf("failed to query Block, err: %s", err)
	//}
	//log.Println(block)
	///////////////////////
	//loop until the last block on blockchain
	for i := startBlockNum; i < endBlockNum; i++ {
		log.Printf("=================================== Sync on block number: %d ===================================", i)
		block, err := ledgerClient.QueryBlock(i)
		if err != nil {
			log.Fatalf("failed to query Block, err: %s", err)
		}
		if block == nil {
			log.Printf("null block : %d ", i)
		}
		blockData := block.Data.Data
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
		if chaincodeEvent.EventName == "" {
			log.Println("event did not happen")
		} else {
			log.Println("############################################################")
			log.Println("###################### Received event ######################")
			log.Printf("################### BlockNum : %d ##########################", block.Header.Number)
			log.Printf("#################### Block event : %v ########### ", chaincodeEvent.EventName)
			log.Printf("#################### Block info - block %v ########### ", chaincodeAction.String())
			log.Println("#############################################################")
		}
	}
}
