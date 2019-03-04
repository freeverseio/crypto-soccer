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
        self.actionsPerMatchday = []

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

    def getActionsForLeagueAndVerse(self, leagueIdx, verse):
        # each action has the form [leagueIdx, actions]
        # this function assumes that all actions for a given leagueIdx are in one single leagueIdx entry
        actionsInThisVerse = pylio.duplicate(self.Accumulator.commitedActions[verse])
        actions = [a for a in actionsInThisVerse if a[0] == leagueIdx]
        assert len(actions)<=1, "Actions for a league should be packed in one single entry"
        return actions


    def getAllActionsForLeague(self, leagueIdx):
        assert self.leagues[leagueIdx].hasLeagueFinished(self.currentVerse), "All actions only available at end of league"
        # actionsPerVerse will have form: verse x nLeagues, the latter in form: [ [ leagueIdx, actions], ...]
        actionsPerVerse = []
        for verse in self.getVersesForLeague(leagueIdx):
            actions = self.getActionsForLeagueAndVerse(leagueIdx, verse)
            actionsPerVerse.append(actions)
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
        leagueIdxAndActionsArray = []
        for leagueIdx in leaguesPlayingInThisVerse:
            if leagueIdx in self.Accumulator.buffer:
                leagueIdxAndActionsArray.append([leagueIdx, self.Accumulator.buffer[leagueIdx]])
                self.leagues[leagueIdx].actionsPerMatchday.append(self.Accumulator.buffer[leagueIdx])

        if leagueIdxAndActionsArray:
            tree, depth = make_tree(leagueIdxAndActionsArray, pylio.serialHash)
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
        self.Accumulator.commitedActions.append(leagueIdxAndActionsArray)
        self.Accumulator.commitedTrees.append(tree)
        self.Accumulator.clearBuffer(leagueIdxAndActionsArray)

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
            values,
            depth
                                ):
        verse = self.leagues[leagueIdx].verseInit + selectedMatchday * self.leagues[leagueIdx].verseStep
        seed  = pylio.getBlockHash(self.VerseCommits[verse].blockNum)
        # TODO: looks like if actions is empty, it does not know how to compare merkle
        assert verify(
            self.VerseCommits[verse].actionsMerkleRoots,
            depth,
            values,
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


    def getMerkleProof(self, leagueIdx, selectedMatchday, actionsAtSelectedMatchday):
        verse = self.leagues[leagueIdx].verseInit + selectedMatchday * self.leagues[leagueIdx].verseStep
        for idx, action in enumerate(self.Accumulator.commitedActions[verse]):
            if action[0] == leagueIdx:
                break
        tree = self.Accumulator.commitedTrees[verse]
        # get the needed hashes and the "values". The latter are simply the corresponding
        # leaf (=actionsAtSelectedMatchday) formated so that is has the form {idx: actionsAtSelectedMatchday}.
        neededHashes, values = pylio.prepareProofForIdxs([idx], tree, actionsAtSelectedMatchday)

        # verify(self.VerseCommits[verse].actionsMerkleRoots, depth, values, merkleProof, pylio.serialHash)
        assert verify(root(tree), get_depth(tree), values, neededHashes, pylio.serialHash), "Generated Merkle proof will not work"

        return neededHashes, values, get_depth(tree)

    def getPlayerIdxFromTeamIdxAndShirt(self, teamIdx, shirtNum):
        # If player has never been sold (virtual team): simple relation between playerIdx and (teamIdx, shirtNum)
        # Otherwise, read what's written in the playerState
        # playerIdx = 0 andt teamdIdx = 0 are the null player and teams
            self.assertTeamIdx(teamIdx)
            isPlayerIdxAssigned = self.teams[teamIdx].playerIdxs[shirtNum] != 0
            if isPlayerIdxAssigned:
                return self.teams[teamIdx].playerIdxs[shirtNum]
            else:
                return 1 + (teamIdx - 1) * NPLAYERS_PER_TEAM + shirtNum

    def assertTeamIdx(self, teamIdx):
        assert teamIdx < len(self.teams), "Team for this playerIdx not created yet!"
        assert teamIdx != 0, "Team 0 is reserved for null team!"

    def getLastWrittenPlayerStateFromPlayerIdx(self, playerIdx):
        prevLeagueIdx, teamPosInPrevLeague = self.getLastPlayedLeagueIdx(playerIdx)
        if prevLeagueIdx == 0:
            # this can be known both by CLIENT and BC
            return self.getPlayerStateBeforePlayingAnyLeague(playerIdx)
        else:
            # this can only be accessed by the CLIENT
            return self.getPlayerStateAtEndOfLeague(prevLeagueIdx, teamPosInPrevLeague, playerIdx)

    def getPlayerStateBeforePlayingAnyLeague(self, playerIdx):
        # this can be called by BC or CLIENT, as both have enough data
        playerStateAtBirth = self.getPlayerStateAtBirth(playerIdx)

        if self.isPlayerVirtual(playerIdx):
            return playerStateAtBirth
        else:
            # if player has been sold before playing any league, it'll conserve skills at birth,
            # but have different metadata in the other fields
            playerState = duplicate(self.playerIdxToPlayerState[playerIdx])
            copySkillsAndAgeFromTo(playerStateAtBirth, playerState)
            return playerState


    def getPlayerStateAtBirth(self, playerIdx):
        # Disregard his current team, just look at the team at moment of birth to build skills
        teamIdx, shirtNum = self.getTeamIdxAndShirtForPlayerIdx(playerIdx, forceAtBirth=True)
        seed = pylio.getPlayerSeedFromTeamAndShirtNum(self.teams[teamIdx].name, shirtNum)
        playerState = pylio.duplicate(pylio.getPlayerStateFromSeed(seed))
        # Once the skills have been added, complete the rest of the player data
        playerState.setPlayerIdx(playerIdx)
        playerState.setCurrentTeamIdx(teamIdx)
        playerState.setCurrentShirtNum(shirtNum)
        return playerState

    # The inverse of the previous relation
    def getTeamIdxAndShirtForPlayerIdx(self, playerIdx, forceAtBirth=False):
        if forceAtBirth or self.isPlayerVirtual(playerIdx):
            teamIdx = 1 + (playerIdx - 1) // NPLAYERS_PER_TEAM
            shirtNum = (playerIdx - 1) % NPLAYERS_PER_TEAM
            return teamIdx, shirtNum
        else:
            return self.playerIdxToPlayerState[playerIdx].getCurrentTeamIdx(), \
                   self.playerIdxToPlayerState[playerIdx].getCurrentShirtNum()

    # if player has never been sold, it will not be in the map playerIdxToPlayerState
    # and his team is derived from a formula
    def isPlayerVirtual(self, playerIdx):
        return not playerIdx in self.playerIdxToPlayerState


    def getPlayerStateAtEndOfLeague(self, prevLeagueIdx, teamPosInPrevLeague, playerIdx):
        selectedStates = [s for s in self.leagues[prevLeagueIdx].statesAtMatchday[-1][teamPosInPrevLeague] if
                          s.getPlayerIdx() == playerIdx]
        assert len(
            selectedStates) == 1, "PlayerIdx not found in previous league final states, or too many with same playerIdx"
        return selectedStates[0]

    def getLastPlayedLeagueIdx(self, playerIdx):
        # if player state has never been written, it played all leagues with current team (obtained from formula)
        # otherwise, we check if it was sold to current team before start of team's previous league
        if self.isPlayerVirtual(playerIdx):
            teamIdx, shirtNum = self.getTeamIdxAndShirtForPlayerIdx(playerIdx)
            return self.teams[teamIdx].prevLeagueIdx, self.teams[teamIdx].teamPosInPrevLeague

        currentTeamIdx = self.playerIdxToPlayerState[playerIdx].getCurrentTeamIdx()
        prevLeagueIdxForCurrentTeam = self.teams[currentTeamIdx].prevLeagueIdx
        didHePlayLastLeagueWithCurrentTeam = self.playerIdxToPlayerState[playerIdx].getLastSaleBlocknum() < \
                                             self.leagues[prevLeagueIdxForCurrentTeam].blockInit
        if didHePlayLastLeagueWithCurrentTeam:
            return prevLeagueIdxForCurrentTeam, self.teams[currentTeamIdx].teamPosInPrevLeague
        else:
            return self.playerIdxToPlayerState[playerIdx].prevLeagueIdx, self.playerIdxToPlayerState[
                playerIdx].prevTeamPosInLeague


    def getInitPlayerStates(self, leagueIdx):
        usersInitData = pylio.duplicate(self.leagues[leagueIdx].usersInitData)
        nTeams = len(usersInitData["teamIdxs"])
        # an array of size [nTeams][NPLAYERS_PER_TEAM]
        initPlayerStates = [[None for playerPosInLeague in range(NPLAYERS_PER_TEAM)] for team in range(nTeams)]
        teamPosInLeague = 0
        for teamIdx, teamOrder in zip(usersInitData["teamIdxs"], usersInitData["teamOrders"]):
            for shirtNum, playerPosInLeague in enumerate(teamOrder):
                playerIdx = self.getPlayerIdxFromTeamIdxAndShirt(teamIdx, shirtNum)
                playerState = self.getLastWrittenPlayerStateFromPlayerIdx(playerIdx)
                initPlayerStates[teamPosInLeague][playerPosInLeague] = playerState
            teamPosInLeague += 1
        return initPlayerStates


    def computeAllMatchdayStates(self, leagueIdx):
        initPlayerStates = self.getInitPlayerStates(leagueIdx)
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

        # the following beast has dimension nMatchdays x nTeams x nPlayersPerTeam
        statesAtMatchday = [pylio.createEmptyPlayerStatesForAllTeams(nTeams) for matchday in range(nMatchdays)]
        tacticsAtMatchDay = []
        teamOrdersAtMatchDay = []

        for matchday in range(nMatchdays):
            self.updateTacticsToVerseNum(leagueIdx, tactics, teamOrders, matchday)
            prevStates = initPlayerStates if matchday == 0 else statesAtMatchday[matchday - 1]
            statesAtMatchday[matchday], scores[matchday] = pylio.computeStatesAtMatchday(
                matchday,
                prevStates,
                tactics,
                teamOrders,
                seedsPerVerse[matchday]
            )
            tacticsAtMatchDay.append(pylio.duplicate(tactics))
            teamOrdersAtMatchDay.append(pylio.duplicate(teamOrders))

        return statesAtMatchday, tacticsAtMatchDay, teamOrdersAtMatchDay, scores

    def updateTacticsToVerseNum(self, leagueIdx, tactics, teamOrders, matchday):
        actionsInThisMatchday = pylio.duplicate(self.leagues[leagueIdx].actionsPerMatchday[matchday])
        for action in actionsInThisMatchday:
            teamPosInLeauge = self.getTeamPosInLeague(action["teamIdx"], leagueIdx)
            tactics[teamPosInLeauge] = action["tactics"]
            teamOrders[teamPosInLeauge] = action["teamOrder"]

    def getTeamPosInLeague(self, teamIdx, leagueIdx):
        for tPos, tIdx in enumerate(self.leagues[leagueIdx].usersInitData["teamIdxs"]):
            if teamIdx == tIdx:
                return tPos
        assert False, "Team not found in league"

