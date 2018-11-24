import random
import numpy as np
import copy
from simulate import *
import matplotlib.pyplot as plt

from testingLib import *

def analyzeAllPlayersEqualButDiffTeams(skillValue1, skillValue2, plotOutFile, nGames):
    skills1 = np.array([skillValue1 for skill in range(5)])
    skills2 = np.array([skillValue2 for skill in range(5)])
    age = 20
    team1 = createAllPlayersEqualTeam(skills1, age, roles433)
    team2 = createAllPlayersEqualTeam(skills2, age, roles433)
    analyzeTeam1AgainstTeam2(nGames, team1, team1, plotOutFile)


def analyzeAllPlayersEqualTeams(skillValue, plotOutFile, nGames):
    skills = np.array([skillValue for skill in range(5)])
    age = 20
    team = createAllPlayersEqualTeam(skills, age, roles433)
    analyzeTeam1AgainstTeam2(nGames, team, team, plotOutFile)

def analyzeTeam1AgainstTeam2(nGames, team1, team2, plotOutFile):
    goals1, goals2 = playNGames(nGames, team1, team2)
    # print "Averages : %s - %s" %(goals1.mean(), goals2.mean())
    print "Average Goals : %.2f,  Distr: %s" %(goals1.mean(), printHistogram(goals1))
    printHistogram(goals1)
    printHistogram(goals2)


def printHistogram(goals):
    str = '|  '
    for (g, hist) in enumerate(getHistogram(goals)):
        str += "%d: %.1f  |  " % (g, hist)
    return str


def getHistogram( goals ):
    hist = np.histogram(goals, [-0.5 + p for p in range(7)])
    histNorm = hist[0]*100.0/sum(hist[0])
    return histNorm




nGames = 1000;
skillValue = 50
analyzeAllPlayersEqualTeams(50, 'allplayes50.png', nGames)
# analyzeAllPlayersEqualTeams(100, 'allplayers100.png', nGames)
analyzeAllPlayersEqualButDiffTeams(50,100, 'allplayers50-100.png', nGames)






