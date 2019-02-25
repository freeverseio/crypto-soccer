import numpy as np
import copy
import datetime
from constants import *
import pylio
from pickle import dumps as serialize



# simple block counter simulator, where the blockhash is just the hash of the blocknumber
class Counter():
    def __init__(self):
        self.currentBlock = 0
        # We start with one commit, a dummy one, so we start with currentVerse = 1.
        self.currentVerse = 1

    def advanceNBlocks(self, deltaN):
        self.advanceToBlock(self.currentBlock + deltaN)

    def advanceToBlock(self, n):
        assert n > self.currentBlock, "Cannot advance... to a block in the past!"
        if self.currentBlock < self.nextVerseBlock() <= n:
            self.advanceNVerses(1)
        self.currentBlock = n


    def advanceNVerses(self, n):
        self.currentVerse += n

    def advanceToVerse(self, n):
        assert n >= self.currentVerse, "Cannot advance... to a verse in the past!"
        self.currentVerse = n



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
    def __init__(self, verseInit, verseStep, usersInitData):
        nTeams = len(usersInitData["teamIdxs"]) if verseInit != 0 else 0
        nMatches = nTeams*(nTeams-1)
        self.nTeams             = nTeams
        self.verseInit          = verseInit
        self.verseStep          = verseStep
        self.usersInitDataHash  = pylio.serialHash(usersInitData)
        # provided in update/challenge game
        self.initStatesHash     = 0
        self.statesAtMatchdayHashes  = 0
        self.scores             = np.zeros(nMatches)
        self.updaterAddr        = 0
        self.blockLastUpdate    = 0

    def isGenesisLeague(self):
        return self.verseInit == 0

    def isLeagueIsAboutToStart(self, verse):
        return verse < self.verseInit

    def hasLeagueStarted(self, verse):
        return verse >= self.verseInit

    def verseFinal(self):
        nMatchdays = 2 * (self.nTeams - 1)
        return self.verseInit + (nMatchdays-1)*self.verseStep

    def hasLeagueFinished(self, verse):
        return verse >= self.verseFinal()

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
            pylio.serialize2str(usersAlongData)
        )

    def updateLeague(self, initStatesHash, statesAtMatchdayHashes, scores, updaterAddr, blocknum):
        assert self.hasLeagueFinished(blocknum), "League cannot be updated before the last matchday finishes"
        assert not self.hasLeagueBeenUpdated(), "League has already been updated"
        self.initStatesHash             = initStatesHash
        self.statesAtMatchdayHashes     = statesAtMatchdayHashes
        self.scores                     = scores
        self.updaterAddr                = updaterAddr
        self.blockLastUpdate            = blocknum

    def challengeMatchdayStates(self, selectedMatchday, prevMatchdayStates, usersInitData, usersAlongData, blocknum):
        assert self.hasLeagueBeenUpdated(), "League has not been updated yet, no need to challenge"
        assert not self.isFullyVerified(blocknum), "You cannot challenge after the challenging period"
        assert pylio.serialHash(usersInitData) == self.usersInitDataHash, "Incorrect provided: usersInitData"
        assert pylio.computeUsersAlongDataHash(usersAlongData) == self.usersAlongDataHash, "Incorrect provided: usersAlongData"
        if selectedMatchday == 0:
            assert pylio.serialHash(prevMatchdayStates) == self.initStatesHash, "Incorrect provided: prevMatchdayStates"
        else:
            assert pylio.serialHash(prevMatchdayStates) == self.statesAtMatchdayHashes[selectedMatchday-1], "Incorrect provided: prevMatchdayStates"

        matchdayBlock = self.blockInit + selectedMatchday * self.blockStep
        tactics = pylio.duplicate(usersInitData["tactics"])
        pylio.updateTacticsToBlockNum(tactics, matchdayBlock, usersAlongData)

        statesAtMatchday, scores = pylio.computeStatesAtMatchday(
            selectedMatchday,
            prevMatchdayStates,
            tactics,
            matchdayBlock
        )

        if not pylio.serialHash(statesAtMatchday) == self.statesAtMatchdayHashes[selectedMatchday]:
            print("Challenger Wins: statesAtMatchday provided by updater are invalid")
            self.resetUpdater()
            return

        if not (self.scores[selectedMatchday] == scores).all():
            print("Challenger Wins: scores provided by updater are invalid")
            self.resetUpdater()
            return

        print("Challenger failed to prove that statesAtMatchday nor scores were wrong")

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
                    print("Challenger Wins: initStates provided by updater are invalid")
                    self.resetUpdater()
                    return

        # TODO: check that the provided state proofs contain the actual player idx!!!!!

        if pylio.serialHash(initPlayerStates) == self.initStatesHash:
            print("Challenger failed to prove that initStates were wrong")
        else:
            print("Challenger Wins: initStates provided by updater are invalid")
            self.resetUpdater()


