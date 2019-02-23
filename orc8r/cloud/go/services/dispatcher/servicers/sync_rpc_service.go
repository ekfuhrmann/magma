/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/dispatcher/broker"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// heartBeat from cloud to gateway
const heartBeatInterval = time.Minute

type SyncRPCService struct {
	// the host at which this service instance is running on
	hostName string
	broker   broker.GatewayRPCBroker
}

func NewSyncRPCService(hostName string, broker broker.GatewayRPCBroker) (*SyncRPCService, error) {

	return &SyncRPCService{
		hostName: hostName,
		broker:   broker,
	}, nil
}

// support for backward compatibility
func (srv *SyncRPCService) SyncRPC(stream protos.SyncRPCService_SyncRPCServer) error {
	return srv.EstablishSyncRPCStream(stream)
}

// rpc call that will be called by gateways
// SyncRPC establishes a bidirectional stream between the gateway and the cloud,
// and the streams will close if it returns.
// every active connection will run this function in its own goroutine.
func (srv *SyncRPCService) EstablishSyncRPCStream(stream protos.SyncRPCService_EstablishSyncRPCStreamServer) error {
	// Check if we can get a valid Gateway identity
	gw, err := identity.GetStreamGatewayId(stream)
	if err != nil {
		return err
	}
	if gw == nil || len(gw.HardwareId) == 0 {
		return status.Errorf(codes.PermissionDenied, "Gateway hardware id is nil")
	}
	return srv.serveGwId(stream, gw.HardwareId)
}

type streamCoordinator struct {
	GwID    string
	ErrChan chan error
	Wg      *sync.WaitGroup
	Ctx     context.Context
	Cancel  context.CancelFunc
}

func makeStreamCoordinator(gwId string, streamCtx context.Context,
) *streamCoordinator {
	errChan := make(chan error, 1)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(streamCtx)
	return &streamCoordinator{gwId, errChan, wg, ctx, cancel}
}

// try send the err to the errChan,
// if nobody is listening on the errChan, log the error and return regardless
func (streamCoordinator *streamCoordinator) sendErrOrLog(err error) {
	select {
	case streamCoordinator.ErrChan <- err:
		return
	case <-time.After(time.Second):
		if err == nil {
			glog.V(2).Infof("Received EOF from client, return\n")
		} else {
			glog.Errorf(err.Error())
		}
		return
	}
}

// called directly by test service
func (srv *SyncRPCService) serveGwId(stream protos.SyncRPCService_EstablishSyncRPCStreamServer, gwId string) error {
	coordinator := makeStreamCoordinator(gwId, stream.Context())
	queue := srv.broker.InitializeGateway(gwId)
	glog.V(2).Infof("Initialized gateway for hwId %v\n", gwId)
	coordinator.Wg.Add(1)
	go srv.receiveFromStream(stream, coordinator)
	coordinator.Wg.Add(1)
	go srv.sendToStream(stream, queue, coordinator)
	// wait on err returned from either sendToStream or receiveFromStream goroutines
	err := <-coordinator.ErrChan
	if err == nil {
		glog.V(2).Infof(
			"SyncRPC return for %v due to client sending EOF\n", gwId)
	} else {
		glog.Infof("SyncRPC error for %v: %v\n", gwId, err)
	}
	coordinator.Cancel()
	coordinator.Wg.Wait()
	srv.broker.CleanupGateway(gwId)
	glog.V(2).Infof("Cleaned up gateway for hwId %v\n", gwId)
	return err
}

func (srv *SyncRPCService) sendToStream(
	stream protos.SyncRPCService_EstablishSyncRPCStreamServer, queue chan *protos.SyncRPCRequest, coordinator *streamCoordinator) {
	defer coordinator.Wg.Done()
	for {
		select {
		case <-coordinator.Ctx.Done():
			coordinator.sendErrOrLog(fmt.Errorf(
				"context cancelled in sendToStream: %v\n",
				coordinator.Ctx.Err()))
			return
		case <-time.After(heartBeatInterval):
			glog.V(2).Infof("sending heartBeat to hwId %v\n",
				coordinator.GwID)
			err := stream.Send(&protos.SyncRPCRequest{HeartBeat: true})
			if err != nil {
				errMsg := fmt.Sprintf("sendHeartBeat err: %v\n", err)
				coordinator.sendErrOrLog(errors.New(errMsg))
				return
			}
		case reqToSend, ok := <-queue:
			if !ok {
				err := fmt.Errorf("Queue is closed for hwId %v\n", coordinator.GwID)
				coordinator.sendErrOrLog(err)
				return
			}
			if reqToSend != nil {
				glog.V(2).Infof(
					"sending req to stream for hwId %v\n",
					coordinator.GwID)
				err := stream.Send(reqToSend)
				if err != nil {
					errMsg := fmt.Sprintf("sendToStream err: %v\n", err)
					coordinator.sendErrOrLog(errors.New(errMsg))
					return
				}
			}
		}
	}
}

func (srv *SyncRPCService) receiveFromStream(
	stream protos.SyncRPCService_EstablishSyncRPCStreamServer, coordinator *streamCoordinator) {
	defer coordinator.Wg.Done()
	for {
		//recv() can be cancelled via ctx
		syncRPCResp, err := RecvWithContext(coordinator.Ctx, func() (*protos.SyncRPCResponse, error) { return stream.Recv() })
		if err == io.EOF {
			coordinator.sendErrOrLog(nil)
			return
		} else if err != nil {
			errMsg := fmt.Sprintf("receiveFromStream err: %v\n", err)
			coordinator.sendErrOrLog(errors.New(errMsg))
			return
		} else {
			glog.V(2).Infof("processing response for hwId %v\n",
				coordinator.GwID)
			err := srv.processSyncRPCResp(syncRPCResp, coordinator.GwID)
			if err != nil {
				errMsg := fmt.Sprintf("procesSyncRPCResp err: %v\n", err)
				coordinator.sendErrOrLog(errors.New(errMsg))
				return
			}
		}
	}
}

// returning err indicates to end the bidirectional stream
func (srv *SyncRPCService) processSyncRPCResp(resp *protos.SyncRPCResponse, hwId string) error {
	if resp.HeartBeat {
		err := directoryd.UpdateHostNameByHwId(hwId, srv.hostName)
		if err != nil {
			// cannot persist <gwId, hostName> so nobody can send things to this gateway use the stream,
			// therefore return err to end the stream
			return err
		}
	} else if resp.ReqId > 0 {
		err := srv.broker.ProcessGatewayResponse(resp)
		if err != nil {
			// no need to end the stream, just log the error
			glog.Errorf("err processing gateway response: %v\n", err)
		}
	} else {
		glog.Errorf("Cannot send a non-heartbeat with invalid ReqId\n")

	}
	return nil
}

type WrappedSyncResponse struct {
	Resp *protos.SyncRPCResponse
	Err  error
}

// RecvWithContext runs f and returns its error.  If the context is cancelled or
// times out first, it returns the context's error instead.
// See https://github.com/grpc/grpc-go/issues/1229#issuecomment-300938770
func RecvWithContext(ctx context.Context, f func() (*protos.SyncRPCResponse, error)) (*protos.SyncRPCResponse, error) {
	wrappedRespchan := make(chan WrappedSyncResponse, 1)
	go func() {
		resp, err := f()
		wrappedRespchan <- WrappedSyncResponse{resp, err}
		close(wrappedRespchan)
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case wrappedResp := <-wrappedRespchan:
		return wrappedResp.Resp, wrappedResp.Err
	}
}
