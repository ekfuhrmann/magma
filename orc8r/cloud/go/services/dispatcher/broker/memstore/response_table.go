/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package memstore

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"magma/orc8r/cloud/go/protos"

	"github.com/golang/glog"
)

type ResponseTable interface {
	InitializeResponse() (chan *protos.GatewayResponse, uint32)
	SendResponse(*protos.SyncRPCResponse) error
}

type ResponseTableImpl struct {
	respByReqId *sync.Map // <uint32, chan *protos.GatewayResponse>
	// strictly increase, making sure all reqIds will be unique.
	reqIdCounter uint32
	timeout      time.Duration
}

func NewResponseTable(timeout time.Duration) *ResponseTableImpl {
	return &ResponseTableImpl{respByReqId: &sync.Map{}, timeout: timeout}
}

// create a request id to bind to the GatewayResponse channel, so when SyncRPCResponse comes back,
// it can be written to the corresponding GatewayResponse channel identified by the request id.
func (table *ResponseTableImpl) InitializeResponse() (chan *protos.GatewayResponse, uint32) {
	reqId := generateReqId(&table.reqIdCounter)
	respChan := make(chan *protos.GatewayResponse)
	table.respByReqId.Store(reqId, respChan)
	return respChan, reqId
}

// send response to the corresponding response channel
func (table *ResponseTableImpl) SendResponse(resp *protos.SyncRPCResponse) error {
	if resp == nil {
		return errors.New("cannot send nil SyncRPCResponse")
	}
	respChanVal, ok := table.respByReqId.Load(resp.ReqId)
	if ok {
		respChan := respChanVal.(chan *protos.GatewayResponse)
		if resp.RespBody == nil {
			glog.Errorf("Nil response body received, forward to httpServer anyways\n")
		}
		select {
		case respChan <- resp.RespBody:
			return nil
		case <-time.After(table.timeout):
			// give up sending, close the channel
			close(respChan)
			return errors.New("sendResponse timed out as respChan is not being actively waited on")
		}
	} else {
		return errors.New(fmt.Sprintf("No response channel found for reqId %v\n", resp.ReqId))
	}
}

func generateReqId(counter *uint32) uint32 {
	return atomic.AddUint32(counter, 1)
}
