/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package gx

import (
	"fmt"
	"time"

	"magma/feg/gateway/policydb"
	"magma/lte/cloud/go/protos"

	"github.com/fiorix/go-diameter/diam"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func (rd *RuleDefinition) ToProto() *protos.PolicyRule {
	monitoringKey := ""
	if rd.MonitoringKey != nil {
		monitoringKey = *rd.MonitoringKey
	}
	var ratingGroup uint32 = 0
	if rd.RatingGroup != nil {
		ratingGroup = *rd.RatingGroup
	}
	flowList := getFlowList(rd.FlowDescriptions, rd.FlowInformations)

	var qos *protos.FlowQos
	if rd.Qos != nil {
		qos = &protos.FlowQos{}
		if rd.Qos.MaxReqBwUL != nil {
			qos.MaxReqBwUl = *rd.Qos.MaxReqBwUL
		}
		if rd.Qos.MaxReqBwDL != nil {
			qos.MaxReqBwDl = *rd.Qos.MaxReqBwDL
		}
	}

	return &protos.PolicyRule{
		Id:            rd.RuleName,
		RatingGroup:   ratingGroup,
		MonitoringKey: monitoringKey,
		Priority:      rd.Precedence,
		Redirect:      rd.getRedirectInfo(),
		FlowList:      flowList,
		Qos:           qos,
		TrackingType:  rd.getTrackingType(),
	}
}

func (rd *RuleDefinition) getTrackingType() protos.PolicyRule_TrackingType {
	if rd.MonitoringKey != nil && rd.RatingGroup != nil {
		return protos.PolicyRule_OCS_AND_PCRF
	} else if rd.MonitoringKey != nil && rd.RatingGroup == nil {
		return protos.PolicyRule_ONLY_PCRF
	} else if rd.MonitoringKey == nil && rd.RatingGroup != nil {
		return protos.PolicyRule_ONLY_OCS
	} else {
		return protos.PolicyRule_NO_TRACKING
	}
}

func (rd *RuleDefinition) getRedirectInfo() *protos.RedirectInformation {
	if rd.RedirectInformation == nil {
		return nil
	}
	return &protos.RedirectInformation{
		Support:       protos.RedirectInformation_Support(rd.RedirectInformation.RedirectSupport),
		AddressType:   protos.RedirectInformation_AddressType(rd.RedirectInformation.RedirectAddressType),
		ServerAddress: rd.RedirectInformation.RedirectServerAddress,
	}
}

func getFlowList(flowStrings []string, flowInfos []*FlowInformation) []*protos.FlowDescription {
	allFlowStrings := flowStrings[:]
	for _, info := range flowInfos {
		allFlowStrings = append(allFlowStrings, info.FlowDescription)
	}
	var flowList []*protos.FlowDescription
	for _, flowString := range allFlowStrings {
		flow, err := policydb.GetFlowDescriptionFromFlowString(flowString)
		if err != nil {
			glog.Errorf("Could not get flow for description %s : %s", flowString, err)
		} else {
			flowList = append(flowList, flow)
		}
	}
	return flowList
}

func (rar *ReAuthRequest) ToProto(imsi, sid string, policyDBClient policydb.PolicyDBClient) *protos.PolicyReAuthRequest {
	var rulesToRemove, baseNamesToRemove []string

	for _, ruleRemove := range rar.RulesToRemove {
		rulesToRemove = append(rulesToRemove, ruleRemove.RuleNames...)
		baseNamesToRemove = append(baseNamesToRemove, ruleRemove.RuleBaseNames...)
	}

	baseNameRuleIDsToRemove := policyDBClient.GetRuleIDsForBaseNames(baseNamesToRemove)
	rulesToRemove = append(rulesToRemove, baseNameRuleIDsToRemove...)

	staticRulesToInstall, dynamicRulesToInstall := ParseRuleInstallAVPs(
		policyDBClient,
		rar.RulesToInstall,
	)

	return &protos.PolicyReAuthRequest{
		SessionId:             sid,
		Imsi:                  imsi,
		RulesToRemove:         rulesToRemove,
		RulesToInstall:        staticRulesToInstall,
		DynamicRulesToInstall: dynamicRulesToInstall,
	}
}

func (raa *ReAuthAnswer) FromProto(sessionID string, answer *protos.PolicyReAuthAnswer) *ReAuthAnswer {
	raa.SessionID = sessionID
	raa.ResultCode = diam.Success
	raa.RuleReports = make([]*ChargingRuleReport, 0, len(answer.FailedRules))
	for ruleName, code := range answer.FailedRules {
		raa.RuleReports = append(
			raa.RuleReports,
			&ChargingRuleReport{RuleNames: []string{ruleName}, FailureCode: RuleFailureCode(code)},
		)
	}
	return raa
}

func ConvertToProtoTimestamp(unixTime *time.Time) (*timestamp.Timestamp, error) {
	if unixTime == nil {
		return nil, nil
	}
	protoTimestamp, err := ptypes.TimestampProto(*unixTime)
	if err != nil {
		return nil, err
	}
	return protoTimestamp, nil
}

func ParseRuleInstallAVPs(
	policyDBClient policydb.PolicyDBClient,
	ruleInstalls []*RuleInstallAVP,
) ([]*protos.StaticRuleInstall, []*protos.DynamicRuleInstall) {
	staticRulesToInstall := make([]*protos.StaticRuleInstall, 0, len(ruleInstalls))
	dynamicRulesToInstall := make([]*protos.DynamicRuleInstall, 0, len(ruleInstalls))
	for _, ruleInstall := range ruleInstalls {
		activationTime, err := ConvertToProtoTimestamp(ruleInstall.RuleActivationTime)
		if err != nil {
			errMsg := fmt.Sprintf("Cannot convert time.Time to "+
				"google.protobuf.Timestamp for rule activation time: %s", err)
			glog.Error(errMsg)
		}

		deactivationTime, err := ConvertToProtoTimestamp(ruleInstall.RuleDeactivationTime)
		if err != nil {
			errMsg := fmt.Sprintf("Cannot convert time.Time to "+
				"google.protobuf.Timestamp for rule deactivation time: %s", err)
			glog.Error(errMsg)
		}

		for _, staticRuleName := range ruleInstall.RuleNames {
			staticRulesToInstall = append(
				staticRulesToInstall,
				&protos.StaticRuleInstall{
					RuleId:           staticRuleName,
					ActivationTime:   activationTime,
					DeactivationTime: deactivationTime,
				},
			)
		}

		if len(ruleInstall.RuleBaseNames) != 0 {
			baseNameRuleIdsToInstall := policyDBClient.GetRuleIDsForBaseNames(ruleInstall.RuleBaseNames)
			for _, baseNameRuleId := range baseNameRuleIdsToInstall {
				staticRulesToInstall = append(
					staticRulesToInstall,
					&protos.StaticRuleInstall{
						RuleId:           baseNameRuleId,
						ActivationTime:   activationTime,
						DeactivationTime: deactivationTime,
					},
				)
			}
		}

		for _, def := range ruleInstall.RuleDefinitions {
			dynamicRulesToInstall = append(
				dynamicRulesToInstall,
				&protos.DynamicRuleInstall{
					PolicyRule:       def.ToProto(),
					ActivationTime:   activationTime,
					DeactivationTime: deactivationTime,
				},
			)
		}
	}
	return staticRulesToInstall, dynamicRulesToInstall
}
