/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"context"
	"testing"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/health/servicers"
	"magma/feg/cloud/go/services/health/storage/mocks"
	"magma/feg/cloud/go/services/health/test_utils"
	orcprotos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/magmad"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	magmad_test_init "magma/orc8r/cloud/go/services/magmad/test_init"

	"github.com/stretchr/testify/assert"
)

func TestHealthServer_GetHealth(t *testing.T) {
	healthStore := &mocks.HealthStorage{}
	clusterStore := &mocks.ClusterStorage{}
	service := servicers.NewTestHealthServer(healthStore, clusterStore)

	gwStatusReq := &protos.GatewayStatusRequest{
		NetworkId: test_utils.TestFegNetwork,
		LogicalId: test_utils.TestFegLogicalId1,
	}
	healthyReq := test_utils.GetHealthyRequest()
	healthStore.On("GetHealth", test_utils.TestFegNetwork, test_utils.TestFegLogicalId1).
		Return(healthyReq.HealthStats, nil).Once()

	stats, err := service.GetHealth(context.Background(), gwStatusReq)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthStatus_HEALTHY, stats.Health.Health)
	assert.Equal(t, healthyReq.HealthStats.SystemStatus, stats.SystemStatus)
	assert.Equal(t, healthyReq.HealthStats.ServiceStatus, stats.ServiceStatus)

	healthStore.AssertExpectations(t)

	unhealthyReq := test_utils.GetUnhealthyRequest()
	healthStore.On("GetHealth", test_utils.TestFegNetwork, test_utils.TestFegLogicalId1).
		Return(unhealthyReq.HealthStats, nil).Once()

	stats, err = service.GetHealth(context.Background(), gwStatusReq)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthStatus_UNHEALTHY, stats.Health.Health)
	assert.Equal(t, unhealthyReq.HealthStats.SystemStatus, stats.SystemStatus)
	assert.Equal(t, unhealthyReq.HealthStats.ServiceStatus, stats.ServiceStatus)

	healthStore.AssertExpectations(t)
}

// Test that a single feg will always remain ACTIVE
func TestHealthServer_UpdateHealth_SingleGateway(t *testing.T) {
	magmad_test_init.StartTestService(t)
	healthStore := &mocks.HealthStorage{}
	clusterStore := &mocks.ClusterStorage{}
	service := servicers.NewTestHealthServer(healthStore, clusterStore)

	testNetworkID, err := magmad.RegisterNetwork(
		&magmad_protos.MagmadNetworkRecord{Name: test_utils.TestFegNetwork},
		test_utils.TestFegNetwork,
	)
	assert.NoError(t, err)

	logicalId, err := magmad.RegisterGatewayWithId(
		testNetworkID,
		&magmad_protos.AccessGatewayRecord{
			HwId: &orcprotos.AccessGatewayID{Id: test_utils.TestFegHwId1},
		},
		test_utils.TestFegLogicalId1,
	)
	assert.NoError(t, err)
	assert.Equal(t, test_utils.TestFegLogicalId1, logicalId)

	// Use Healthy Request metrics
	healthyRequest := test_utils.GetHealthyRequest()
	clusterState := getClusterState(logicalId)

	// Ensure FeG is active and receives SYSTEM_UP
	healthStore.On("UpdateHealth", testNetworkID, logicalId, healthyRequest.HealthStats).
		Return(nil).Once()
	clusterStore.On("DoesKeyExist", testNetworkID, testNetworkID).Return(false, nil)
	clusterStore.On("UpdateClusterState", testNetworkID, testNetworkID, logicalId).Return(nil)
	clusterStore.On("GetClusterState", testNetworkID, testNetworkID).Return(clusterState, nil)

	res, err := service.UpdateHealth(context.Background(), healthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_UP, res.Action)
	healthStore.AssertExpectations(t)
	clusterStore.AssertExpectations(t)

	// Ensure we stay active with only one feg, even if it is unhealthy
	unhealthyRequest := test_utils.GetUnhealthyRequest()

	healthStore.On("UpdateHealth", testNetworkID, logicalId, unhealthyRequest.HealthStats).
		Return(nil)
	clusterStore.On("DoesKeyExist", testNetworkID, testNetworkID).Return(true, nil)
	clusterStore.On("GetClusterState", testNetworkID, testNetworkID).Return(clusterState, nil)

	res, err = service.UpdateHealth(context.Background(), unhealthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_UP, res.Action)
	healthStore.AssertExpectations(t)
	clusterStore.AssertExpectations(t)
}

