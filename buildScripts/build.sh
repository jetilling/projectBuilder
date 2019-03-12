#!/bin/bash

# APP_NAME=$1
# UNIQUE_ID=$2

# DIR_NAME=$APP_NAME
# DIR_NAME+="_$UNIQUE_ID"

# cd ~/projects

# mkdir $DIR_NAME

# cd $DIR_NAME

# echo "---------------------------------------------------"
# echo "| Downloading Laravel in $DIR_NAME"
# echo "---------------------------------------------------"

# curl -L https://github.com/laravel/laravel/archive/v5.8.3.tar.gz | tar xz

# echo "---------------------------------------------------"
# echo "| Renaming laravel-5.8.3 to $APP_NAME"
# echo "---------------------------------------------------"

# mv laravel-5.8.3 $APP_NAME

# cd $APP_NAME

# echo "---------------------------------------------------"
# echo "| Running composer install"
# echo "---------------------------------------------------"

# docker run --rm -v $(pwd):/app composer install --ignore-platform-reqs

# echo "---------------------------------------------------"
# echo "| Setting permissions"
# echo "---------------------------------------------------"

# chmod -R ugo+rw storage/

# echo "---------------------------------------------------"
# echo "| Updating ENV file"
# echo "---------------------------------------------------"

# cp ~/go/src/github.com/jetilling/projectBuilder/.env.example .env
# rm .env.example

# echo "---------------------------------------------------"
# echo "| Project is all setup! "
# echo "| "
# echo "| Please note that docker has not been built locally"
# echo "| This probably a good time to try writing to files"
# echo "| or try creating a new repo in GitHub and pushing"
# echo "| to that."
# echo "| It could then be pushed to Heroku"
# echo "---------------------------------------------------"
