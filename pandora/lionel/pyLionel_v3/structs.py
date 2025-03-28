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

    def advanceToVerse(self, n):
        assert n >= self.currentVerse, "Cannot advance... to a verse in the past!"
        self.currentVerse = n


class MinimalPlayerState():
    def __init__(self, playerState = None):
        if playerState:
            self.skills                  = playerState.skills
            self.monthOfBirthInUnixTime  = playerState.monthOfBirthInUnixTime
            self.playerIdx               = playerState.playerIdx
        else:
            self.skills                  = np.zeros(N_SKILLS)
            self.monthOfBirthInUnixTime  = 0
            self.playerIdx               = 0

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

# In Solidity, PlayerState will be just a uin256, serializing the data shown here,
# and there'll be associated read/write functions
# playerIdx = 0 is the null player
class PlayerState(MinimalPlayerState):
    def __init__(self):
        MinimalPlayerState.__init__(self)
        self.currentTeamIdx          = 0
        self.currentShirtNum         = 0
        self.prevLeagueIdx          = 0
        self.prevTeamPosInLeague    = 0
        self.prevShirtNumInLeague   = 0
        self.lastSaleBlocknum        = 0

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
    def __init__(self, name, nowInMonthsUnixTime):
        self.name = name
        self.monthOfTeamCreationInUnixTime = nowInMonthsUnixTime
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
        self.initSkillsHash     = 0
        self.dataAtMatchdayHashes = 0
        self.scores             = np.zeros(nMatches)
        self.updaterAddr        = 0
        self.blockLastUpdate    = 0

    def isGenesisLeague(self):
        return self.verseInit == 0

    def verseFinal(self):
        nMatchdays = 2 * (self.nTeams - 1)
        return self.verseInit + (nMatchdays-1)*self.verseStep

    def hasLeagueBeenUpdated(self):
        return self.blockLastUpdate != 0

    def resetUpdater(self):
        self.blockLastUpdate = 0


    def updateLeague(self, initSkillsHash, dataAtMatchdayHashes, scores, updaterAddr, blocknum):
        self.initSkillsHash             = initSkillsHash
        self.dataAtMatchdayHashes       = dataAtMatchdayHashes
        self.scores                     = scores
        self.updaterAddr                = updaterAddr
        self.blockLastUpdate            = blocknum


# client leagues inherit from leagues, and extend to include the data pre-hash
class LeagueClient(League):
    def __init__(self, verseInit, verseStep, usersInitData):
        League.__init__(self, verseInit, verseStep, usersInitData)
        self.usersInitData      = usersInitData
        self.initPlayerStates   = None
        self.skillsAtMatchday   = None
        self.lastDayTree        = None
        self.tacticsAtMatchday  = None
        self.scores             = None
        self.actionsPerMatchday = []
        self.dataToChallengeInitSkills = None

    def storeDataAtMatchdays(self, dataAtMatchdays, scores):
        self.dataAtMatchdays    = pylio.duplicate(dataAtMatchdays)
        self.scores             = pylio.duplicate(scores)

    def writeInitState(self, initPlayerStates):
        self.initPlayerStates = pylio.duplicate(initPlayerStates)

    def writeDataToChallengeInitSkills(self, dataToChallengeInitSkills):
        self.dataToChallengeInitSkills = pylio.duplicate(dataToChallengeInitSkills)

    def getInitPlayerSkills(self):
        initSkills = []
        for team in self.initPlayerStates:
            initSkills.append([pylio.duplicate(MinimalPlayerState(state)) for state in team])
        return initSkills


# The VerseCommit basically stores the merkle roots of all actions corresponding to a league starting at that moment
# The Merkle Roots are computed from the leafs:
#
#  leafs = [ [leagueIdx0, allActionsInLeagueIdx0], ..., ]
#
#  where allActionsInLeagueIdx0 = [action0, action1, action2,...]

class VerseCommit:
    def __init__(self, actionsMerkleRoots = 0, blockNum = 0):
        self.actionsMerkleRoots = actionsMerkleRoots
        self.blockNum = blockNum


class VerseCommitClient(VerseCommit):
    def __init__(self):
        VerseCommit.__init__(self)
        self.actions = 0