# client leagues inherit from leagues, and extend to include the data pre-hash
class LeagueClient(League):
    def __init__(self, verseInit, verseStep, usersInitData):
        League.__init__(self, verseInit, verseStep, usersInitData)
        self.usersInitData      = usersInitData
        self.initPlayerStates   = None
        self.statesAtMatchday   = None
        self.scores             = None

    def updateStatesAtMatchday(self, statesAtMatchday, scores):
        self.statesAtMatchday   = statesAtMatchday
        self.scores             = scores

    def updateInitState(self, initPlayerStates):
        self.initPlayerStates = initPlayerStates

class VerseCommit:
    def __init__(self, actionsHashes = 0, blockNum = 0, blockHash = 0):
        self.actionsHashes = actionsHashes
        self.blockNum = blockNum
        self.blockHash = blockHash


class VerseCommitClient(VerseCommit):
    def __init__(self):
        VerseCommit.__init__(self)
        self.actions = 0


class ActionsAccumulator():
    def __init__(self):
        self.buffer                     = {}
        self.commitedActions            = [0]

    def accumulateAction(self, action, blocknum):
        if blocknum in self.buffer:
            self.buffer[blocknum].append(action)
        else:
            self.buffer[blocknum] = [action]

    def removeActionsBeforeBlockNumFromBuffer(self, blockNum):
        allActions = pylio.duplicate(self.buffer)
        for (block, actions) in self.buffer.items():
            if block < blockNum:
                del allActions[block]
        self.buffer = allActions



class Storage(Counter):
    def __init__(self):

        Counter.__init__(self)

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


        self.blocksBetweenVerses = 360
        self.VerseCommits = [VerseCommit()]

    def lastVerseBlock(self):
        return self.VerseCommits[-1].blockNum

    def nextVerseBlock(self):
        return self.lastVerseBlock() + self.blocksBetweenVerses

    def commit(self, actionsHash, commitBlockNum, commitBlockHash, actionsPrehash = None):
        self.VerseCommits.append(VerseCommit(actionsHash, commitBlockNum, commitBlockHash))

    def addAccumulator(self):
        self.Accumulator = ActionsAccumulator()


    def accumulateAction(self, action):
        assert self.currentBlock >= self.lastVerseBlock(), "Weird, blocknum for action received that belonged to past commit"
        self.Accumulator.accumulateAction(action, self.currentBlock)

    def getAllActionsBeforeBlock(self, blockNum):
        actions2commit = []
        for (block, actions) in self.Accumulator.buffer.items():
            if block < blockNum:
                for a in actions:
                    actions2commit.append(a)
        return actions2commit

    def getAllActionsForLeague(self, leagueIdx):
        nMatchdays = 2*(self.leagues[leagueIdx].nTeams-1)
        verseInit = self.leagues[leagueIdx].verseInit
        verseEnd = self.leagues[leagueIdx].verseFinal()
        verseStep = self.leagues[leagueIdx].verseStep
        versesThisLeague = range(verseInit, verseEnd, self.leagues[leagueIdx])

        for matchday in range(nMatchdays):
            verse = self.leagues[leagueIdx].verseInit
            self.Accumulator.commitedActions

    def syncActions(self, ST):
        assert self.currentBlock == ST.currentBlock, "Client and BC are out of sync in blocknum!"
        if self.currentBlock >= self.nextVerseBlock():
            actions2commit = self.getAllActionsBeforeBlock(self.nextVerseBlock())
            # leaguesPlayingInThisVerse
            ST.commit(
                pylio.serialize2str(actions2commit),
                self.nextVerseBlock(),
                pylio.getBlockhashForBlock(self.nextVerseBlock())
            )
            self.commit(
                pylio.serialize2str(actions2commit),
                self.nextVerseBlock(),
                pylio.getBlockhashForBlock(self.nextVerseBlock()),
            )
            self.Accumulator.commitedActions.append(actions2commit)
            self.Accumulator.removeActionsBeforeBlockNumFromBuffer(self.nextVerseBlock())