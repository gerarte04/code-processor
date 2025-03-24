#!/bin/bash

if [[ $TRANSLATOR == "gcc" ]]; then
    g++ $FILE_NAME -std=c++20 -O2 -o exe && ./exe
elif [[ $TRANSLATOR == "clang" ]]; then
    clang++ $FILE_NAME -std=c++20 -O2 -o exe && ./exe
elif [[ $TRANSLATOR == "python3" ]]; then
    python3 $FILE_NAME
else
    echo "No supported translators"
fi
