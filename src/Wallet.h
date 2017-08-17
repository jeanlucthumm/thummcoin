//
// Created by Jean-Luc Thumm on 8/17/17.
//

#ifndef THUMMCOIN_WALLET_H
#define THUMMCOIN_WALLET_H

#include <string>
#include "Config.h"

class Wallet {
public:
  Wallet(const Config &config);

private:
  double balance;
  std::string address;
};


#endif //THUMMCOIN_WALLET_H
