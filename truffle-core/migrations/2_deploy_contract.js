const CryptoPlayers = artifacts.require("CryptoPlayers");
const CryptoTeams = artifacts.require("CryptoTeams");
const Horizon = artifacts.require("Horizon");

module.exports = function (deployer) {
  deployer.deploy(CryptoPlayers)
    .then(cryptoPlayers => {
      console.log(`CryptoPlayers deployed at address: ${cryptoPlayers.address}`);
      console.log(`CryptoPlayers transaction at hash: ${cryptoPlayers.transactionHash}`);
      return deployer.deploy(CryptoTeams)
        .then(cryptoTeams => {
          console.log(`CryptoTeams deployed at address: ${cryptoTeams.address}`);
          console.log(`CryptoTeams transaction at hash: ${cryptoTeams.transactionHash}`);
          return deployer.deploy(Horizon, cryptoPlayers.address, cryptoTeams.address)
            .then(horizon => {
              console.log(`Horizon deployed at address: ${horizon.address}`);
              console.log(`Horizon transaction at hash: ${horizon.transactionHash}`);
            })
            .catch(console.error);
        })
        .catch(console.error);
    })
    .catch(console.error);
};

