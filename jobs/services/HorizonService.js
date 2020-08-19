const { request, gql } = require('graphql-request')
const { horizonConfig } = require('../config.js')

class HorizonService {
  constructor() {
    this.endpoint = horizonConfig.url
  }

  async getAllUsersTeam() {
    const query = gql`
    { 
      allTeams(filter: { owner: { notEqualTo:"0x0000000000000000000000000000000000000000" }}){
        nodes{
          teamId
          name
          managerName
          owner
        }
      }
    }
    `
    const result = await request(this.endpoint, query)
    
    return result && result.allTeams && result.allTeams.nodes ? result.allTeams.nodes : ''
  }

}

module.exports = new HorizonService();
