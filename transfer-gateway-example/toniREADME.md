## Install

First YARN: I messed up Yarn and Npm, until I only used brew for them all, and installed Yarn WITHOUT node support: brew install yarn --without-node

Second Metamask with the private account...(TODO specify)

## Run the standard example (you need to install )

Then, just do: ./transfer_gateway setup, and then ./transfer_gateway start

## Run the under-construction setup to prepare tests:

./prepareForTesting.sh start


## WEBPACK

The web interface (in ./webclient) uses Webpack to dynamically pack all JS libraries. 
It creates a ./dist folder with an dummy index.html which calls a dynamically generated bundle.js.

The config is in webpack.config.js contains:

const webpack = require('webpack')

webpack is called directly as specified in package.json:
    "start": "webpack-dev-server --hot --content-base ./dist"





