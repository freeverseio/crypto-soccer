const HorizonService = require("../services/HorizonService.js");
const { updatePlayerName } = require("../repositories/index.js");

const fillGameDbPlayerNames = async () => {
  const timezoneIdx = 10;
  const countryIdx = 0;
  const teams = await HorizonService.getAllPlayersFromTimezoneAndCountry({
    timezoneIdx,
    countryIdx,
  });
  let teamsCount = 0;
  for (let team of teams) {
    console.log(`Processing ${teamsCount} of total of ${teams.length} teams`);
    if (team.playersByTeamId && team.playersByTeamId.nodes) {
      for (let player of team.playersByTeamId.nodes) {
        await updatePlayerName({
          playerId: player.name,
          playerName: player.playerId,
        });
      }
    }

    teamsCount++;
  }
};

fillGameDbPlayerNames();
