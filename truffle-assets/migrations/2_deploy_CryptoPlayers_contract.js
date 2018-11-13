const CryptoPlayers = artifacts.require('CryptoPlayers')

module.exports = (deployer) => {
  const name = "name";
  const symbol = "symbol";

  deployer.deploy(CryptoPlayers, name, symbol)
    .then(instance => {
      console.log(`CryptoPlayers deployed at address: ${instance.address}`)
      console.log(`CryptoPlayers transaction at hash: ${instance.transactionHash}`)
    });
}