# The Accumulator is responsible for receving user actions and committing them in the correct verse.
# An action is a struct:
#    action00 = {"teamIdx": 2, "teamOrder": [0,4,2,3...,NPLAYERS_PER_TEAM], "tactics": 3}
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



# The MAIN CLASS that manages all BC & CLIENT storage
class Storage(Counter):
    def __init__(self, isClient):

        Counter.__init__(self)

        # this bool is just to understand if the created BC is actually a client
        # it allows to assert that some funcions should only be run by the client
        self.isClient = isClient

        # an array of Team structs, the first entry being the null team
        self.teams = [Team("",0)]

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

    def assertIsClient(self):
        assert self.isClient, "This code should only be run by CLIENTS, not the BC"

    # ------------------------------------------------------------------------
    # ----------      Functions common to both BC and CLIENT      ------------
    # ------------------------------------------------------------------------

    def lastVerseBlock(self):
        return self.VerseCommits[-1].blockNum

    def nextVerseBlock(self):
        return self.lastVerseBlock() + self.blocksBetweenVerses

    def commit(self, actionsRoot):
        self.VerseCommits.append(VerseCommit(actionsRoot, self.currentBlock))

    def updateLeague(self, leagueIdx, initSkillsHash, dataAtMatchdayHashes, scores, updaterAddr):
        assert self.hasLeagueFinished(leagueIdx), "League cannot be updated before the last matchday finishes"
        assert not self.leagues[leagueIdx].hasLeagueBeenUpdated(), "League has already been updated"
        self.leagues[leagueIdx].updateLeague(
            initSkillsHash,
            dataAtMatchdayHashes,
            scores,
            updaterAddr,
            self.currentBlock,
        )


    # note that values = actionsAtSelectedMatchday, formated so that is has the form
    # {idx: actionsAtSelectedMatchday}, where idx is the leaf idx.
    # so it should happen that both things coincide.

    def challengeMatchdayStates(self,
            leagueIdx,
            selectedMatchday,
            dataAtPrevMatchday,
            usersInitData,
            merkleProofDataForMatchday
        ):
        assert self.leagues[leagueIdx].hasLeagueBeenUpdated(), "League has not been updated yet, no need to challenge"
        assert not self.isFullyVerified(leagueIdx), "You cannot challenge after the challenging period"
        assert pylio.serialHash(usersInitData) == self.leagues[leagueIdx].usersInitDataHash, "Incorrect provided: usersInitData"
        assert merkleProofDataForMatchday.leaf[0] == leagueIdx, "Deverr: The actions do not belong to this league"
        verse = self.leagues[leagueIdx].verseInit + selectedMatchday * self.leagues[leagueIdx].verseStep

        # Validate that the provided actions where in the verse MerkleRoot
        assert pylio.verifyMerkleProof(
            self.VerseCommits[verse].actionsMerkleRoots,
            merkleProofDataForMatchday,
            pylio.serialHash,
        ), "Actions are not part of the corresponding commit"

        # Validate "dataAtPrevMatchday"
        # - if day =0, validate only that skills coincide with initSkillsHash,
        #              and initialize tactics and orders from usersInitData
        # - if day!=0, validate that the entire hash of dataAtPrevMatchday coincides with
        #               the hashes that the updater provided
        if selectedMatchday == 0:
            assert pylio.serialHash(dataAtPrevMatchday.skillsAtMatchday) == self.leagues[leagueIdx].initSkillsHash, "Incorrect provided: prevMatchdayStates"
            # initialize tactics and teams as written in league creation:
            assert dataAtPrevMatchday.tacticsAtMatchday == 0, "Incorrect provided: prevMatchdayStates"
            assert dataAtPrevMatchday.teamOrdersAtMatchday == 0, "Incorrect provided: prevMatchdayStates"
            dataAtPrevMatchday.tacticsAtMatchday = usersInitData["tactics"]
            dataAtPrevMatchday.teamOrdersAtMatchday = usersInitData["teamOrders"]
        else:
            assert self.leagues[leagueIdx].dataAtMatchdayHashes[selectedMatchday-1] == self.prepareOneMatchdayHash(dataAtPrevMatchday),\
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
            self.getSeedForVerse(verse)
        )

        dataAtMatchdayHash = self.prepareOneMatchdayHash(dataAtPrevMatchday)

        if not dataAtMatchdayHash == self.leagues[leagueIdx].dataAtMatchdayHashes[selectedMatchday]:
            print("Challenger Wins: skillsAtMatchday provided by updater are invalid")
            self.leagues[leagueIdx].resetUpdater()
            return

        if not (self.leagues[leagueIdx].scores[selectedMatchday] == scores).all():
            print("Challenger Wins: scores provided by updater are invalid")
            self.leagues[leagueIdx].resetUpdater()
            return

        print("Challenger failed to prove that skillsAtMatchday nor scores were wrong")

    def getPlayerIdxFromTeamIdxAndShirt(self, teamIdx, shirtNum):
        # If player has never been sold (virtual team): simple relation between playerIdx and (teamIdx, shirtNum)
        # Otherwise, read what's written in the playerState
        # playerIdx = 0 and teamdIdx = 0 are the null player and teams
            self.assertTeamIdx(teamIdx)
            isPlayerIdxAssigned = self.teams[teamIdx].playerIdxs[shirtNum] != 0
            if isPlayerIdxAssigned:
                return self.teams[teamIdx].playerIdxs[shirtNum]
            else:
                return 1 + (teamIdx - 1) * NPLAYERS_PER_TEAM + shirtNum

    def assertTeamIdx(self, teamIdx):
        assert teamIdx < len(self.teams), "Team for this playerIdx not created yet!"
        assert teamIdx != 0, "Team 0 is reserved for null team!"

    def getLastWrittenInBCPlayerStateFromPlayerIdx(self, playerIdx):
        if self.isPlayerVirtual(playerIdx):
            return self.getPlayerStateBeforePlayingAnyLeague(playerIdx)
        else:
            return self.playerIdxToPlayerState[playerIdx]

    def getPlayerStateBeforePlayingAnyLeague(self, playerIdx):
        # this can be called by BC or CLIENT, as both have enough data
        playerStateAtBirth = self.getPlayerStateAtBirth(playerIdx)

        if self.isPlayerVirtual(playerIdx):
            return playerStateAtBirth
        else:
            # if player has been sold before playing any league, it'll conserve skills at birth,
            # but have different metadata in the other fields
            playerState = pylio.duplicate(self.playerIdxToPlayerState[playerIdx])
            return self.copySkillsAndAgeFromTo(playerStateAtBirth, playerState)

    def copySkillsAndAgeFromTo(self, playerStateOrig, playerStateDest):
        newPlayerState = pylio.duplicate(playerStateDest)
        newPlayerState.setSkills(pylio.duplicate(playerStateOrig.getSkills()))
        newPlayerState.setMonth(pylio.duplicate(playerStateOrig.getMonth()))
        return newPlayerState

    # the skills of a player are determined by concat of teamName and shirtNum
    def getPlayerSeedFromTeamAndShirtNum(self, teamName, shirtNum):
        return pylio.limitSeed(pylio.intHash(teamName + str(shirtNum)))

    # Given a seed, returns a balanced player.
    # It only deals with skills & age, not playerIdx.
    def getPlayerStateFromSeed(self, seed, monthOfTeamCreationInUnixTime):
        newPlayerState = PlayerState()
        np.random.seed(seed)
        monthsWhenTeamWasCreated = np.random.randint(MIN_PLAYER_AGE, MAX_PLAYER_AGE) * 12
        newPlayerState.setMonth(monthOfTeamCreationInUnixTime-monthsWhenTeamWasCreated)
        skills = np.random.randint(0, AVG_SKILL - 1, N_SKILLS)
        excess = int((AVG_SKILL * N_SKILLS - skills.sum()) / N_SKILLS)
        skills += excess
        newPlayerState.setSkills(skills)
        return newPlayerState


    def getPlayerStateAtBirth(self, playerIdx):
        # Disregard his current team, just look at the team at moment of birth to build skills
        teamIdx, shirtNum = self.getTeamIdxAndShirtForPlayerIdx(playerIdx, forceAtBirth=True)
        seed = self.getPlayerSeedFromTeamAndShirtNum(self.teams[teamIdx].name, shirtNum)
        playerState = pylio.duplicate(self.getPlayerStateFromSeed(
            seed,
            self.teams[teamIdx].monthOfTeamCreationInUnixTime
        ))
        # Once the skills have been added, complete the rest of the player data
        playerState.setPlayerIdx(playerIdx)
        playerState.setCurrentTeamIdx(teamIdx)
        playerState.setCurrentShirtNum(shirtNum)
        return playerState

    # The inverse of the previous relation
    def getTeamIdxAndShirtForPlayerIdx(self, playerIdx, forceAtBirth=False):
        if forceAtBirth or self.isPlayerVirtual(playerIdx):
            teamIdx = int(1 + (playerIdx - 1) // NPLAYERS_PER_TEAM)
            shirtNum = int((playerIdx - 1) % NPLAYERS_PER_TEAM)
            return teamIdx, shirtNum
        else:
            return self.playerIdxToPlayerState[playerIdx].getCurrentTeamIdx(), \
                   self.playerIdxToPlayerState[playerIdx].getCurrentShirtNum()

    # if player has never been sold, it will not be in the map playerIdxToPlayerState
    # and his team is derived from a formula
    def isPlayerVirtual(self, playerIdx):
        return not playerIdx in self.playerIdxToPlayerState

    def verse2blockNum(self, verse):
        return self.VerseCommits[verse].blockNum

    def getLastPlayedLeagueIdx(self, playerIdx):
        # if player state has never been written, it played all leagues with current team (obtained from formula)
        # otherwise, we check if it was sold to current team before start of team's previous league
        if self.isPlayerVirtual(playerIdx):
            teamIdx, shirtNum = self.getTeamIdxAndShirtForPlayerIdx(playerIdx)
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
    def verifyInitSkillsMerkleProofs(self, usersInitData, dataToChallengeInitSkills):
        nTeams = len(usersInitData["teamIdxs"])
        # an array of size [nTeams][NPLAYERS_PER_TEAM]
        initPlayerSkills = pylio.createEmptyPlayerStatesForAllTeams(nTeams)
        for teamPosInLeague, teamIdx in enumerate(usersInitData["teamIdxs"]):
            for shirtNum in range(NPLAYERS_PER_TEAM):
                playerIdx = self.getPlayerIdxFromTeamIdxAndShirt(teamIdx, shirtNum)
                playerSkills = dataToChallengeInitSkills[teamPosInLeague][shirtNum].leaf
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

    def areLatestSkills(self, merkleProof):
        # If player has never played a league, we can compute the playerSkills directly in the BC
        # It basically is equal to the birth skills, with ,potentially, a few team changes via sales.
        # If not, we can just compare the hash of the dataToChallengePlayerState with the stored hash in the prev league
        playerSkills = merkleProof.leaf
        playerIdx = playerSkills.getPlayerIdx()
        prevLeagueIdx, teamPosInPrevLeague = self.getLastPlayedLeagueIdx(playerIdx)
        if prevLeagueIdx == 0:
            return pylio.areEqualStructs(
                playerSkills,
                MinimalPlayerState(self.getPlayerStateBeforePlayingAnyLeague(playerIdx))
            )
        else:
            return pylio.verifyMerkleProof(
                self.leagues[prevLeagueIdx].dataAtMatchdayHashes[-1],
                merkleProof,
                pylio.serialHash
            ), "Provided Merkle proof is invalid"


    def challengeInitSkills(self, leagueIdx, usersInitData, dataToChallengeInitSkills):
        assert self.leagues[leagueIdx].hasLeagueBeenUpdated(), "League has not been updated yet, no need to challenge"
        assert not self.isFullyVerified(leagueIdx), "You cannot challenge after the challenging period"
        assert pylio.serialHash(usersInitData) == self.leagues[leagueIdx].usersInitDataHash, "Incorrect provided: usersInitData"

        initSkills = self.verifyInitSkillsMerkleProofs(usersInitData, dataToChallengeInitSkills)
        # if None is returned, it means that at least one player had incorrect challenge data
        if not initSkills:
            print("Challenger Wins: initSkills provided by updater are invalid")
            self.leagues[leagueIdx].resetUpdater()
            return

        # We now know that the initSkills were correct. We just check that
        # the updater had not provided exactly the same correct skills!
        if pylio.serialHash(initSkills) == self.leagues[leagueIdx].initSkillsHash:
            print("Challenger failed to prove that initStates were wrong")
        else:
            print("Challenger Wins: initStates provided by updater are invalid")
            self.leagues[leagueIdx].resetUpdater()

    def getBlockNumForLastLeagueOfTeam(self, teamIdx):
        return self.verse2blockNum(self.leagues[self.teams[teamIdx].currentLeagueIdx].verseInit)




    # quick solution to simulate changing teams.
    # for the purpose of Lionel, we'll start with a simple exchange, instead
    # of the more convoluted sell, assign, etc.
    def exchangePlayers(self, playerIdx1, address1, playerIdx2, address2):
        assert not self.isPlayerBusy(playerIdx1), "Player sale failed: player is busy playing a league, wait until it finishes"
        assert not self.isPlayerBusy(playerIdx2), "Player sale failed: player is busy playing a league, wait until it finishes"

        teamIdx1, shirtNum1 = self.getTeamIdxAndShirtForPlayerIdx(playerIdx1)
        teamIdx2, shirtNum2 = self.getTeamIdxAndShirtForPlayerIdx(playerIdx2)

        # check ownership!
        assert self.getOwnerAddrFromTeamIdx(teamIdx1) == address1, "Exchange Failed, owner not correct"
        assert self.getOwnerAddrFromTeamIdx(teamIdx2) == address2, "Exchange Failed, owner not correct"

        # get states from BC in memory to do changes, and only write back once at the end
        state1 = pylio.duplicate(self.getLastWrittenInBCPlayerStateFromPlayerIdx(playerIdx1))
        state2 = pylio.duplicate(self.getLastWrittenInBCPlayerStateFromPlayerIdx(playerIdx2))

        # a player should change his prevLeagueIdx only if the current team played
        # a last league that started AFTER the last sale
        if self.getBlockNumForLastLeagueOfTeam(teamIdx1) > state1.getLastSaleBlocknum():
            state1.prevLeagueIdx = self.teams[teamIdx1].currentLeagueIdx
            state1.prevTeamPosInLeague = self.teams[teamIdx1].teamPosInCurrentLeague

        if self.getBlockNumForLastLeagueOfTeam(teamIdx2) > state2.getLastSaleBlocknum():
            state2.prevLeagueIdx = self.teams[teamIdx2].currentLeagueIdx
            state2.prevTeamPosInLeague = self.teams[teamIdx2].teamPosInCurrentLeague

        state1.setCurrentTeamIdx(teamIdx2)
        state2.setCurrentTeamIdx(teamIdx1)

        state1.setCurrentShirtNum(shirtNum2)
        state2.setCurrentShirtNum(shirtNum1)

        state1.setLastSaleBlocknum(self.currentBlock)
        state2.setLastSaleBlocknum(self.currentBlock)

        self.teams[teamIdx1].playerIdxs[shirtNum1] = playerIdx2
        self.teams[teamIdx2].playerIdxs[shirtNum2] = playerIdx1

        self.playerIdxToPlayerState[playerIdx1] = pylio.duplicate(state1)
        self.playerIdxToPlayerState[playerIdx2] = pylio.duplicate(state2)

    def isPlayerBusy(self, playerIdx1):
        return self.areTeamsBusyInPrevLeagues(
            [self.getTeamIdxAndShirtForPlayerIdx(playerIdx1)[0]])



    def areTeamsBusyInPrevLeagues(self, teamIdxs):
        for teamIdx in teamIdxs:
            if not self.isFullyVerified(self.teams[teamIdx].currentLeagueIdx):
                return True
        return False

    # Creates the league in the BC, storing only the hash of usersInitData
    # It signs teams in League, which allows the BC to now that they're busy
    # without 'seeing' the pre-hash usersInitData
    def createLeague(self, verseInit, verseStep, usersInitData):
        assert verseInit > self.currentVerse, "League cannot start in the past"
        assert not self.areTeamsBusyInPrevLeagues(usersInitData["teamIdxs"]), "League cannot create: some teams involved in prev leagues"
        assert len(usersInitData["teamIdxs"]) % 2 == 0, "Currently we only support leagues with even nTeams"
        leagueIdx = len(self.leagues)
        self.leagues.append(League(verseInit, verseStep, usersInitData))
        self.signTeamsInLeague(usersInitData["teamIdxs"], leagueIdx)
        return leagueIdx



    def signTeamsInLeague(self, teamIdxs, leagueIdx):
        for teamPosInLeague, teamIdx in enumerate(teamIdxs):
            self.teams[teamIdx].prevLeagueIdx             = pylio.duplicate(self.teams[teamIdx].currentLeagueIdx)
            self.teams[teamIdx].teamPosInPrevLeague       = pylio.duplicate(self.teams[teamIdx].teamPosInCurrentLeague)

            self.teams[teamIdx].currentLeagueIdx          = leagueIdx
            self.teams[teamIdx].teamPosInCurrentLeague    = teamPosInLeague


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
        return self.leagues[leagueIdx].hasLeagueBeenUpdated() and \
               (self.currentBlock > self.leagues[leagueIdx].blockLastUpdate + CHALLENGING_PERIOD_BLKS)

    def getPlayerSkillsFromChallengeData(self, playerIdx, dataToChallengePlayerState):
        # dataToChallengePlayerState is either:
        #  - just a player state
        #  - a merkle proof for that player
        # In the latter case, we can extract the state just from the values (leafs)
        if type(dataToChallengePlayerState) == type(MinimalPlayerState()):
            assert dataToChallengePlayerState.getPlayerIdx() == playerIdx, "This data does not contain the required playerIdx"
            return dataToChallengePlayerState
        else:
            playerSkills = dataToChallengePlayerState.leaf
            assert playerSkills.getPlayerIdx() == playerIdx, "This data does not contain the required player"
            return pylio.duplicate(playerSkills)

    def getOwnerAddrFromTeamIdx(self, teamIdx):
        return self.teamNameHashToOwnerAddr[pylio.intHash(self.teams[teamIdx].name)]

    def getOwnerAddrFromPlayerIdx(self, playerIdx):
        currentTeamIdx = self.getLastWrittenInBCPlayerStateFromPlayerIdx(playerIdx).currentTeamIdx
        return self.getOwnerAddrFromTeamIdx(currentTeamIdx)


    # A mockup of how to obtain the block hash for a given blocknum.
    # This is a function that is available in Ethereum after Constatinople
    def getBlockHash(self, blockNum):
        return pylio.intHash('salt' + str(blockNum))

    def getSeedForVerse(self, verse):
        return self.getBlockHash(self.VerseCommits[verse].blockNum)


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

    def certifyPlayerState(self, playerState, neededHashes, depth):
        # check that the skills inside playerState match the end of last league:
        playerSkills = MinimalPlayerState(playerState)
        skillsMerkleProof = MerkleProof(neededHashes, depth, playerSkills, 0)
        assert self.areLatestSkills(skillsMerkleProof), "Computed player state by CLIENT is not recognized by BC.."
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


    # ------------------------------------------------------------------------
    # ------------      Functions only for CLIENT                 ------------
    # ------------------------------------------------------------------------

    # return state of a player at end of a certain league
    # note that these do not contain potential sales done after the league
    def getPlayerSkillsAtEndOfLeague(self, leagueIdx, teamPosInLeague, playerIdx):
        self.assertIsClient()
        if leagueIdx == 0:
            return MinimalPlayerState(self.getPlayerStateBeforePlayingAnyLeague(playerIdx))

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
    def storePreHashDataInClientAtEndOfLeague(self, leagueIdx, dataAtMatchdays, lastDayTree, scores):
        self.assertIsClient()
        self.leagues[leagueIdx].storeDataAtMatchdays(dataAtMatchdays, scores)
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
        assert len(usersInitData["teamIdxs"]) % 2 == 0, "Currently we only support leagues with even nTeams"
        leagueIdx = len(self.leagues)
        self.leagues.append(LeagueClient(verseInit, verseStep, usersInitData))
        self.signTeamsInLeague(usersInitData["teamIdxs"], leagueIdx)
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
            self.VerseCommits[verse].actionsMerkleRoots,
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
                if playerIdx == 0: # if never written in teams.playerIdxs array
                    correctPlayerIdx = self.getPlayerIdxFromTeamIdxAndShirt(teamIdx, shirtNum)
                else:
                    correctPlayerIdx = playerIdx
                dataToChallengeInitSkills[teamPos][shirtNum] = self.computeDataToChallengePlayerSkills(correctPlayerIdx)
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
            return MerkleProof([], 0, self.getPlayerSkillsAtEndOfLastLeague(playerIdx), 0)
        else:
            skillsAllTeamsAtEndOfPrevLeague = self.leagues[prevLeagueIdx].dataAtMatchdays[-1].skillsAtMatchday
            playerSkills, playerPosInPrevLeague = self.getPlayerFromTeamStates(playerIdx, skillsAllTeamsAtEndOfPrevLeague[teamPosInPrevLeague])
            idxInFlattenedSkills = teamPosInPrevLeague*NPLAYERS_PER_TEAM+playerPosInPrevLeague

            lastDayTree = self.leagues[prevLeagueIdx].lastDayTree
            merkleProof = lastDayTree.prepareProofForLeaf(playerSkills, idxInFlattenedSkills)

            assert pylio.verifyMerkleProof(
                self.leagues[prevLeagueIdx].dataAtMatchdayHashes[-1],
                merkleProof,
                pylio.serialHash
            ), "Generated Merkle proof will not work"

            return merkleProof

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
        lastStatesFlattened = pylio.flatten(dataAtMatchdays[-1].skillsAtMatchday)
        lastDayTree = MerkleTree(lastStatesFlattened)
        dataAtMatchdayHashes.append(lastDayTree.root)
        return dataAtMatchdayHashes, lastDayTree


    # The CLIENT:
    # - computes all games of the league,
    # - in particular, all DataAtMatchdays => for every matchday: all teams states, tactics, teamOrders.
    # - stores both the pre-hash and the hashed DataAtMatchdays
    # - returns the hashed data so that it can then be send to the BC
    def updateLeagueInClient(self, leagueIdx, ADDR):
        self.assertIsClient()
        assert self.hasLeagueFinished(leagueIdx), "cannot update a league that is not finished"
        assert not self.leagues[leagueIdx].hasLeagueBeenUpdated(), "League has already been updated"
        dataAtMatchdays, scores = self.computeAllMatchdayStates(leagueIdx)
        initSkillsHash          = pylio.serialHash(self.leagues[leagueIdx].getInitPlayerSkills())
        dataAtMatchdayHashes, lastDayTree = self.prepareHashesForDataAtMatchdays(dataAtMatchdays)

        self.updateLeague(
            leagueIdx,
            initSkillsHash,
            dataAtMatchdayHashes,
            scores,
            ADDR
        )
        # and additionally, stores the league pre-hash data, and updates every player involved
        self.storePreHashDataInClientAtEndOfLeague(leagueIdx, dataAtMatchdays, lastDayTree, scores)
        assert initSkillsHash == pylio.serialHash(self.leagues[leagueIdx].getInitPlayerSkills()), "InitSkillsHash changed unintentionally"
        assert self.leagues[leagueIdx].hasLeagueBeenUpdated(), "League not detected as already updated"
        return initSkillsHash, dataAtMatchdayHashes, scores

    # returns states of all teams at start of a league. These include skills from previous league, and possible
    # sales after end of that league
    def getInitPlayerStates(self, leagueIdx):
        self.assertIsClient()
        usersInitData = pylio.duplicate(self.leagues[leagueIdx].usersInitData)
        nTeams = len(usersInitData["teamIdxs"])
        # an array of size [nTeams][NPLAYERS_PER_TEAM]
        initPlayerStates = pylio.createEmptyPlayerStatesForAllTeams(nTeams)
        for teamPosInLeague, teamIdx in enumerate(usersInitData["teamIdxs"]):
            for shirtNum in range(NPLAYERS_PER_TEAM):
                playerIdx = self.getPlayerIdxFromTeamIdxAndShirt(teamIdx, shirtNum)
                playerState = self.getCurrentPlayerState(playerIdx)
                initPlayerStates[teamPosInLeague][shirtNum] = playerState
        return initPlayerStates

