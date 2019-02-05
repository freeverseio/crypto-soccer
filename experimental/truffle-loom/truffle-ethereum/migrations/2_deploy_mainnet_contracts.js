const { writeFileSync } = require('fs')

const Players = artifacts.require('Players')
const Teams = artifacts.require('Teams');
const Gateway = artifacts.require('Gateway')

module.exports = (deployer, _network, accounts) => {
  const [_, user] = accounts
  const validator = accounts[9]
  deployer.deploy(Gateway, [validator], 3, 4).then(async () => {
    const gatewayInstance = await Gateway.deployed()

    console.log(`Gateway deployed at address: ${gatewayInstance.address}`)

    const PlayersContract = await deployer.deploy(Players, gatewayInstance.address)
    const PlayersInstance = await Players.deployed()

    console.log(`Players deployed at address: ${PlayersInstance.address}`)
    console.log(`Players transaction at hash: ${PlayersContract.transactionHash}`)

    const TeamsContract = await deployer.deploy(Teams, gatewayInstance.address)
    const TeamsInstance = await Teams.deployed()

    console.log(`Teams deployed at address: ${TeamsInstance.address}`)
    console.log(`Teams transaction at hash: ${TeamsContract.transactionHash}`)

    await gatewayInstance.toggleToken(PlayersInstance.address, { from: validator })
    await gatewayInstance.toggleToken(TeamsInstance.address, { from: validator })
    await PlayersInstance.register(user)

    writeFileSync('../gateway_address', gatewayInstance.address)
    writeFileSync('../crypto_players_address', PlayersInstance.address)
    writeFileSync('../crypto_players_tx_hash', PlayersContract.transactionHash)
    writeFileSync('../crypto_teams_address', TeamsInstance.address)
    writeFileSync('../crypto_teams_tx_hash', TeamsContract.transactionHash)
  })
}
