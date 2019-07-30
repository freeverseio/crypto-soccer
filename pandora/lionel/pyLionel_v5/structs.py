import numpy as np
# we need 'copy' because Python passes ALWAYS by reference, and sometimes, we want to avoid unwated modification inside a function
import copy
import datetime
# These are the main constants that govern the whole thing.
from constants import *
# We use the serialization function from the Pickle library. In Solidity, we'll serialize our way.
from pickle import dumps as serialize
# Currently, we use this library for building Merkle Trees and Proofs. Any other would work:
from merkle_tree import *
# Pylio is a set of useful handmade functions:
import pylio


# ------------------------------------------------------------------------
# ----------      Classes common to both ST and CLIENT        ------------
# ------------------------------------------------------------------------


class Country():
    def __init__(self, timeZone):
        self.nDivisions = 1
        self.timeZone = timeZone
        self.nFrozenDivisions = 0
        self.divisonIdxToRound = {1: 1} # divId = 1 starts at the very first round = 1
        self.teamIdxInCountryToTeam = {}


class TimeZoneUpdate():
    def __init__(self):
        self.teamOrgMap = [0, 0]
        self.newestOrgMap = 0
        self.teamsSkills = 0
        self.scores = 0
        self.updateCycleIdx = 0
        self.lastBlockUpdate = 0

    def isTimeZoneMarketOpen(self, nowBlock):
        if self.updateCycleIdx > TZ_IDX_MARKET_OPENS:
            return True
        if self.isLastUpdateSettled(nowBlock) and (
            self.updateCycleIdx == TZ_IDX_MARKET_OPENS or
            self.updateCycleIdx == 0
        ):
            return True
        return False

    def isLastUpdateSettled(self, nowBlock):
        if self.lastBlockUpdate == 0:
            return True
        return nowBlock > self.lastBlockUpdate + CHALLENGING_PERIOD_BLKS

    # Todo: implement do something with updateData
    def newUpdate(self, nowBlock, updateData):
        assert self.isLastUpdateSettled(nowBlock), "cannot update until settled!"
        self.updateCycleIdx = (self.updateCycleIdx + 1) % TZ_IDX_LAST_CYCLE_IDX
        isInFreezePeriod = self.updateCycleIdx > TZ_IDX_DRAW_NEXT_LEAGUES
        if not isInFreezePeriod:
            # do something with update data
            self.lastBlockUpdate = nowBlock


# In Solidity, PlayerState will be just a uin256, serializing the data shown here,
# ...and there'll be associated read/write functions
# Note: playerIdx = 0 is the null player

# PlayerSkills contains the data regarding the player's DNA (not about the team it belongs, etc)
class PlayerSkills():
    def __init__(self, skills = None, monthOfBirthInUnixTime = None, playerIdx = None):
        self.skills                  = skills
        self.monthOfBirthInUnixTime  = monthOfBirthInUnixTime
        self.playerIdx               = playerIdx

    def setSkills(self, skills):
        self.skills = skills

    def getSkills(self):
        return self.skills

    def setMonth(self, month):
        self.monthOfBirthInUnixTime = month

    def getMonth(self):
        return self.monthOfBirthInUnixTime

    def setPlayerIdx(self, playerIdx):
        self.playerIdx = playerIdx

    def getPlayerIdx(self):
        return self.playerIdx


class PlayerState():
    def __init__(self):
        self.playerIdx              = 0
        self.currentTeamIdx         = 0
        self.currentShirtNum        = 0
        self.prevPlayedTeamIdx      = 0
        self.lastSaleBlocknum       = 0

    def getCurrentTeamIdx(self):
        return self.currentTeamIdx

    def setCurrentTeamIdx(self, currentTeamIdx):
        self.currentTeamIdx = currentTeamIdx

    def getPrevPlayedTeamIdx(self):
        return self.prevPlayedTeamIdx

    def setPrevPlayedTeamIdx(self, prevPlayedTeamIdx):
        self.prevPlayedTeamIdx = prevPlayedTeamIdx

    def setCurrentShirtNum(self, currentShirtNum):
        self.currentShirtNum = currentShirtNum

    def getCurrentShirtNum(self):
        return self.currentShirtNum

    def setLastSaleBlocknum(self, blocknum):
        self.lastSaleBlocknum = blocknum

    def getLastSaleBlocknum(self):
        return self.lastSaleBlocknum


# Note: teamIdx = 0 is the null team
# The Team struct contains the array playerIdxs.
# - if playerIdx = 0, it is considered a virtual player
# - if playerIdx = UINTMINUS1, this place in the team is free
class Team():
    def __init__(self, addr):
        self.teamOwner = addr
        self.playerIdxs             = np.append(
            np.zeros(PLAYERS_PER_TEAM_INIT, int),
            UINTMINUS1*np.ones(PLAYERS_PER_TEAM_MAX-PLAYERS_PER_TEAM_INIT, int)
        )
        assert len(self.playerIdxs) == PLAYERS_PER_TEAM_MAX


class League():
    def __init__(self, verseInit, verseStep, nTeams):
        self.nTeams             = nTeams
        self.verseInit          = verseInit
        self.verseStep          = verseStep
        self.usersInitDataHash  = 0

    # simulates what would happen when users sign up, one by one
    def signTeamInLeague(self, teamIdx, teamOrders, tactics):
        self.usersInitDataHash = pylio.serialHash([self.usersInitDataHash, teamIdx, teamOrders, tactics])

    def isGenesisLeague(self):
        return self.verseInit == 0

    def verseFinal(self):
        nMatchdays = 2 * (self.nTeams - 1)
        return self.verseInit + (nMatchdays-1)*self.verseStep


# The main class that rules the update/challenge process
# Level 1: VerseRoot provided (one single hash)
# Level 2: SuperRoots provided (up to 200 Hashes, indexed by 'subverse')
# Level 3: LeagueRoots provided (up to 200 LeagueRoots that challenge a particular subverse; indexed by 'posInSubverse')
# Level 4: OneLeague provided (one hash per mathchday that challenge a particular LeagueRoot)
class VerseUpdate():
    # Constructed directly at Level 1 by an updater that provides the verseRoot.
    def __init__(self, verseRoot, addr, blocknum):
        self.verseRoot                      = pylio.duplicate(verseRoot)
        self.verseRootAddr                  = pylio.duplicate(addr)
        self.lastWriteBlocknum              = pylio.duplicate(blocknum)

        # Levels 2, 3, 4 start at zero.
        self.initLevel2()
        self.initLevel3()
        self.initLevel4()

    def initLevel2(self):
        self.superRoots                     = None
        self.superRootsOwner                = None
        self.superRootsVerseRoot            = None

    def initLevel3(self):
        self.subVerse                       = None
        self.leagueRoots                    = None
        self.leagueRootsOwner               = None

    def initLevel4(self):
        self.posInSubVerse                  = None
        self.dataToChallengeLeague          = None
        self.oneLeagueDataOwner             = None

    # Challenge to Level 1, moves to Level 2
    def writeLevel2(self, superRoots, superRootsVerseRoot, ownerAddr, blocknum):
        self.superRoots                     = superRoots
        self.superRootsOwner                = ownerAddr
        self.lastWriteBlocknum              = blocknum
        self.superRootsVerseRoot            = superRootsVerseRoot

    # Challenge to Level 2, moves to Level 3
    def writeLevel3(self, subVerse, leagueRoots, ownerAddr, blocknum):
        self.subVerse           = subVerse
        self.leagueRoots        = leagueRoots
        self.leagueRootsOwner   = ownerAddr
        self.lastWriteBlocknum  = blocknum

    # Challenge to Level 3, moves to Level 4
    def writeLevel4(self, posInSubVerse, dataToChallengeLeague, addr, blocknum):
        self.posInSubVerse              = posInSubVerse
        self.dataToChallengeLeague      = dataToChallengeLeague
        self.oneLeagueDataOwner         = addr
        self.lastWriteBlocknum          = blocknum

    # slashing basically resets data below the given update, and resets timer.
    def slashLevel2(self, blocknum):
        self.lastWriteBlocknum = blocknum
        self.initLevel2()
        self.initLevel3()
        self.initLevel4()

    def slashLevel3(self, blocknum):
        self.lastWriteBlocknum          = blocknum
        self.initLevel3()
        self.initLevel4()

    def slashLevel4(self, blocknum):
        self.lastWriteBlocknum          = blocknum
        self.initLevel4()


# ------------------------------------------------------------------------
# ----------      Classes used only by CLIENT                 ------------
# ------------------------------------------------------------------------

# LeagueClient extends League to store pre-hash stuff, etc.
class LeagueClient(League):
    def __init__(self, verseInit, verseStep, usersInitData):
        nTeams = len(usersInitData["teamIdxs"])
        League.__init__(self, verseInit, verseStep, nTeams)
        self.usersInitData      = usersInitData
        self.initPlayerStates   = None
        self.lastDayTree        = None
        self.actionsPerMatchday = []
        self.dataToChallengeInitSkills = None
        self.dataAtMatchdays    = 0
        self.dataToChallengeLeague = DataToChallengeLeague(0,0,0)

    def storeDataAtMatchdays(self, dataAtMatchdays):
        self.dataAtMatchdays    = pylio.duplicate(dataAtMatchdays)

    def writeInitState(self, initPlayerStates):
        self.initPlayerStates = pylio.duplicate(initPlayerStates)

    def writeDataToChallengeInitSkills(self, dataToChallengeInitSkills):
        self.dataToChallengeInitSkills = pylio.duplicate(dataToChallengeInitSkills)

    def writeDataToChallengeLeague(self, dataToChallengeLeague):
        self.dataToChallengeLeague = pylio.duplicate(dataToChallengeLeague)

    def getInitPlayerSkills(self):
        initSkills = []
        for team in self.initPlayerStates:
            initSkills.append([pylio.duplicate(PlayerSkills(state)) for state in team])
        return initSkills


# The VerseActionsCommit basically stores the merkle roots of all actions corresponding to a league starting at that moment
# The Merkle Roots are computed from the leafs:
#   - leafs = [ [leagueIdx0, allActionsInLeagueIdx0], ..., ]
#   - where allActionsInLeagueIdx0 = [action0, action1, action2,...]

class VerseActionsCommit:
    def __init__(self, actionsMerkleRoots = 0, blockNum = 0):
        self.actionsMerkleRoots = actionsMerkleRoots
        self.blockNum = blockNum


# Data that the CLIENT builds to challenge a league at level 4.
class DataToChallengeLeague():
    def __init__(self, initSkillsHash, dataAtMatchdayHashes, scores):
        self.initSkillsHash         = initSkillsHash
        self.dataAtMatchdayHashes   = dataAtMatchdayHashes
        self.scores                 = scores


