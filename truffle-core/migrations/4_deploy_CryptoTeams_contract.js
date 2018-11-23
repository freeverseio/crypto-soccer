const CryptoTeams = artifacts.require('CryptoTeams');

module.exports = (deployer) => {
  const name = "CryptoSoccerTeams";
  const symbol = "CSP";

  deployer.deploy(CryptoTeams, name, symbol)
    .then(instance => {
      console.log(`CryptoTeams deployed at address: ${instance.address}`)
      console.log(`CryptoTeams transaction at hash: ${instance.transactionHash}`)
    });
}
