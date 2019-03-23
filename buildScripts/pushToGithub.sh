#!/bin/bash

cd ~/projects/$1/$2

git init

git add .

git commit -m "initial commit from Launch"


JSON_FMT='{"name":"%s"}\n'
json_string=$(printf "$JSON_FMT" "$2")

curl -H "Authorization: token $3" https://api.github.com/user/repos -d $json_string


git remote add origin git@github.com:jetilling/$2.git

git push -u origin master