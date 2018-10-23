const { writeFileSync } = require('fs')

const CryptoPlayers = artifacts.require('CryptoPlayers')
const CryptoTeams = artifacts.require('CryptoTeams');
const Gateway = artifacts.require('Gateway')

module.exports = (deployer, _network, accounts) => {
  const [_, user] = accounts
  const validator = accounts[9]
  deployer.deploy(Gateway, [validator], 3, 4).then(async () => {
    const gatewayInstance = await Gateway.deployed()

    console.log(`Gateway deployed at address: ${gatewayInstance.address}`)

    const cryptoPlayersContract = await deployer.deploy(CryptoPlayers, gatewayInstance.address)
    const cryptoPlayersInstance = await CryptoPlayers.deployed()

    console.log(`CryptoPlayers deployed at address: ${cryptoPlayersInstance.address}`)
    console.log(`CryptoPlayers transaction at hash: ${cryptoPlayersContract.transactionHash}`)

    const cryptoTeamsContract = await deployer.deploy(CryptoTeams, gatewayInstance.address)
    const cryptoTeamsInstance = await CryptoTeams.deployed()

    console.log(`CryptoTeams deployed at address: ${cryptoTeamsInstance.address}`)
    console.log(`CryptoTeams transaction at hash: ${cryptoTeamsContract.transactionHash}`)

    await gatewayInstance.toggleToken(cryptoPlayersInstance.address, { from: validator })
    await gatewayInstance.toggleToken(cryptoTeamsInstance.address, { from: validator })
    await cryptoPlayersInstance.register(user)

    writeFileSync('../gateway_address', gatewayInstance.address)
    writeFileSync('../crypto_players_address', cryptoPlayersInstance.address)
    writeFileSync('../crypto_players_tx_hash', cryptoPlayersContract.transactionHash)
    writeFileSync('../crypto_teams_address', cryptoTeamsInstance.address)
    writeFileSync('../crypto_teams_tx_hash', cryptoTeamsContract.transactionHash)
  })
}
