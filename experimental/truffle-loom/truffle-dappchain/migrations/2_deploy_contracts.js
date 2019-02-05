const { writeFileSync, readFileSync } = require('fs')

const PlayersDappChain = artifacts.require('PlayersDappChain')
const TeamsDappChain = artifacts.require('TeamsDappChain')

module.exports = (deployer, network, accounts) => {
  const gatewayAddress = readFileSync('../gateway_dappchain_address', 'utf-8')

  deployer.deploy(PlayersDappChain, gatewayAddress).then(async () => {
    const PlayersDAppChainInstance = await PlayersDappChain.deployed()
    console.log(`PlayersDAppChain deployed at address: ${PlayersDAppChainInstance.address}`)
    writeFileSync('../crypto_cards_dappchain_address', PlayersDAppChainInstance.address)
  })

  deployer.deploy(TeamsDappChain, gatewayAddress).then(async () => {
    const TeamsDAppChainInstance = await TeamsDappChain.deployed()
    console.log(`TeamsDAppChain deployed at address: ${TeamsDAppChainInstance.address}`)
    writeFileSync('../crypto_cards_dappchain_address', TeamsDAppChainInstance.address)
  })
}
