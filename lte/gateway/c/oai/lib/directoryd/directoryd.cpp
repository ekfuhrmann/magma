/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include <string>
#include <thread>

#include "DirectorydClient.h"
#include "directoryd.h"

static void directoryd_rpc_call_done(const grpc::Status &status);

bool directoryd_report_location(table_id_t table, char *imsi)
{
  // Actual GW_ID will be filled in the cloud
  magma::DirectoryServiceClient::UpdateLocation(
    static_cast<magma::TableID>(table),
    "IMSI" + std::string(imsi),
    std::string("GW_ID"),
    [&](grpc::Status status, magma::Void response) {
      directoryd_rpc_call_done(status);
    });
  return true;
}

bool directoryd_remove_location(table_id_t table, char *imsi)
{
  magma::DirectoryServiceClient::DeleteLocation(
    static_cast<magma::TableID>(table),
    "IMSI" + std::string(imsi),
    [&](grpc::Status status, magma::Void response) {
      directoryd_rpc_call_done(status);
    });
  return true;
}

bool directoryd_update_location(table_id_t table, char *imsi, char *location)
{
  magma::DirectoryServiceClient::UpdateLocation(
    static_cast<magma::TableID>(table),
    "IMSI" + std::string(imsi),
    std::string(location),
    [&](grpc::Status status, magma::Void response) {
      directoryd_rpc_call_done(status);
    });
  return true;
}

void directoryd_rpc_call_done(const grpc::Status &status)
{
  if (!status.ok()) {
    std::cerr << "Directoryd RPC failed with code " << status.error_code()
              << ", msg: " << status.error_message() << std::endl;
  }
}
