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
  const schema = await generateJSON({ playersContract, teamsContract, teamId });
  res.send(schema);
});

const generateJSON = async ({ playersContract, teamsContract, teamId }) => {
  let speed = 0;
  let defence = 0;
  let endurance = 0;
  let shoot = 0;
  let pass = 0;

  try {
    var name = await teamsContract.methods.getName(teamId).call();
    var players = await teamsContract.methods.getPlayers(teamId).call();
    var playersName = [];
    for (let i = 0; i < players.length; i++) {
      const playerId = players[i];
      const playerName = await playersContract.methods.getName(playerId).call();
      playersName.push({
        "trait_type": "player",
        "value": playerName
      });
      speed += Number(await playersContract.methods.getSpeed(playerId).call());
      defence += Number(await playersContract.methods.getDefence(playerId).call());
      endurance += Number(await playersContract.methods.getEndurance(playerId).call());
      shoot += Number(await playersContract.methods.getShoot(playerId).call());
      pass += Number(await playersContract.methods.getPass(playerId).call());
    }
    speed = Math.floor(speed / players.length);
    defence = Math.floor(defence /players.length);
    endurance = Math.floor(endurance / players.length);
    shoot = Math.floor(shoot / players.length);
    pass = Math.floor(pass / players.length);
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
    "attributes": [
      {
        "trait_type": "speed",
        "value": Number(speed),
        "max_value": 100
      },
      {
        "trait_type": "defence",
        "value": Number(defence),
        "max_value": 100
      },
      {
        "trait_type": "endurance",
        "value": Number(endurance),
        "max_value": 100
      },
      {
        "trait_type": "shoot",
        "value": Number(shoot),
        "max_value": 100
      },
      {
        "trait_type": "pass",
        "value": Number(pass),
        "max_value": 100
      }
    ]
  };

  schema.attributes = schema.attributes.concat(playersName);

  return schema;
};

module.exports = router;
module.exports.generateJSON = generateJSON;

