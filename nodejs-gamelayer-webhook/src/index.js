const express = require('express');
const bodyParser = require('body-parser');
const routes = require('./routes');
const port = 5000;
const app = express();

app.use(bodyParser.json());
app.use('/', routes);

const start = async () => {
  try {
    app.listen(port, () => {
      console.log(`Running on ${port}`);
    });
  } catch (e) {
    console.error(e);
    process.exit(1);
  }
};

start();
