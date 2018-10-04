const { writeFileSync, readFileSync } = require('fs')

const CryptoPlayersDappChain = artifacts.require('CryptoPlayersDappChain')
const CryptoTeamsDappChain = artifacts.require('CryptoTeamsDappChain')

module.exports = (deployer, network, accounts) => {
  const gatewayAddress = readFileSync('../gateway_dappchain_address', 'utf-8')

  deployer.deploy(CryptoPlayersDappChain, gatewayAddress).then(async () => {
    const cryptoPlayersDAppChainInstance = await CryptoPlayersDappChain.deployed()
    console.log(`CryptoPlayersDAppChain deployed at address: ${cryptoPlayersDAppChainInstance.address}`)
    writeFileSync('../crypto_cards_dappchain_address', cryptoPlayersDAppChainInstance.address)
  })

  deployer.deploy(CryptoTeamsDappChain, gatewayAddress).then(async () => {
    const cryptoTeamsDAppChainInstance = await CryptoTeamsDappChain.deployed()
    console.log(`CryptoTeamsDAppChain deployed at address: ${cryptoTeamsDAppChainInstance.address}`)
    writeFileSync('../crypto_cards_dappchain_address', cryptoTeamsDAppChainInstance.address)
  })
}
