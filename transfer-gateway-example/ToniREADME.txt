webpack.config.js contains:

const webpack = require('webpack')



loads dist/index.html --> which calls ./bundle.js

this is because in package.json we have:
    "start": "webpack-dev-server --hot --content-base ./dist",



