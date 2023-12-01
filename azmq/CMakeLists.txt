cmake_minimum_required(VERSION 3.18)

project(async-bench)

add_executable(async-bench src/main.cpp)
target_compile_features(async-bench PRIVATE cxx_std_20)
target_compile_options(async-bench PRIVATE
    -Wall -Wextra -fno-exceptions -fno-rtti)
