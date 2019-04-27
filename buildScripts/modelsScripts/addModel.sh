#!/bin/bash

cd ~/projects/$1/$2

php artisan make:model $3 -mc