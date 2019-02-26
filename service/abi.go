package service

const lionelAbiJson = `
[
	{
		"constant": false,
		"inputs": [
			{
				"name": "_leagueNo",
				"type": "uint256"
			},
			{
				"name": "_value",
				"type": "bytes32"
			}
		],
		"name": "challange",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [],
		"name": "next",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "_leagueNo",
				"type": "uint256"
			},
			{
				"name": "_value",
				"type": "bytes32"
			}
		],
		"name": "update",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "leagueNo",
				"type": "uint256"
			},
			{
				"indexed": false,
				"name": "value",
				"type": "bytes32"
			}
		],
		"name": "LeagueChallangeAvailable",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"name": "leagueNo",
				"type": "uint256"
			}
		],
		"name": "LeagueChallangeSucessfull",
		"type": "event"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "_leagueNo",
				"type": "uint256"
			}
		],
		"name": "canLeagueBeChallanged",
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
				"name": "_leagueNo",
				"type": "uint256"
			}
		],
		"name": "canLeagueBeUpdated",
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
		"name": "legueCount",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "pure",
		"type": "function"
	}
]
`
