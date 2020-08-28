const { request, gql } = require('graphql-request');
const { horizonConfig } = require('../config.js');

class HorizonService {
  constructor() {
    this.endpoint = horizonConfig.url;
  }

  async getTeamOwner({ teamId }) {
    const query = gql`
    {
        teamByTeamId(teamId: "${teamId}"){ owner }
    }
    `;
    const result = await request(this.endpoint, query);

    return result && result.teamByTeamId && result.teamByTeamId.owner
      ? result.teamByTeamId.owner
      : '';
  }

  async getPlayerOwner({ playerId }) {
    const query = gql`
    {
        playerByPlayerId(playerId: "${playerId}"){ 
            teamId
              teamByTeamId {
              owner
            }
        }
    }
    `;
    const result = await request(this.endpoint, query);

    return result &&
      result.playerByPlayerId &&
      result.playerByPlayerId.teamByTeamId &&
      result.playerByPlayerId.teamByTeamId.owner
      ? result.playerByPlayerId.teamByTeamId.owner
      : '';
  }

  async getAllLeaguesLeaderboard() {
    const query = gql`
      {
        allLeagues(orderBy: LEAGUE_IDX_DESC) {
          nodes {
            countryIdx
            leagueIdx
            timezoneIdx
            teamsByTimezoneIdxAndCountryIdxAndLeagueIdx(
              orderBy: LEADERBOARD_POSITION_ASC
            ) {
              nodes {
                teamId
                name
                owner
                w
                d
                l
                goalsForward
                goalsAgainst
                points
                rankingPoints
              }
            }
            counter: matchesByTimezoneIdxAndCountryIdxAndLeagueIdx(
              filter: {
                or: [
                  { state: { equalTo: END } }
                  { state: { equalTo: CANCELLED } }
                ]
              }
            ) {
              totalCount
            }
            matchesByTimezoneIdxAndCountryIdxAndLeagueIdx(
              first: 1
              orderBy: START_EPOCH_DESC
            ) {
              nodes {
                startEpoch
              }
            }
            lastlast: matchesByTimezoneIdxAndCountryIdxAndLeagueIdx(
              first: 1
              offset: 4
              orderBy: START_EPOCH_DESC
            ) {
              nodes {
                startEpoch
              }
            }
          }
        }
      }
    `;
    const result = await request(this.endpoint, query);

    return result && result.allLeagues && result.allLeagues.nodes
      ? result.allLeagues.nodes
      : '';
  }
}

module.exports = new HorizonService();
