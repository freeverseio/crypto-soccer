const CryptoPlayers = artifacts.require('CryptoPlayers')

module.exports = (deployer) => {
  deployer.deploy(CryptoPlayers)
    .then(instance => {
      console.log(`CryptoPlayers deployed at address: ${instance.address}`)
      console.log(`CryptoPlayers transaction at hash: ${instance.transactionHash}`)
    });
}