func TestHealthServer_UpdateHealth_DualFeg_HealthyActive(t *testing.T) {
	magmad_test_init.StartTestService(t)
	healthStore := &mocks.HealthStorage{}
	clusterStore := &mocks.ClusterStorage{}
	service := servicers.NewTestHealthServer(healthStore, clusterStore)

	testNetworkID, logicalId, logicalId2 := registerTwoFegs(t)

	healthyRequest := test_utils.GetHealthyRequest()
	clusterState := getClusterState(logicalId)
	healthStore.On("UpdateHealth", testNetworkID, logicalId, healthyRequest.HealthStats).Return(nil)
	clusterStore.On("DoesKeyExist", testNetworkID, testNetworkID).Return(true, nil)
	clusterStore.On("GetClusterState", testNetworkID, testNetworkID).Return(clusterState, nil)
	healthStore.On("GetHealth", testNetworkID, logicalId2).Return(healthyRequest.HealthStats, nil)

	res, err := service.UpdateHealth(context.Background(), healthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_UP, res.Action)
	healthStore.AssertExpectations(t)
	clusterStore.AssertExpectations(t)

	// Update test servicer to simulate like this request is coming from standby feg
	service.Feg1 = false
	healthyRequest2 := test_utils.GetHealthyRequest()
	healthStore.On("UpdateHealth", testNetworkID, logicalId2, healthyRequest2.HealthStats).Return(nil)
	clusterStore.On("DoesKeyExist", testNetworkID, testNetworkID).Return(true, nil)
	clusterStore.On("GetClusterState", testNetworkID, testNetworkID).Return(clusterState, nil)
	healthStore.On("GetHealth", testNetworkID, logicalId).Return(healthyRequest.HealthStats, nil)
	healthStore.On("GetHealth", testNetworkID, logicalId2).Return(healthyRequest.HealthStats, nil)

	// Standby FeG receives SYSTEM_DOWN, since active is still healthy
	res, err = service.UpdateHealth(context.Background(), healthyRequest2)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_DOWN, res.Action)
	healthStore.AssertExpectations(t)
	clusterStore.AssertExpectations(t)
}

func TestNewHealthServer_UpdateHealth_FailoverFromActive(t *testing.T) {
	magmad_test_init.StartTestService(t)
	healthStore := &mocks.HealthStorage{}
	clusterStore := &mocks.ClusterStorage{}
	service := servicers.NewTestHealthServer(healthStore, clusterStore)

	testNetworkID, logicalId, logicalId2 := registerTwoFegs(t)

	// Simulate an unhealthy active FeG, and thus a failover
	unhealthyRequest := test_utils.GetUnhealthyRequest()
	healthyRequest := test_utils.GetHealthyRequest()

	clusterState := getClusterState(logicalId)
	healthStore.On("UpdateHealth", testNetworkID, logicalId, unhealthyRequest.HealthStats).Return(nil)
	clusterStore.On("DoesKeyExist", testNetworkID, testNetworkID).Return(true, nil)
	clusterStore.On("GetClusterState", testNetworkID, testNetworkID).Return(clusterState, nil)
	healthStore.On("GetHealth", testNetworkID, logicalId2).Return(healthyRequest.HealthStats, nil)
	clusterStore.On("UpdateClusterState", testNetworkID, testNetworkID, logicalId2).Return(nil)

	res, err := service.UpdateHealth(context.Background(), unhealthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_DOWN, res.Action)
	healthStore.AssertExpectations(t)
	clusterStore.AssertExpectations(t)
}

