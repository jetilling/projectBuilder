#!/bin/bash

echo "---------------------------------------------------"
echo "| Running composer install"
echo "---------------------------------------------------"

cd ~/projects/$1/$2

docker run --rm -v $(pwd):/app composer install --ignore-platform-reqs