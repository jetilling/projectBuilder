#!/bin/bash

cd ~/projects/$1

echo "---------------------------------------------------"
echo "| Downloading Laravel in $1"
echo "---------------------------------------------------"

curl -L https://github.com/laravel/laravel/archive/v5.8.5.tar.gz | tar xz
