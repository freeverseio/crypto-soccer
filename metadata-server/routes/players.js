const express = require('express');
const Web3 = require('web3');
const jsonInterface = require('../../truffle-core/build/contracts/CryptoPlayers.json').abi;
const teamsJSONInterface = require('../../truffle-core/build/contracts/CryptoTeams.json').abi;
const config = require('../config.json');
const spawn = require("child_process").spawn;

const router = express.Router();

const web3 = new Web3(config.provider);
const playersContract = new web3.eth.Contract(jsonInterface, config.crypto_player_address);
const teamsContract = new web3.eth.Contract(teamsJSONInterface, config.crypto_teams_contract);

/* GET JSON schema for players with id. */
router.get('/:id', async (req, res, next) => {
  const playerId = req.params.id;
  const schema = await generateJSON({ playersContract, teamsContract, playerId });
  res.send(schema);
});

const generateJSON = async ({ playersContract, teamsContract, playerId }) => {
  try {
    const pythonProcess = spawn('python',["../extract_svg/player_composer.py", '-n', playerId, '-o', 'public/images/' + playerId]);
    pythonProcess.stdout.on('data', data => console.log(data.toString())) 
    pythonProcess.stderr.on('data', data => console.log(data.toString()))
    var name = await playersContract.methods.getName(playerId).call();
    var image = 'http://metadata.busyverse.com:3000/images/' + playerId + '.svg';
    var speed = await playersContract.methods.getSpeed(playerId).call();
    var defence = await playersContract.methods.getDefence(playerId).call();
    var endurance = await playersContract.methods.getEndurance(playerId).call();
    var shoot = await playersContract.methods.getShoot(playerId).call();
    var pass = await playersContract.methods.getPass(playerId).call();
    const teamId = await playersContract.methods.getTeam(playerId).call();
    var teamName = teamId == 0 ? "" : await teamsContract.methods.getName(teamId).call();
  }
  catch (err) {
    console.error(err);
    return {};
  }

  const schema = {
    "name": name,
    "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
    "image": image,
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
      },
      {
        "trait_type": "team",
        "value": teamName
      },
      {
        "display_type": "boost_number",
        "trait_type": "shoot_power",
        "value": 10
      },
      {
        "display_type": "boost_percentage",
        "trait_type": "pass_increase",
        "value": 5
      },
      {
        "display_type": "number",
        "trait_type": "generation",
        "value": 0
      }
    ]
  };

  return schema;
};

module.exports = router;
module.exports.generateJSON = generateJSON;
