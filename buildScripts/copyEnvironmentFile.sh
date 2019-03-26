#!/bin/bash

cd ~/projects/$1/$2

echo "---------------------------------------------------"
echo "| Setting permissions"
echo "---------------------------------------------------"

chmod -R ugo+rw storage/

echo "---------------------------------------------------"
echo "| Copying over ENV template file"
echo "---------------------------------------------------"

cp ~/go/src/github.com/jetilling/projectBuilder/laravel_files/.env.template .env.template

# Laravel comes with an .env.example file - we want to remove that
rm .env.example

echo "---------------------------------------------------"
echo "| Copying .gitignore"
echo "---------------------------------------------------"

cp ~/go/src/github.com/jetilling/projectBuilder/laravel_files/.gitignore .gitignore