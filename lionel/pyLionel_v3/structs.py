import numpy as np
import copy
import datetime
from constants import *
import pylio
from pickle import dumps as serialize
from merkle_tree import *


# simple block counter simulator, where the blockhash is just the hash of the blocknumber
class Counter():
    def __init__(self):
        self.currentBlock = 0
        self.currentVerse = 0

    def advanceNBlocks(self, deltaN):
        self.advanceToBlock(self.currentBlock + deltaN)

    def advanceToBlock(self, n):
        assert n > self.currentBlock, "Cannot advance... to a block in the past!"
        verseWasCrossed = False
        if self.currentBlock < self.nextVerseBlock() <= n:
            self.advanceNVerses(1)
            verseWasCrossed = True
        self.currentBlock = n
        return verseWasCrossed


    def advanceNVerses(self, n):
        self.currentVerse += n
        # TONI self.lastBlock

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
        self.dataAtMatchdayHashes = 0
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

    def updateLeague(self, initStatesHash, dataAtMatchdayHashes, scores, updaterAddr, blocknum, verse):
        assert self.hasLeagueFinished(verse), "League cannot be updated before the last matchday finishes"
        assert not self.hasLeagueBeenUpdated(), "League has already been updated"
        self.initStatesHash             = initStatesHash
        self.dataAtMatchdayHashes     = dataAtMatchdayHashes
        self.scores                     = scores
        self.updaterAddr                = updaterAddr
        self.blockLastUpdate            = blocknum

    def challengeMatchdayStates(self,
                                selectedMatchday,
                                prevMatchdayStates,
                                prevMatchdayTactics,
                                prevMatchdayTeamOrders,
                                usersInitData,
                                actionsAtSelectedMatchday,
                                seed,
                                currentBlocknum):

        assert self.hasLeagueBeenUpdated(), "League has not been updated yet, no need to challenge"
        assert not self.isFullyVerified(currentBlocknum), "You cannot challenge after the challenging period"
        assert pylio.serialHash(usersInitData) == self.usersInitDataHash, "Incorrect provided: usersInitData"
        if selectedMatchday == 0:
            assert pylio.serialHash(prevMatchdayStates) == self.initStatesHash, "Incorrect provided: prevMatchdayStates"
            # TODO: the next two are a bit useless, rethink
            assert pylio.serialHash(prevMatchdayTactics) == pylio.serialHash(usersInitData["tactics"]), "Incorrect provided: prevMatchdayStates"
            assert pylio.serialHash(prevMatchdayTeamOrders) == pylio.serialHash(usersInitData["teamOrders"]), "Incorrect provided: prevMatchdayStates"
        else:
            # TODO: sum of hashes is not secure, hash the result!
            assert self.dataAtMatchdayHashes[selectedMatchday-1] == \
                pylio.serialHash(prevMatchdayStates) + \
                pylio.serialHash(prevMatchdayTactics) +\
                pylio.serialHash(prevMatchdayTeamOrders), \
                "Incorrect provided: prevMatchdayStates"

        # assert pylio.computeUsersAlongDataHash(usersAlongData) == self.usersAlongDataHash, "Incorrect provided: usersAlongData"

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
        self.tacticsAtMatchDay  = None
        self.scores             = None

    def updateDataAtMatchday(self, statesAtMatchday, tacticsAtMatchDay, teamOrdersAtMatchDay, scores):
        self.statesAtMatchday   = statesAtMatchday
        self.tacticsAtMatchDay  = tacticsAtMatchDay
        self.teamOrdersAtMatchDay  = teamOrdersAtMatchDay
        self.scores             = scores

    def updateInitState(self, initPlayerStates):
        self.initPlayerStates = initPlayerStates

class VerseCommit:
    def __init__(self, actionsMerkleRoots = 0, blockNum = 0):
        self.actionsMerkleRoots = actionsMerkleRoots
        self.blockNum = blockNum


class VerseCommitClient(VerseCommit):
    def __init__(self):
        VerseCommit.__init__(self)
        self.actions = 0


