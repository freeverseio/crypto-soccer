import random
import numpy as np
import copy
from simulate import *
import matplotlib.pyplot as plt

from testingLib import *

def analyzeAllPlayersEqualButDiffTeams(skillValue1, skillValue2, plotOutFile, nGames):
    skills1 = skillValue1 * np.ones(5)
    skills2 = skillValue2 * np.ones(5)
    age = 20
    team1 = createAllPlayersEqualTeam(skills1, age, roles433)
    team2 = createAllPlayersEqualTeam(skills2, age, roles433)
    # showTeam(team1)
    # showTeam(team2)
    return analyzeTeam1AgainstTeam2(nGames, team1, team2, plotOutFile)


def analyzeAllPlayersEqualTeams(skillValue, plotOutFile, nGames):
    skills = skillValue * np.ones(5)
    age = 20
    team = createAllPlayersEqualTeam(skills, age, roles433)
    analyzeTeam1AgainstTeam2(nGames, team, team, plotOutFile)


def analyzeTeam1AgainstTeam2(nGames, team1, team2, plotOutFile):
    goals1, goals2 = playNGames(nGames, team1, team2)
    print "Average Goals T1: %.2f,  Distr: %s" %(goals1.mean(), printHistogram(goals1))
    print "Average Goals T2: %.2f,  Distr: %s" %(goals2.mean(), printHistogram(goals2))
    return goals1, goals2


def printHistogram(goals):
    str = '|  '
    for (g, hist) in enumerate(getHistogram(goals)):
        str += "%d: %.1f  |  " % (g, hist)
    return str


def getHistogram(goals):
    edges = [-0.5 + p for p in range(7)]
    edges.append(2000)  # last edge just captures all goals >= 6
    hist = np.histogram(goals, edges)
    histNorm = hist[0]*100.0/sum(hist[0])
    return histNorm



np.random.seed(0)
nGames = 100;
skillValueT1 = 50

skillValuesT2 = [10 + s*2 for s in range(45)]

for skillValue2 in skillValuesT2:
    print "\n%d against %d" %(skillValueT1, skillValue2)
    goals1, goals2 = analyzeAllPlayersEqualButDiffTeams(skillValueT1,skillValue2, 'allplayers50-100.png', nGames)






