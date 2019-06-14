#!/bin/bash

APP_NAME=$1
UNIQUE_ID=$2

DIR_NAME=$APP_NAME
DIR_NAME+="_$UNIQUE_ID"

cd ~/projects

mkdir $DIR_NAME