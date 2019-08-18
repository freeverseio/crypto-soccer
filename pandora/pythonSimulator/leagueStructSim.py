import random
import numpy as np
import copy
import datetime
from os import listdir, makedirs
from os.path import isfile, join, exists
# import matplotlib.pyplot as plt
# from mpl_toolkits.mplot3d import Axes3D  # noqa: F401 unused import
import matplotlib.pyplot as plt

N_ROUNDS = 100
INIT_SORT = 0
ALPHA_INTERTIA = 0.4
WEIGHT_SKILLS = 100
WEIGHT_PERF = 1
PLAYERS_PER_TEAM = 18
TEAMS_PER_LEAGUE = 8
GAMES_PER_LEAGUE = TEAMS_PER_LEAGUE * (TEAMS_PER_LEAGUE-1)
N_LEAGUES = 16*5
N_TEAMS = N_LEAGUES * TEAMS_PER_LEAGUE
SK_START    = 1000 * PLAYERS_PER_TEAM
SK_LOW      = int(SK_START / 12)
SK_HIGH     = int(SK_START * 12)
# SK_LOW      = int(SK_START / 2)
# SK_HIGH     = int(SK_START * 2)
RESULT_WINS1 = 0
RESULT_WINS2 = 2
RESULT_TIE = 1
PERF = 1
PERF_POINTS = PERF * np.array([-8, -5, -3, 0, 2, 5, 8, 10])
MAX_LEAGUES_PLAYER = 18
MAX_GAMES_PLAYER = MAX_LEAGUES_PLAYER * GAMES_PER_LEAGUE
MAX_PERPOINTS_PLAYER = 12000
AVG_PERPOINTS_PER_GAME = MAX_PERPOINTS_PLAYER // MAX_GAMES_PLAYER * 3
PERPOINTS_RANGE = int(0.8*AVG_PERPOINTS_PER_GAME)
MAX_DICE_RAND = 16383

def probOverResults(sk1, sk2):
    probWins1 = min(0.9, 1/3*sk1/sk2)
    probWins2 = min(0.9, 1/3*sk2/sk1)
    probTie = 1 - probWins1 - probWins2
    return np.array([probWins1, probTie, probWins2])

# result = 0 if sk1 wins, 1 if tie, 2 if sk2 wins
def updateSkills(sk1, sk2, result):
    if result == RESULT_TIE:
        return sk1, sk2

    if sk1 >= sk2:
        extra = min(3, sk1/sk2)
        winnerWasBetter = (result == RESULT_WINS1)
    else:
        extra = min(3, sk2/sk1)
        winnerWasBetter = (result == RESULT_WINS2)

    if winnerWasBetter:
        perfPoints = AVG_PERPOINTS_PER_GAME - PERPOINTS_RANGE * min(1, 0.3 * extra)
    else:
        perfPoints = AVG_PERPOINTS_PER_GAME + PERPOINTS_RANGE * min(1, 0.3 * extra)

    if result == RESULT_WINS1:
        return sk1 + perfPoints, sk2
    else:
        return sk1, sk2 + perfPoints


def throwDiceArray(weights, rndNum, maxRndNum):
    uniformRndInSumOfWeights = sum(weights) * rndNum
    cumSum = 0
    for w in range(len(weights)):
        cumSum += weights[w]
        if ( uniformRndInSumOfWeights < ( cumSum * (maxRndNum-1) )):
            return w
    return w

def playGameAndUpdateSkills(sk1, sk2):
    results = [RESULT_WINS1, RESULT_TIE, RESULT_WINS2]
    probResults = probOverResults(sk1, sk2)
    result = throwDiceArray(probResults, np.random.randint(0, MAX_DICE_RAND), MAX_DICE_RAND)
    (newSk1, newSk2) = updateSkills(sk1, sk2, result)
    if result == RESULT_WINS1:
        leaguePoints1 = 3
        leaguePoints2 = 0
    elif result == RESULT_WINS2:
        leaguePoints1 = 0
        leaguePoints2 = 3
    else:
        leaguePoints1 = 1
        leaguePoints2 = 1
    return newSk1, newSk2, leaguePoints1, leaguePoints2


def playGameAndUpdateSkillsAvg(sk1, sk2):
    results = [RESULT_WINS1, RESULT_TIE, RESULT_WINS2]
    probResults = probOverResults(sk1, sk2)
    newSk1 = 0
    newSk2 = 0
    leaguePoints1 = 3 * probResults[RESULT_WINS1] + probResults[RESULT_TIE]
    leaguePoints2 = 3 * probResults[RESULT_WINS2] + probResults[RESULT_TIE]
    for r in results:
        (nSk1, nSk2) = updateSkills(sk1, sk2, r)
        newSk1 += probResults[r] * nSk1
        newSk2 += probResults[r] * nSk2
    return newSk1, newSk2, leaguePoints1, leaguePoints2


def shiftBack(t, nTeams):
    if (t < nTeams):
        return t
    else:
        return t-(nTeams-1)


def getTeamsInMatchFirstHalf(matchday, match, nTeams):
    team1 = 0
    if (match > 0):
        team1 = shiftBack(nTeams-match+matchday, nTeams)

    team2 = shiftBack(match+1+matchday, nTeams)
    if ( (matchday % 2) == 0):
        return team1, team2
    else:
        return team2, team1


