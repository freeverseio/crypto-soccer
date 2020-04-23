from gql import gql, Client
from gql.transport.requests import RequestsHTTPTransport
import json

devURL  = 'https://k8s.gorengine.com/auth/'
prodURL = 'https://k8s.goalrevolution.live/auth/'

def setClient(url):
    sample_transport = RequestsHTTPTransport(
        url=url,
        use_json=True,
        headers={
            "Content-type": "application/json",
            "Authorization": "Bearer joshua",
        },
        verify=False
    )
    return Client(
        retries=3,
        transport=sample_transport,
        fetch_schema_from_transport=True,
    )

def saveAllPlayers(client, filename):
    # allPlayers(condition: { teamId: "1099511627776" }, orderBy: SHOOT_DESC) {
    query = gql('''
        query {
            allPlayers {
                nodes {
                  name
                  preferredPosition
                  potential
                  shoot   
                  defence
                  pass
                  endurance
                  speed
                  teamId
                }
            }
        }
    ''')
    result = client.execute(query)
    with open(filename, 'w') as outfile:
        json.dump(result, outfile)

def getSavedPlayers(filename):
    with open(filename) as json_file:
        result = json.load(json_file)
    return result

client = setClient(prodURL)

# saveAllPlayers(client, "allplayersprod.txt")

players = getSavedPlayers("allplayersprod.txt")






