/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <iostream>
#include <fstream>
#include <cstdlib>
#include <json.hpp> // JSON library
#include <google/protobuf/message.h>
#include <google/protobuf/util/json_util.h>

#include "MConfigLoader.h"
#include "magma_logging.h"

using json = nlohmann::json;

namespace magma {

static bool check_file_exists(const std::string filename) {
  std::ifstream f(filename.c_str());
  return f.is_open();
}

bool MConfigLoader::load_service_mconfig(
    const std::string& service_name,
    google::protobuf::Message* message) {
  std::ifstream file;
  get_mconfig_file(&file);
  if (!file.is_open()) {
    MLOG(MERROR) << "Couldn't load mconfig file";
    return false;
  }

  json mconfig_json;
  mconfig_json << file;
  file.close();

  // config is located at mconfig_json["configs_by_key"][service_name]
  auto configs_it = mconfig_json.find("configs_by_key");
  if (configs_it == mconfig_json.end()) {
    configs_it = mconfig_json.find("configsByKey");
    if (configs_it == mconfig_json.end()) {
      MLOG(MERROR) << "Could not find configs_by_key in mconfig";
      return false;
    }
  }

  // Check if service exists
  auto service_it = configs_it->find(service_name);
  if (service_it == configs_it->end()) {
    MLOG(MERROR) << "Couldn't find pipelined config";
    return false;
  }
  service_it->erase("@type"); // @type param makes parsing fail

  // Parse to message and return
  auto status = google::protobuf::util::JsonStringToMessage(
    service_it->dump(), message);
  if (!status.ok()) {
    MLOG(MERROR) << "Couldn't parse pipelined config";
  }
  return status.ok();
}

void MConfigLoader::get_mconfig_file(std::ifstream* file) {
  // Load from /var/opt/magma if config exists, else read from /etc/magma
  if (check_file_exists(MConfigLoader::DYNAMIC_MCONFIG_PATH)) {
    file->open(MConfigLoader::DYNAMIC_MCONFIG_PATH);
    return;
  }
  const char* cfg_dir = std::getenv("MAGMA_CONFIG_LOCATION");
  if (cfg_dir == nullptr) {
    cfg_dir = MConfigLoader::CONFIG_DIR;
  }
  auto file_path = std::string(cfg_dir) + "/"
    + std::string(MConfigLoader::MCONFIG_FILE_NAME);
  file->open(file_path.c_str());
  return;
}

}
