const HorizonService = require('../services/HorizonService.js')
const { updateTeamManagerName, updateTeamName } = require('../repositories/index.js')

const fillGameDb = async () => {
    const teams = await HorizonService.getAllUsersTeam()
    let teamsCount = 0
    for(let team of teams) {
        console.log(`Processing ${teamsCount} of total of ${teams.length} teams`)
        if (team.managerName) {
            await updateTeamManagerName({ teamId: team.teamId, teamManagerName: team.managerName })
        }

        if (team.name) {
            await updateTeamName({ teamId: team.teamId, teamName: team.name})
        }
        teamsCount++
    }
};

fillGameDb()