func TestNewHealthServer_UpdateHealth_FailoverFromStandby(t *testing.T) {
	magmad_test_init.StartTestService(t)
	healthStore := &mocks.HealthStorage{}
	clusterStore := &mocks.ClusterStorage{}
	service := servicers.NewTestHealthServer(healthStore, clusterStore)

	testNetworkID, logicalId, logicalId2 := registerTwoFegs(t)

	// Update test servicer to act as though this request is coming from the standby feg
	service.Feg1 = false

	// Simulate that the active's last update was too long ago
	healthyRequestTooLongAgo := test_utils.GetHealthyRequest()
	healthyRequestTooLongAgo.HealthStats.Time = 0
	healthyRequest := test_utils.GetHealthyRequest()
	clusterState := getClusterState(logicalId)

	healthStore.On("UpdateHealth", testNetworkID, logicalId2, healthyRequest.HealthStats).Return(nil)
	clusterStore.On("DoesKeyExist", testNetworkID, testNetworkID).Return(true, nil)
	clusterStore.On("GetClusterState", testNetworkID, testNetworkID).Return(clusterState, nil)
	healthStore.On("GetHealth", testNetworkID, logicalId).Return(healthyRequestTooLongAgo.HealthStats, nil)
	clusterStore.On("UpdateClusterState", testNetworkID, testNetworkID, logicalId2).Return(nil)

	res, err := service.UpdateHealth(context.Background(), healthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_UP, res.Action)
	healthStore.AssertExpectations(t)
	clusterStore.AssertExpectations(t)
}

func TestNewHealtherServer_UpdateHealth_AllUnhealthy(t *testing.T) {
	magmad_test_init.StartTestService(t)
	healthStore := &mocks.HealthStorage{}
	clusterStore := &mocks.ClusterStorage{}
	service := servicers.NewTestHealthServer(healthStore, clusterStore)

	testNetworkID, logicalId, logicalId2 := registerTwoFegs(t)

	// Simulate that both the active and standby are unhealthy
	unhealthyRequest := test_utils.GetUnhealthyRequest()
	clusterState := getClusterState(logicalId)
	healthStore.On("UpdateHealth", testNetworkID, logicalId, unhealthyRequest.HealthStats).Return(nil)
	clusterStore.On("DoesKeyExist", testNetworkID, testNetworkID).Return(true, nil)
	clusterStore.On("GetClusterState", testNetworkID, testNetworkID).Return(clusterState, nil)
	healthStore.On("GetHealth", testNetworkID, logicalId2).Return(unhealthyRequest.HealthStats, nil)

	res, err := service.UpdateHealth(context.Background(), unhealthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_UP, res.Action)
	healthStore.AssertExpectations(t)
	clusterStore.AssertExpectations(t)

	// Update test servicer to simulate like this request is coming from standby feg
	service.Feg1 = false

	healthStore.On("UpdateHealth", testNetworkID, logicalId2, unhealthyRequest.HealthStats).Return(nil)
	clusterStore.On("DoesKeyExist", testNetworkID, testNetworkID).Return(true, nil)
	clusterStore.On("GetClusterState", testNetworkID, testNetworkID).Return(clusterState, nil)
	healthStore.On("GetHealth", testNetworkID, logicalId).Return(unhealthyRequest.HealthStats, nil)

	res, err = service.UpdateHealth(context.Background(), unhealthyRequest)
	assert.NoError(t, err)
	assert.Equal(t, protos.HealthResponse_SYSTEM_DOWN, res.Action)
	healthStore.AssertExpectations(t)
	clusterStore.AssertExpectations(t)
}

func registerTwoFegs(t *testing.T) (string, string, string) {
	testNetworkID := test_utils.RegisterNetwork(t, test_utils.TestFegNetwork)
	logicalId := test_utils.RegisterGateway(
		t,
		testNetworkID,
		test_utils.TestFegHwId1,
		test_utils.TestFegLogicalId1,
	)
	logicalId2 := test_utils.RegisterGateway(
		t,
		testNetworkID,
		test_utils.TestFegHwId2,
		test_utils.TestFegLogicalId2,
	)
	return testNetworkID, logicalId, logicalId2
}

func getClusterState(logicalID string) *protos.ClusterState {
	return &protos.ClusterState{
		ActiveGatewayLogicalId: logicalID,
		Time:                   uint64(time.Now().Unix()),
	}
}
