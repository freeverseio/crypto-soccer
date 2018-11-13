const CryptoTeams = artifacts.require('CryptoTeams');

module.exports = (deployer) => {
  deployer.deploy(CryptoTeams)
    .then(instance => {
      console.log(`CryptoTeams deployed at address: ${instance.address}`)
      console.log(`CryptoTeams transaction at hash: ${instance.transactionHash}`)
    });
}
