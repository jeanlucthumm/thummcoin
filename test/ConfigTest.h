//
// Created by Jean-Luc Thumm on 8/17/17.
//

#ifndef THUMMCOIN_CONFIGTEST_H
#define THUMMCOIN_CONFIGTEST_H

#include <gtest/gtest.h>
#include <string>
#include "../src/Config.h"

TEST(ConfigTest, validate) {
  Config good{"test/good_data"};
  Config bad{"test/bad_data"};

  EXPECT_TRUE(good.validate());
  EXPECT_FALSE(bad.validate());
}

TEST(ConfigTest, walletConfig) {
  Config c{"test/good_data"};
  c.validate();

  std::ifstream walletConfig = c.walletConfig();
  std::ifstream peerList = c.peerList();

  std::ifstream actWalletConfig{"test/good_data/wallet_config.cfg"};
  std::ifstream actPeerList{"test/good_data/peer_list.cfg"};

  std::string actual, line;

  std::getline(actWalletConfig, actual);
  std::getline(walletConfig, line);
  EXPECT_EQ(actual, line);
  std::getline(actPeerList, actual);
  std::getline(peerList, line);
  EXPECT_EQ(actual, line);
}



#endif //THUMMCOIN_CONFIGTEST_H
