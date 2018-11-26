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

def analyzeAllPlayersEqualExplicitSkills(skills1, skills2, doHist, nGames):
    age = 20
    team1 = createAllPlayersEqualTeam(skills1, age, roles433)
    team2 = createAllPlayersEqualTeam(skills2, age, roles433)
    # showTeam(team1)
    # showTeam(team2)
    return analyzeTeam1AgainstTeam2(nGames, team1, team2, doHist)



def analyzeAllPlayersEqualTeams(skillValue, plotOutFile, nGames):
    skills = skillValue * np.ones(5)
    age = 20
    team = createAllPlayersEqualTeam(skills, age, roles433)
    analyzeTeam1AgainstTeam2(nGames, team, team, plotOutFile)


def analyzeTeam1AgainstTeam2(nGames, team1, team2, doHist):
    goals1, goals2 = playNGames(nGames, team1, team2)
    if doHist:
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

def getProbabilityOfWinning(goals1, goals2):
    wins = 1.*sum( goals1 > goals2)
    ties = 1.*sum( goals1 == goals2)
    loss = 1.*sum( goals1 < goals2)
    return wins/goals1.size, ties/goals1.size, loss/goals1.size


def createPlot1(nGames):
    #
    # fixes Team1 to all 50
    # varies Team2 from 20 to 100
    # output: plot of pWin, pTie for team1.
    #
    np.random.seed(0)
    skillValueT1 = 50

    # skillValuesT2 = [10 + s*2 for s in range(46)]
    skillValuesT2 = np.array([20 + s*2 for s in range(40)])
    pWin = np.empty(skillValuesT2.size)
    pTie = np.empty(skillValuesT2.size)
    pLoss = np.empty(skillValuesT2.size)

    for s, skillValue2 in enumerate(skillValuesT2):
        print "\n%d against %d" %(skillValueT1, skillValue2)
        goals1, goals2 = analyzeAllPlayersEqualButDiffTeams(skillValueT1,skillValue2, False, nGames)
        pWin[s], pTie[s], pLoss[s] = getProbabilityOfWinning(goals1, goals2)

    fig, ax = plt.subplots()
    line1, = ax.plot(skillValuesT2, pWin, label='Win prob')
    line2, = ax.plot(skillValuesT2, pTie, label='Tie prob')
    setPlot(ax, 'Skills of all players in Team 2', 'Probability','Team 1: all-50.  Team 2: all-N')
    plt.savefig('team1-all50_team2-allVarying.png')

def setPlot(ax, xlabel, ylabel, title):
    ax.legend()
    ax.set_xlabel(xlabel)
    ax.set_ylabel(ylabel)
    ax.set_title(title)
    ax.grid(True,'both')
    ax.set_ylim(0,1)
    ax.set_yticks([0.1*t for t in range(11)], True)

def createPlot2(nGames):
    #
    # fixes Team1 to all 50
    # varies Team2, skill by skill, from 20 to 100
    # output: plot of pWin, pTie for team1, one curve per skill varied.
    #
    np.random.seed(0)
    nSkills = 6  # the last skill is the goalie
    skills1 = 50 * np.ones(nSkills-1)
    skillValuesT2 = np.array([20 + s * 4 for s in range(20)])

    pWin = np.empty([skillValuesT2.size, nSkills])
    pTie = np.empty([skillValuesT2.size, nSkills])
    pLoss = np.empty([skillValuesT2.size, nSkills])
    skNames.append('Goalie')

    for sk in range(nSkills):
        print "Varying skill: %s" %skNames[sk]
        skills2 = 50 * np.ones(nSkills)
        for s, skillValue2 in enumerate(skillValuesT2):
            skills2[sk] = skillValue2
            goals1, goals2 = analyzeAllPlayersEqualExplicitSkills(skills1,skills2, False, nGames)
            pWin[s, sk], pTie[s, sk], pLoss[s,sk] = getProbabilityOfWinning(goals1, goals2)

    fig, ax = plt.subplots()
    for sk in range(nSkills):
        ax.plot(skillValuesT2, pWin[:,sk], label=skNames[sk])
    setPlot(ax, 'Value of specific skill varied', 'Probability of Win','Team 1: all-50.  Team 2: all-50-except-1-skill')
    plt.savefig('team1-all50_team2-oneSkillVariedAtATime-probWin.png')

    plt.clf()

    fig, ax = plt.subplots()
    for sk in range(nSkills):
        ax.plot(skillValuesT2, pTie[:,sk], label=skNames[sk])
    setPlot(ax, 'Value of specific skill varied', 'Probability of Tie','Team 1: all-50.  Team 2: all-50-except-1-skill')
    plt.savefig('team1-all50_team2-oneSkillVariedAtATime-probTie.png')


createPlot1(10000)  # ideal plot: nGames = 1e4
# createPlot2(10000)  # ideal plot: nGames = 1e4



