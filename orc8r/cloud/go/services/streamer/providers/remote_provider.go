/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package providers

import (
	"context"
	"strings"

	streamer_protos "magma/orc8r/cloud/go/services/streamer/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/any"
)

type remoteProvider struct {
	// service name of the provider
	// should always be uppercase to match service registry convention
	service string
	// stream name
	stream string
}

// NewRemoteProvider returns a stream provider that forwards its methods to the
// remote stream provider servicer.
func NewRemoteProvider(serviceName, stream string) StreamProvider {
	return &remoteProvider{service: strings.ToUpper(serviceName), stream: stream}
}

func (r *remoteProvider) GetStreamName() string {
	return r.stream
}

func (r *remoteProvider) GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error) {
	c, err := r.getProviderClient()
	if err != nil {
		return nil, err
	}
	res, err := c.GetUpdates(context.Background(), &protos.StreamRequest{
		GatewayId:  gatewayId,
		StreamName: r.GetStreamName(),
		ExtraArgs:  extraArgs,
	})
	if err != nil {
		return nil, err
	}
	return res.Updates, nil
}

func (r *remoteProvider) getProviderClient() (streamer_protos.StreamProviderClient, error) {
	conn, err := registry.GetConnection(r.service)
	if err != nil {
		initErr := merrors.NewInitError(err, r.service)
		glog.Error(initErr)
		return nil, initErr
	}
	return streamer_protos.NewStreamProviderClient(conn), nil
}