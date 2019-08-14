import random
import numpy as np
import copy
import datetime
from os import listdir, makedirs
from os.path import isfile, join, exists

N_ROUNDS = 100
ALPHA = 0.75
PLAYERS_PER_TEAM = 18
TEAMS_PER_LEAGUE = 8
N_LEAGUES = 16*10
N_TEAMS = N_LEAGUES * TEAMS_PER_LEAGUE
SK_START    = 50 * PLAYERS_PER_TEAM
SK_LOW      = int(SK_START / 1.2)
SK_HIGH     = int(SK_START * 1.2)
RESULT_WINS1 = 0
RESULT_WINS2 = 2
RESULT_TIE = 1
PERF_POINTS = [-8, -5, -3, 0, 2, 5, 8, 10]

def probOverResults(sk1, sk2):
    if sk1/sk2 > 1.1:
        probWins1   = 0.85
        probTie     = 0.05
    elif sk2/sk1 > 1.1:
        probWins1   = 1-0.85
        probTie     = 0.05
    else:
        probWins1   = 0.3
        probTie     = 1 - 2*probWins1
    probWins2 = 1 - probWins1 - probTie
    return np.array([probWins1, probTie, probWins2])

# result = 0 if sk1 wins, 1 if tie, 2 if sk2 wins
def updateSkills(sk1, sk2, result):
    if result == RESULT_TIE:
        return sk1, sk2

    ratingDiff = sk1 - sk2
    winnerWasBetter = (ratingDiff > 0 and result == RESULT_WINS1) or (ratingDiff < 0 and result == RESULT_WINS2)

    if ratingDiff == 0:
        perfPoints = 5
    elif winnerWasBetter:
        perfPoints = 4
    else:
        perfPoints = 6

    if result == RESULT_WINS1:
        return sk1 + perfPoints, sk2
    else:
        return sk1, sk2 + perfPoints
    assert False, "we should never reach here"

def playGameAndUpdateSkills(sk1, sk2):
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

# testing recursively getting better than the other guy
sk1 = 1.2
sk2 = 1
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

allSkills = 1.0 * np.random.randint(low=40, high=70, size= N_TEAMS)
avgPoints = 100.0 * np.ones(N_TEAMS)

for round in range(N_ROUNDS):
    for league in range(N_LEAGUES):
        perfPoints = playLeague(allSkills[league*TEAMS_PER_LEAGUE:(league+1)*TEAMS_PER_LEAGUE])
        avgPoints[league*TEAMS_PER_LEAGUE:(league+1)*TEAMS_PER_LEAGUE] += perfPoints
    newOrder = np.argsort(avgPoints)
    allSkills = allSkills[newOrder]
    avgPoints = avgPoints[newOrder]


print(allSkills)
print(avgPoints)








