#!/bin/bash

cd ~/projects/$1/$2

echo "---------------------------------------------------"
echo "| Setting permissions"
echo "---------------------------------------------------"

chmod -R ugo+rw storage/

echo "---------------------------------------------------"
echo "| Updating ENV file"
echo "---------------------------------------------------"

cp ~/go/src/github.com/jetilling/projectBuilder/.env.example .env

rm .env.example