def getTeamsInMatch(matchday, match, nTeams):
    assert matchday < 2 * (nTeams - 1), "This league does not have so many matchdays"
    if (matchday < (nTeams - 1)):
        (team1, team2) = getTeamsInMatchFirstHalf(matchday, match, nTeams)
    else:
        (team2, team1) = getTeamsInMatchFirstHalf(matchday - (nTeams - 1), match, nTeams);
    return team1, team2



print(probOverResults(2,1))
print(probOverResults(1,1))
print(probOverResults(1,2))
print(probOverResults(1,4))
print(probOverResults(4,1))
print(sum(probOverResults(4,1)))
print(playGameAndUpdateSkills(1,1))
print(playGameAndUpdateSkills(1.2,1))

print(" testing recursively getting better than the other guy")
sk1 = 1010
sk2 = 1000
for t in range(10):
    newSk1, newSk2, leaguePoints1, leaguePoints2 = playGameAndUpdateSkills(sk1, sk2)
    print(newSk1, newSk2)
    sk1 = newSk1
    sk2 = newSk2


def playLeague(skills):
    nMatchDays = 2*(TEAMS_PER_LEAGUE-1)
    gamesPerMatchday = TEAMS_PER_LEAGUE//2
    leaguePoints = np.ones(TEAMS_PER_LEAGUE, int)
    for matchday in range(nMatchDays):
        for match in range(gamesPerMatchday):
            t1, t2 = getTeamsInMatch(matchday, match, TEAMS_PER_LEAGUE)
            skills[t1], skills[t2], leaguePoints1, leaguePoints2 = playGameAndUpdateSkills(skills[t1], skills[t2])
            leaguePoints[t1] += leaguePoints1
            leaguePoints[t2] += leaguePoints2

    classification = np.argsort(leaguePoints)
    perfPoints = np.zeros(TEAMS_PER_LEAGUE)
    perfPoints[classification] = PERF_POINTS
    return perfPoints

def setPlot(ax, xlabel, ylabel, title):
    ax.legend()
    ax.set_xlabel(xlabel)
    ax.set_ylabel(ylabel)
    ax.set_title(title, fontsize = 7)
    # ax.text(0, 1200, subtitle)
    ax.grid(True,'both')
    # ax.set_ylim(0,1)
    # ax.set_yticks([0.1*t for t in range(11)], True)

def computeQualities(allSkills):
    quality = np.zeros(N_LEAGUES)
    for league in range(N_LEAGUES):
        quality[league] = sum(allSkills[league*TEAMS_PER_LEAGUE:(league+1)*TEAMS_PER_LEAGUE])
    return quality

def overallRating(allSkills, avgPoints):
    return WEIGHT_SKILLS * allSkills/SK_START + WEIGHT_PERF * avgPoints

allSkills = 1.0 * np.random.randint(low=SK_LOW, high=SK_START, size= N_TEAMS)
avgPoints = np.zeros(N_TEAMS)
prevOrder = range(N_TEAMS)
if INIT_SORT == 1:
    newOrder = np.argsort(overallRating(allSkills, avgPoints))
    allSkills = allSkills[newOrder]
else:
    newOrder = np.argsort(prevOrder)
rankings = np.zeros([N_ROUNDS+1, N_TEAMS])
qualities = np.zeros([N_ROUNDS+1, N_LEAGUES])
rankings[0, :] = newOrder[prevOrder]
qualities[0, :] = computeQualities(allSkills)

for round in range(N_ROUNDS):
    for league in range(N_LEAGUES):
        perfPoints = playLeague(allSkills[league*TEAMS_PER_LEAGUE:(league+1)*TEAMS_PER_LEAGUE])
        avgPoints[league*TEAMS_PER_LEAGUE:(league+1)*TEAMS_PER_LEAGUE] = ALPHA_INTERTIA * perfPoints + (1-ALPHA_INTERTIA) * avgPoints[league*TEAMS_PER_LEAGUE:(league+1)*TEAMS_PER_LEAGUE]
    newOrder = np.argsort(overallRating(allSkills, avgPoints))
    allSkills = allSkills[newOrder]
    avgPoints = avgPoints[newOrder]
    rankings[round+1, :] = newOrder[prevOrder]
    qualities[round+1, :] = computeQualities(allSkills)
    prevOrder = newOrder


fig, ax = plt.subplots()
for round in range(0, N_ROUNDS+1,N_ROUNDS//10):
    ax.plot(qualities[round,:]/TEAMS_PER_LEAGUE/PLAYERS_PER_TEAM)
plotname = '(initSort,inertia)=(%s,%s),(low,high)=(%s,%s),(Perf,wSk,wPerf)=(%s,%s,%s)' %(ALPHA_INTERTIA, INIT_SORT, int(SK_LOW/PLAYERS_PER_TEAM),int(SK_HIGH/PLAYERS_PER_TEAM), PERF, WEIGHT_SKILLS, WEIGHT_PERF)

setPlot(ax,
        'league number',
        'Average quality per player',
        'Quality at various rounds,' + plotname
)
plt.savefig('qualities'+plotname+'.png')
plt.close(fig)


fig, ax = plt.subplots()
for team in range(0, N_TEAMS, N_TEAMS//3):
    ax.plot(rankings[:,team])
ax.plot(rankings[:, N_TEAMS-1])
setPlot(ax,
        'league number',
        'League pos',
        'Evolution for various teams' + plotname
        )
plt.savefig('leagueEvoPerTeam'+plotname+'.png')
plt.close(fig)



a=2+2








