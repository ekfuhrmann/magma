/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package identity

import (
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/certifier"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const CLIENT_CERT_SN_KEY = "x-magma-client-cert-serial"

// NewOperator returns Operator identity corresponding to given opId
// see protos/identity_helper.go
func NewOperator(opId string) *protos.Identity {
	return protos.NewOperatorIdentity(opId)
}

// NewNetwork returns Network identity corresponding to given networkId
// see protos/identity_helper.go
func NewNetwork(networkId string) *protos.Identity {
	return protos.NewNetworkIdentity(networkId)
}

// NewGateway returns Gateway identity corresponding to given hardware network &
// logical gateway IDs
func NewGateway(hwId, networkId, logicalId string) *protos.Identity {
	return protos.NewGatewayIdentity(hwId, networkId, logicalId)
}

// NewGatewayWildcard returns Network Wildcard identity
// see protos/identity_helper.go
func NewGatewayWildcard() *protos.Identity {
	return protos.NewGatewayWildcardIdentity()
}

// NewOperatorWildcard returns Operator Wildcard identity
// see protos/identity_helper.go
func NewOperatorWildcard() *protos.Identity {
	return protos.NewOperatorWildcardIdentity()
}

// NewNetworkWildcard returns Network Wildcard identity
// see protos/identity_helper.go
func NewNetworkWildcard() *protos.Identity {
	return protos.NewNetworkWildcardIdentity()
}

// IsOperator Checks if it's an Identity of Operator and returns true if it is
func IsOperator(id *protos.Identity) bool {
	if id != nil {
		_, ok := id.Value.(*protos.Identity_Operator)
		return ok
	}
	return false
}

// IsGateway Checks if it's an Identity of Gateway and returns true if it is
func IsGateway(id *protos.Identity) bool {
	if id != nil {
		_, ok := id.Value.(*protos.Identity_Gateway_)
		return ok
	}
	return false
}

//GetStreamGatewayId returns a valid, non nil Gateway identity based on the
//stream's metadata CTX or error if no GW Identity can be found/verified
func GetStreamGatewayId(stream grpc.Stream) (*protos.Identity_Gateway, error) {
	ctx := stream.Context()
	if ctx == nil {
		msg := "Missing Stream Context"
		glog.Errorf(msg)
		return nil, status.Error(codes.Unauthenticated, msg)
	}
	ctxMetadata, ok := metadata.FromIncomingContext(ctx)
	if !ok || ctxMetadata == nil {
		glog.Errorf("Missing Metadata from Stream Ctx: %+v", ctx)
		return nil, status.Error(codes.Unauthenticated, "Missing CTX Metadata")
	}
	// Find the caller's identity
	snlist, snok := ctxMetadata[CLIENT_CERT_SN_KEY]
	if !snok || len(snlist) == 0 {
		err := status.Error(codes.Unauthenticated, "Missing Certificate SN")
		glog.Errorf("%s in stream CTX Metadata: %+v", err, ctxMetadata)
		return nil, err
	}
	serialNum := snlist[0]
	id, err := certifier.GetVerifiedCertificateIdentity(serialNum)
	if err != nil {
		glog.Errorf(err.Error())
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	gwIdentity := id.GetGateway()
	if gwIdentity == nil {
		glog.Errorf(
			"Identity (%s) of Cert SN %s from metadata %+v is not a Gateway",
			id.HashString(), serialNum, ctxMetadata)

		return nil, status.Error(codes.PermissionDenied, "Invalid Identity Type")
	}
	return gwIdentity, nil
}
