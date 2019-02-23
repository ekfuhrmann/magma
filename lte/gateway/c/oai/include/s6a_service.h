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

#pragma once

#include "bstrlib.h"
#include "mme_default_values.h"

typedef struct s6a_service_data_s {
  bstring server_address;
} s6a_service_data_t;

/*
  Init a grpc s6a_service
*/
int s6a_service_init(void);

#ifdef __cplusplus
extern "C" {
#endif
/**
 * Start the s6a_service Server and blocks
 *
 * @param server_address: the address and port to bind to ex "0.0.0.0:50051"
 */
void start_s6a_service_server(bstring server_address);

/**
 * Stop the server and clean up
 *
 */
void stop_s6a_service_server(void);

#ifdef __cplusplus
}
#endif
