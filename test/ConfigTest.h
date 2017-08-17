//
// Created by Jean-Luc Thumm on 8/17/17.
//

#ifndef THUMMCOIN_CONFIGTEST_H
#define THUMMCOIN_CONFIGTEST_H

#include <gtest/gtest.h>
#include "../src/Config.h"

class ConfigTest : public ::testing::Test {
public:
  virtual SetUp() {
    good = Config{"data/data1"};
    bad = Config{"test/bad_data"};
  }

  Config good;
  Config bad;
};


#endif //THUMMCOIN_CONFIGTEST_H
