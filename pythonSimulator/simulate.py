import random
import numpy as np
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



class Player:
    age     = None
    skills  = None
    role    = None

# starts with an int seed (e.g. intSeed = 3) and returns an unpredictable int
# as if a hash had taken place.
def intSeed2RndSeed(intSeed):
    np.random.seed(intSeed)
    return np.random.randint(0,1000000000,1)

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
    newTeam = []
    rndSeed = intSeed2RndSeed(intSeed)
    for p in range(nPlayers):
        newTeam.append(createRandomPlayer(rndSeed+p, roles[p]))
    return newTeam

def showTeam(team):
    for player in team:
        print str(player.role) + " - " + str(player.skills)

def getDefenders(team, skill):
    return [p.skills[skill] for p in team if p.role==DEFENDER]

def getMids(team, skill):
    return [p.skills[skill] for p in team if p.role==MIDFIELD]

def getAttackers(team, skill):
    return [p.skills[skill] for p in team if p.role==ATTACKER]

def getGoalie(team, skill):
    return [p.skills[skill] for p in team if p.role==GOALIE]


def getTeamGlobalSkills(team):
    endurance = sum([player.skills[EN] for player in team])

    #createShoot = speed(attackers) + pass(attackers)
    createShoot = sum( getAttackers(team,SP) ) + sum( getAttackers(team,PA))

    #defendShoot =    speed(defenders) + defence(defenders)
    defendShoot = sum( getDefenders(team,SP) ) + sum( getDefenders(team,DE))

    blockShoot = sum( getGoalie(team,SH))

    # move2attack = defence(defenders + 2 * midfields + attackers) +
    #               speed(defenders + 2 * midfields) + pass(defenders + 3 * midfields)
    move2attack =  sum(getDefenders(team,DE)) + 2*sum(getMids(team,DE)) + sum(getAttackers(team,DE))
    move2attack += sum(getDefenders(team,SP)) + 2*sum(getMids(team,SP))
    move2attack += sum(getDefenders(team,PA)) + 3*sum(getMids(team,PA))
    return endurance, createShoot, defendShoot, blockShoot, move2attack

def playGame(t1, t2, intSeed):
    rndSeed = intSeed2RndSeed(intSeed)
    endurance1, createShoot1, defendShoot1, blockShoot1, move2attack1 = getTeamGlobalSkills(t1)
    endurance2, createShoot2, defendShoot2, blockShoot2, move2attack2 = getTeamGlobalSkills(t2)
    


barca = createRandomTeam(0,roles433)

showTeam(barca)
playGame(barca, barca, 0)