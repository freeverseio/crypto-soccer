const CryptoPlayers = artifacts.require("CryptoPlayers");
const CryptoTeams = artifacts.require("CryptoTeams");

module.exports = function (deployer) {
  deployer.deploy(CryptoPlayers)
    .then(cryptoPlayers => {
      console.log(`CryptoPlayers deployed at address: ${cryptoPlayers.address}`);
      console.log(`CryptoPlayers transaction at hash: ${cryptoPlayers.transactionHash}`);
      deployer.deploy(CryptoTeams)
        .then(cryptoTeams => {
          console.log(`CryptoTeams deployed at address: ${cryptoTeams.address}`);
          console.log(`CryptoTeams transaction at hash: ${cryptoTeams.transactionHash}`);
        })
        .catch(console.error);
    })
    .catch(console.error);
};

