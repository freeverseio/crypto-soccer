//Configuration paths and urls
const rpcURL = "ws://127.0.0.1:8545"
const listeningPort = 3000
const addressRefreshTime = 600000 //10 minutes
const contractsPath = '/home/ignasi/blockchain/cryptosoccer/web3/cryptoSoccerServer/contracts/'
//const contractsPath = '/Users/ignasi/Documents/blockchain/projecteFinal/cryptosoccer/web3/cryptoSoccerServer/contracts/'
//Variables
let contractAddress
let mainContract
const numOfPlayers = 11
let createdTeams = 0

//Required libs
const Web3 = require('web3')
const ganache = require('ganache-cli')
const web3 = new Web3(rpcURL);
const express = require('express')
const app = express()
const solc = require('solc')
const fs = require('fs');

//Get ganache addresses
let addresses
retrieveAddresses()
function retrieveAddresses() {
    console.log('Checking ganache addresses')
    web3.eth.getAccounts(function (error, result) {

        addresses = result
        compileContracts()

    })
    console.log('Ganache addresses checked! :D')
}
setTimeout(retrieveAddresses, addressRefreshTime);

function compileContracts() {
    //Compile player_factory and its imports (synchronous)
    console.log('Compiling contracts')

    try {

        var input = {
            'Math.sol': fs.readFileSync(contractsPath + 'Math.sol', 'utf8'),
            'SafeMath.sol': fs.readFileSync(contractsPath + 'SafeMath.sol', 'utf8'),
            'helper_functions.sol': fs.readFileSync(contractsPath + 'helper_functions.sol', 'utf8'),
            'Ownable.sol': fs.readFileSync(contractsPath + 'Ownable.sol', 'utf8'),
            'player_factory.sol': fs.readFileSync(contractsPath + 'player_factory.sol', 'utf8'),
            'team_factory.sol': fs.readFileSync(contractsPath + 'team_factory.sol', 'utf8'),
            'game_engine.sol': fs.readFileSync(contractsPath + 'game_engine.sol', 'utf8')

        };

        var output = solc.compile({ sources: input }, 1)
        let abi = JSON.parse(output.contracts['game_engine.sol:GameEngine'].interface);
        let bytecode = '0x' + output.contracts['game_engine.sol:GameEngine'].bytecode;

        const contract = new web3.eth.Contract(abi)

        contract.deploy({
            data: bytecode
        })
            .send({
                from: addresses[0],
                gas: 6721975,
                gasPrice: '20000000000'
            }, function (error, transactionHash) { console.error(error) })
            .on('error', function (error) { console.error(error) })
            .on('transactionHash', function (transactionHash) { console.log(transactionHash) })
            .on('receipt', function (receipt) {
                console.log(receipt.contractAddress) // contains the new contract address
            })
            .on('confirmation', function (confirmationNumber, receipt) { console.log('Confirmed deploy') })
            .then(function (newContractInstance) {
                console.log(newContractInstance.options.address)
                contractAddress = newContractInstance.options.address // instance with the new contract address
                mainContract = newContractInstance
            });

    } catch (e) {
        console.error('Problem compiling:' + e);
    }

}

//Creates a team with the players
function createTeam(res) {

    var randomTeamName = Math.random().toString(36).substring(2, 5)

    mainContract.methods.createTeam(randomTeamName).send({ from: addresses[0], gasPrice: 20000000000, gas: 6721975 }, function (error, transactionHash) {

        //Get created team Id //TODO: get team id dinamically
        let teamId = 0

        //Create the players for the created team
        createRandomPlayer(teamId, 1, 1, res)

    })
}

//Creates team players
function createRandomPlayer(teamId, randomTeamName, playerNumberInTeam, res) {

    var playerRole = 0
    var randomName = Math.random().toString(36).substring(2, 5)
    if (playerNumberInTeam > 0) playerRole = (playerNumberInTeam % 3) + 1;

    mainContract.methods.createRandomPlayer(randomName, teamId, randomTeamName, playerNumberInTeam, playerRole)
        .send({ from: addresses[0], gasPrice: 20000000000, gas: 6721975 }, function (error, transactionHash) {

            if (playerNumberInTeam < numOfPlayers) {
                createRandomPlayer(teamId, randomTeamName, playerNumberInTeam + 1, res)
            } else {
 
                //Check everyting was correct
                //Get the team
                mainContract.methods.getNCreatedTeams().call({ from: addresses[0] }, function (error, result) {
                    console.log("Created teams:" + result)
                    createdTeams = result
                    res.send(createdTeams)
                })

                //Get the players
                mainContract.methods.getNCreatedPlayers().call({ from: addresses[0] }, function (error, result) {
                    console.log("Created players:" + result)
                })

                //Print players of the created Team
                mainContract.methods.getSkillsOfPlayersInTeam(teamId).call({ from: addresses[0], gasPrice: 20000000000, gas: 6721975 }, function (error, result) {
                    console.log("Player's skills:" + result)
                })
            }
        })
}

//Plays a match and return the results
function playMatch(local, visitor, res) {

    var randomSeed = 5 //Ask Toni for random seed
    mainContract.methods.playGame(local, visitor, randomSeed).call({ from: addresses[0], gasPrice: 20000000000, gas: 6721975 }, function (error, result) {
        if (error) {
            console.log(error)
        }
        console.log("Result match:" + result)
        res.send(result)
    })
}

function getPlayersOfTeam(teamId, res) {

    mainContract.methods.getSkillsOfPlayersInTeam(teamId).call({ from: addresses[0], gasPrice: 20000000000, gas: 6721975 }, function (error, result) {
        if (error) {
            console.log(error)
        }
        console.log(result)
        res.send(result)
    })

}

//Start express server
app.listen(listeningPort, () => console.log('CryptoServer listening on port ' + listeningPort + '!'))

//Routing handlers
app.get('/', (req, res) => res.send('The Server is running!'))

app.get('/ganache', function (req, res) {

    console.log('recieved ganache call')
    setResponse(res)
    res.send(addresses)
})

app.get('/teamMembers', function (req, res) {
    console.log('recieved team members call')
})

app.get('/createTeam', function (req, res) {
    setResponse(res)
    createTeam(res)
})

app.get('/createdTeams', function (req, res) {
    setResponse(res)
    res.send(createdTeams.toString())
})

app.get('/playersOfTeam', function (req, res) {
    setResponse(res)
    var teamId = req.query.teamId
    getPlayersOfTeam(teamId, res)
})

app.get('/deployedContractAddress', function (req, res) {
    setResponse(res)
    res.send(contractAddress)
})

app.get('/playMatch', function (req, res) {

    setResponse(res)

    var local = parseInt(req.query.localTeam)
    var visitor = parseInt(req.query.visitorTeam)

    playMatch(local, visitor, res)
})

function setResponse(res) {

    res.set('Access-Control-Allow-Origin', '*');

}