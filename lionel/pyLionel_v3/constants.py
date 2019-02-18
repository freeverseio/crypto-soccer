import random
import numpy as np
import copy
import datetime
from os import listdir, makedirs
from os.path import isfile, join, exists

CHALLENGING_PERIOD_BLKS = 60

AVG_SKILL = 50
NPLAYERS_PER_TEAM = 16
MAX_NTEAMS_PER_LEAGUE = 10

MIN_PLAYER_AGE = 16
MAX_PLAYER_AGE = 35 # max age at time of creation, of course

ADDR1 = 0x5eD8Cee6b63b1c6AFce3AD7c92f4fD7E1B8fAd9F
ADDR2 = 0x01D4950B1Ed0cDAc801973EA8968785148a9E006
ADDR3 = 0x38aa48A49034c7AF5C6b04b3AF39F2BaAFe9fc3a


# defence, speed, pass, shoot, endurance
DE = 0
SP = 1
PA = 2
SH = 3
EN = 4

N_SKILLS    = 5
SK_NAMES    = ["Defense", "Speed", "Pass", "Shoot", "Stamina"]


MAX_RND_SEED_ALLOWED_BY_NUMPY = 2**32 - 1

GOALIE = 0
DEFENDER = 1
MIDFIELD = 2
ATTACKER = 3

# TACTICS BITS
TACTICS = {
    "433": 0,
    "442": 1,
    # ... fill up to 7
    # after 8, same but with pressing
    "433pressing": 8,
    "442pressing": 9,
    # ....
}

# When a team is created, there is no particular order of the players, which are
# generated randomly. We need to define an order, from most defensive to more attacking
# so that afterwards, user actions when changing formation are either not needed, or minimal.
# For example, if a user formation is 433, then, the first 4 are defenders, the next 3 are mids, etc.
# So when we create a league, we provide a maps between:
#    a list of NPLAYER_PER_TEAM  --->  a list of NPLAYER_PER_TEAM
#
# To test, we will either use two. One is the "I'm ok with the original order", which
# is just a map from each number to itself; and the extreme "I want the reverse order", where
# the last created player in the team will act as goalie, then prev-to-last is defender, etc:
DEFAULT_ORDER = np.arange(NPLAYERS_PER_TEAM)
REVERSE_ORDER = np.arange(NPLAYERS_PER_TEAM, 0, -1) - 1


