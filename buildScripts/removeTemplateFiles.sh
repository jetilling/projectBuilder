#!/bin/bash

cd ~/projects/$1/$2

rm .env.template
rm create-testing-db_template.sql
rm docker-compose_template.yml