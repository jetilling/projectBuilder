#!/bin/bash

cd ~/projects/$1/$2

# Copy over component files
cp ~/go/src/github.com/jetilling/projectBuilder/react_files/main_container.js resources/js/components/main_container.js
cp ~/go/src/github.com/jetilling/projectBuilder/react_files/main.js.template resources/js/components/main.js.template

# Copy over serviceWorker, webpack, initialState, actions, package.json
cp ~/go/src/github.com/jetilling/projectBuilder/react_files/serviceWorker.js resources/js/serviceWorker.js
cp ~/go/src/github.com/jetilling/projectBuilder/react_files/initialState.js resources/js/initialState.js
cp ~/go/src/github.com/jetilling/projectBuilder/react_files/actions.js resources/js/actions.js
cp ~/go/src/github.com/jetilling/projectBuilder/react_files/index.js.template resources/js/index.js.template
cp ~/go/src/github.com/jetilling/projectBuilder/react_files/reducer.js.template resources/js/reducer.js.template
cp ~/go/src/github.com/jetilling/projectBuilder/react_files/webpack.mix.js webpack.mix.js
cp ~/go/src/github.com/jetilling/projectBuilder/react_files/package.json package.json

# Add styles directory
mkdir resources/styles
cp ~/go/src/github.com/jetilling/projectBuilder/react_files/index.css resources/styles/index.css

rm resources/js/components/ExampleComponent.vue
rm resources/js/app.js