class ActionsAccumulator():
    def __init__(self):
        self.buffer                     = {}
        self.commitedActions            = [0]
        self.commitedTrees              = [0]

    def accumulateAction(self, action, leagueIdx):
        if leagueIdx in self.buffer:
            self.buffer[leagueIdx].append(action)
        else:
            self.buffer[leagueIdx] = [action]

    def clearBuffer(self, actions2remove):
        for action in actions2remove:
            leagueIdx = action[0]
            del self.buffer[leagueIdx]



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

    def commit(self, actionsHash, commitBlockNum, actionsPrehash = None):
        self.VerseCommits.append(VerseCommit(actionsHash, commitBlockNum))

    def addAccumulator(self):
        self.Accumulator = ActionsAccumulator()


    def accumulateAction(self, action):
        assert self.currentBlock >= self.lastVerseBlock(), "Weird, blocknum for action received that belonged to past commit"
        self.Accumulator.accumulateAction(action, self.getLeagueForAction(action))

    def getAllActionsBeforeBlock(self, blockNum):
        actions2commit = []
        for (block, actions) in self.Accumulator.buffer.items():
            if block < blockNum:
                for a in actions:
                    actions2commit.append(a)
        return actions2commit

    def getVersesForLeague(self, leagueIdx):
        nMatchdays = 2*(self.leagues[leagueIdx].nTeams-1)
        verses = []
        for matchday in range(nMatchdays):
            verses.append(self.leagues[leagueIdx].verseInit + matchday * self.leagues[leagueIdx].verseStep)
        return verses

    def getSeedForVerse(self, verse):
        return pylio.getBlockHash(self.VerseCommits[verse].blockNum)

    def getAllSeedsForLeague(self, leagueIdx):
        assert self.leagues[leagueIdx].hasLeagueFinished(self.currentVerse), "All seeds only available at end of league"
        seedsPerVerse = []
        for verse in self.getVersesForLeague(leagueIdx):
            seedsPerVerse.append(self.getSeedForVerse(verse))
        return seedsPerVerse

    def getAllActionsForLeague(self, leagueIdx):
        assert self.leagues[leagueIdx].hasLeagueFinished(self.currentVerse), "All actions only available at end of league"
        actionsPerVerse = []
        for verse in self.getVersesForLeague(leagueIdx):
            actionsInThisVerse = pylio.duplicate(self.Accumulator.commitedActions[verse])
            if leagueIdx in actionsInThisVerse:
                actionsInThisVerse = actionsInThisVerse[leagueIdx]
            else:
                actionsInThisVerse = []
            # Convert teamIdx -> teamPos
            for a in actionsInThisVerse:
                teamPosInLeague = pylio.getTeamPosInLeague(a["teamIdx"], self.leagues[leagueIdx] )
                a["teamIdx"] = teamPosInLeague
            actionsPerVerse.append(actionsInThisVerse)
        return actionsPerVerse

    def getLeagueForAction(self, action):
        return self.teams[action["teamIdx"]].currentLeagueIdx


    def getLeaguesPlayingInThisVerse(self, verse):
        # TODO: make this less terribly slow
        leagueIdxs = []
        nLeagues = len(self.leagues)
        for leagueIdx  in range(1,nLeagues): # bypass the first (dummy) league
            if verse in self.getVersesForLeague(leagueIdx):
                leagueIdxs.append(leagueIdx)
        return leagueIdxs

    def syncActions(self, ST):
        assert self.currentBlock == ST.currentBlock, "Client and BC are out of sync in blocknum!"
        leaguesPlayingInThisVerse = self.getLeaguesPlayingInThisVerse(ST.currentVerse)
        leagueIdxVsActionsMatrix = []
        for leagueIdx in leaguesPlayingInThisVerse:
            if leagueIdx in self.Accumulator.buffer:
                leagueIdxVsActionsMatrix.append([leagueIdx, self.Accumulator.buffer[leagueIdx]])

        if leagueIdxVsActionsMatrix:
            tree, depth = make_tree(leagueIdxVsActionsMatrix, pylio.serialHash)
            rootTree    = root(tree)
        else:
            tree        = 0
            rootTree    = 0

        ST.commit(
            rootTree,
            self.nextVerseBlock(),
        )
        self.commit(
            rootTree,
            self.nextVerseBlock(),
        )
        self.Accumulator.commitedActions.append(leagueIdxVsActionsMatrix)
        self.Accumulator.commitedTrees.append(tree)
        self.Accumulator.clearBuffer(leagueIdxVsActionsMatrix)

    def updateLeague(self, leagueIdx, initStatesHash, dataAtMatchdayHashes, scores, updaterAddr):
        self.leagues[leagueIdx].updateLeague(
            initStatesHash,
            dataAtMatchdayHashes,
            scores,
            updaterAddr,
            self.currentBlock,
            self.currentVerse
        )


    def challengeMatchdayStates(self,
            leagueIdx,
            selectedMatchday,
            prevMatchdayStates,
            prevMatchdayTactics,
            prevMatchdayTeamOrders,
            usersInitData,
            actionsAtSelectedMatchday,
            merkleProof,
            depth
                                ):
        verse = self.leagues[leagueIdx].verseInit + selectedMatchday * self.leagues[leagueIdx].verseStep
        seed  = pylio.getBlockHash(self.VerseCommits[verse].blockNum)

        # TODO: looks like if actions is empty, it does not know how to compare merkle
        assert verify(
            self.VerseCommits[verse].actionsMerkleRoots,
            depth,
            actionsAtSelectedMatchday,
            merkleProof,
            pylio.serialHash,
            debug_print=False
        ), "Actions are not part of the corresponding commit"


        # TODO: verify that actionsAtSelectedMatchday are such that hash is correct
        assert self.actionsArePartOfCommit(leagueIdx, actionsAtSelectedMatchday, verse, merkleProof), "Actions are not part of a Verse Commit"


        self.leagues[leagueIdx].challengeMatchdayStates(
            selectedMatchday,
            prevMatchdayStates,
            prevMatchdayTactics,
            prevMatchdayTeamOrders,
            usersInitData,
            actionsAtSelectedMatchday,
            seed,
            self.currentBlock
        )


    def getMerkleProof(self, leagueIdx, selectedMatchday):
        verse = self.leagues[leagueIdx].verseInit + selectedMatchday * self.leagues[leagueIdx].verseStep
        for idx, action in enumerate(self.Accumulator.commitedActions[verse]):
            if action[0] == leagueIdx:
                break
        tree = self.Accumulator.commitedTrees[verse]
        return proof(tree, [idx]), get_depth(tree)


