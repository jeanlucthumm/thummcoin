//
// Created by Jean-Luc Thumm on 8/17/17.
//

#include <vector>
#include <string>
#include "Config.h"

using namespace std;
namespace fs = boost::filesystem;

#define WALLET_CONFIG_NAME "wallet_config.cfg"
#define PEER_LIST_NAME "peer_list.cfg"

Config::Config(const boost::filesystem::path dataDir)
  : dataDir{dataDir} {}

bool Config::validate() {
  if (!fs::exists(dataDir) || !fs::is_directory(dataDir)) {
    return false;
  }

  // look for config files
  fs::directory_iterator begin{dataDir}, end;
  vector<fs::directory_entry> entries(begin, end);

  for (auto &entry : entries) {
    const string &name = entry.path().filename().string();
    if (name == WALLET_CONFIG_NAME) {
      walletConfigPath = entry.path();
    } else if (name == PEER_LIST_NAME) {
      peerListPath = entry.path();
    }
  }

  return !walletConfigPath.empty() &&
         !peerListPath.empty();
}

std::ifstream Config::walletConfig() {
  return ifstream{walletConfigPath.string()};
}

std::ifstream Config::peerList() {
  return ifstream{peerListPath.string()};
}
