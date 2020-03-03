import { check } from "k6";
import http  from "k6/http";

const url = "https://k8s.goalrevolution.live/auth";
const query = `
query getTrainingStatus($teamId: String!) {
  allTeams(condition: { teamId: $teamId }) {
    nodes {
      trainingPoints
      playersByTeamId(orderBy: PREFERRED_POSITION_DESC) {
        nodes {
          playerId
          name
          preferredPosition
          shirtNumber
        }
      }
    }
  }
  allTrainings(
    orderBy: PRIMARY_KEY_DESC
    first: 1
    condition: { teamId: $teamId }
  ) {
    nodes {
      specialPlayerShirt
      goalkeepersDefence
      goalkeepersSpeed
      goalkeepersPass
      goalkeepersShoot
      goalkeepersEndurance
      defendersDefence
      defendersSpeed
      defendersPass
      defendersShoot
      defendersEndurance
      midfieldersDefence
      midfieldersSpeed
      midfieldersPass
      midfieldersShoot
      midfieldersEndurance
      attackersDefence
      attackersSpeed
      attackersPass
      attackersShoot
      attackersEndurance
      specialPlayerDefence
      specialPlayerSpeed
      specialPlayerPass
      specialPlayerShoot
      specialPlayerEndurance
    }
  }
}`;

export function graphql(query, vars) {
  let payload = JSON.stringify({ query : query , variables : vars });
  let params = { headers : { 
	"Content-Type" : "application/json",
        "Authorization": "Bearer joshua" 
  }};
  let res = http.post(url, payload, params);
  check(res, {
    "is status 200": (r) => r.status === 200
  });
  try {
    let gres = JSON.parse(res.body)
    check(gres, {
      "success" : (r) => !("errors" in r)
    })
    return gres.data;
  } catch (e) {
    console.log("RESP.BODY=[",res.body,"]")
    check(res.body, {
      "success" : false
    })
  }
} 

export function setup() {
  let nodes = graphql("{allTeams{nodes{teamId}}}",{}).allTeams.nodes;
  let teamIds = []
  for (let node of nodes) {
     teamIds.push(node.teamId)
  }
  return { teamIds : teamIds }
}

export default function(data) {
  let teamId = data.teamIds[Math.floor(Math.random()*data.teamIds.length)]; 
  graphql(query,{ "teamId" : teamId })
};
