/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package mconfig provides gateway Go support for cloud managed configuration (mconfig)
package mconfig

import (
	"errors"

	"magma/orc8r/cloud/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

func GetServiceConfigs(service string, result proto.Message) error {
	current := localConfig.Load().(*protos.GatewayConfigs)
	anyCfg, found := current.ConfigsByKey[service]
	if !found {
		return errors.New("No configs found for service: " + service)
	}
	return ptypes.UnmarshalAny(anyCfg, result)
}

func GetGatewayConfigs() *protos.GatewayConfigs {
	return localConfig.Load().(*protos.GatewayConfigs)
}