# The Accumulator is responsible for receving user actions and committing them in the correct verse.
# An action is a struct:
#    action00 = {"teamIdx": 2, "teamOrder": [0,4,2,3...,PLAYERS_PER_TEAM_INIT], "tactics": 3}
#       where "tactics" is an int < 16. For example, tactics = 1 means (4,4,2).
#
# The buffer is an array that maintains:  buffer[leagueIdx] = [action0, action1, ...]
class ActionsAccumulator():
    def __init__(self):
        self.buffer                     = {}
        self.commitedActions            = [0] # The genesis commit is a dummy one, as always
        self.commitedTrees              = [0]

    def accumulateAction(self, action, leagueIdx):
        if leagueIdx in self.buffer:
            self.buffer[leagueIdx].append(action)
        else:
            self.buffer[leagueIdx] = [action]

    def clearBuffer(self, actions2remove):
        for action in actions2remove:
            leagueIdx = action[0]
            if leagueIdx in self.buffer:
                del self.buffer[leagueIdx]
            else:
                assert action[1] == 0, "Tried to remove from buffer the actions in a league that was not present"


# Simple struct that stores the data that is computed/updated every matchday
class DataAtMatchday():
    def __init__(self, skillsAtMatchday, tacticsAtMatchday, teamOrdersAtMatchday):
        self.skillsAtMatchday       = pylio.duplicate(skillsAtMatchday)
        self.tacticsAtMatchday      = pylio.duplicate(tacticsAtMatchday)
        self.teamOrdersAtMatchday   = pylio.duplicate(teamOrdersAtMatchday)


# Simple struct that stores the data needed to proof that a certain leaf belongs to a Merkle tree
# "Values" is just the pair [ leafIdx, leafValue ]
class MerkleProof():
    def __init__(self, neededHashes, depth, leaf, leafIdx):
        self.neededHashes   = pylio.duplicate(neededHashes)
        self.depth          = pylio.duplicate(depth)
        self.leaf           = pylio.duplicate(leaf)
        self.leafIdx        = pylio.duplicate(leafIdx)


# Simple struct that stores a tree
class MerkleTree():
    def __init__(self, leafs):
        self.tree    = None
        self.root    = None
        self.depth   = None
        self.makeTree(pylio.duplicate(leafs), pylio.serialHash)

    def makeTree(self, leafs, hashFunction):
        tree, depth = make_tree(leafs, hashFunction)
        self.tree  = tree
        self.root  = root(tree)
        self.depth = depth

    # Merkle proof: given a tree, and its leafs,
    # it creates the hashes required to prove that a given idx in the leave belongs to the tree.
    # it returns the neededHashes
    def prepareProofForLeaf(self, leaf, leafIdx):
        neededHashes = proof(self.tree, [leafIdx])
        return MerkleProof(neededHashes, self.depth, leaf, leafIdx)


# Simple block counter simulator, where the blockhash is just the hash of the blocknumber. Not serious, who cares here.
class Counter():
    def __init__(self, nowInSecsOfADay):
        self.currentBlock = 0
        self.currentVerse = 0
        self.deployTimeInSecsOfADay = nowInSecsOfADay

    def isCrossingVerse(self):
        if self.currentBlock == 0:
            return False
        return self.currentBlock % self.blocksBetweenVerses == 0

    def incrementBlock(self):
        if self.isCrossingVerse():
            self.currentVerse += 1
        self.currentBlock += 1

    def getDeployHoursMinutes(self):
        hours = self.deployTimeInSecsOfADay // 3600
        minutes = (self.deployTimeInSecsOfADay - hours*3600) // 60
        return hours, minutes

class DataToChallengePlayerSkills():
    def __init__(self, merkleProofStates, merkleProofLeague, merkleProofLeagueRoots, merkleProofSuperRoots):
        self.merkleProofStates      = merkleProofStates
        self.merkleProofLeague      = merkleProofLeague
        self.merkleProofLeagueRoots = merkleProofLeagueRoots
        self.merkleProofSuperRoots  = merkleProofSuperRoots




    # ------------------------------------------------------------------------
    # ------     THE MAIN CLASS. PART THAT IS COMMON TO BC AND CLIENT---------
    # ------------------------------------------------------------------------

