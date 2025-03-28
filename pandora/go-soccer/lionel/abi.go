package lionel

const stateAbiJson = `[{"constant":true,"inputs":[{"name":"playerState","type":"uint256"}],"name":"getSkills","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"leagueState","type":"uint256[]"},{"name":"teamState","type":"uint256[]"}],"name":"leagueStateAppend","outputs":[{"name":"state","type":"uint256[]"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"state","type":"uint256"}],"name":"isValidPlayerState","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"leagueState","type":"uint256[]"},{"name":"teamIdx","type":"uint256"},{"name":"teamState","type":"uint256[]"}],"name":"leagueStateUpdate","outputs":[{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"}],"name":"getEndurance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"leagueState","type":"uint256[]"},{"name":"idx","type":"uint256"}],"name":"leagueStateAt","outputs":[{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"teamState","type":"uint256[]"},{"name":"idx","type":"uint256"}],"name":"teamStateAt","outputs":[{"name":"playerState","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"teamState","type":"uint256[]"}],"name":"teamStateSize","outputs":[{"name":"count","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"state","type":"uint256"},{"name":"value","type":"uint256"}],"name":"setPrevLeagueId","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"}],"name":"getSpeed","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"}],"name":"getDefence","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"defence","type":"uint256"},{"name":"speed","type":"uint256"},{"name":"pass","type":"uint256"},{"name":"shoot","type":"uint256"},{"name":"endurance","type":"uint256"},{"name":"monthOfBirthInUnixTime","type":"uint256"},{"name":"playerId","type":"uint256"},{"name":"currentTeamId","type":"uint256"},{"name":"currentShirtNum","type":"uint256"},{"name":"prevLeagueId","type":"uint256"},{"name":"prevTeamPosInLeague","type":"uint256"},{"name":"prevShirtNumInLeague","type":"uint256"},{"name":"lastSaleBlock","type":"uint256"}],"name":"playerStateCreate","outputs":[{"name":"state","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"}],"name":"getPass","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"}],"name":"getPrevTeamPosInLeague","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"teamState","type":"uint256[]"},{"name":"delta","type":"uint8"}],"name":"teamStateEvolve","outputs":[{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"}],"name":"getShoot","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"}],"name":"getPrevShirtNumInLeague","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[],"name":"leagueStateCreate","outputs":[{"name":"state","type":"uint256[]"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"state","type":"uint256[]"}],"name":"isValidLeagueState","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[],"name":"teamStateCreate","outputs":[{"name":"state","type":"uint256[]"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"state","type":"uint256"},{"name":"value","type":"uint256"}],"name":"setPrevTeamPosInLeague","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"}],"name":"getMonthOfBirthInUnixTime","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"},{"name":"delta","type":"uint16"}],"name":"playerStateEvolve","outputs":[{"name":"evolvedState","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"leagueState","type":"uint256[]"}],"name":"leagueStateGetSkills","outputs":[{"name":"skills","type":"uint256[]"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"state","type":"uint256"},{"name":"currentShirtNum","type":"uint256"}],"name":"setCurrentShirtNum","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"state","type":"uint256[]"}],"name":"isValidTeamState","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"state","type":"uint256"},{"name":"lastSaleBlock","type":"uint256"}],"name":"setLastSaleBlock","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"teamState","type":"uint256[]"}],"name":"computeTeamRating","outputs":[{"name":"rating","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"leagueState","type":"uint256[]"}],"name":"leagueStateSize","outputs":[{"name":"count","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"},{"name":"teamId","type":"uint256"}],"name":"setCurrentTeamId","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"}],"name":"getLastSaleBlock","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"}],"name":"getSkillsVec","outputs":[{"name":"skills","type":"uint16[5]"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"}],"name":"getCurrentTeamId","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"teamState","type":"uint256[]"},{"name":"playerState","type":"uint256"}],"name":"teamStateAppend","outputs":[{"name":"state","type":"uint256[]"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"}],"name":"getCurrentShirtNum","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"}],"name":"getPlayerId","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"playerState","type":"uint256"}],"name":"getPrevLeagueId","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"}]`
const assetsAbiJson = `[{"constant":false,"inputs":[{"name":"teamId","type":"uint256"},{"name":"leagueId","type":"uint256"},{"name":"posInLeague","type":"uint8"}],"name":"signToLeague","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"teamId","type":"uint256"}],"name":"getTeamPlayerIds","outputs":[{"name":"playerIds","type":"uint256[11]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"countTeams","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"NUM_SKILLS","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"playerId0","type":"uint256"},{"name":"playerId1","type":"uint256"}],"name":"exchangePlayersTeams","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"teamId","type":"uint256"}],"name":"getTeamName","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"teamId","type":"uint256"},{"name":"posInTeam","type":"uint8"}],"name":"getPlayerIdFromTeamIdAndPos","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"name","type":"string"},{"name":"owner","type":"address"}],"name":"createTeam","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"PLAYERS_PER_TEAM","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"name","type":"string"}],"name":"getTeamOwner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"playerId","type":"uint256"}],"name":"getPlayerState","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"teamId","type":"uint256"}],"name":"getTeamCurrentHistory","outputs":[{"name":"currentLeagueId","type":"uint256"},{"name":"posInCurrentLeague","type":"uint8"},{"name":"prevLeagueId","type":"uint256"},{"name":"posInPrevLeague","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"playerState","type":"address"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"name":"teamName","type":"string"},{"indexed":false,"name":"teamId","type":"uint256"}],"name":"TeamCreation","type":"event"}]`

