import numpy as np
import copy
import datetime
from constants import *
import pylio
from pickle import dumps as serialize

# simple block counter simulator, where the blockhash is just the hash of the blocknumber
class BlockCounter():
    def __init__(self):
        self.currentBlock = 0

    def advanceNBlocks(self, n):
        self.currentBlock += n

    def advanceToBlock(self, n):
        assert n >= self.currentBlock, "Cannot advance... to a block in the past!"
        self.currentBlock = n


class Storage(BlockCounter):
    def __init__(self):

        BlockCounter.__init__(self)

        # an array of Team structs, the first entry being the null team
        self.teams = [Team("")]

        # a map from playerIdx to playerState, only available for players already sold once,
        # or for 'promo players' not created directly from team creation.
        # In Python, maps are closer to 'dictionaries'
        self.playerIdxToPlayerState = {}

        # the obvious ownership map:
        self.teamNameHashToOwnerAddr = {}

        # an array of leagues, first entry is dummy
        self.leagues = [League(0,0,0)]



# In Solidity, PlayerState will be just a uin256, serializing the data shown here,
# and there'll be associated read/write functions
# playerIdx = 0 is the null player
class PlayerState():
    def __init__(self):
        self.skills                  = np.zeros(N_SKILLS)
        self.monthOfBirthInUnixTime  = 0
        self.playerIdx               = 0
        self.currentTeamIdx          = 0
        self.currentShirtNum         = 0
        self.prevLeagueIdx          = 0
        self.prevTeamPosInLeague    = 0
        self.prevShirtNumInLeague   = 0
        self.lastSaleBlocknum        = 0

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

    def setCurrentTeamIdx(self, currentTeamIdx):
        self.currentTeamIdx = currentTeamIdx

    def getPrevLeagueIdx(self):
        return self.prevLeagueIdx

    def setPrevTeamPosInLeague(self, prevTeamPosInLeague):
        self.prevTeamPosInLeague = prevTeamPosInLeague

    def getCurrentTeamIdx(self):
        return self.currentTeamIdx

    def getPrevTeamIdxLeague(self):
        return self.prevTeamPosInLeague

    def setCurrentShirtNum(self, currentShirtNum):
        self.currentShirtNum = currentShirtNum

    def getCurrentShirtNum(self):
        return self.currentShirtNum

    def setLastSaleBlocknum(self, blocknum):
        self.lastSaleBlocknum = blocknum

    def getLastSaleBlocknum(self):
        return self.lastSaleBlocknum



# teamIdx = 0 is the null team
class Team():
    def __init__(self, name):
        self.name = name
        self.playerIdxs             = np.zeros(NPLAYERS_PER_TEAM, int)
        self.currentLeagueIdx       = 0
        self.teamPosInCurrentLeague = 0
        self.prevLeagueIdx          = 0
        self.teamPosInPrevLeague    = 0


