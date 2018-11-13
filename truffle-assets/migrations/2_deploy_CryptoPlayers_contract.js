const CryptoPlayers = artifacts.require('CryptoPlayers')

module.exports = (deployer) => {
  const name = "CryptoSoccerPlayers";
  const symbol = "CSP";

  deployer.deploy(CryptoPlayers, name, symbol)
    .then(instance => {
      console.log(`CryptoPlayers deployed at address: ${instance.address}`)
      console.log(`CryptoPlayers transaction at hash: ${instance.transactionHash}`)
    });
}
