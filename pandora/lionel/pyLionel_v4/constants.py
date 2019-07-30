import random
import numpy as np
import copy
import datetime
from os import listdir, makedirs
from os.path import isfile, join, exists


# Blocks after an update/challenge during which it can still be challenged.
SECS_BETWEEN_BLOCKS = 10 # 10 for ETH, 5 for POA
CHALLENGING_PERIOD_MINS = 3
CHALLENGING_PERIOD_BLKS = int(CHALLENGING_PERIOD_MINS * 60 / SECS_BETWEEN_BLOCKS)

# Main league organization
MAX_DIV_PER_COUNTRY = 100000
TEAMS_PER_LEAGUE = 8
LEAGUES_PER_DIVISON = 16
VERSES_PER_DAY = 4*24
VERSES_PER_MONTH = 2920 # exactly 365/12 in verses
VERSES_PER_TIMEZONE = 4
MATCHDAYS_PER_LEAGUE = (TEAMS_PER_LEAGUE-1)*2
DAYS_BETWEEN_LEAGUES = 2
DAYS_PER_ROUND = MATCHDAYS_PER_LEAGUE + DAYS_BETWEEN_LEAGUES
VERSES_PER_ROUND = DAYS_PER_ROUND * VERSES_PER_DAY
BLOCKS_BETWEEN_VERSES = 360
assert VERSES_PER_ROUND == 1536

# TIMEZONE UPDATE CODES
TZ_NULL = 24 # (one more than allowed)
TZ_IDX_MARKET_OPENS = 55
TZ_IDX_DRAW_NEXT_LEAGUES = 56
TZ_IDX_LAST_CYCLE_IDX = 62
# Time evolution
DEPLOYMENT_IN_UNIX_MONTHS = 602
GAME_MONTH_IN_VERSES = 7 * VERSES_PER_DAY / 12 # 7 VERSES_PER_DAY = 1 Year => 1 month = 7/12*VERSES_PER_DAY = 7 * 4 * 24 / 12 = 56
assert GAME_MONTH_IN_VERSES == 56


# BITS TO ENCODE COUNTRY
BITS_PER_COUNTRYIDX = 10
BITS_PER_TEAMIDX    = 28
BITS_PER_PLAYERIDX  = 28

# Verse update status. Only one at a time is possible. The further down we are, the more details have been provided
UPDT_NONE    = 0 # no update at all
UPDT_LEVEL1  = 1 # only verseRoot, but nothing else. Aka LEVEL 1
UPDT_LEVEL2  = 2 # only superRoot, but nothing else. Aka LEVEL 2
UPDT_LEVEL3  = 3 # allLeaguesRoot, but nothing else. Aka LEVEL 3
UPDT_LEVEL4  = 4 # matchdayHashes provided. Aka LEVEL 4

# Players and teams constants
AVG_SKILL = 50
PLAYERS_PER_TEAM_INIT  = 18
PLAYERS_PER_TEAM_MAX   = 25
MIN_PLAYER_AGE = 16
MAX_PLAYER_AGE = 35 # max age at time of creation, of course

# Default addresses
ALICE = 0x5eD8Cee6b63b1c6AFce3AD7c92f4fD7E1B8fAd9F
BOB = 0x01D4950B1Ed0cDAc801973EA8968785148a9E006
CAROL = 0x38aa48A49034c7AF5C6b04b3AF39F2BaAFe9fc3a
FREEVERSE = 0x000000000000000000000000000000000000001


# Indices to refer to player skills
N_SKILLS    = 5
SK_NAMES    = ["Defense", "Speed", "Pass", "Shoot", "Stamina"]
DE = 0
SP = 1
PA = 2
SH = 3
EN = 4
GOALIE      = 0
DEFENDER    = 1
MIDFIELD    = 2
ATTACKER    = 3


# the largest number that can be used for playerIdx. To be decided. If we went for 26 bits:
UINTMINUS1 = 2**26-1
MAX_RND_SEED_ALLOWED_BY_NUMPY = 2**32 - 1

# For the time being, selecting "tactics" is choosing a number between 0...MAX, where 0 = 442, etc.
TACTICS = {
    "433": 0,
    "442": 1,
    # ... fill up to 7
    # after 8, same but with pressing
    "433pressing": 8,
    "442pressing": 9,
    # ....
}

# For the time being, selecting a "team order" is choosing an order for the PLAYERS_PER_TEAM_MAX
# So if each player has a shirt number (0,...PLAYERS_PER_TEAM_MAX-1), then we need to order them, for example:
#   order = [2, 4, 1, 6, 7, 15, ...]
# The first position plays as goalkeeper.
# If the chosen tactics if 433, then the next players [4, 1, 6, 7] would be the defenders, etc.
# This will change when game design is completed.

# a few orders we will play with:
DEFAULT_ORDER = np.arange(PLAYERS_PER_TEAM_MAX)
REVERSE_ORDER = np.arange(PLAYERS_PER_TEAM_MAX, 0, -1) - 1

ORDER1 = np.array([0,5,4,3,2,1])
ORDER1 = np.append(ORDER1, range(6,PLAYERS_PER_TEAM_MAX))
ORDER2 = np.array([3,2,1,4,5,0])
ORDER2 = np.append(ORDER2, range(6,PLAYERS_PER_TEAM_MAX))

POSSIBLE_ORDERS  = [DEFAULT_ORDER, REVERSE_ORDER, ORDER1, ORDER2]
POSSIBLE_TACTICS = [TACTICS["433"], TACTICS["442"], TACTICS["433pressing"], TACTICS["442pressing"]]

# Max number of superroots that fit in a single verse
SUPERROOTS_PER_VERSE = 200

LEAGUE_INIT_SKILLS_ID = -1