class League():
    def __init__(self, blockInit, blockStep, usersInitData):
        nTeams = len(usersInitData["teamIdxs"]) if blockInit != 0 else 0
        nMatches = nTeams*(nTeams-1)
        self.nTeams             = nTeams
        self.blockInit          = blockInit
        self.blockStep          = blockStep
        self.usersInitDataHash  = pylio.serialHash(usersInitData)
        self.usersAlongDataHash = 0
        # provided in update/challenge game
        self.initStatesHash     = 0
        self.finalStatesHashes  = 0
        self.scores             = np.zeros(nMatches)
        self.updaterAddr        = 0
        self.blockLastUpdate    = 0

    def isGenesisLeague(self):
        return self.blockInit == 0

    def isLeagueIsAboutToStart(self, blocknum):
        return blocknum < self.blockInit

    def hasLeagueStarted(self, blocknum):
        return blocknum >= self.blockInit

    def hasLeagueFinished(self, blocknum):
        nMatchdays = 2 * (self.nTeams-1)
        return blocknum >= self.blockInit + (nMatchdays-1) * self.blockStep

    def hasLeagueBeenUpdated(self):
        return self.blockLastUpdate != 0

    def resetUpdater(self):
        self.blockLastUpdate = 0

    def isFullyVerified(self, blocknum):
        if self.isGenesisLeague():
            return True
        return self.hasLeagueBeenUpdated() and (blocknum > self.blockLastUpdate + CHALLENGING_PERIOD_BLKS)

    def updateUsersAlongDataHash(self, usersAlongData):
        self.usersAlongDataHash = pylio.intHash(
            str(self.usersAlongDataHash) +
            serialize(usersAlongData)
        )

    def updateLeague(self, initStatesHash, finalStatesHashes, scores, updaterAddr, blocknum):
        assert self.hasLeagueFinished(blocknum), "League cannot be updated before the last matchday finishes"
        assert not self.hasLeagueBeenUpdated(), "League has already been updated"
        self.initStatesHash     = initStatesHash
        self.finalStatesHashes  = finalStatesHashes
        self.scores             = scores
        self.updaterAddr        = updaterAddr
        self.blockLastUpdate    = blocknum

    def challengeFinalStates(self, selectedTeam, initPlayerStates, usersInitData, usersAlongData, blocknum):
        assert self.hasLeagueBeenUpdated(), "League has not been updated yet, no need to challenge"
        assert not self.isFullyVerified(blocknum), "You cannot challenge after the challenging period"
        assert pylio.serialHash(usersInitData) == self.usersInitDataHash, "Incorrect provided: usersInitData"
        assert pylio.computeUsersAlongDataHash(usersAlongData) == self.usersAlongDataHash, "Incorrect provided: usersAlongData"
        assert pylio.serialHash(initPlayerStates) == self.initStatesHash, "Incorrect provided: initPlayerStates"
        selectedTeamFinalState, selectedMatchInMatchday, selectedScores  = pylio.computeTeamFinalState(selectedTeam, self.blockInit, self.blockStep, initPlayerStates, usersInitData, usersAlongData)

        if not pylio.serialHash(selectedTeamFinalState) == self.finalStatesHashes[selectedTeam]:
            print "Challenger Wins: finalStates provided by updater are invalid"
            self.resetUpdater()
            return

        if not pylio.areUpdaterScoresCorrect(selectedMatchInMatchday, selectedScores, self.scores):
            print "Challenger Wins: scores provided by updater are invalid"
            self.resetUpdater()
            return

        print "Challenger failed to prove that finalStates nor scores were wrong"

    def challengeInitStates(self, usersInitData, dataToChallengeInitStates, ST, blocknum):
        assert self.hasLeagueBeenUpdated(), "League has not been updated yet, no need to challenge"
        assert not self.isFullyVerified(blocknum), "You cannot challenge after the challenging period"
        assert pylio.serialHash(usersInitData) == self.usersInitDataHash, "Incorrect provided: usersInitData"

        nTeams = len(usersInitData["teamIdxs"])
        # dimensions: [team, nPlayersInTeam]
        #   if that a given player is virtual, then it contains just its state
        #   if not, it contains all states of prev league's team
        initPlayerStates = [[None for playerPosInLeague in range(NPLAYERS_PER_TEAM)] for team in range(nTeams)]
        for teamPos, teamIdx in enumerate(usersInitData["teamIdxs"]):
            for playerShirt, playerIdx in enumerate(ST.teams[teamIdx].playerIdxs):
                isOK = pylio.isCorrectStateForPlayerIdx(
                    pylio.getPlayerStateFromChallengeData(playerIdx, dataToChallengeInitStates[teamPos][playerShirt]),
                    dataToChallengeInitStates[teamPos][playerShirt],
                    ST
                )
                if isOK:
                    initPlayerStates[teamPos][playerShirt] = pylio.getPlayerStateFromChallengeData(
                        playerIdx,
                        dataToChallengeInitStates[teamPos][playerShirt]
                    )
                else:
                    print "Challenger Wins: initStates provided by updater are invalid"
                    self.resetUpdater()
                    return

        # TODO: check that the provided state proofs contain the actual player idx!!!!!

        if pylio.serialHash(initPlayerStates) == self.initStatesHash:
            print "Challenger failed to prove that initStates were wrong"
        else:
            print "Challenger Wins: initStates provided by updater are invalid"
            self.resetUpdater()


# client leagues inherit from leagues, and extend to include the data pre-hash
class LeagueClient(League):
    def __init__(self, blockInit, blockStep, usersInitData):
        League.__init__(self, blockInit, blockStep, usersInitData)
        self.usersInitData = usersInitData
        self.usersAlongData = []  # this list must be ordered!
        self.initPlayerStates = None
        self.finalStates = None
        self.scores      = None

    def updateUsersAlongData(self, usersAlongData):
        self.updateUsersAlongDataHash(usersAlongData)
        self.usersAlongData.append(usersAlongData)

    def updateFinalState(self, finalStates, scores):
        self.finalStates = finalStates
        self.scores      = scores

    def updateInitState(self, initPlayerStates):
        self.initPlayerStates = initPlayerStates


