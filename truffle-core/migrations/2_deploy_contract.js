const CryptoPlayers = artifacts.require("CryptoPlayers");
const CryptoTeams = artifacts.require("CryptoTeams");
const Horizon = artifacts.require("Horizon");

module.exports = function (deployer) {
  deployer.then(async () => {
      const cryptoPlayers = await deployer.deploy(CryptoPlayers);
      console.log(`CryptoPlayers deployed at address: ${cryptoPlayers.address}`);
      console.log(`CryptoPlayers transaction at hash: ${cryptoPlayers.transactionHash}`);

      const basePlayersURI = "https://www.freeverse.io/api/player/";
      await cryptoPlayers.setBaseTokenURI(basePlayersURI);
      console.log("CryptoPlayers base URI: " + basePlayersURI);

      const cryptoTeams = await deployer.deploy(CryptoTeams, cryptoPlayers.address);
      console.log(`CryptoTeams deployed at address: ${cryptoTeams.address}`);
      console.log(`CryptoTeams transaction at hash: ${cryptoTeams.transactionHash}`);

      const baseTeamsURI = "https://www.freeverse.io/api/team/";
      await cryptoTeams.setBaseTokenURI(baseTeamsURI);
      console.log("CryptoTeams base URI: " + baseTeamsURI);

      await cryptoPlayers.addTeamsContract(cryptoTeams.address);
      console.log("CryptoTeams can change players ownership");

      await cryptoPlayers.renounceTeamsContract();
      console.log("Deployer can't change players ownership");

      const horizon = await deployer.deploy(Horizon, cryptoTeams.address);
      console.log(`Horizon deployed at address: ${horizon.address}`);
      console.log(`Horizon transaction at hash: ${horizon.transactionHash}`);

      await cryptoPlayers.addMinter(horizon.address);
      console.log("Horizon can mint players");

      await cryptoPlayers.renounceMinter();
      console.log("Deployer can't mint players");

      await cryptoTeams.addMinter(horizon.address);
      console.log("Horizon can mint teams");

      await cryptoTeams.renounceMinter();
      console.log("Deployer can't mint teams");
    })
    .catch(console.error);
};

