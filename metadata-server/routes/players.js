var express = require('express');
const Web3 = require('web3');
const schema = require('./example_schema.json');
const jsonInterface = require('../../truffle-core/build/contracts/CryptoPlayers.json').abi;
const config = require('../config.json');

var router = express.Router();

// "Web3.providers.givenProvider" will be set if in an Ethereum supported browser.
const web3 = new Web3(Web3.givenProvider || config.provider);
const instance = new web3.eth.Contract(jsonInterface, config.cryptoPlayerAddress);

/* GET JSON schema for players with id. */
router.get('/:id', async (req, res, next) => {
  const id = req.params.id;
  res.send(schema);
});

module.exports = router;
