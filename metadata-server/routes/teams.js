var express = require('express');
const Web3 = require('web3');
const jsonInterface = require('../../truffle-core/build/contracts/CryptoPlayers.json').abi;
const teamsJSONInterface = require('../../truffle-core/build/contracts/CryptoTeams.json').abi;
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

const generateJSON = async ({playersContract, teamsContract, teamId}) => {
    try {
        var name = await teamsContract.methods.getName(teamId).call();
        var players = await teamsContract.methods.getPlayers(teamId).call();
        var playersName = [];
        for (let i = 0; i < players.length; i++) {
            const playerName = await playersContract.methods.getName(players[i]).call();
            playersName.push(playerName);
        }
    }
    catch (err) {
        console.error(err);
        return {};
    }

    const schema = {
        "name": name,
        "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
        "image": "http://www.monkers.net/fotos/lucasartsDead.jpg",
        "external_url": "https://www.freeverse.io/",
        "attributes":
            playersName.map(name => ({
                "trait_type": "player",
                "value": name
            }))
    };

    return schema;
};

module.exports = router;
module.exports.generateJSON = generateJSON;

