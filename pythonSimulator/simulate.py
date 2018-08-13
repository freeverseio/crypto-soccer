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
    skills = np.random.randint(0,49,5)          # states[sk] = states[sk] % 50
    excess = int( (250-skills.sum())/5 )        # excess = (250 - excess)/5
    skills += excess
    return newPlayer


def createRandomTeam(intSeed, roles):
    newTeam = []
    rndSeed = intSeed2RndSeed(intSeed)
    for p in range(nPlayers):
        newTeam.append(createRandomPlayer(rndSeed+p, roles[p]))
    return newTeam


barca = createRandomTeam(0,roles433)
