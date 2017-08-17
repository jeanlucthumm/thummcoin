//
// Created by Jean-Luc Thumm on 8/17/17.
//

#ifndef THUMMCOIN_CONFIG_H
#define THUMMCOIN_CONFIG_H


#include <fstream>
#include <boost/filesystem.hpp>

class Config {
public:
  Config(const boost::filesystem::path dataDir);

  bool validate(); // this must be called right after constructor

  std::ifstream walletConfig();

  std::ifstream peerList();

private:
  boost::filesystem::path dataDir;
  boost::filesystem::path walletConfigPath;
  boost::filesystem::path peerListPath;
};


#endif //THUMMCOIN_CONFIG_H
