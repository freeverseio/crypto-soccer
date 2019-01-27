var express = require('express');
const Web3 = require('web3');
const jsonInterface = require('../../truffle-core/build/contracts/CryptoPlayers.json').abi;
const teamsJSONInterface = require('../../truffle-core/build/contracts/CryptoTeams.json').abi;
const teamsJSON = require('./teamsJSON');
const config = require('../config.json');

var router = express.Router();

const web3 = new Web3(config.provider);
const playersContract = new web3.eth.Contract(jsonInterface, config.crypto_player_address);
const teamsContract = new web3.eth.Contract(teamsJSONInterface, config.crypto_teams_contract);

/* GET JSON schema for teams with id. */
router.get('/:id', async (req, res, next) => {
    const teamId = req.params.id;
    const schema = await teamsJSON({ playersContract, teamsContract, teamId });
    res.send(schema);
});

module.exports = router;

