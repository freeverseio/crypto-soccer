import random
import numpy as np
import copy
from simulate import *

def playNGames(nGames, team1, team2):
    goals1 = np.zeros(nGames);
    goals2 = np.zeros(nGames);
    for game in range(nGames):
        goals1[game], goals2[game] = playGame(team1, team2)
        # print 'Result: %s - %s' % (goals1, goals2)
    return goals1, goals2