class Storage(Counter):
    def __init__(self, nowInSecsOfADay, isClient):

        # The Blockchain does not need this fake counter :-)
        Counter.__init__(self, nowInSecsOfADay)

        self.countries = []
        self.timeZoneForRound1, self.verseForRound1 = self.initFirstRound()

        self.countries.append(Country(0)) # countries[0] is dummy
        self.timeZoneToCountries = {}
        self.timeZoneUpdates = {}
        self.createCountry(self.timeZoneForRound1)

        # this bool is just to understand if the created BC is actually a client
        # it allows us, in this simulation, to ensure that the functions that are
        # only to be used by the CLIENT are actually used only by the CLIENT :-)
        self.isClient = isClient

        # a map from playerIdx to playerState, only available for players already sold once,
        # or for 'promo players' not created directly from team creation.
        # In Python, maps are closer to 'dictionaries'
        self.playerIdxToPlayerState = {}

        # the obvious ownership map:
        self.teamNameHashToOwnerAddr = {}

        # an array of leagues, first entry is dummy
        self.leagues = [League(0,0,0)]

        self.blocksBetweenVerses = BLOCKS_BETWEEN_VERSES
        self.VerseActionsCommits = [VerseActionsCommit()]
        self.verseToLeagueCommits = {}

        self.verseToFinishingLeagueIdxs = {}

        if isClient:
            self.forceVerseRootLie = False

    def assertIsClient(self):
        assert self.isClient, "This code should only be run by CLIENTS, not the BC"



    # ------------------------------------------------------------------------
    # ----------      Functions common to both BC and CLIENT      ------------
    # ------------------------------------------------------------------------
    def initFirstRound(self):
        hours, minutes = self.getDeployHoursMinutes()
        quarter = minutes // 15 # = 0, 1, 2, 3
        if quarter == 3:
            return hours + 1, 4
        else:
            return hours, 3 - quarter

    def currentRound(self):
        # verse starts at 0, rounds at 1.
        return self.verseToRound(self.currentVerse)

    def addCountryToTimeZone(self, countryIdx, timeZone):
        if timeZone in self.timeZoneToCountries:
            self.timeZoneToCountries[timeZone].append(countryIdx)
        else:
            self.timeZoneToCountries[timeZone] = [countryIdx]
            self.timeZoneUpdates[timeZone] = TimeZoneUpdate()

    def createCountry(self, timeZone):
        countryIdx = len(self.countries)
        self.countries.append(Country(timeZone))
        self.addCountryToTimeZone(countryIdx, timeZone)
        return countryIdx

    def addDivision(self, countryIdx):
        self.countries[countryIdx].nDivisions += 1
        divisionIdx = self.countries[countryIdx].nDivisions
        self.countries[countryIdx].divisonIdxToRound[divisionIdx] = self.currentRound() + 1
        return divisionIdx

    def verseToUnixMonths(self, verse):
        return DEPLOYMENT_IN_UNIX_MONTHS + int(verse/VERSES_PER_MONTH)

    def getDivisionCreationDay(self, countryIdx, divisionIdx):
        # disregards the offset introduced by timeZone, and thanks to this, avoids requiring country.timeZone
        creationRound = self.countries[countryIdx].divisonIdxToRound[divisionIdx]
        return (creationRound - 1)* DAYS_PER_ROUND

    def getNDivisionsInCountry(self, countryIdx):
        return self.countries[countryIdx].nDivisions

    def getNLeaguesInCountry(self, countryIdx):
        return 1 + (self.getNDivisionsInCountry(countryIdx) -1) * LEAGUES_PER_DIVISON

    def getNTeamsInCountry(self, countryIdx):
        return self.getNLeaguesInCountry(countryIdx) * TEAMS_PER_LEAGUE

    def getTeamIdxInCountryFromLeagueAndPos(self, divisionIdx, leaguePosInDiv, teamPosInLeague):
        # posInDiv and posInLeague start at zero.
        assert divisionIdx > 0, "divs start at idx = 1"
        if divisionIdx == 1:
            assert leaguePosInDiv == 0, "First division only has 1 league"
            nLeaguesAbove = 0
        else:
            nLeaguesAbove = 1 + (divisionIdx - 2) * LEAGUES_PER_DIVISON + leaguePosInDiv
        return 1 + nLeaguesAbove * TEAMS_PER_LEAGUE + teamPosInLeague

    def getDisivionIdxFromTeamIdxInCountry(self, teamIdxInCountry):
        if teamIdxInCountry <= TEAMS_PER_LEAGUE:
            return 1
        else:
            return 2 + int( (teamIdxInCountry - 1 - TEAMS_PER_LEAGUE) / (LEAGUES_PER_DIVISON * TEAMS_PER_LEAGUE) )


    def getTeamIdxInCountryFromPlayerIdxInCountry(self, playerIdxInCountry):
        # posInDiv and posInLeague start at zero.
        return 1 + int((playerIdxInCountry-1) / PLAYERS_PER_TEAM_INIT)

    def getTeamIdxInCountryAndShirtNumFromPlayerIdxInCountry(self, playerIdxInCountry):
        teamIdxInCountry = self.getTeamIdxInCountryFromPlayerIdxInCountry(playerIdxInCountry)
        shirtNum = (playerIdxInCountry - 1) - (teamIdxInCountry - 1) * PLAYERS_PER_TEAM_INIT
        return (teamIdxInCountry, shirtNum)

    def encode(self, val1, val2, bits1, bits2):
        assert val1 < 2**bits1 - 1, "val too big"
        assert val2 < 2**bits2 - 1, "val too big"
        return val1 * 2**bits2 + val2

    def decode(self, val, bits1, bits2):
        assert val < 2**(bits1+bits2) - 1, "val too big"
        val2 = val % 2**bits2
        val1 = int( (val - val2)/2**bits2 )
        return (val1, val2)

    def encodeCountryAndVal(self, val1, val2):
        return self.encode(val1, val2, BITS_PER_COUNTRYIDX, BITS_PER_TEAMIDX)

    def decodeCountryAndVal(self, val):
        return self.decode(val, BITS_PER_COUNTRYIDX, BITS_PER_TEAMIDX)

    def countryExists(self, countryIdx):
        return countryIdx <= len(self.countries) - 1

    def teamExists(self, teamIdx):
        (countryIdx, teamIdxInCountry) = self.decodeCountryAndVal(teamIdx)
        if not self.countryExists(countryIdx):
            return False
        return teamIdxInCountry <= self.getNTeamsInCountry(countryIdx)

    def playerExists(self, playerIdx):
        (countryIdx, playerIdxInCountry) = self.decodeCountryAndVal(playerIdx)
        if not self.countryExists(countryIdx):
            return False
        return playerIdxInCountry <= self.getNTeamsInCountry(countryIdx) * PLAYERS_PER_TEAM_INIT

    def isBotTeam(self, teamIdx):
        (countryIdx, teamIdxInCountry) = self.decodeCountryAndVal(teamIdx)
        return (teamIdxInCountry not in self.countries[countryIdx].teamIdxInCountryToTeam)

    def acquireBot(self, teamIdx, addr):
        assert self.isBotTeam(teamIdx), "cannot acquire a team that is not a Bot"
        (countryIdx, teamIdxInCountry) = self.decodeCountryAndVal(teamIdx)
        self.countries[countryIdx].teamIdxInCountryToTeam[teamIdxInCountry] = Team(addr)

    def getCountryTimeZone(self, countryIdx):
        return self.countries[countryIdx].timeZone

    def getVerseLeaguesStartFromTimeZoneAndRound(self, timeZone, round):
        assert round > 0, "league has never started"
        return self.verseForRound1 + 4 * (timeZone - self.timeZoneForRound1) + (round-1) * VERSES_PER_ROUND

    def verseToRound(self, verse):
        if verse < self.verseForRound1:
            return 0
        else:
            return 1 + (verse - self.verseForRound1) // VERSES_PER_ROUND

    def isPlayerTransferable(self, playerIdx):
        (countryIdx, playerIdxInCountry) = self.decodeCountryAndVal(playerIdx)
        return self.isCountryMarketOpen(countryIdx)

    def isCountryMarketOpen(self, countryIdx):
        timeZone = self.getCountryTimeZone(countryIdx)
        return self.timeZoneUpdates[timeZone].isTimeZoneMarketOpen(self.currentBlock)

    # DeltaVerse < VERSES_PER_DAY:
    # timeZoneToUpdate = timeZoneForRound1 + DeltaVerse // 4 % 24
    # posInUpdate = DeltaVerse % 4
    #
    # if DeltaVerse == VERSES_PER_DAY:
    #     timeZoneToUpdate = None
    #
    # if DeltaVerse > VERSES_PER_DAY and
    #     timeZoneToUpdate = timeZoneForRound1 + (DeltaVerse - 1) // 4 % 24
    # posInUpdate = (DeltaVerse - 1) % 4
    #
    # Note that VERSES_PER_DAY // 4 % 4 = 0, so:
    # DeltaVerse = VPD + 1 = > timeZoneToUpdate = timeZoneForRound1(o'clock)
    # DeltaVerse = VPD + 2 = > timeZoneToUpdate = timeZoneForRound1(past 15)
    # DeltaVerse = VPD + 3 = > timeZoneToUpdate = timeZoneForRound1(past 30)
    # DeltaVerse = VPD + 4 = > timeZoneToUpdate = timeZoneForRound1(past 45)
    # DeltaVerse = VPD + 5 = > timeZoneToUpdate = timeZoneForRound1 + 1(o'clock)

    def verseToTimeZoneToUpdate(self, verse):
        if verse < self.verseForRound1:
            return TZ_NULL, TZ_NULL

        deltaVerse = ( verse - self.verseForRound1 ) % VERSES_PER_ROUND
        if deltaVerse < VERSES_PER_DAY:
            timeZone    = (self.timeZoneForRound1 + deltaVerse//4) % 24
            posInZone   = deltaVerse % 4
        elif deltaVerse == VERSES_PER_DAY:
            timeZone    = TZ_NULL
            posInZone   = TZ_NULL
        else:
            timeZone    = (self.timeZoneForRound1 + (deltaVerse - 1)//4) % 24
            posInZone   = (deltaVerse - 1) % 4
        return timeZone, posInZone

    def currentTimeZoneToUpdate(self):
        return self.verseToTimeZoneToUpdate(self.currentVerse)


    # toni

    def lastVerseBlock(self):
        return self.VerseActionsCommits[-1].blockNum

    def nextVerseBlock(self):
        return self.lastVerseBlock() + self.blocksBetweenVerses

    def commit(self, actionsRoot):
        self.VerseActionsCommits.append(VerseActionsCommit(actionsRoot, self.currentBlock))

    def hasLeagueBeenUpdated(self, leagueIdx):
        verse = self.leagues[leagueIdx].verseFinal()
        verseStatus, isVerseSettled, needsSlash = self.getVerseUpdateStatus(verse)
        return verseStatus != UPDT_NONE


    def haveNPeriodsPassed(self, verse, nPeriods):
        return (self.currentBlock - self.verseToLeagueCommits[verse].lastWriteBlocknum) > nPeriods*CHALLENGING_PERIOD_BLKS


    def getVerseSettledVerseRoot(self, verse):
        verseStatus, isVerseSettled, needsSlash = self.getVerseUpdateStatus(verse)
        assert isVerseSettled, "Asking for a settled superRoot of a not-settled verse"
        if verseStatus == UPDT_LEVEL1:
            return self.verseToLeagueCommits[verse].verseRoot
        if verseStatus == UPDT_LEVEL2:
            return self.verseToLeagueCommits[verse].superRootsVerseRoot
        assert False, "We should never be in this verse state"


    # getVerseUpdateStatus - returns:
    # - Level at which the current is (from no update to Level 1,2,3,4)
    # - Should someone be slashed?
    # - is verse settled
    def getVerseUpdateStatus(self, verse):
        # If verse was never updated, return immediately
        if not (verse in self.verseToLeagueCommits):
            verseStatus     = UPDT_NONE
            needsSlash      = UPDT_NONE
            isVerseSettled  = False
            return verseStatus, isVerseSettled, needsSlash

        # Start from the bottom. If there is Level 4 data:
        if self.verseToLeagueCommits[verse].oneLeagueDataOwner:
            if self.haveNPeriodsPassed(verse, 2):   # successful, since time passed, and settled
                verseStatus     = UPDT_LEVEL2     # so move to Level 2
                needsSlash      = UPDT_LEVEL3      # and report slash for Level 3
                isVerseSettled  = True
            elif self.haveNPeriodsPassed(verse, 1): # successful, since time passed, but not settled yet
                verseStatus     = UPDT_LEVEL2     # so move to Level 2
                needsSlash      = UPDT_LEVEL3      # and report slash for Level 3
                isVerseSettled  = False
            else:                                   # not sure if successful yet, need more time
                verseStatus     = UPDT_LEVEL4    # so, still at Level 4
                needsSlash      = UPDT_NONE
                isVerseSettled  = False
            return verseStatus, isVerseSettled, needsSlash

        # If we're here, there's not Level 4 data.
        # If there's Level 3 data:
        if self.verseToLeagueCommits[verse].leagueRootsOwner:
            if self.haveNPeriodsPassed(verse, 2):   # successful, since time passed, and settled
                verseStatus     = UPDT_LEVEL1        # so move to Level 1
                needsSlash      = UPDT_LEVEL2     # and report slash for Level 2
                isVerseSettled  = True
            elif self.haveNPeriodsPassed(verse, 1): # successful, since time passed, but not settled yet
                verseStatus     = UPDT_LEVEL1        # so move to Level 1
                needsSlash      = UPDT_LEVEL2     # and report slash for Level 2
                isVerseSettled  = False
            else:                                   # not sure if successful yet, need more time
                verseStatus     = UPDT_LEVEL3      # so, still at Level 3
                needsSlash      = UPDT_NONE
                isVerseSettled  = False
            return verseStatus, isVerseSettled, needsSlash

        # If we're here, there's not Level 3 nor Level 4 data.
        # If there's Level 2 data:
        if self.verseToLeagueCommits[verse].superRootsOwner:
            if self.haveNPeriodsPassed(verse, 1):   # successful, since time passed, and settled
                verseStatus     = UPDT_LEVEL2     # so stay at Level 2
                needsSlash      = UPDT_LEVEL1        # and slash the guy at Level 1
                isVerseSettled  = True
            else:                                   # not sure if successful yet, need more time
                verseStatus     = UPDT_LEVEL2     # so stay at Level 2
                needsSlash      = UPDT_NONE
                isVerseSettled  = False
            return verseStatus, isVerseSettled, needsSlash

        # If we're here, there's not Level 2, 3 nor Level 4 data.
        # And there is only Level 1 data.
        # So, isSettled here is just a checking time.
        verseStatus     = UPDT_LEVEL1
        needsSlash      = UPDT_NONE
        isVerseSettled  = self.haveNPeriodsPassed(verse, 1)
        return verseStatus, isVerseSettled, needsSlash

    # fail unless the status of a verse is as expected, and not settled yet.
    # if does not fail, returns if someone needs slashing
    def assertCanChallengeStatus(self, verse, status):
        verseStatus, isVerseSettled, needsSlash = self.getVerseUpdateStatus(verse)
        assert not isVerseSettled, "Verse Settled already, cannot challenge it"
        assert verseStatus == status, "Verse not ready to challenge this status"
        return needsSlash

    def updateVerseRoot(self, verse, verseRoot, addr):
        self.assertCanChallengeStatus(verse, UPDT_NONE)
        self.verseToLeagueCommits[verse] = VerseUpdate(verseRoot, addr, self.currentBlock)

    def updateLeague(self, leagueIdx, initSkillsHash, dataAtMatchdayHashes, scores, updaterAddr):
        assert self.hasLeagueFinished(leagueIdx), "League cannot be updated before the last matchday finishes"
        assert not self.hasLeagueBeenUpdated(leagueIdx), "League has already been updated"
        self.leagues[leagueIdx].updateLeague(
            initSkillsHash,
            dataAtMatchdayHashes,
            scores,
            updaterAddr,
            self.currentBlock,
        )

    def computeUsersInitDataHash(self, usersInitData):
        hash = 0
        nTeams = len(usersInitData["teamIdxs"])
        assert nTeams == len(usersInitData["teamOrders"]), "init data not consistent"
        assert nTeams == len(usersInitData["tactics"]), "init data not consistent"
        for team in range(nTeams):
            teamIdx = usersInitData["teamIdxs"][team]
            teamOrders = usersInitData["teamOrders"][team]
            tactics = usersInitData["tactics"][team]
            hash = pylio.serialHash([hash, teamIdx, teamOrders, tactics])
        return hash

    # note that values = actionsAtSelectedMatchday, formated so that is has the form
    # {idx: actionsAtSelectedMatchday}, where idx is the leaf idx.
    # so it should happen that both things coincide.
    def challengeMatchdayStates(self,
            verse,
            selectedMatchday,
            dataAtPrevMatchday,
            usersInitData,
            merkleProofDataForMatchday
        ):
        posInSubVerse = self.verseToLeagueCommits[verse].posInSubVerse
        leagueIdx = self.getLeagueIdxFromPosInSubverse(verse, posInSubVerse)

        assert self.hasLeagueBeenUpdated(leagueIdx), "League has not been updated yet, no need to challenge"
        # TODO: re-put isFullyVerified in next line
        # assert not self.isFullyVerified(leagueIdx), "You cannot challenge after the challenging period"
        assert self.computeUsersInitDataHash(usersInitData) == self.leagues[leagueIdx].usersInitDataHash, "Incorrect provided: usersInitData"
        assert merkleProofDataForMatchday.leaf[0] == leagueIdx, "Deverr: The actions do not belong to this league"
        verseActions = self.leagues[leagueIdx].verseInit + selectedMatchday * self.leagues[leagueIdx].verseStep

        # Validate that the provided actions where in the verse MerkleRoot
        assert pylio.verifyMerkleProof(
            self.VerseActionsCommits[verseActions].actionsMerkleRoots,
            merkleProofDataForMatchday,
            pylio.serialHash,
        ), "Actions are not part of the corresponding commit"

        # Validate "dataAtPrevMatchday"
        # - if day =0, validate only that skills coincide with initSkillsHash,
        #              and initialize tactics and orders from usersInitData
        # - if day!=0, validate that the entire hash of dataAtPrevMatchday coincides with
        #               the hashes that the updater provided
        if selectedMatchday == 0:
            assert pylio.serialHash(dataAtPrevMatchday.skillsAtMatchday) == self.verseToLeagueCommits[verse].dataToChallengeLeague.initSkillsHash, "Incorrect provided: prevMatchdayStates"
            # initialize tactics and teams as written in league creation:
            assert dataAtPrevMatchday.tacticsAtMatchday == 0, "Incorrect provided: prevMatchdayStates"
            assert dataAtPrevMatchday.teamOrdersAtMatchday == 0, "Incorrect provided: prevMatchdayStates"
            dataAtPrevMatchday.tacticsAtMatchday = usersInitData["tactics"]
            dataAtPrevMatchday.teamOrdersAtMatchday = usersInitData["teamOrders"]
        else:
            assert self.verseToLeagueCommits[verse].dataToChallengeLeague.dataAtMatchdayHashes[selectedMatchday-1] == self.prepareOneMatchdayHash(dataAtPrevMatchday),\
                "Incorrect provided: dataAtPrevMatchday"

        actionsAtSelectedMatchday = merkleProofDataForMatchday.leaf[1]
        self.updateTactics(
            dataAtPrevMatchday.tacticsAtMatchday,
            dataAtPrevMatchday.teamOrdersAtMatchday,
            actionsAtSelectedMatchday,
            usersInitData
        )

        dataAtPrevMatchday.skillsAtMatchday, scores = pylio.computeStatesAtMatchday(
            selectedMatchday,
            pylio.duplicate(dataAtPrevMatchday.skillsAtMatchday),
            pylio.duplicate(dataAtPrevMatchday.tacticsAtMatchday),
            pylio.duplicate(dataAtPrevMatchday.teamOrdersAtMatchday),
            self.getSeedForVerse(verseActions)
        )

        dataAtMatchdayHash = self.prepareOneMatchdayHash(dataAtPrevMatchday)

        if not dataAtMatchdayHash == self.verseToLeagueCommits[verse].dataToChallengeLeague.dataAtMatchdayHashes[selectedMatchday]:
            print("Challenger Wins: skillsAtMatchday provided by updater are invalid")
            self.verseToLeagueCommits[verse].slashLevel4(self.currentBlock)
            return

        if not (scores == self.verseToLeagueCommits[verse].dataToChallengeLeague.scores[selectedMatchday]).all():
            print("Challenger Wins: scores provided by updater are invalid")
            self.verseToLeagueCommits[verse].slashLevel4(self.currentBlock)
            return

        print("Challenger failed to prove that skillsAtMatchday nor scores were wrong")


    def isShirtNumFree(self, countryIdx, teamIdxInCountry, shirtNum):
        return self.countries[countryIdx].teamIdxInCountryToTeam[teamIdxInCountry].playerIdxs[shirtNum] == UINTMINUS1

    def getPlayerIdxFromTeamIdxAndShirt(self, teamIdx, shirtNum):
        # If player has never been sold (virtual team): simple relation between playerIdx and (teamIdx, shirtNum)
        # Otherwise, read what's written in the playerState
        # playerIdx = 0 and teamdIdx = 0 are the null player and teams
        self.assertTeamIdx(teamIdx)
        if self.isShirtNumFree(teamIdx, shirtNum):
            return UINTMINUS1
        isPlayerIdxAssigned = self.teams[teamIdx].playerIdxs[shirtNum] != 0
        if isPlayerIdxAssigned:
            return self.teams[teamIdx].playerIdxs[shirtNum]
        else:
            return 1 + (teamIdx - 1) * PLAYERS_PER_TEAM_MAX + shirtNum

    def assertTeamIdx(self, teamIdx):
        assert teamIdx < len(self.teams), "Team for this playerIdx not created yet!"
        assert teamIdx != 0, "Team 0 is reserved for null team!"

    def getLastWrittenInBCPlayerStateFromPlayerIdx(self, playerIdx):
        if self.isPlayerVirtual(playerIdx):
            return self.getPlayerStateAtBirth(playerIdx)
        else:
            return self.playerIdxToPlayerState[playerIdx]

    def getPlayerStateAtBirth(self, playerIdx):
        (countryIdx, playerIdxInCountry) = self.decodeCountryAndVal(playerIdx)
        (teamIdxInCountry, shirtNum) = self.getTeamIdxInCountryAndShirtNumFromPlayerIdxInCountry(playerIdxInCountry)
        playerState = PlayerState()
        playerState.setCurrentTeamIdx(self.encodeCountryAndVal(countryIdx, teamIdxInCountry))
        playerState.setCurrentShirtNum(shirtNum)
        return playerState


    # def getPlayerStateBeforePlayingAnyLeague(self, playerIdx):
    #     # this can be called by BC or CLIENT, as both have enough data
    #     playerStateAtBirth = self.getPlayerStateAtBirth(playerIdx)
    #
    #     if self.isPlayerVirtual(playerIdx):
    #         return playerStateAtBirth
    #     else:
    #         # if player has been sold before playing any league, it'll conserve skills at birth,
    #         # but have different metadata in the other fields
    #         playerState = pylio.duplicate(self.playerIdxToPlayerState[playerIdx])
    #         return self.copySkillsAndAgeFromTo(playerStateAtBirth, playerState)

    def copySkillsAndAgeFromTo(self, playerStateOrig, playerStateDest):
        newPlayerState = pylio.duplicate(playerStateDest)
        newPlayerState.setSkills(pylio.duplicate(playerStateOrig.getSkills()))
        newPlayerState.setMonth(pylio.duplicate(playerStateOrig.getMonth()))
        return newPlayerState


    # Given a seed, returns a balanced player.
    # It only deals with skills & age, not playerIdx.
    def getPlayerSkillsFromSeed(self, seed, monthOfTeamCreationInUnixTime):
        newPlayerState = PlayerSkills()
        np.random.seed(seed % 2**32) # we need mod(.,32) due to numpy limitation
        monthsAtBirth = np.random.randint(MIN_PLAYER_AGE, MAX_PLAYER_AGE) * 12
        newPlayerState.setMonth(monthOfTeamCreationInUnixTime-monthsAtBirth)
        skills = np.random.randint(0, AVG_SKILL - 1, N_SKILLS)
        excess = int((AVG_SKILL * N_SKILLS - skills.sum()) / N_SKILLS)
        skills += excess
        newPlayerState.setSkills(skills)
        return newPlayerState



    def getPlayerSkillsAtBirth(self, playerIdx):
        # Disregard his current team, just look at the team at moment of birth to build skills
        (countryIdx, playerIdxInCountry) = self.decodeCountryAndVal(playerIdx)
        (teamIdxInCountry, shirtNum) = self.getTeamIdxInCountryAndShirtNumFromPlayerIdxInCountry(playerIdxInCountry)
        playerDNA = pylio.serialHash([teamIdxInCountry, shirtNum])

        divisionIdx = self.getDisivionIdxFromTeamIdxInCountry(teamIdxInCountry)
        creationDay = self.getDivisionCreationDay(countryIdx, divisionIdx)
        monthOfTeamCreationInUnixTime = self.verseToUnixMonths(creationDay * VERSES_PER_DAY)

        minimalPlayerState = pylio.duplicate(self.getPlayerSkillsFromSeed(
            playerDNA,
            monthOfTeamCreationInUnixTime
        ))
        # Once the skills have been added, complete the rest of the player data
        minimalPlayerState.setPlayerIdx(playerIdx)
        return minimalPlayerState


    # The inverse of the previous relation
    def getCurrentTeamIdxAndShirtForPlayerIdx(self, playerIdx):
        if self.isPlayerVirtual(playerIdx):
            (countryIdx, playerIdxInCountry) = self.decodeCountryAndVal(playerIdx)
            (teamdIdxInCountry, shirtNum) = self.getTeamIdxInCountryAndShirtNumFromPlayerIdxInCountry(playerIdxInCountry)
            return self.encodeCountryAndVal(countryIdx, teamdIdxInCountry), shirtNum
        else:
            return self.playerIdxToPlayerState[playerIdx].getCurrentTeamIdx(), \
                   self.playerIdxToPlayerState[playerIdx].getCurrentShirtNum()

    # if player has never been sold, it will not be in the map playerIdxToPlayerState
    # and his team is derived from a formula
    def isPlayerVirtual(self, playerIdx):
        return not playerIdx in self.playerIdxToPlayerState

    def verse2blockNum(self, verse):
        return self.VerseActionsCommits[verse].blockNum

    def getLastPlayedLeagueIdx(self, playerIdx):
        # if player state has never been written, it played all leagues with current team (obtained from formula)
        # otherwise, we check if it was sold to current team before start of team's previous league
        if self.isPlayerVirtual(playerIdx):
            teamIdx, shirtNum = self.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx)
            return self.teams[teamIdx].prevLeagueIdx, self.teams[teamIdx].teamPosInPrevLeague

        currentTeamIdx = self.playerIdxToPlayerState[playerIdx].getCurrentTeamIdx()
        prevLeagueIdxForCurrentTeam = self.teams[currentTeamIdx].prevLeagueIdx
        didHePlayLastLeagueWithCurrentTeam = \
            self.playerIdxToPlayerState[playerIdx].getLastSaleBlocknum() < \
                                             self.verse2blockNum(self.leagues[prevLeagueIdxForCurrentTeam].verseInit)
        if didHePlayLastLeagueWithCurrentTeam:
            return prevLeagueIdxForCurrentTeam, self.teams[currentTeamIdx].teamPosInPrevLeague
        else:
            return self.playerIdxToPlayerState[playerIdx].prevLeagueIdx, self.playerIdxToPlayerState[
                playerIdx].prevTeamPosInLeague

    # returns skills of all teams at start of a league, basically equal to skills at end of previous league,
    # from the provided dataToChallengeInitSkills.
    # It does an extra check to make sure that the dataToChallengeInitSkills matches the previous league final matchday hash
    def verifyInitSkills(self, usersInitData, dataToChallengeInitSkills):
        nTeams = len(usersInitData["teamIdxs"])
        # an array of size [nTeams][PLAYERS_PER_TEAM_INIT]
        initPlayerSkills = pylio.createEmptyPlayerStatesForAllTeams(nTeams)
        for teamPosInLeague, teamIdx in enumerate(usersInitData["teamIdxs"]):
            for shirtNum in range(PLAYERS_PER_TEAM_MAX):
                if self.isShirtNumFree(teamIdx, shirtNum):
                    initPlayerSkills[teamPosInLeague][shirtNum] = PlayerSkills()
                    continue
                playerIdx = self.getPlayerIdxFromTeamIdxAndShirt(teamIdx, shirtNum)
                playerSkills = dataToChallengeInitSkills[teamPosInLeague][shirtNum].merkleProofStates.leaf
                assert playerSkills.getPlayerIdx() == playerIdx, "The playerIdx provided does not agree with what the BC expects"
                # it makes sure that the state matches what the BC says about that player
                if not self.areLatestSkills(dataToChallengeInitSkills[teamPosInLeague][shirtNum]):
                    return None
                initPlayerSkills[teamPosInLeague][shirtNum] = playerSkills
        return pylio.duplicate(initPlayerSkills)

    def getTeamPosInLeague(self, teamIdx, leagueUsersInitData):
        for tPos, tIdx in enumerate(leagueUsersInitData["teamIdxs"]):
            if teamIdx == tIdx:
                return tPos
        assert False, "Team not found in league"

    def areLatestSkills(self, dataToChallengeLatestSkills):
        # If player has never played a league, we can compute the playerSkills directly in the BC
        # It basically is equal to the birth skills, with ,potentially, a few team changes via sales.
        # If not, we can just compare the hash of the dataToChallengePlayerState with the stored hash in the prev league
        playerSkills = dataToChallengeLatestSkills.merkleProofStates.leaf
        playerIdx = playerSkills.getPlayerIdx()
        prevLeagueIdx, teamPosInPrevLeague = self.getLastPlayedLeagueIdx(playerIdx)
        if prevLeagueIdx == 0:
            return pylio.areEqualStructs(
                playerSkills,
                PlayerSkills(self.getPlayerStateAtBirth(playerIdx))
            )
        else:
            # First verify that the data provided match with the prevLeague SuperRoot:
            #   OneLeagueData
            #   => leads to leagueRoot which is included in the provided allLeagueRoots
            #   => which leads to a superRoot which matchs the one provided in the verse update
            leagueFinalVerse = self.leagues[prevLeagueIdx].verseFinal()
            settledVerseRoot = self.getVerseSettledVerseRoot(leagueFinalVerse)

            if not pylio.verifyMerkleProof(
                settledVerseRoot,
                dataToChallengeLatestSkills.merkleProofSuperRoots,
                pylio.serialHash
            ):
                print("SuperRoot not part of VerseRoot MerkleTree")
                return False

            if not pylio.verifyMerkleProof(
                dataToChallengeLatestSkills.merkleProofSuperRoots.leaf,
                dataToChallengeLatestSkills.merkleProofLeagueRoots,
                pylio.serialHash
            ):
                print("LeagueRoot not part of SuperRoot MerkleTree")
                return False

            if not pylio.verifyMerkleProof(
                dataToChallengeLatestSkills.merkleProofLeagueRoots.leaf,
                dataToChallengeLatestSkills.merkleProofLeague,
                pylio.serialHash
            ):
                print("LeagueData not part of League MerkleTree")
                return False

            if not pylio.verifyMerkleProof(
                dataToChallengeLatestSkills.merkleProofLeague.leaf,
                dataToChallengeLatestSkills.merkleProofStates,
                pylio.serialHash
            ):
                print("LeagueStates not part of LeagueData MerkleTree")
                return False
            return True



    def challengeLevel4InitSkills(self, verse, usersInitData, dataToChallengeInitSkills):
        self.assertCanChallengeStatus(verse, UPDT_LEVEL4)

        posInSubVerse = self.verseToLeagueCommits[verse].posInSubVerse
        leagueIdx = self.getLeagueIdxFromPosInSubverse(verse, posInSubVerse)

        leagueRoot = self.verseToLeagueCommits[verse].leagueRoots[posInSubVerse]

        assert leagueRoot != 0, "You cannot challenge a league that is not part of the verse commit"
        assert self.hasLeagueBeenUpdated(leagueIdx), "League has not been updated yet, no need to challenge"
        assert self.computeUsersInitDataHash(usersInitData) == self.leagues[leagueIdx].usersInitDataHash, "Incorrect provided: usersInitData"

        # it first makes sure that the provided initSkills are certified as the last ones.
        initSkills = self.verifyInitSkills(usersInitData, dataToChallengeInitSkills)
        # if None is returned, it means that at least one player had incorrect challenge data
        if not initSkills:
            print("Challenger failed to provide certified initSkills")
            return

        # We now know that the initSkills were correct. We just check that
        # the updater had not provided exactly the same correct skills!
        if pylio.serialHash(initSkills) == self.verseToLeagueCommits[verse].dataToChallengeLeague.initSkillsHash:
            print("Challenger failed to prove that initStates were wrong")
        else:
            print("Challenger Wins: initStates provided by updater are invalid")
            self.verseToLeagueCommits[verse].slashLevel4(self.currentBlock)

    def getBlockNumForLastLeagueOfTeam(self, teamIdx):
        (countryIdx, teamIdxInCountry) = self.decodeCountryAndVal(teamIdx)
        timeZone = self.getCountryTimeZone(countryIdx)
        verseStart = self.getVerseLeaguesStartFromTimeZoneAndRound(timeZone, self.currentRound())
        return self.verse2blockNum(verseStart)

    def getFreeShirtNum(self, teamIdx):
        (countryIdx, teamIdxInCountry) = self.decodeCountryAndVal(teamIdx)
        for shirtNum in range(PLAYERS_PER_TEAM_MAX-1, -1, -1):
            if self.isShirtNumFree(countryIdx, teamIdxInCountry, shirtNum):
                return shirtNum
        assert "Team is already full"

    # does not check ownership
    def movePlayerToTeam(self, playerIdx, buyerTeamIdx):
        assert not self.isBotTeam(buyerTeamIdx), "cannot transfer players to Bot teams"
        assert self.isPlayerTransferable(playerIdx), "Player sale failed: player is busy playing a league, wait until it finishes"
        sellerTeamIdx, sellerShirtNum = self.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx)
        assert not self.isBotTeam(sellerTeamIdx), "cannot transfer players from Bot teams"
        buyerShirtNum = self.getFreeShirtNum(buyerTeamIdx)

        # get states from BC in memory to do changes, and only write back once at the end
        state = pylio.duplicate(self.getLastWrittenInBCPlayerStateFromPlayerIdx(playerIdx))

        # a player should change his prevLeagueIdx only if the current team played
        # a last league that started AFTER the last sale
        if self.currentRound() == 0 or self.getBlockNumForLastLeagueOfTeam(sellerTeamIdx) > state.getLastSaleBlocknum():
            state.setPrevPlayedTeamIdx(sellerTeamIdx)

        state.setCurrentTeamIdx(buyerTeamIdx)
        state.setCurrentShirtNum(buyerShirtNum)
        state.setLastSaleBlocknum(self.currentBlock)

        sellerCountryIdx, sellerTeamIdxInCountry    = self.decodeCountryAndVal(sellerTeamIdx)
        buyerCountryIdx, buyerTeamIdxInCountry      = self.decodeCountryAndVal(buyerTeamIdx)
        self.countries[sellerCountryIdx].teamIdxInCountryToTeam[sellerTeamIdxInCountry].playerIdxs[sellerShirtNum] = UINTMINUS1
        self.countries[buyerCountryIdx].teamIdxInCountryToTeam[buyerTeamIdxInCountry].playerIdxs[buyerShirtNum] = playerIdx

        self.playerIdxToPlayerState[playerIdx] = pylio.duplicate(state)



    # quick solution to simulate changing teams.
    # for the purpose of Lionel, we'll start with a simple exchange, instead
    # of the more convoluted sell, assign, etc.
    def exchangePlayers(self, playerIdx1, address1, playerIdx2, address2):
        teamIdx1, shirtNum1 = self.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx1)
        teamIdx2, shirtNum2 = self.getCurrentTeamIdxAndShirtForPlayerIdx(playerIdx2)

        # check ownership!
        assert self.getOwnerAddrFromTeamIdx(teamIdx1) == address1, "Exchange Failed, owner not correct"
        assert self.getOwnerAddrFromTeamIdx(teamIdx2) == address2, "Exchange Failed, owner not correct"

        self.movePlayerToTeam(playerIdx1, teamIdx2)
        self.movePlayerToTeam(playerIdx2, teamIdx1)





    def addLeagueToVerse(self, leagueIdx, verse):
        if verse in self.verseToFinishingLeagueIdxs:
            self.verseToFinishingLeagueIdxs[verse].append(leagueIdx)
        else:
            self.verseToFinishingLeagueIdxs[verse] = [leagueIdx]

    # Creates the league in the BC, storing only the hash of usersInitData
    # It signs teams in League, which allows the BC to now that they're busy
    # without 'seeing' the pre-hash usersInitData
    def createLeague(self, verseInit, verseStep, usersInitData):
        assert verseInit > self.currentVerse, "League cannot start in the past"
        assert not self.areTeamsBusyInPrevLeagues(usersInitData["teamIdxs"]), "League cannot create: some teams involved in prev leagues"
        nTeams = len(usersInitData["teamIdxs"])
        assert nTeams % 2 == 0, "Currently we only support leagues with even nTeams"
        leagueIdx = len(self.leagues)
        self.leagues.append(League(verseInit, verseStep, nTeams))
        self.addLeagueToVerse(leagueIdx, self.leagues[leagueIdx].verseFinal())
        self.signTeamsInLeague(usersInitData, leagueIdx)
        return leagueIdx



    def signTeamsInLeague(self, usersInitData, leagueIdx):
        nTeams = len(usersInitData["teamIdxs"])
        assert nTeams == len(usersInitData["teamOrders"]), "init data not consistent"
        for team in range(nTeams):
            teamIdx     = usersInitData["teamIdxs"][team]
            teamOrders  = usersInitData["teamOrders"][team]
            tactics     = usersInitData["tactics"][team]
            self.leagues[leagueIdx].signTeamInLeague(teamIdx, teamOrders, tactics)

            self.teams[teamIdx].prevLeagueIdx             = pylio.duplicate(self.teams[teamIdx].currentLeagueIdx)
            self.teams[teamIdx].teamPosInPrevLeague       = pylio.duplicate(self.teams[teamIdx].teamPosInCurrentLeague)

            self.teams[teamIdx].currentLeagueIdx          = leagueIdx
            self.teams[teamIdx].teamPosInCurrentLeague    = team

    # Minimal (virtual) team creation. The Name could be the concat of the given name, and user int choice
    # e.g. teamName = "Barcelona5443"
    def createTeam(self, teamName, ownerAddr):
        assert pylio.intHash(teamName) not in self.teamNameHashToOwnerAddr, "You cannot create to teams with equal name!"
        teamIdx = len(self.teams)
        nowInMonthsUnixTime = 602
        self.teams.append(Team(teamName, nowInMonthsUnixTime))
        self.teamNameHashToOwnerAddr[pylio.intHash(teamName)] = ownerAddr
        return teamIdx


    # ------------ LEAGUE STATUS --------------
    def isLeagueIsAboutToStart(self, leagueIdx):
        return self.currentVerse < self.leagues[leagueIdx].verseInit

    def hasLeagueStarted(self, leagueIdx):
        return self.currentVerse >= self.leagues[leagueIdx].verseInit

    def hasLeagueFinished(self, leagueIdx):
        return self.currentVerse >= self.leagues[leagueIdx].verseFinal()

    def isFullyVerified(self, leagueIdx):
        if self.leagues[leagueIdx].isGenesisLeague():
            return True
        verse = self.leagues[leagueIdx].verseFinal()
        verseStatus, isVerseSettled, needsSlash = self.getVerseUpdateStatus(verse)
        return isVerseSettled

    def getPlayerSkillsFromChallengeData(self, playerIdx, dataToChallengePlayerState):
        # dataToChallengePlayerState is either:
        #  - just a player state
        #  - a merkle proof for that player
        # In the latter case, we can extract the state just from the values (leafs)
        if type(dataToChallengePlayerState) == type(PlayerSkills()):
            assert dataToChallengePlayerState.getPlayerIdx() == playerIdx, "This data does not contain the required playerIdx"
            return dataToChallengePlayerState
        else:
            playerSkills = dataToChallengePlayerState.leaf
            assert playerSkills.getPlayerIdx() == playerIdx, "This data does not contain the required player"
            return pylio.duplicate(playerSkills)

    def getOwnerAddrFromTeamIdx(self, teamIdx):
        if self.isBotTeam(teamIdx):
            return FREEVERSE
        (countryIdx, teamIdxInCountry) = self.decodeCountryAndVal(teamIdx)
        return self.countries[countryIdx].teamIdxInCountryToTeam[teamIdxInCountry].teamOwner

    def getOwnerAddrFromPlayerIdx(self, playerIdx):
        currentTeamIdx = self.getLastWrittenInBCPlayerStateFromPlayerIdx(playerIdx).currentTeamIdx
        return self.getOwnerAddrFromTeamIdx(currentTeamIdx)


    # A mockup of how to obtain the block hash for a given blocknum.
    # This is a function that is available in Ethereum after Constatinople
    def getBlockHash(self, blockNum):
        return pylio.intHash('salt' + str(blockNum))

    def getSeedForVerse(self, verse):
        return self.getBlockHash(self.VerseActionsCommits[verse].blockNum)


    # From the states at a given matchday, we just need to store the hash... of the skills,
    # ... disregarding other side info, like lastSaleBlock...
    # This is important, because otherwise, it's impossible to use these hashes for challenges once
    # sales have taken place.
    def prepareOneMatchdayHash(self, dataAtMatchday):
        # note that the matrix has size: statesAtOneMatchday[team][player]
        # we basically convert from 'states' to 'skills':
        #   dataAtMatchday.statesAtOneMatchday[team][player] --> dataAtMatchday.skillsAtOneMatchday[team][player]
        skillsAtOneMatchday = []
        for teams in dataAtMatchday.skillsAtMatchday:
            allTeamSkills = [s.getSkills() for s in teams]
            skillsAtOneMatchday.append(pylio.duplicate(allTeamSkills))

        updatedData = pylio.duplicate(dataAtMatchday)
        updatedData.skillsAtMatchday = skillsAtOneMatchday

        return pylio.serialHash(updatedData)

    def certifyPlayerState(self, playerState, dataToChallengePlayerSkills):
        # check that the skills inside playerState match the end of last league:
        playerSkills = PlayerSkills(playerState)
        dataToChallengePlayerSkills.merkleProof = MerkleProof(
            dataToChallengePlayerSkills.merkleProofStates.neededHashes,
            dataToChallengePlayerSkills.merkleProofStates.depth,
            playerSkills,
            0
        )
        assert self.areLatestSkills(dataToChallengePlayerSkills), "Computed player state by CLIENT is not recognized by BC.."
        # evolve skills to last written state in the BC
        currentStateCertified = self.skillsToLastWrittenState(playerSkills)
        return pylio.areEqualStructs(playerState, currentStateCertified)

    # If we start from the state at the end of last played league, then only the skills remain unchanged.
    # In general, the player can have been sold many times up to the current time.
    # So we start with whatever state is currently written, and insert the skills from end of last league
    def skillsToLastWrittenState(self, playerSkills):
        lastWrittenPlayerState = self.getLastWrittenInBCPlayerStateFromPlayerIdx(playerSkills.getPlayerIdx())
        return self.copySkillsAndAgeFromTo(playerSkills, lastWrittenPlayerState)

    def updateTactics(self, tactics, teamOrders, actions, usersInitData):
        if actions == 0:
            return
        for action in actions:
            teamPosInLeague = self.getTeamPosInLeague(action["teamIdx"], usersInitData)
            tactics[teamPosInLeague] = action["tactics"]
            teamOrders[teamPosInLeague] = action["teamOrder"]

    def challengeLevel1(self, verse, superRoots, addr):
        needsSlash = self.assertCanChallengeStatus(verse, UPDT_LEVEL1)
        if needsSlash == UPDT_LEVEL2:
            self.verseToLeagueCommits[verse].slashLevel2(self.currentBlock)
        tree = MerkleTree(superRoots)

        assert tree.root != self.verseToLeagueCommits[verse].verseRoot, \
            "The superRoots provided lead to the same verseRoot as already provided by updater"

        self.verseToLeagueCommits[verse].writeLevel2(
            superRoots,
            tree.root,
            addr,
            self.currentBlock
        )



    def challengeLevel2(self, verse, subVerse, leagueRoots, addr):
        needsSlash = self.assertCanChallengeStatus(verse, UPDT_LEVEL2)
        if needsSlash == UPDT_LEVEL3:
            self.verseToLeagueCommits[verse].slashLevel3(self.currentBlock)

        tree = MerkleTree(leagueRoots)
        assert tree.root != self.verseToLeagueCommits[verse].superRoots[subVerse], \
            "The leagueRoots provided lead to the same superRoot as already provided by updated"


        self.verseToLeagueCommits[verse].writeLevel3(
            subVerse,
            leagueRoots,
            addr,
            self.currentBlock
        )


    def getPosInSubverse(self, verse, leagueIdx):
        challengedSubVerse = self.verseToLeagueCommits[verse].subVerse
        nLeagues, nSubVerses, leagueIdxsInVerse = self.getSubVerseData(verse)
        return self.getLeaguesInSubVerse(leagueIdxsInVerse, challengedSubVerse).index(leagueIdx)

    def getLeagueIdxFromPosInSubverse(self, verse, posInSubVerse):
        challengedSubVerse = self.verseToLeagueCommits[verse].subVerse
        return self.verseToFinishingLeagueIdxs[verse][challengedSubVerse * SUPERROOTS_PER_VERSE + posInSubVerse]


    def isLeagueIdxInVerseCommit(self, verse, leagueIdx):
        for leaguePair in self.verseToLeagueCommits[verse].leagueRoots:
            if leaguePair[0] == leagueIdx:
                return True
        return False

    def challengeLevel3(self, verse, posInSubVerse, dataToChallengeLeague, addr):
        self.assertCanChallengeStatus(verse, UPDT_LEVEL3)
        leagueRoot = self.verseToLeagueCommits[verse].leagueRoots[posInSubVerse]
        assert leagueRoot != 0, "You cannot challenge a league that is not part of the verse commit"
        assert self.computeLeagueRoot(dataToChallengeLeague) != leagueRoot, \
            "Your data coincides with the updater. Nothing to challenge."
        self.verseToLeagueCommits[verse].writeLevel4(
            posInSubVerse,
            dataToChallengeLeague,
            addr,
            self.currentBlock
        )

    def flattenLeagueData(self, dataToChallengeLeague):
        leafs = [dataToChallengeLeague.initSkillsHash]
        for hash in dataToChallengeLeague.dataAtMatchdayHashes:
            leafs.append(hash)
        for score in dataToChallengeLeague.scores:
            leafs.append(score)
        return leafs

    def computeLeagueRoot(self, dataToChallengeLeague):
        leagueTree = MerkleTree(self.flattenLeagueData(dataToChallengeLeague))
        return leagueTree.root

    def isLeagueSettled(self, leagueIdx):
        verse = self.leagues[leagueIdx].verseFinal()
        verseStatus, isVerseSettled, needsSlash = self.getVerseUpdateStatus(verse)
        return isVerseSettled

    def getLeaguesInSubVerse(self, leagueIdxsInVerse, subVerse):
        firstLeague = subVerse * SUPERROOTS_PER_VERSE
        lastLeague = (subVerse + 1) * SUPERROOTS_PER_VERSE - 1
        lastLeague = min(len(leagueIdxsInVerse), lastLeague - 1)
        return leagueIdxsInVerse[firstLeague:lastLeague]

    def getSubVerseData(self, verse):
        leagueIdxsInVerse = self.getLeaguesFinishingInVerse(verse)
        nLeagues = len(leagueIdxsInVerse)
        nSubVerses = math.ceil(nLeagues/SUPERROOTS_PER_VERSE)
        return nLeagues, nSubVerses, leagueIdxsInVerse

    def getLeagueSubVerse(self, verse, leagueIdx):
        posInVerse = self.getLeaguePosInVerse(verse, leagueIdx)
        subVerse = math.floor(posInVerse / SUPERROOTS_PER_VERSE)
        posInSubVerse = posInVerse - subVerse * SUPERROOTS_PER_VERSE
        return subVerse, posInSubVerse

    def getLeaguePosInVerse(self, verse, leagueIdx):
        leaguesFinishingInVerse = self.getLeaguesFinishingInVerse(verse)
        for leaguePos, finishingLeagueIdx in enumerate(leaguesFinishingInVerse):
            if finishingLeagueIdx == leagueIdx:
                return leaguePos
        assert False, "league not found in verse!"


    def getLeaguesFinishingInVerse(self, verse):
        if not verse in self.verseToFinishingLeagueIdxs:
            return []
        else:
            return self.verseToFinishingLeagueIdxs[verse]

    def updateTimeZone(self, updateData):
        timeZoneToUpdate, posInUpdate = self.currentTimeZoneToUpdate()
        assert not timeZoneToUpdate == TZ_NULL, "nothing to update in this verse"
        assert timeZoneToUpdate in self.timeZoneUpdates, "nothing to update in this particular timeZone"
        self.timeZoneUpdates[timeZoneToUpdate].newUpdate(self.currentBlock, updateData)


    # ------------------------------------------------------------------------
    # ------------      Functions only for CLIENT                 ------------
    # ------------------------------------------------------------------------

    # return state of a player at end of a certain league
    # note that these do not contain potential sales done after the league
    def getPlayerSkillsAtEndOfLeague(self, leagueIdx, teamPosInLeague, playerIdx):
        self.assertIsClient()
        if leagueIdx == 0:
            return PlayerSkills(self.getPlayerStateAtBirth(playerIdx))

        selectedSkills = [s for s in self.leagues[leagueIdx].dataAtMatchdays[-1].skillsAtMatchday[teamPosInLeague] if
                          s.getPlayerIdx() == playerIdx]
        assert len(selectedSkills) == 1, "PlayerIdx not found in previous league final states, or too many with same playerIdx"
        return selectedSkills[0]

    def getCurrentPlayerState(self, playerIdx):
        self.assertIsClient()
        currentSkills = pylio.duplicate(self.getPlayerSkillsAtEndOfLastLeague(playerIdx))
        lastBCState = pylio.duplicate(self.getLastWrittenInBCPlayerStateFromPlayerIdx(playerIdx))
        return self.copySkillsAndAgeFromTo(currentSkills, lastBCState)

    def getPlayerSkillsAtEndOfLastLeague(self, playerIdx):
        self.assertIsClient()
        prevLeagueIdx, teamPosInPrevLeague = self.getLastPlayedLeagueIdx(playerIdx)
        return self.getPlayerSkillsAtEndOfLeague(prevLeagueIdx, teamPosInPrevLeague, playerIdx)

    # Stores the data, pre-hash, in the CLIENT
    def storePreHashDataInClientAtEndOfLeague(self, leagueIdx, dataAtMatchdays, lastDayTree):
        self.assertIsClient()
        self.leagues[leagueIdx].storeDataAtMatchdays(dataAtMatchdays)
        self.leagues[leagueIdx].lastDayTree = lastDayTree

        # the last matchday gives the final skills used to update all players:
        # After the end of the league, there could be other things, like sales, so we need to update
        # those (while keeping the skills as of last league's end)
        for skillsAtEndOfLeaguePerTeam in dataAtMatchdays[-1].skillsAtMatchday:
            for playerSkills in skillsAtEndOfLeaguePerTeam:
                self.playerIdxToPlayerState[playerSkills.getPlayerIdx()] = \
                    pylio.duplicate(self.skillsToLastWrittenState(playerSkills))

    def getPrevMatchdayData(self, leagueIdx, selectedMatchday):
        self.assertIsClient()
        if selectedMatchday == 0:
            return DataAtMatchday(
                self.leagues[leagueIdx].getInitPlayerSkills(),
                0,
                0
            )
        else:
            return pylio.duplicate(self.leagues[leagueIdx].dataAtMatchdays[selectedMatchday-1])


    # Besides creating the league, it also:
    # - computes the init states and stores them
    # - computes the data needed to challenge those init states and stores them
    def createLeagueClient(self, verseInit, verseStep, usersInitData):
        self.assertIsClient()
        assert not self.areTeamsBusyInPrevLeagues(usersInitData["teamIdxs"]), "League cannot create: some teams involved in prev leagues"
        nTeams = len(usersInitData["teamIdxs"])
        assert nTeams % 2 == 0, "Currently we only support leagues with even nTeams"
        leagueIdx = len(self.leagues)
        self.leagues.append(LeagueClient(verseInit, verseStep, usersInitData))
        self.signTeamsInLeague(usersInitData, leagueIdx)
        self.addLeagueToVerse(leagueIdx, self.leagues[leagueIdx].verseFinal())
        self.leagues[leagueIdx].writeInitState(self.getInitPlayerStates(leagueIdx))
        self.leagues[leagueIdx].writeDataToChallengeInitSkills(self.prepareDataToChallengeLeagueInitSkills(leagueIdx))
        return leagueIdx

    def addAccumulator(self):
        self.assertIsClient()
        self.Accumulator = ActionsAccumulator()

    def accumulateAction(self, action):
        self.assertIsClient()
        assert self.currentBlock >= self.lastVerseBlock(), "Weird, blocknum for action received that belonged to past commit"
        leagueIdx = self.getLeagueForAction(action)
        if self.hasLeagueFinished(leagueIdx):
            print("Cannot accept actions for leagues that already finished! Action discarded")
        else:
            self.Accumulator.accumulateAction(action, leagueIdx)

    # returns all verses were matchdays of a league took/take place
    def getVersesForLeague(self, leagueIdx):
        self.assertIsClient()
        nMatchdays = 2*(self.leagues[leagueIdx].nTeams-1)
        verses = []
        for matchday in range(nMatchdays):
            verses.append(self.leagues[leagueIdx].verseInit + matchday * self.leagues[leagueIdx].verseStep)
        return verses

    # returns all seeds used for all matchdays of a league
    def getAllSeedsForLeague(self, leagueIdx):
        self.assertIsClient()
        assert self.hasLeagueFinished(leagueIdx), "All seeds only available at end of league"
        seedsPerVerse = []
        for verse in self.getVersesForLeague(leagueIdx):
            seedsPerVerse.append(self.getSeedForVerse(verse))
        return seedsPerVerse

    # returns which league did this action refer to
    def getLeagueForAction(self, action):
        self.assertIsClient()
        return self.teams[action["teamIdx"]].currentLeagueIdx

    def getLeaguesPlayingInThisVerse(self, verse):
        self.assertIsClient()
        # TODO: make this less terribly slow
        leagueIdxs = []
        nLeagues = len(self.leagues)
        for leagueIdx in range(1,nLeagues): # bypass the first (dummy) league
            if verse in self.getVersesForLeague(leagueIdx):
                leagueIdxs.append(leagueIdx)
        return leagueIdxs


    # Sends the actions acummulated in the buffer to the BC, by sending the Merkle Root first.
    # It only sends the actions corresponding to leagues that play games at the current verse.
    # Before computing the Merkler Root, it first orders all the actions in the form:
    # [leagueIdx0, allActionsInLeagueIdx0], [leagueIdx1, allActionsInLeagueIdx1], ...
    # So each leaf has the form [leagueIdx, allActionsInLeagueIdx]
    def syncActions(self, ST):
        self.assertIsClient()
        assert self.currentBlock == ST.currentBlock, "Client and BC are out of sync in blocknum!"
        leaguesPlayingInThisVerse = self.getLeaguesPlayingInThisVerse(ST.currentVerse)
        leagueIdxAndActionsArray = []
        # Builds two quantities.
        #   - leagueIdxAndActionsArray, used to get the Merkle game
        #   - self.leagues[leagueIdx].actionsPerMatchday to store the pre-hashes
        for leagueIdx in leaguesPlayingInThisVerse:
            if leagueIdx in self.Accumulator.buffer:
                leagueIdxAndActionsArray.append([leagueIdx, self.Accumulator.buffer[leagueIdx]])
                self.leagues[leagueIdx].actionsPerMatchday.append(self.Accumulator.buffer[leagueIdx])
            else:
                leagueIdxAndActionsArray.append([leagueIdx, 0])
                self.leagues[leagueIdx].actionsPerMatchday.append(0)

        if leagueIdxAndActionsArray:
            tree = MerkleTree(leagueIdxAndActionsArray)
            rootTree    = tree.root
        else:
            tree        = 0
            rootTree    = 0

        ST.commit(rootTree)
        self.commit(rootTree)

        self.Accumulator.commitedActions.append(leagueIdxAndActionsArray)
        self.Accumulator.commitedTrees.append(tree)
        self.Accumulator.clearBuffer(pylio.duplicate(leagueIdxAndActionsArray))


    # It gathers all actions that were sent for a given selectedMatchday of a given league.
    # It then builds a proof that can be used by someone who:
    # - receives those actions (pre-hash)
    # - knows the MerkleRoot of that verse
    # ...to prove that the actions were included in that Merkle proof verse
    def getActionsMerkleProofForMatchday(self, leagueIdx, selectedMatchday):
        self.assertIsClient()
        verse = self.leagues[leagueIdx].verseInit + selectedMatchday * self.leagues[leagueIdx].verseStep
        if not self.Accumulator.commitedActions[verse]:
            return MerkleProof(0, 0, 0, 0)

        for idx, action in enumerate(self.Accumulator.commitedActions[verse]):
            if action[0] == leagueIdx:
                break

        tree = self.Accumulator.commitedTrees[verse]
        merkleProof = tree.prepareProofForLeaf(action, idx)

        assert pylio.verifyMerkleProof(
            self.VerseActionsCommits[verse].actionsMerkleRoots,
            merkleProof,
            pylio.serialHash
        ), "Generated Merkle proof will not work"

        return merkleProof


    def computeAllMatchdayStates(self, leagueIdx):
        self.assertIsClient()
        initPlayerSkills = self.leagues[leagueIdx].getInitPlayerSkills()
        usersInitData = pylio.duplicate(self.leagues[leagueIdx].usersInitData)
        seedsPerVerse = self.getAllSeedsForLeague(leagueIdx)

        # In this initial implementation, evolution happens at the end of the league only
        tactics     = pylio.duplicate(usersInitData["tactics"])
        teamOrders  = pylio.duplicate(usersInitData["teamOrders"])

        nTeams = len(usersInitData["teamIdxs"])
        nMatchdays = 2*(nTeams-1)
        assert nMatchdays == len(seedsPerVerse), "We should have as many matchdays as verses"
        nMatchesPerMatchday = nTeams//2
        scores = np.zeros([nMatchdays, nMatchesPerMatchday, 2], int)

        dataAtMatchdays = []

        skillsAtMatchday = initPlayerSkills
        for matchday in range(nMatchdays):
            actionsInThisMatchday = pylio.duplicate(self.leagues[leagueIdx].actionsPerMatchday[matchday])
            self.updateTactics(tactics, teamOrders, actionsInThisMatchday, self.leagues[leagueIdx].usersInitData)
            skillsAtMatchday, scores[matchday] = pylio.computeStatesAtMatchday(
                matchday,
                pylio.duplicate(skillsAtMatchday),
                pylio.duplicate(tactics),
                pylio.duplicate(teamOrders),
                seedsPerVerse[matchday]
            )
            dataAtMatchdays.append(DataAtMatchday(skillsAtMatchday, tactics, teamOrders))

        return dataAtMatchdays, scores

    # Data needed to challenge the init states of a league. If the player has never played before,
    # it's easy, otherwise, it needs to prove that his state is in the final states of a previous league...
    # In all cases it returns an array [N_PLAyERS, nTeams] where each entry is a MerkleProof
    def prepareDataToChallengeLeagueInitSkills(self, leagueIdx):
        self.assertIsClient()
        thisLeague = pylio.duplicate(self.leagues[leagueIdx])
        nTeams = len(thisLeague.usersInitData["teamIdxs"])
        dataToChallengeInitSkills = pylio.createEmptyPlayerStatesForAllTeams(nTeams)
        for teamPos, teamIdx in enumerate(thisLeague.usersInitData["teamIdxs"]):
            for shirtNum, playerIdx in enumerate(self.teams[teamIdx].playerIdxs):
                if self.isShirtNumFree(teamIdx, shirtNum):
                    dataToChallengeInitSkills[teamPos][shirtNum] = DataToChallengePlayerSkills(0, 0, 0, 0)
                elif playerIdx == 0: # if never written in teams.playerIdxs array
                    dataToChallengeInitSkills[teamPos][shirtNum] = self.computeDataToChallengePlayerSkills(
                        self.getPlayerIdxFromTeamIdxAndShirt(teamIdx, shirtNum)
                    )
                else:
                    dataToChallengeInitSkills[teamPos][shirtNum] = self.computeDataToChallengePlayerSkills(playerIdx)
        return dataToChallengeInitSkills

    # This function uses CLIENT data to return what is needed to then be able to challenge the player skills.
    # If it has already played leagues, it returns the playerSkills and the MerkleProof that it belongs to last leagues' matchday.
    # If not, then the birth skills.
    # note: skillsAtEndOfPrevLeague does not obviously take into account possible evolution/sales after the league
    # note: yes, it returns either a playerSkills, or a MerkleProof
    def computeDataToChallengePlayerSkills(self, playerIdx):
        self.assertIsClient()
        prevLeagueIdx, teamPosInPrevLeague = self.getLastPlayedLeagueIdx(playerIdx)
        if prevLeagueIdx == 0:
            merkleProofStates = MerkleProof([], 0, self.getPlayerSkillsAtEndOfLastLeague(playerIdx), 0)
            return DataToChallengePlayerSkills(merkleProofStates, 0, 0, 0)
        else:
            # construct merkle proofs for:
            # - leagues states in prevLeague last day's hash
            # - prevLeague data in prevLeague root
            # - prevLeague root in verse superRoot
            # For each proof, we need the idx in the tree.
            # For each proof we double check that it'd pass it in the BC
            #
            # ----- leagues states in prevLeague last day's hash ------
            skillsAllTeamsAtEndOfPrevLeague = self.leagues[prevLeagueIdx].dataAtMatchdays[-1].skillsAtMatchday
            playerSkills, playerPosInPrevLeague = self.getPlayerFromTeamStates(playerIdx, skillsAllTeamsAtEndOfPrevLeague[teamPosInPrevLeague])
            idxInFlattenedSkills = teamPosInPrevLeague*PLAYERS_PER_TEAM_MAX+playerPosInPrevLeague

            lastDayTree = self.leagues[prevLeagueIdx].lastDayTree
            merkleProofStates = lastDayTree.prepareProofForLeaf(playerSkills, idxInFlattenedSkills)

            assert pylio.verifyMerkleProof(
                self.leagues[prevLeagueIdx].dataToChallengeLeague.dataAtMatchdayHashes[-1],
                merkleProofStates,
                pylio.serialHash
            ), "Generated Merkle proof will not work"

            # ----- prevLeague data in prevLeague root ------
            leagueData = self.leagues[prevLeagueIdx].dataToChallengeLeague
            leafs = self.flattenLeagueData(
                leagueData
            )
            treeLeague = MerkleTree(leafs)
            idxInFlattenedLeagueData = len(leagueData.dataAtMatchdayHashes)

            merkleProofLeague = treeLeague.prepareProofForLeaf(
                self.leagues[prevLeagueIdx].dataToChallengeLeague.dataAtMatchdayHashes[-1],
                idxInFlattenedLeagueData
            )

            assert pylio.verifyMerkleProof(
                treeLeague.root,
                merkleProofLeague,
                pylio.serialHash
            ), "Generated Merkle proof will not work"

            # ----- prevLeague root in verse superRoot ------
            verse = self.leagues[prevLeagueIdx].verseFinal()
            subVerse, posInSubVerse = self.getLeagueSubVerse(verse, prevLeagueIdx)

            superRoots, leagueRoots = self.computeLeagueHashesForVerse(verse)

            treeLeagueRoots = MerkleTree(leagueRoots[posInSubVerse])
            assert treeLeagueRoots.root == superRoots[subVerse], "Computed leagueRoots inconsistent"

            merkleProofLeagueRoots = treeLeagueRoots.prepareProofForLeaf(
                treeLeague.root,
                posInSubVerse
            )

            assert pylio.verifyMerkleProof(
                treeLeagueRoots.root,
                merkleProofLeagueRoots,
                pylio.serialHash
            ), "Generated Merkle proof will not work"

            # ----- superRoot in VerseRoot ------

            treeSuperRoots = MerkleTree(superRoots)
            assert treeSuperRoots.root == self.verseToLeagueCommits[verse].verseRoot, "Computed superRoots inconsistent"

            merkleProofSuperRoots = treeSuperRoots.prepareProofForLeaf(
                treeLeagueRoots.root,
                subVerse
            )

            assert pylio.verifyMerkleProof(
                self.verseToLeagueCommits[verse].verseRoot,
                merkleProofSuperRoots,
                pylio.serialHash
            ), "Generated Merkle proof will not work"


            return DataToChallengePlayerSkills(
                merkleProofStates,
                merkleProofLeague,
                merkleProofLeagueRoots,
                merkleProofSuperRoots
            )


    # Given all states of players in a team, returns the state corresponding to
    # the required playerIdx, as well as its position in the team.
    def getPlayerFromTeamStates(self, playerIdx, statesInTeam):
        self.assertIsClient()
        playerState = None
        playerPos   = None
        for pos, state in enumerate(statesInTeam):
            if playerIdx == state.getPlayerIdx():
                if playerState:
                    assert False, "Same player appears twice in a team!!!"
                else:
                    playerState             = pylio.duplicate(state)
                    playerPos   = pylio.duplicate(pos)
        return playerState, playerPos


    # for all days in the league, except for the last one, it basically hashes
    # the struct dataAtMatchday, which contains, besides states, the tactics and teamOrders.
    # for the last day, you just need the states, and you actuall do a MerkleRoot, so
    # that it can be later used for Merkle Proofs
    def prepareHashesForDataAtMatchdays(self, dataAtMatchdays):
        self.assertIsClient()
        # hash all except for last day:
        # dataAtMatchdayHashes = [pylio.serialHash(d) for d in dataAtMatchdays[:-1]]
        dataAtMatchdayHashes = [self.prepareOneMatchdayHash(dataAtMatchday) for dataAtMatchday in dataAtMatchdays]

        # compute MerkleRoot for last day:
        lastSkillsFlattened = pylio.flatten(dataAtMatchdays[-1].skillsAtMatchday)
        lastDayTree = MerkleTree(lastSkillsFlattened)
        dataAtMatchdayHashes.append(lastDayTree.root)
        return dataAtMatchdayHashes, lastDayTree



    # The CLIENT:
    # - computes all games of the league,
    # - in particular, all DataAtMatchdays => for every matchday: all teams states, tactics, teamOrders.
    # - stores both the pre-hash and the hashed DataAtMatchdays
    # - returns the hashed data so that it can then be send to the BC
    def updateLeagueInClient(self, leagueIdx):
        self.assertIsClient()
        assert self.hasLeagueFinished(leagueIdx), "cannot update a league that is not finished"
        assert not self.hasLeagueBeenUpdated(leagueIdx), "League has already been updated"
        dataAtMatchdays, scores = self.computeAllMatchdayStates(leagueIdx)
        initSkillsHash          = pylio.serialHash(self.leagues[leagueIdx].getInitPlayerSkills())
        dataAtMatchdayHashes, lastDayTree = self.prepareHashesForDataAtMatchdays(dataAtMatchdays)
        dataToChallengeLeague = DataToChallengeLeague(
            initSkillsHash,
            dataAtMatchdayHashes,
            scores,
        )

        self.leagues[leagueIdx].writeDataToChallengeLeague(dataToChallengeLeague)

        # and additionally, stores the league pre-hash data, and updates every player involved
        self.storePreHashDataInClientAtEndOfLeague(leagueIdx, dataAtMatchdays, lastDayTree)
        assert initSkillsHash == pylio.serialHash(self.leagues[leagueIdx].getInitPlayerSkills()), "InitSkillsHash changed unintentionally"
        assert self.hasLeagueBeenUpdatedByClient(leagueIdx), "League not detected as already updated in client"
        # return initSkillsHash, dataAtMatchdayHashes, scores

    def hasLeagueBeenUpdatedByClient(self, leagueIdx):
        return self.leagues[leagueIdx].dataToChallengeLeague.initSkillsHash != 0


    # returns states of all teams at start of a league. These include skills from previous league, and possible
    # sales after end of that league
    def getInitPlayerStates(self, leagueIdx):
        self.assertIsClient()
        usersInitData = pylio.duplicate(self.leagues[leagueIdx].usersInitData)
        nTeams = len(usersInitData["teamIdxs"])
        # an array of size [nTeams][PLAYERS_PER_TEAM_INIT]
        initPlayerStates = pylio.createEmptyPlayerStatesForAllTeams(nTeams)
        for teamPosInLeague, teamIdx in enumerate(usersInitData["teamIdxs"]):
            for shirtNum in range(PLAYERS_PER_TEAM_MAX):
                if self.isShirtNumFree(teamIdx, shirtNum):
                    initPlayerStates[teamPosInLeague][shirtNum] = PlayerState()
                else:
                    playerIdx = self.getPlayerIdxFromTeamIdxAndShirt(teamIdx, shirtNum)
                    playerState = self.getCurrentPlayerState(playerIdx)
                    initPlayerStates[teamPosInLeague][shirtNum] = playerState
        return initPlayerStates


    def computeLeagueRootFromLeagueIdx(self, leagueIdx):
        self.assertIsClient()
        return self.computeLeagueRoot(
            self.leagues[leagueIdx].dataToChallengeLeague
        )


    def computeLeagueHashesForVerse(self, verse):
        self.assertIsClient()
        nLeagues, nSubVerses, leagueIdxsInVerse = self.getSubVerseData(verse)
        superRoots = []
        leagueRoots = []

        for subVerse in range(nSubVerses):
            leagueIdxsInSubVerse = self.getLeaguesInSubVerse(leagueIdxsInVerse, subVerse)
            thisSuperRoot, thisleagueRoots = self.computeHashesForSubverse(leagueIdxsInSubVerse)
            superRoots.append(thisSuperRoot)
            leagueRoots.append(thisleagueRoots)
        return superRoots, leagueRoots


    def computeHashesForSubverse(self,leagueIdxsInVerse):
        leagueRoots = []
        for leagueIdx in leagueIdxsInVerse:
            leagueRoots.append(self.computeLeagueRootFromLeagueIdx(leagueIdx))
        tree = MerkleTree(leagueRoots)
        superRoot = tree.root
        return superRoot, leagueRoots

    def updateAllLeaguesForVerseInClient(self, verse):
        self.assertIsClient()
        leagueIdxsForThisCommit = self.getLeaguesFinishingInVerse(verse)
        if len(leagueIdxsForThisCommit) > 0:
            for leagueIdx in leagueIdxsForThisCommit:
                self.updateLeagueInClient(leagueIdx)
        return leagueIdxsForThisCommit

    def syncLeagueCommits(self, ST):
        self.assertIsClient()
        leagueIdxsForThisCommit = self.updateAllLeaguesForVerseInClient(self.currentVerse)
        if len(leagueIdxsForThisCommit) == 0:
            return
        superRoots, leagueRoots = self.computeLeagueHashesForVerse(self.currentVerse)
        tree = MerkleTree(superRoots)
        verseRoot = tree.root
        self.updateVerseRoot(self.currentVerse, verseRoot, ALICE)
        # only lie (if forced) in the BC, not locally
        verseRootFinal = pylio.duplicate(verseRoot)
        if self.forceVerseRootLie:
            verseRootFinal = verseRootFinal * 2
        ST.updateVerseRoot(self.currentVerse, verseRootFinal, ALICE)


    # TODO: implement
    def getDataForUpdate(self, timeZoneToUpdate, posInUpdate):
        self.assertIsClient()
        return 0


    def syncTimeZoneCommits(self, ST):
        self.assertIsClient()
        timeZoneToUpdate, posInUpdate = ST.currentTimeZoneToUpdate()
        if timeZoneToUpdate == TZ_NULL:
            return
        if timeZoneToUpdate in ST.timeZoneUpdates:
            dataForUpdate = self.getDataForUpdate(timeZoneToUpdate, posInUpdate)
            ST.updateTimeZone(dataForUpdate)


