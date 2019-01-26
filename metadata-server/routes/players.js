const express = require('express');
const Web3 = require('web3');
const jsonInterface = require('../../truffle-core/build/contracts/CryptoPlayers.json').abi;
const teamsJSONInterface = require('../../truffle-core/build/contracts/CryptoTeams.json').abi;
const playersJSON = require('./playersJSON');
const config = require('../config.json');

const router = express.Router();

const web3 = new Web3(config.provider);
const playersContract = new web3.eth.Contract(jsonInterface, config.crypto_player_address);
const teamsContract = new web3.eth.Contract(teamsJSONInterface, config.crypto_teams_contract);

/* GET JSON schema for players with id. */
router.get('/:id', async (req, res, next) => {
  const id = req.params.id;
  const schema = await playersJSON({ playersContract, teamsContract, playerId });
  res.send(schema);
});

module.exports = router;
