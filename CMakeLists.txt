cmake_minimum_required(VERSION 3.9)
project(thummcoin)

set(CMAKE_CXX_STANDARD 14)

enable_testing()

find_package(Boost COMPONENTS filesystem system thread REQUIRED)
find_package(GTest)

file(GLOB SOURCES src/*.cpp src/*.h)
file(GLOB TEST_SOURCES test/*.h)

add_executable(thummcoin ${SOURCES})
target_link_libraries(thummcoin Boost::boost ${Boost_LIBRARIES})


add_executable(test_thummcoin test/main.cpp ${SOURCES} ${TEST_SOURCES})
target_link_libraries(
        test_thummcoin Boost::boost ${Boost_LIBRARIES}
        GTest::GTest GTest::Main)
add_test(all_tests test_thummcoin)
