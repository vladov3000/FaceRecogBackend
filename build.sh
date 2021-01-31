#!/usr/bin/env bash

# compile cpp to library
if [[ $1 == "prod" ]]; then
        echo "building prod..."

    if [[ -d build ]]; then
        rm -rf build
    fi
    mkdir build
    cd build

    cmake -DCMAKE_BUILD_TYPE=Release ..
    cmake --build .

else 
    echo "building dev..."

    if [[ ! -d build ]]; then
        echo "build folder not found, creating new one"

        mkdir build
        cd build 
        cmake -DCMAKE_BUILD_TYPE=Debug ..

    else
        cd build
    fi
    
    cmake --build .
fi

if [[ -d bin ]]; then
    rm -rf bin
fi
mkdir bin

# compile go
cd ../src
go build -o start-server
echo $PWD
mv start-server ../build/bin
