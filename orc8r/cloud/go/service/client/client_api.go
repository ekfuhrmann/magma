/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package client

import (
	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/registry"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func getClient(service string) (protos.Service303Client, *grpc.ClientConn, error) {
	conn, err := registry.GetConnection(service)
	if err != nil {
		initErr := errors.NewInitError(err, "SERVICE303")
		glog.Error(initErr)
		return nil, nil, initErr
	}
	return protos.NewService303Client(conn), conn, nil
}

func Service303GetServiceInfo(service string) (*protos.ServiceInfo, error) {
	client, conn, err := getClient(service)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return client.GetServiceInfo(context.Background(), new(protos.Void))
}

func Service303GetMetrics(service string) (*protos.MetricsContainer, error) {
	client, conn, err := getClient(service)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return client.GetMetrics(context.Background(), new(protos.Void))
}

func Service303StopService(service string) error {
	client, conn, err := getClient(service)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.StopService(context.Background(), new(protos.Void))
	return err
}

func Service303SetLogLevel(service string, in *protos.LogLevelMessage) error {
	client, conn, err := getClient(service)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.SetLogLevel(context.Background(), in)
	return err
}

func Service303SetLogVerbosity(service string, in *protos.LogVerbosity) error {
	client, conn, err := getClient(service)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = client.SetLogVerbosity(context.Background(), in)
	return err
}
