#!/bin/bash

cd ~/projects/$1/$2

echo "---------------------------------------------------"
echo "| Removing Template Files"
echo "---------------------------------------------------"

rm .env.template
rm create-testing-db_template.sql
rm docker-compose_template.yml