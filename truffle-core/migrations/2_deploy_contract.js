const TeamFactory = artifacts.require("TeamFactory");
const CryptoPlayers = artifacts.require("CryptoPlayers");
const CryptoTeams = artifacts.require("CryptoTeams");
const Testing = artifacts.require("Testing");
const League = artifacts.require("League");

module.exports = function (deployer) {
  deployer.deploy(TeamFactory)
    .then(teamFactory => {
      console.log(`TeamFactory deployed at address: ${teamFactory.address}`);
      console.log(`TeamFactory transaction at hash: ${teamFactory.transactionHash}`);
      deployer.deploy(CryptoPlayers)
        .then(cryptoPlayers => {
          console.log(`CryptoPlayers deployed at address: ${cryptoPlayers.address}`);
          console.log(`CryptoPlayers transaction at hash: ${cryptoPlayers.transactionHash}`);
        })
        .catch(console.error);
      deployer.deploy(CryptoTeams)
        .then(cryptoTeams => {
          console.log(`CryptoTeams deployed at address: ${cryptoTeams.address}`);
          console.log(`CryptoTeams transaction at hash: ${cryptoTeams.transactionHash}`);
        })
        .catch(console.error);
      deployer.deploy(League, teamFactory.address)
        .then(league => {
          console.log(`League deployed at address: ${league.address}`);
          console.log(`League transaction at hash: ${league.transactionHash}`);
        })
        .catch(console.error);
      deployer.deploy(Testing, teamFactory.address)
        .then(testing => {
          console.log(`Testing deployed at address: ${testing.address}`);
          console.log(`Testing transaction at hash: ${testing.transactionHash}`);
        })
        .catch(console.error);
    })
    .catch(console.error);
};

