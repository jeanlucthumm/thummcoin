//
// Created by Jean-Luc Thumm on 8/17/17.
//

#ifndef THUMMCOIN_EVENT_H
#define THUMMCOIN_EVENT_H

#include <cstddef>
#include <vector>

class Event {
public:
  enum Type {
    Transaction
  };

  Event(const std::vector<std::byte> &data);

  Type type();

private:
  virtual bool interpret() = 0;
};


#endif //THUMMCOIN_EVENT_H