const leaguesAbiJson = `[
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getScores",
      "outputs": [
        {
          "name": "",
          "type": "uint16[]"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "teamIds",
          "type": "uint256[]"
        },
        {
          "name": "tactics",
          "type": "uint8[3][]"
        },
        {
          "name": "blocks",
          "type": "uint256[]"
        }
      ],
      "name": "computeUsersAlongDataHash",
      "outputs": [
        {
          "name": "",
          "type": "bytes32"
        }
      ],
      "payable": false,
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getNTeams",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "hasFinished",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "resetUpdater",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getUpdateBlock",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "leaguesCount",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "scores",
          "type": "uint16[]"
        },
        {
          "name": "score",
          "type": "uint16"
        }
      ],
      "name": "scoresAppend",
      "outputs": [
        {
          "name": "",
          "type": "uint16[]"
        }
      ],
      "payable": false,
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "isVerified",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "tactics",
          "type": "uint256[3][]"
        }
      ],
      "name": "hashTactics",
      "outputs": [
        {
          "name": "",
          "type": "bytes32"
        }
      ],
      "payable": false,
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        },
        {
          "name": "day",
          "type": "uint256"
        }
      ],
      "name": "scoresGetDay",
      "outputs": [
        {
          "name": "dayScores",
          "type": "uint16[]"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "hasStarted",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getTactics",
      "outputs": [
        {
          "name": "",
          "type": "uint8[]"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getDayStateHashes",
      "outputs": [
        {
          "name": "",
          "type": "bytes32[]"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "home",
          "type": "uint8"
        },
        {
          "name": "visitor",
          "type": "uint8"
        }
      ],
      "name": "encodeScore",
      "outputs": [
        {
          "name": "score",
          "type": "uint16"
        }
      ],
      "payable": false,
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "leagueId",
          "type": "uint256"
        },
        {
          "name": "leagueDay",
          "type": "uint256"
        },
        {
          "name": "initLeagueState",
          "type": "uint256[]"
        },
        {
          "name": "tactics",
          "type": "uint8[3][]"
        }
      ],
      "name": "computeDay",
      "outputs": [
        {
          "name": "scores",
          "type": "uint16[]"
        },
        {
          "name": "finalLeagueState",
          "type": "uint256[]"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        },
        {
          "name": "initBlock",
          "type": "uint256"
        },
        {
          "name": "step",
          "type": "uint256"
        },
        {
          "name": "teamIds",
          "type": "uint256[]"
        },
        {
          "name": "tactics",
          "type": "uint8[3][]"
        }
      ],
      "name": "create",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "getChallengePeriod",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getStep",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        },
        {
          "name": "usersInitDataTeamIds",
          "type": "uint256[]"
        },
        {
          "name": "usersInitDataTactics",
          "type": "uint8[3][]"
        },
        {
          "name": "usersAlongDataTeamIds",
          "type": "uint256[]"
        },
        {
          "name": "usersAlongDataTactics",
          "type": "uint8[3][]"
        },
        {
          "name": "usersAlongDataBlocks",
          "type": "uint256[]"
        },
        {
          "name": "leagueDay",
          "type": "uint256"
        },
        {
          "name": "prevMatchdayStates",
          "type": "uint256[]"
        }
      ],
      "name": "challengeMatchdayStates",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getMatchPerDay",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        },
        {
          "name": "matchday",
          "type": "uint256"
        },
        {
          "name": "matchIdx",
          "type": "uint256"
        }
      ],
      "name": "getTeamsInMatch",
      "outputs": [
        {
          "name": "homeIdx",
          "type": "uint256"
        },
        {
          "name": "visitorIdx",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getEndBlock",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        },
        {
          "name": "initStateHash",
          "type": "bytes32"
        },
        {
          "name": "dayStateHashes",
          "type": "bytes32[]"
        },
        {
          "name": "scores",
          "type": "uint16[]"
        },
        {
          "name": "isLie",
          "type": "bool"
        }
      ],
      "name": "updateLeague",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "getEngineContract",
      "outputs": [
        {
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getUsersAlongDataHash",
      "outputs": [
        {
          "name": "",
          "type": "bytes32"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "target",
          "type": "uint16[]"
        },
        {
          "name": "scores",
          "type": "uint16[]"
        }
      ],
      "name": "scoresConcat",
      "outputs": [
        {
          "name": "",
          "type": "uint16[]"
        }
      ],
      "payable": false,
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getUpdater",
      "outputs": [
        {
          "name": "",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "countLeagueDays",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        },
        {
          "name": "teamIds",
          "type": "uint256[]"
        },
        {
          "name": "tactics",
          "type": "uint8[3][]"
        },
        {
          "name": "dataToChallengeInitStates",
          "type": "uint256[]"
        }
      ],
      "name": "challengeInitStates",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        },
        {
          "name": "teamIds",
          "type": "uint256[]"
        },
        {
          "name": "tactics",
          "type": "uint8[3][]"
        },
        {
          "name": "blocks",
          "type": "uint256[]"
        }
      ],
      "name": "updateUsersAlongDataHash",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getInitStateHash",
      "outputs": [
        {
          "name": "",
          "type": "bytes32"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "teamIds",
          "type": "uint256[]"
        },
        {
          "name": "tactics",
          "type": "uint8[3][]"
        }
      ],
      "name": "hashUsersInitData",
      "outputs": [
        {
          "name": "",
          "type": "bytes32"
        }
      ],
      "payable": false,
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        },
        {
          "name": "day",
          "type": "uint256"
        }
      ],
      "name": "getMatchDayBlockHash",
      "outputs": [
        {
          "name": "",
          "type": "bytes32"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getUsersInitDataHash",
      "outputs": [
        {
          "name": "",
          "type": "bytes32"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "stakersContract",
          "type": "address"
        }
      ],
      "name": "setStakersContract",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "state",
          "type": "uint256[]"
        }
      ],
      "name": "hashInitState",
      "outputs": [
        {
          "name": "",
          "type": "bytes32"
        }
      ],
      "payable": false,
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getLastChallengeBlock",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "state",
          "type": "uint256[]"
        }
      ],
      "name": "hashDayState",
      "outputs": [
        {
          "name": "",
          "type": "bytes32"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "isUpdated",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getIsLie",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "scoresCreate",
      "outputs": [
        {
          "name": "",
          "type": "uint16[]"
        }
      ],
      "payable": false,
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getInitBlock",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        },
        {
          "name": "teamIds",
          "type": "uint256[]"
        },
        {
          "name": "tactics",
          "type": "uint8[3][]"
        },
        {
          "name": "dataToChallengeInitStates",
          "type": "uint256[]"
        }
      ],
      "name": "getInitPlayerStates",
      "outputs": [
        {
          "name": "state",
          "type": "uint256[]"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "id",
          "type": "uint256"
        }
      ],
      "name": "getTeams",
      "outputs": [
        {
          "name": "",
          "type": "uint256[]"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "score",
          "type": "uint16"
        }
      ],
      "name": "decodeScore",
      "outputs": [
        {
          "name": "home",
          "type": "uint8"
        },
        {
          "name": "visitor",
          "type": "uint8"
        }
      ],
      "payable": false,
      "stateMutability": "pure",
      "type": "function"
    },
    {
      "inputs": [
        {
          "name": "engine",
          "type": "address"
        },
        {
          "name": "state",
          "type": "address"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "constructor"
    }
  ]`
