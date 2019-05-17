import random
import numpy as np
import copy
from simulate import *
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D  # noqa: F401 unused import
from testingLib import *
from matplotlib import rc
# import pandas as pd


def analyzeAllPlayersEqualButDiffTeams(skillValue1, skillValue2, plotOutFile, nGames, roles1, roles2):
    skills1 = skillValue1 * np.ones(5)
    skills2 = skillValue2 * np.ones(5)
    age = 20
    team1 = createAllPlayersEqualTeam(skills1, age, roles1)
    team2 = createAllPlayersEqualTeam(skills2, age, roles2)
    # showTeam(team1)
    # showTeam(team2)
    return analyzeTeam1AgainstTeam2(nGames, team1, team2, plotOutFile)

def analyzeAllPlayersEqualExplicitSkills(skills1, skills2, doHist, nGames, roles1, roles2):
    age = 20
    team1 = createAllPlayersEqualTeam(skills1, age, roles1)
    team2 = createAllPlayersEqualTeam(skills2, age, roles2)
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


def create3dplot(histo, nHistVals, skillValuesT2, roles1, roles2):
    dz = []
    cols = ['red', 'blue', 'yellow', 'green', 'black', 'orange', 'pink', 'magenta']
    colours = []
    xx, yy = np.meshgrid(range(nHistVals),skillValuesT2)
    x, y = xx.ravel(), yy.ravel()

    for g, goals in enumerate(histo):
        for sk, frequency in enumerate(goals):
            dz.append(histo[g,sk])
            colours.append(cols[sk])

    z = np.zeros(len(y))
    dx = 3*np.ones(len(y))
    dy = 0.7*np.ones(len(y))

    ax3d = plt.figure().gca(projection='3d')
    ax3d.bar3d(x, y, z, dy, dx, dz, color=colours)
    ax3d.set_xlabel('Num of Goals Team 1')
    ax3d.set_ylabel('Skill Value Team 2')
    ax3d.set_zlabel('Percentage of games')
    plt.savefig('team1-%s-all50_team2-%s-allVarying_hist.png' % (roles1[0], roles2[0]))

def createHistPlot(histo, nHistVals, skillValuesT2, roles1, roles2):

    barWidth = 1.0
    x_names = skillValuesT2

    nSkillVals = len(skillValuesT2)
    x_range = range(nSkillVals)
    accumulated = np.zeros(nSkillVals)

    cols = ['red', 'blue', 'yellow', 'green', 'black', 'orange', 'pink', 'magenta']
    plots = []

    for histVal in range(nHistVals):
        bars = copy.deepcopy(histo[:,histVal])
        plots.append(plt.bar(x_range, bars, bottom = accumulated, color=cols[histVal], edgecolor='white', width=barWidth))
        accumulated += bars

    plt.xticks(x_range, x_names)
    plt.yticks([10*i for i in range(11)])
    plt.grid(False, alpha=0.4)
    return plots

def createPlot1(nGames, roles1, roles2):
    #
    # fixes Team1 to all 50
    # varies Team2 from 20 to 100
    # output: plot of pWin, pTie for team1.
    #
    np.random.seed(0)
    skillValueT1 = 50

    # skillValuesT2 = [10 + s*2 for s in range(46)]
    skillValuesT2 = np.array([20 + s*4 for s in range(20)])
    pWin = np.zeros(skillValuesT2.size)
    pTie = np.zeros(skillValuesT2.size)
    pLoss = np.zeros(skillValuesT2.size)

    roles1array = roles1[1]
    roles2array = roles2[1]

    nHistVals = 7
    histo1 = np.zeros([skillValuesT2.size, nHistVals])
    histo2 = np.zeros([skillValuesT2.size, nHistVals])
    for s, skillValue2 in enumerate(skillValuesT2):
        goals1, goals2 = analyzeAllPlayersEqualButDiffTeams(skillValueT1,skillValue2, False, nGames, roles1array, roles2array)
        pWin[s], pTie[s], pLoss[s] = getProbabilityOfWinning(goals1, goals2)
        histo1[s,:] = getHistogram(goals1)
        histo2[s,:] = getHistogram(goals2)

    fig, ax = plt.subplots()
    line1, = ax.plot(skillValuesT2, pWin, label='Win prob')
    line2, = ax.plot(skillValuesT2, pTie, label='Tie prob')
    setPlot(ax,
            'Skills of all players in Team 2',
            'Probability','Team 1: %s, all-50.  Team 2: %s, all-N' % (roles1[0], roles2[0])
            )
    plt.savefig('team1-%s-all50_team2-%s-allVarying.png' % (roles1[0], roles2[0]))
    plt.close(fig)

    # create3dplot(histo, nHistVals, skillValuesT2, roles1, roles2)

    fig, axes = plt.subplots(nrows=2, ncols=1)
    # Set the ticks and ticklabels for all axes
    plt.setp(axes, xticks=skillValuesT2, xticklabels=skillValuesT2, yticks=[10*i for i in range(11)])


    fig.suptitle("Team1: %s-all50 vs Team2: %s-allVarying.png" % (roles1[0], roles2[0]))

    plt.subplot(211)
    createHistPlot(histo1, nHistVals, skillValuesT2, roles1, roles2)
    plt.title('Top: Team 1,    Bottom: Team 2')
    plt.subplot(212)
    plots = createHistPlot(histo2, nHistVals, skillValuesT2, roles1, roles2)

    plt.xlabel("Skills team 2")
    fig.text(0.04, 0.5, 'Percentage of games with these numGoals', va='center', rotation='vertical')
    fig.legend([p[0] for p in plots], range(nHistVals))
    plt.savefig('team1-%s-all50_team2-%s-allVarying_hist.png' % (roles1[0], roles2[0]))

