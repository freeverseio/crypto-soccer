import random
import numpy as np
import copy
import datetime
from os import listdir, makedirs
from os.path import isfile, join, exists

# defence, speed, pass, shoot, endurance
DE = 0
SP = 1
PA = 2
SH = 3
EN = 4

nPlayers = 11

GOALIE = 0
DEFENDER = 1
MIDFIELD = 2
ATTACKER = 3

roles433 = [0,1,1,1,1,2,2,2,3,3,3]
roles442 = [0,1,1,1,1,2,2,2,2,3,3]
roles541 = [0,1,1,1,1,1,2,2,2,2,3]
roles631 = [0,1,1,1,1,1,1,2,2,2,3]
roles640 = [0,1,1,1,1,1,1,2,2,2,2]
roles451 = [0,1,1,1,1,2,2,2,2,2,3]

ROUNDS = 18
MAX_DICE_RAND = 1000-1 # basically, discretization used to determine who wins the dice


class Player:
    age     = None
    skills  = None
    role    = None

class Team:
    players         = None
    endurance       = None
    createShoot     = None
    defendShoot     = None
    blockShoot      = None
    move2attack     = None
    goals           = None

# starts with an int seed (e.g. intSeed = 3) and returns an unpredictable int
# as if a hash had taken place.
def intSeed2RndSeed(intSeed, maxRnd = 1000000000):
    np.random.seed(intSeed)
    return np.random.randint(0,maxRnd,1)

def createRandomPlayer(rndSeed, role):
    np.random.seed(rndSeed)
    newPlayer = Player()
    newPlayer.role = role
    newPlayer.age = 16 + random.random()*19     # states[0] = 16 + (states[0] % 20)
    newPlayer.skills = np.random.randint(0,49,5)          # states[sk] = states[sk] % 50
    excess = int( (250-newPlayer.skills.sum())/5 )        # excess = (250 - excess)/5
    newPlayer.skills += excess
    return newPlayer


def createRandomTeam(intSeed, roles):
    newTeam = Team()
    newTeam.players = []
    rndSeed = intSeed2RndSeed(intSeed)
    for p in range(nPlayers):
        newTeam.players.append(createRandomPlayer(rndSeed+p, roles[p]))
    return newTeam

def showTeam(team):
    for player in team.players:
        print str(player.role) + " - " + str(player.skills)

def getDefenders(team, skill):
    return [p.skills[skill] for p in team.players if p.role==DEFENDER]

def getMids(team, skill):
    return [p.skills[skill] for p in team.players if p.role==MIDFIELD]

def getAttackers(team, skill):
    return [p.skills[skill] for p in team.players if p.role==ATTACKER]

def getGoalie(team, skill):
    return [p.skills[skill] for p in team.players if p.role==GOALIE]

def convertEndurance2Percentage(endurance):
    # endurance is converted to a percentage that will be maintained:
    # 100 is super - endurant(1500), 70 is bad. For an avg starting team (550).
    if (endurance < 500):
        endurance = 70
    elif (endurance < 1400):
        endurance = 100 - (1400-endurance) / 30
    else:
        endurance = 100
    return endurance

def throwDice(w1, w2, rndNum, maxRndNum):
    if (((w1 + w2) * rndNum) < (w1 * (maxRndNum - 1))):
        return 0
    else:
        return 1

def throwDiceArray(weights, rndNum, maxRndNum):
    uniformRndInSumOfWeights = sum(weights) * rndNum
    cumSum = 0
    for w in range(len(weights)):
        cumSum += weights[w]
        if ( uniformRndInSumOfWeights < ( cumSum * (maxRndNum-1) )):
            return w
    return w


def computeTeamGlobalSkills(team):
    endurance = sum([player.skills[EN] for player in team.players])
    team.endurance = convertEndurance2Percentage(endurance)

    #createShoot = speed(attackers) + pass(attackers)
    #defendShoot =    speed(defenders) + defence(defenders)
    # move2attack = defence(defenders + 2 * midfields + attackers) +
    #               speed(defenders + 2 * midfields) + pass(defenders + 3 * midfields)
    team.createShoot = sum( getAttackers(team,SP) ) + sum( getAttackers(team,PA))
    team.defendShoot = sum( getDefenders(team,SP) ) + sum( getDefenders(team,DE))
    team.blockShoot = sum( getGoalie(team,SH))

    move2attack =  sum(getDefenders(team,DE)) + 2*sum(getMids(team,DE)) + sum(getAttackers(team,DE))
    move2attack += sum(getDefenders(team,SP)) + 2*sum(getMids(team,SP))
    move2attack += sum(getDefenders(team,PA)) + 3*sum(getMids(team,PA))
    team.move2attack = move2attack

def manages2Shoot(defendShoot, createShoot, rndNum1, maxRndNum):
    return throwDice(defendShoot, 0.6*createShoot, rndNum1, maxRndNum) == 1

def manages2Score(teamThatAttacks, teamThatDefends, rndNum1, rndNum2, maxRndNum):
    #attacker who actually shoots is selected weighted by his speed
    attackers = [p for p in teamThatAttacks.players if p.role == ATTACKER]
    attackersSpeed = [p.skills[SP] for p in attackers]

    shooterIdx = throwDiceArray(attackersSpeed, rndNum1, maxRndNum)

    # a goal is scored by confronting his shoot skill to the goalkeeper block skill
    goalieBlock = getGoalie(teamThatDefends, SH)

    return throwDice( attackers[shooterIdx].skills[SH] * 0.7, goalieBlock, rndNum2, maxRndNum) == 0


def playGame(team1, team2, intSeed):
    t1 = copy.deepcopy(team1)
    t2 = copy.deepcopy(team2)
    computeTeamGlobalSkills(t1)
    computeTeamGlobalSkills(t2)
    rndSeed = intSeed2RndSeed(intSeed)
    np.random.seed(rndSeed)

    t1.goals = 0
    t2.goals = 0
    for round in range(ROUNDS):
        teamAttacks = throwDice(t1.move2attack, t2.move2attack, np.random.randint(0,MAX_DICE_RAND), MAX_DICE_RAND)
        if teamAttacks == 0:
            teamThatAttacks = t1
            teamThatDefends = t2
        else:
            teamThatAttacks = t2
            teamThatDefends = t1

        if manages2Shoot(teamThatDefends.defendShoot, teamThatAttacks.createShoot, np.random.randint(0,MAX_DICE_RAND), MAX_DICE_RAND):
            if manages2Score(teamThatAttacks, teamThatDefends, np.random.randint(0,MAX_DICE_RAND), np.random.randint(0,MAX_DICE_RAND), MAX_DICE_RAND):
                teamThatAttacks.goals += 1

    print 'Result: %s - %s' % (t1.goals, t2.goals)


barca = createRandomTeam(0,roles433)

showTeam(barca)

for game in range(10):
    playGame(barca, barca, game)