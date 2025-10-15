#!/bin/sh


init()
{
    # cmake -DCMAKE_PREFIX_PATH=~/Qt/6.8.1/gcc_64/ -S ./src/ -B ./build/
    ~/Qt/6.8.1/gcc_64/bin/qt-cmake -S ./src/ -B ./build/
}