def setPlot(ax, xlabel, ylabel, title):
    ax.legend()
    ax.set_xlabel(xlabel)
    ax.set_ylabel(ylabel)
    ax.set_title(title)
    ax.grid(True,'both')
    ax.set_ylim(0,1)
    ax.set_yticks([0.1*t for t in range(11)], True)

def createPlot2(nGames, roles1, roles2):
    #
    # fixes Team1 to all 50
    # varies Team2, skill by skill, from 20 to 100
    # output: plot of pWin, pTie for team1, one curve per skill varied.
    #
    np.random.seed(0)
    nSkills = 6  # the last skill is the goalie
    skills1 = 50 * np.ones(nSkills-1)
    skillValuesT2 = np.array([20 + s * 4 for s in range(20)])

    pWin = np.zeros([skillValuesT2.size, nSkills])
    pTie = np.zeros([skillValuesT2.size, nSkills])
    pLoss = np.zeros([skillValuesT2.size, nSkills])
    skNames.append('Goalie')

    roles1array = roles1[1]
    roles2array = roles2[1]

    for sk in range(nSkills):
        print "Varying skill: %s" %skNames[sk]
        skills2 = 50 * np.ones(nSkills)
        for s, skillValue2 in enumerate(skillValuesT2):
            skills2[sk] = skillValue2
            goals1, goals2 = analyzeAllPlayersEqualExplicitSkills(skills1,skills2, False, nGames, roles1array, roles2array)
            pWin[s, sk], pTie[s, sk], pLoss[s,sk] = getProbabilityOfWinning(goals1, goals2)

    fig, ax = plt.subplots()
    for sk in range(nSkills):
        ax.plot(skillValuesT2, pWin[:,sk], label=skNames[sk])
    setPlot(ax,
            'Value of specific skill varied',
            'Probability of Win','Team 1: %s, all-50.  Team 2: %s, all-50-except-1-skill' %(roles1[0], roles2[0])
            )
    plt.savefig('team1-%s-all50_team2-%s-oneSkillVariedAtATime-probWin.png' %( roles1[0], roles2[0]))

    plt.close(fig)

    fig, ax = plt.subplots()
    for sk in range(nSkills):
        ax.plot(skillValuesT2, pTie[:,sk], label=skNames[sk])
    setPlot(ax,
            'Value of specific skill varied',
            'Probability of Tie','Team 1: %s, all-50.  Team 2: %s, all-50-except-1-skill' %(roles1[0], roles2[0])
            )
    plt.savefig('team1-%s-all50_team2-%s-oneSkillVariedAtATime-probTie.png' %( roles1[0], roles2[0]))
    plt.close(fig)



nGames = 2000
nRoles = len(roles)

for r1 in range(nRoles):
    for r2 in range(r1,nRoles):
        print "Roles: %s vs %s" % (roles[r1][0], roles[r2][0])
        createPlot1(nGames, roles[r1], roles[r2])  # ideal plot: nGames = 1e4

for r1 in range(nRoles):
    for r2 in range(r1,nRoles):
        print "Roles: %s vs %s" % (roles[r1][0], roles[r2][0])
        createPlot2(nGames, roles[r1], roles[r2])  # ideal plot: nGames = 1e4





