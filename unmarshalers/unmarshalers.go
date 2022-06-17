package unmarshalers

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
)

func GetChaincodeAction(data []byte) (*peer.ChaincodeAction, error) {
	chaincodeAction := &peer.ChaincodeAction{}
	if err := proto.Unmarshal(data, chaincodeAction); err != nil {
		return nil, err
	}
	return chaincodeAction, nil
}

func GetChaincodeActionPayload(data []byte) (*peer.ChaincodeActionPayload, error) {
	chaincodeActionPayload := &peer.ChaincodeActionPayload{}
	if err := proto.Unmarshal(data, chaincodeActionPayload); err != nil {
		return nil, err
	}
	return chaincodeActionPayload, nil
}

func GetChaincodeEvent(data []byte) (*peer.ChaincodeEvent, error) {
	chaincodeEvent := &peer.ChaincodeEvent{}
	if err := proto.Unmarshal(data, chaincodeEvent); err != nil {
		return nil, err
	}
	return chaincodeEvent, nil
}

func GetChaincodeResults(data []byte) (*peer.ChaincodeAction, error) {
	chaincodeResults := &peer.ChaincodeAction{}
	if err := proto.Unmarshal(data, chaincodeResults); err != nil {
		return nil, err
	}
	return chaincodeResults, nil
}

func GetEnvelopeFromBlock(data []byte) (*common.Envelope, error) {
	env := &common.Envelope{}
	if err := proto.Unmarshal(data, env); err != nil {
		return nil, err
	}
	return env, nil
}

func GetPayloadFromEnv(data []byte) (*common.Payload, error) {
	payload := &common.Payload{}
	if err := proto.Unmarshal(data, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func GetProposalResponsePayload(data []byte) (*peer.ProposalResponsePayload, error) {
	proposalResponsePayload := &peer.ProposalResponsePayload{}
	if err := proto.Unmarshal(data, proposalResponsePayload); err != nil {
		return nil, err
	}
	return proposalResponsePayload, nil
}

func GetTransaction(data []byte) (*peer.Transaction, error) {
	transaction := &peer.Transaction{}
	if err := proto.Unmarshal(data, transaction); err != nil {
		return nil, err
	}
	return transaction, nil
}
