#!/bin/bash

cd ~/projects/$1/$2

echo "---------------------------------------------------"
echo "| Copying Over Docker Files"
echo "---------------------------------------------------"

cp ~/go/src/github.com/jetilling/projectBuilder/laravel_files/app.dockerfile app.dockerfile
cp ~/go/src/github.com/jetilling/projectBuilder/laravel_files/web.dockerfile web.dockerfile
cp ~/go/src/github.com/jetilling/projectBuilder/laravel_files/postgres.dockerfile postgres.dockerfile

cp ~/go/src/github.com/jetilling/projectBuilder/laravel_files/vhost.conf vhost.conf

cp ~/go/src/github.com/jetilling/projectBuilder/laravel_files/create-testing-db_template.sql create-testing-db_template.sql
cp ~/go/src/github.com/jetilling/projectBuilder/laravel_files/docker-compose_template.yml docker-compose_template.yml
