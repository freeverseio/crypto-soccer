const TeamFactory = artifacts.require('TeamFactory');
const CryptoPlayers = artifacts.require('CryptoPlayers')
const CryptoTeams = artifacts.require('CryptoTeams');

module.exports = (deployer) => {
  deployer.deploy(TeamFactory)
    .then(teamFactory => {
      console.log(`TeamFactory deployed at address: ${teamFactory.address}`);
      console.log(`TeamFactory transaction at hash: ${teamFactory.transactionHash}`);
      deployer.deploy(CryptoPlayers)
        .then(cryptoPlayers => {
          console.log(`CryptoPlayers deployed at address: ${cryptoPlayers.address}`);
          console.log(`CryptoPlayers transaction at hash: ${cryptoPlayers.transactionHash}`);
        });
      deployer.deploy(CryptoTeams, teamFactory.address)
        .then(cryptoTeams => {
          console.log(`CryptoTeams deployed at address: ${cryptoTeams.address}`);
          console.log(`CryptoTeams transaction at hash: ${cryptoTeams.transactionHash}`);
        });
    });
};
