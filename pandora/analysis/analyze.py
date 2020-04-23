from gql import gql, Client
from gql.transport.requests import RequestsHTTPTransport
import json
import matplotlib.pyplot as plt
import numpy as np

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

def plotSkill(name):
    skill = [p[name] for p in players["allPlayers"]["nodes"]]
    fig = plt.figure()
    plt.hist(skill, bins = 50)
    fig.suptitle('All players at end of League 1', fontsize=20)
    plt.xlim(xmin=0, xmax=3500)
    # plt.ylim(ymin=0, ymax=100)
    plt.xlabel(name, fontsize=18)
    plt.ylabel('Num Players', fontsize=18)
    fig.savefig(name + '.png')
    # plt.show()

def getTopN(name, n):
    skill = np.array([p[name] for p in players["allPlayers"]["nodes"]])
    sortedIdxs = np.argsort(skill)
    return (skill[sortedIdxs[-n:]], sortedIdxs)


client = setClient(prodURL)
# saveAllPlayers(client, "allplayersprod.txt")

players = getSavedPlayers("allplayersprod.txt")

skillNames = ["shoot", "pass", "defence", "speed", "endurance"]

########################
# PLOT Skills Histograms
########################
# for name in skillNames:
#     plotSkill(name)


########################
# PLOT Top20 players
########################
fig = plt.figure()
for name in skillNames:
    fig.suptitle('Top 20 players, per skill, at end of League 1', fontsize=20)
    (top, topIdxs) = getTopN(name, 20)
    plt.plot(top)
plt.ylabel('skill value', fontsize=18)
plt.legend((skillNames),
           loc='upper left')
fig.savefig('tops' + '.png')

########################
# PLOT Level vs potential
########################

potential = np.array([p["potential"] for p in players["allPlayers"]["nodes"]])
sumSkills = 0 * potential

for name in skillNames:
    sumSkills += np.array([p[name] for p in players["allPlayers"]["nodes"]])
from gql import gql, Client
from gql.transport.requests import RequestsHTTPTransport
import json
import matplotlib.pyplot as plt
import numpy as np

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

def saveAllPlayers(url, filename):
    # allPlayers(condition: { teamId: "1099511627776" }, orderBy: SHOOT_DESC) {
    client = setClient(url)
    query = gql('''
        query {
            allPlayers {
                nodes {
                  name
                  dayOfBirth
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

def plotSkill(name, fixXrange = True):
    skill = [p[name] for p in players["allPlayers"]["nodes"]]
    fig = plt.figure()
    plt.hist(skill, bins = 50)
    fig.suptitle('All players at end of League 1', fontsize=20)
    if fixXrange:
        plt.xlim(xmin=0, xmax=3500)
    # plt.ylim(ymin=0, ymax=100)
    plt.xlabel(name, fontsize=18)
    plt.ylabel('Num Players', fontsize=18)
    fig.savefig(name + '.png')
    # plt.show()

def getTopN(name, n):
    skill = np.array([p[name] for p in players["allPlayers"]["nodes"]])
    sortedIdxs = np.argsort(skill)
    return (skill[sortedIdxs[-n:]], sortedIdxs)


# saveAllPlayers(prodURL, "allplayersprod.txt")

players = getSavedPlayers("allplayersprod.txt")

skillNames = ["shoot", "pass", "defence", "speed", "endurance"]

########################
# PLOT Skills Histograms
########################
if False:
    plotSkill("dayOfBirth", False)
    for name in skillNames:
        plotSkill(name)


########################
# PLOT Top20 players
########################
if False:
    fig = plt.figure()
    for name in skillNames:
        (top, topIdxs) = getTopN(name, 20)
        plt.plot(top)
    plt.ylabel('skill value', fontsize=18)
    plt.legend((skillNames),
               loc='upper left')
    fig.suptitle('Top 20 players, per skill, at end of League 1', fontsize=20)
    fig.savefig('tops' + '.png')

########################
# PLOT Level vs potential
########################

avgSkills = np.zeros(10)

for potential in range(10):
    playersThisPot = [p for  p in players["allPlayers"]["nodes"] \
                      if (p["potential"] == potential) and (p["dayOfBirth"] > 17000) \
    ]
    sumSkills = 0
    for name in skillNames:
        sumSkills += np.array([p[name] for p in playersThisPot]).sum()
    avgSkills[potential] = sumSkills/(len(playersThisPot)*5)

fig = plt.figure()

potential = np.array(range(10))
plt.plot(potential, avgSkills)
plt.ylabel('skill value', fontsize=18)
fig.suptitle('Average players skill as function of potential, at end of League 1', fontsize=20)
fig.savefig('skillVsPotential' + '.png')


















