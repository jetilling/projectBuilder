#!/bin/bash

cd ~/projects/$1/$2

# Setting permissions

chmod -R ugo+rw storage/

# Copying over ENV template file

cp ~/go/src/github.com/jetilling/projectBuilder/laravel_files/.env.template .env.template

# Laravel comes with an .env.example file - we want to remove that
rm .env.example

# Copying .gitignore

cp ~/go/src/github.com/jetilling/projectBuilder/laravel_files/.gitignore .gitignore