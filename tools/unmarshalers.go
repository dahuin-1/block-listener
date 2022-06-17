package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
)

func GetEnvelopeFromBlock(data []byte) (*common.Envelope, error) {
	env := &common.Envelope{}
	err := proto.Unmarshal(data, env)
	if err != nil {
		return nil, err
	}
	return env, nil
}

func GetPayloadFromEnv(data []byte) (*common.Payload, error) {
	payload := &common.Payload{}
	err := proto.Unmarshal(data, payload)
	if err != nil {
		return nil, err
	}
	return payload, nil
}
