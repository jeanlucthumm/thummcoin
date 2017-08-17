//
// Created by Jean-Luc Thumm on 8/17/17.
//

#ifndef THUMMCOIN_CONFIGTEST_H
#define THUMMCOIN_CONFIGTEST_H

#include <gtest/gtest.h>
#include "../src/Config.h"

class ConfigTest : public ::testing::Test {
public:
  virtual void SetUp() {
    good = Config{"data/data1"};
    bad = Config{"test/bad_data"};
  }

  Config good{"test/good_data"};
  Config bad{"test/bad_data"};
};

TEST_F(ConfigTest, validate) {
  EXPECT_TRUE(good.validate());
  EXPECT_FALSE(bad.validate());
}

#endif //THUMMCOIN_CONFIGTEST_H
