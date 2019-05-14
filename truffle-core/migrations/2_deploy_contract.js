// const Players = artifacts.require("Players");
// const Teams = artifacts.require("Teams");
// const Horizon = artifacts.require("Horizon");

module.exports = function (deployer) {
  deployer.then(async () => {
      return;

      const players = await deployer.deploy(Players);
      console.log(`Players deployed at address: ${players.address}`);
      console.log(`Players transaction at hash: ${players.transactionHash}`);

      const basePlayersURI = "https://www.freeverse.io/api/player/";
      await players.setBaseTokenURI(basePlayersURI);
      console.log("Players base URI: " + basePlayersURI);

      const teams = await deployer.deploy(Teams, players.address);
      console.log(`Teams deployed at address: ${teams.address}`);
      console.log(`Teams transaction at hash: ${teams.transactionHash}`);

      const baseTeamsURI = "https://www.freeverse.io/api/team/";
      await teams.setBaseTokenURI(baseTeamsURI);
      console.log("Teams base URI: " + baseTeamsURI);

      await players.addTeamsContract(Teams.address);
      console.log("Teams can change players ownership");

      await players.renounceTeamsContract();
      console.log("Deployer can't change players ownership");

      const horizon = await deployer.deploy(Horizon, Teams.address);
      console.log(`Horizon deployed at address: ${horizon.address}`);
      console.log(`Horizon transaction at hash: ${horizon.transactionHash}`);

      await players.addMinter(horizon.address);
      console.log("Horizon can mint players");

      await players.renounceMinter();
      console.log("Deployer can't mint players");

      await teams.addMinter(horizon.address);
      console.log("Horizon can mint teams");

      await teams.renounceMinter();
      console.log("Deployer can't mint teams");
    })
    .catch(console.error);
};

