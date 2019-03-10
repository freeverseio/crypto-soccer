import numpy as np
import copy
import datetime
import sha3
from copy import deepcopy as duplicate
from constants import *
from structs import *

# Numpy accepts a max possible seed
def limitSeed(seed):
    return seed % MAX_RND_SEED_ALLOWED_BY_NUMPY

# Returns kekkack of string in hex format
def hexHash(str):
    return sha3.keccak_256(str).hexdigest()

# Returns kekkack of string in decimal format
def intHash(str):
    return int(hexHash(str), 16)

def serialHash(obj):
    return intHash(serialize(obj))

# Minimal (virtual) team creation. The Name could be the concat of the given name, and user int choice
# e.g. teamName = "Barcelona5443"
def createTeam(teamName, ownerAddr, ST):
    assert intHash(teamName) not in ST.teamNameHashToOwnerAddr, "You cannot create to teams with equal name!"
    teamIdx = len(ST.teams)
    ST.teams.append(Team(teamName))
    ST.teamNameHashToOwnerAddr[intHash(teamName)] = ownerAddr
    return teamIdx

# Given a seed, returns a balanced player.
# It only deals with skills & age, not playerIdx.
def getPlayerStateFromSeed(seed):
    newPlayerState = PlayerState()
    np.random.seed(seed)
    years = np.random.randint(MIN_PLAYER_AGE, MAX_PLAYER_AGE)
    newPlayerState.setMonth(years*12)
    skills = np.random.randint(0,AVG_SKILL-1,N_SKILLS)
    excess = int( (AVG_SKILL*N_SKILLS-skills.sum())/N_SKILLS )
    skills += excess
    newPlayerState.setSkills(skills)
    return newPlayerState


def assertTeamIdx(teamIdx, ST):
    assert teamIdx < len(ST.teams), "Team for this playerIdx not created yet!"
    assert teamIdx != 0, "Team 0 is reserved for null team!"


# If player has never been sold (virtual team): simple relation between playerIdx and (teamIdx, shirtNum)
# Otherwise, read what's written in the playerState
# playerIdx = 0 andt teamdIdx = 0 are the null player and teams
def getPlayerIdxFromTeamIdxAndShirt(teamIdx, shirtNum, ST):
    assertTeamIdx(teamIdx, ST)
    isPlayerIdxAssigned = ST.teams[teamIdx].playerIdxs[shirtNum] != 0
    if isPlayerIdxAssigned:
        return ST.teams[teamIdx].playerIdxs[shirtNum]
    else:
        return 1 + (teamIdx-1)*NPLAYERS_PER_TEAM + shirtNum

# The inverse of the previous relation
def getTeamIdxAndShirtForPlayerIdx(playerIdx, ST, forceAtBirth = False):
    if forceAtBirth or isPlayerVirtual(playerIdx, ST):
        teamIdx     = 1 + (playerIdx-1)/NPLAYERS_PER_TEAM
        shirtNum    = (playerIdx-1) % NPLAYERS_PER_TEAM
        return teamIdx, shirtNum
    else:
        return ST.playerIdxToPlayerState[playerIdx].getCurrentTeamIdx(), \
               ST.playerIdxToPlayerState[playerIdx].getCurrentShirtNum()

# the skills of a player are determined by concat of teamName and shirtNum
def getPlayerSeedFromTeamAndShirtNum(teamName, shirtNum):
    return limitSeed(intHash(teamName + str(shirtNum)))

# if player has never been sold, it will not be in the map playerIdxToPlayerState
# and his team is derived from a formula
def isPlayerVirtual(playerIdx, ST):
    return not playerIdx in ST.playerIdxToPlayerState

def getLastPlayedLeagueIdx(playerIdx, ST):
    # if player state has never been written, it played all leagues with current team (obtained from formula)
    # otherwise, we check if it was sold to current team before start of team's previous league
    if isPlayerVirtual(playerIdx, ST):
        teamIdx, shirtNum = getTeamIdxAndShirtForPlayerIdx(playerIdx, ST)
        return ST.teams[teamIdx].prevLeagueIdx, ST.teams[teamIdx].teamPosInPrevLeague

    currentTeamIdx  = ST.playerIdxToPlayerState[playerIdx].getCurrentTeamIdx()
    prevLeagueIdxForCurrentTeam = ST.teams[currentTeamIdx].prevLeagueIdx
    didHePlayLastLeagueWithCurrentTeam = ST.playerIdxToPlayerState[playerIdx].getLastSaleBlocknum() < \
                                         ST.leagues[prevLeagueIdxForCurrentTeam].blockInit
    if didHePlayLastLeagueWithCurrentTeam:
        return prevLeagueIdxForCurrentTeam, ST.teams[currentTeamIdx].teamPosInPrevLeague
    else:
        return ST.playerIdxToPlayerState[playerIdx].prevLeagueIdx, ST.playerIdxToPlayerState[playerIdx].prevTeamPosInLeague

def getPlayerStateAtEndOfLeague(prevLeagueIdx, teamPosInPrevLeague, playerIdx, ST):
    selectedStates =[s for s in ST.leagues[prevLeagueIdx].statesAtMatchday[-1][teamPosInPrevLeague] if s.getPlayerIdx() == playerIdx]
    assert len(selectedStates)==1, "PlayerIdx not found in previous league final states, or too many with same playerIdx"
    return selectedStates[0]

def getPlayerStateAtBirth(playerIdx, ST):
    # Disregard his current team, just look at the team at moment of birth to build skills
    teamIdx, shirtNum = getTeamIdxAndShirtForPlayerIdx(playerIdx, ST, forceAtBirth=True)
    seed = getPlayerSeedFromTeamAndShirtNum(ST.teams[teamIdx].name, shirtNum)
    playerState = duplicate(getPlayerStateFromSeed(seed))
    # Once the skills have been added, complete the rest of the player data
    playerState.setPlayerIdx(playerIdx)
    playerState.setCurrentTeamIdx(teamIdx)
    playerState.setCurrentShirtNum(shirtNum)
    return playerState


def copySkillsAndAgeFromTo(playerStateOrig, playerStateDest):
    playerStateDest.setSkills(duplicate(playerStateOrig.getSkills()))
    playerStateDest.setMonth(duplicate(playerStateOrig.getMonth()))


def getPlayerStateBeforePlayingAnyLeague(playerIdx, ST):
    # this can be called by BC or CLIENT, as both have enough data
    playerStateAtBirth = getPlayerStateAtBirth(playerIdx, ST)

    if isPlayerVirtual(playerIdx, ST):
        return playerStateAtBirth
    else:
        # if player has been sold before playing any league, it'll conserve skills at birth,
        # but have different metadata in the other fields
        playerState = duplicate(ST.playerIdxToPlayerState[playerIdx])
        copySkillsAndAgeFromTo(playerStateAtBirth, playerState)
        return playerState


def getLastWrittenPlayerStateFromPlayerIdx(playerIdx, ST_CLIENT):
    prevLeagueIdx, teamPosInPrevLeague = getLastPlayedLeagueIdx(playerIdx, ST_CLIENT)

    if prevLeagueIdx == 0:
        # this can be known both by CLIENT and BC
        return getPlayerStateBeforePlayingAnyLeague(playerIdx, ST_CLIENT)
    else:
        # this can only be accessed by the CLIENT
        return getPlayerStateAtEndOfLeague(prevLeagueIdx, teamPosInPrevLeague, playerIdx, ST_CLIENT)


# Simple player print
def printPlayer(playerState):
    toPrint =  "PlayerIdx: %s\n" % str(playerState.getPlayerIdx())
    toPrint += "Age      : %s\n" % str(playerState.getMonth()/12)
    toPrint += "Skills   : %s\n" % str(playerState.getSkills())
    toPrint += "TeamIdx  : %s\n" % str(playerState.getCurrentTeamIdx())
    toPrint += "ShirtNum : %s\n" % str(playerState.getCurrentShirtNum())
    toPrint += "SaleBlock: %s\n" % str(playerState.getLastSaleBlocknum())
    print "%s" % toPrint
    return intHash(toPrint) % 1000

# Simple team print
def printTeam(teamIdx, ST_CLIENT):
    hash = 0
    print "Player for teamIdx %d, with teamName %s: " %(teamIdx, ST_CLIENT.teams[teamIdx].name)
    for shirtNum in range(NPLAYERS_PER_TEAM):
        playerIdx = getPlayerIdxFromTeamIdxAndShirt(teamIdx, shirtNum, ST_CLIENT)
        playerState = getLastWrittenPlayerStateFromPlayerIdx(playerIdx,ST_CLIENT)
        playerChallengeData = computeDataToChallengePlayerIdx(playerState.getPlayerIdx(), ST_CLIENT)
        assert isCorrectStateForPlayerIdx(playerState, playerChallengeData, ST_CLIENT), "Player state not correctly in sync"
        hash += printPlayer(playerState)
    return hash


# quick solution to simulate changing teams.
# for the purpose of Lionel, we'll start with a simple exchange, instead
# of the more convoluted sell, assign, etc.
def exchangePlayers(playerIdx1, address1, playerIdx2, address2, ST):
    assert not isPlayerBusy(playerIdx1, ST), "Player sale failed: player is busy playing a league, wait until it finishes"
    assert not isPlayerBusy(playerIdx2, ST), "Player sale failed: player is busy playing a league, wait until it finishes"

    teamIdx1, shirtNum1 = getTeamIdxAndShirtForPlayerIdx(playerIdx1, ST)
    teamIdx2, shirtNum2 = getTeamIdxAndShirtForPlayerIdx(playerIdx2, ST)

    # check ownership!
    assert ST.teamNameHashToOwnerAddr[intHash(ST.teams[teamIdx1].name)] == address1, "Exchange Failed, owner not correct"
    assert ST.teamNameHashToOwnerAddr[intHash(ST.teams[teamIdx2].name)] == address2, "Exchange Failed, owner not correct"

    # get states from BC in memory to do changes, and only write back once at the end
    state1 = copy.deepcopy(getLastWrittenPlayerStateFromPlayerIdx(playerIdx1, ST))
    state2 = copy.deepcopy(getLastWrittenPlayerStateFromPlayerIdx(playerIdx2, ST))



    state1.prevLeagueIdx        = ST.teams[teamIdx1].currentLeagueIdx
    state1.prevTeamPosInLeague  = ST.teams[teamIdx1].teamPosInCurrentLeague

    state2.prevLeagueIdx        = ST.teams[teamIdx2].currentLeagueIdx
    state2.prevTeamPosInLeague  = ST.teams[teamIdx2].teamPosInCurrentLeague


    state1.setCurrentTeamIdx(teamIdx2)
    state2.setCurrentTeamIdx(teamIdx1)


    state1.setCurrentShirtNum(shirtNum2)
    state2.setCurrentShirtNum(shirtNum1)

    state1.setLastSaleBlocknum(ST.currentBlock)
    state2.setLastSaleBlocknum(ST.currentBlock)

    ST.teams[teamIdx1].playerIdxs[shirtNum1] = playerIdx2
    ST.teams[teamIdx2].playerIdxs[shirtNum2] = playerIdx1

    ST.playerIdxToPlayerState[playerIdx1] = duplicate(state1)
    ST.playerIdxToPlayerState[playerIdx2] = duplicate(state2)


def isValidOrdering(playerOrders):
    # TODO: check all nums are different and in [0, NPLAYERS_PER_TEAM]
    return True

def computeUsersInitDataHash(teamIdxs, playerOrders, tactics):
    # Consider changing an ordering map by a set of permutations
    # The reason is that it is a correct ordering by construction, no need to check
    #
    # inputs:  teamIdxs[nTeams], playerOrders[nTeams][NPLAYERS_PER_TEAM], tactics[nTeams]
    assert isValidOrdering(playerOrders), "The provided ordering of players is not valid"
    serialization = ""
    for (t, teamIdx) in enumerate(teamIdxs):
        serialization += str(teamIdx) + "-" + str(tactics[t]) + "-"
        for order in playerOrders[t]:
            serialization += str(order) + "-"
    return intHash(serialization)

def isPlayerBusy(playerIdx1, ST):
    return areTeamsBusyInPrevLeagues(
        [getTeamIdxAndShirtForPlayerIdx(playerIdx1, ST)[0]],
        ST)



def areTeamsBusyInPrevLeagues(teamIdxs, ST):
    for teamIdx in teamIdxs:
        if not ST.leagues[ST.teams[teamIdx].currentLeagueIdx].isFullyVerified(ST.currentBlock):
            return True
    return False

def signTeamsInLeague(teamIdxs, leagueIdx, ST):
    for teamPosInLeague, teamIdx in enumerate(teamIdxs):
        ST.teams[teamIdx].prevLeagueIdx             = duplicate(ST.teams[teamIdx].currentLeagueIdx)
        ST.teams[teamIdx].teamPosInPrevLeague       = duplicate(ST.teams[teamIdx].teamPosInCurrentLeague)

        ST.teams[teamIdx].currentLeagueIdx          = leagueIdx
        ST.teams[teamIdx].teamPosInCurrentLeague    = teamPosInLeague


def createLeague(blocknumber, blockStep, usersInitData, ST):
    assert not areTeamsBusyInPrevLeagues(usersInitData["teamIdxs"], ST), "League cannot create: some teams involved in prev leagues"
    assert len(usersInitData["teamIdxs"]) % 2 == 0, "Currently we only support leagues with even nTeams"
    leagueIdx = len(ST.leagues)
    ST.leagues.append( League(blocknumber, blockStep, usersInitData) )
    signTeamsInLeague(usersInitData["teamIdxs"], leagueIdx, ST)
    return leagueIdx

def createLeagueClient(blocknumber, blockStep, usersInitData, ST_CLIENT):
    assert not areTeamsBusyInPrevLeagues(usersInitData["teamIdxs"], ST_CLIENT), "League cannot create: some teams involved in prev leagues"
    leagueIdx = len(ST_CLIENT.leagues)
    ST_CLIENT.leagues.append( LeagueClient(blocknumber, blockStep, usersInitData) )
    signTeamsInLeague(usersInitData["teamIdxs"], leagueIdx, ST_CLIENT)
    return leagueIdx


def getInitPlayerStates(leagueIdx, ST, usersInitData = None, dataToChallengeInitStates = None):
    if not usersInitData:
        usersInitData = duplicate(ST.leagues[leagueIdx].usersInitData)
    nTeams = len(usersInitData["teamIdxs"])
    # an array of size [nTeams][NPLAYERS_PER_TEAM]
    initPlayerStates = [[None for playerPosInLeague in range(NPLAYERS_PER_TEAM)] for team in range(nTeams)]
    teamPosInLeague = 0
    for teamIdx, teamOrder in zip(usersInitData["teamIdxs"], usersInitData["teamOrders"]):
        for shirtNum, playerPosInLeague in enumerate(teamOrder):
            playerIdx = getPlayerIdxFromTeamIdxAndShirt(teamIdx, shirtNum, ST)
            if dataToChallengeInitStates:
                # gets the playerState from the challengeData
                # TODO TONI: there could be an error with player exhcnage...
                playerState = pylio.getPlayerStateFromChallengeData(
                    playerIdx,
                    dataToChallengeInitStates[teamPosInLeague][shirtNum],
                    ST
                )
                if not isCorrectStateForPlayerIdx(
                    playerState,
                    dataToChallengeInitStates[teamPosInLeague][shirtNum],
                    ST
                ):
                    return None
            else:
                # if no dataToChallenge is provided, it means this is a request
                # from a Client, so just read whatever pre-hash data you have
                playerState = getLastWrittenPlayerStateFromPlayerIdx(playerIdx, ST)
            initPlayerStates[teamPosInLeague][playerPosInLeague] = playerState
        teamPosInLeague += 1
    return initPlayerStates

def updateTacticsToBlockNum(tactics, blockNum, usersAlongData):
    for userData in [data for data in usersAlongData if data["block"] < blockNum]:
        for teamIdx, tact in zip(userData["teamIdxsWithinLeague"], userData["tactics"]):
            tactics[teamIdx] = tact

def getBlockHash(blockNum):
    return intHash('salt' + str(blockNum))


def shiftBack(t, nTeams):
    if (t < nTeams):
        return t
    else:
        return t-(nTeams-1)


def getTeamsInMatchFirstHalf(matchday, match, nTeams):
    team1 = 0
    if (match > 0):
        team1 = shiftBack(nTeams-match+matchday, nTeams)

    team2 = shiftBack(match+1+matchday, nTeams)
    if ( (matchday % 2) == 0):
        return team1, team2
    else:
        return team2, team1


def getTeamsInMatch(matchday, match, nTeams):
    assert matchday < 2 * (nTeams - 1), "This league does not have so many matchdays"
    if (matchday < (nTeams - 1)):
        (team1, team2) = getTeamsInMatchFirstHalf(matchday, match, nTeams)
    else:
        (team2, team1) = getTeamsInMatchFirstHalf(matchday - (nTeams - 1), match, nTeams);
    return team1, team2


def playMatch(initPlayerStates1, initPlayerStates2, tactics1, tactics2, MatchSeed):
    hash1 = intHash(str(MatchSeed)+serialize(initPlayerStates1)+serialize(tactics1))
    hash2 = intHash(str(MatchSeed)+serialize(initPlayerStates2)+serialize(tactics2))
    return hash1 % 4, hash2 % 4

def computeTeamRating(playerStates):
    return sum([sum(thisPlayerState.getSkills()) for thisPlayerState in playerStates])

def addFixedPointsToAllPlayers(playerStates, points):
    for playerState in playerStates:
        playerState.setSkills(playerState.getSkills() + points)

def updatePlayerStatesAfterMatch(playerStates1, playerStates2, goals1, goals2):
    ps1 = duplicate(playerStates1)
    ps2 = duplicate(playerStates2)

    if goals1 == goals2:
        return ps1, ps2

    pointsWon = computePointsWon(ps1, ps2, goals1, goals2)
    if goals1 > goals2:
        addFixedPointsToAllPlayers(ps1, pointsWon)
    else:
        addFixedPointsToAllPlayers(ps2, pointsWon)

    return ps1, ps2


def computePointsWon(playerState1, playerState2, goals1, goals2):
    ratingDiff              = computeTeamRating(playerState1) - computeTeamRating(playerState2)
    winnerWasBetter         = (ratingDiff > 0 and goals1>goals2) or (ratingDiff < 0 and goals1<goals2)

    if ratingDiff == 0:
        return 5
    else:
        return (2 if winnerWasBetter else 10)


def computeStatesAtMatchday(matchday, prevStates, tactics, matchdayBlock):
    nTeams = len(prevStates)
    nMatchesPerMatchday = nTeams/2
    scores = np.zeros([nMatchesPerMatchday, 2], int)
    statesAtMatchday = createEmptyPlayerStatesForAllTeams(nTeams)
    matchdaySeed = getBlockHash(matchdayBlock * 3)  # TODO: remove this *3

    for match in range(nMatchesPerMatchday):
        team1, team2 = getTeamsInMatch(matchday, match, nTeams)

        goals1, goals2 = playMatch(
            prevStates[team1],
            prevStates[team2],
            tactics[team1],
            tactics[team2],
            matchdaySeed
        )
        scores[match] = [goals1, goals2]
        statesAtMatchday[team1], statesAtMatchday[team2] = \
            updatePlayerStatesAfterMatch(
                    prevStates[team1],
                    prevStates[team2],
                    goals1,
                    goals2
                )
    return statesAtMatchday, scores



def computeAllMatchdayStates(blockInit, blockStep, initPlayerStates, usersInitData, usersAlongData):
    # In this initial implementation, evolution happens at the end of the league only
    tactics = duplicate(usersInitData["tactics"])
    nTeams = len(usersInitData["teamIdxs"])
    nMatchdays = 2*(nTeams-1)
    nMatchesPerMatchday = nTeams/2
    scores = np.zeros([nMatchdays, nMatchesPerMatchday, 2], int)
    matchdayBlock = duplicate(blockInit)

    # the following beast has dimension nMatchdays x nTeams x nPlayersPerTeam
    statesAtMatchday = [createEmptyPlayerStatesForAllTeams(nTeams) for matchday in range(nMatchdays)]

    for matchday in range(nMatchdays):
        updateTacticsToBlockNum(tactics, matchdayBlock, usersAlongData)
        prevStates = initPlayerStates if matchday == 0 else statesAtMatchday[matchday - 1]
        statesAtMatchday[matchday], scores[matchday] = computeStatesAtMatchday(
            matchday,
            prevStates,
            tactics,
            matchdayBlock
        )
        matchdayBlock += blockStep


    return statesAtMatchday, scores

def computeUsersAlongDataHash(usersAlongData):
    usersAlongDataHash = 0
    for entry in usersAlongData:
        usersAlongDataHash = intHash(str(usersAlongDataHash) + serialize(entry))
    return usersAlongDataHash

def getMatchsPlayerByTeam(selectedTeam, nTeams):
    matchdayMatch = []
    nMatchdays = 2*(nTeams-1)
    nMatchesPerMatchday = nTeams/2
    for matchday in range(nMatchdays):
        for match in range(nMatchesPerMatchday):
            team1, team2 = getTeamsInMatch(matchday, match, nTeams)
            if (team1==selectedTeam) or (team2==selectedTeam):
                matchdayMatch.append([matchday,match])
    return matchdayAndMatch


def areUpdaterScoresCorrect(selectedMatchInMatchday, selectedScores, updaterScores):
    for matchday, match, score in zip(range(len(selectedMatchInMatchday)), selectedMatchInMatchday, selectedScores):
        if any(updaterScores[matchday][match] != score):
            return False
    return True


def updateClientAtEndOfLeague(leagueIdx, initPlayerStates, statesAtMatchday, scores, ST_CLIENT):
    ST_CLIENT.leagues[leagueIdx].updateInitState(initPlayerStates)
    ST_CLIENT.leagues[leagueIdx].updateStatesAtMatchday(statesAtMatchday, scores)
    # the last matchday gives the final states used to update all players:
    for allPlayerStatesInTeam in statesAtMatchday[-1]:
        for playerState in allPlayerStatesInTeam:
            ST_CLIENT.playerIdxToPlayerState[playerState.getPlayerIdx()] = playerState



def getTeamIdxInLeague(currentTeamIdx, lastLeagueIdx, ST_CLIENT):
    for idxInLeague, teamIdx in enumerate(ST_CLIENT.leagues[lastLeagueIdx].usersInitData["teamIdxs"]):
        if teamIdx == currentTeamIdx:
            return idxInLeague
    assert False, "The team is not in this league"



def areEqualStructs(st1, st2):
    return serialHash(st1) == serialHash(st2)

def computeDataToChallengePlayerIdx(playerIdx, ST_CLIENT):
    prevLeagueIdx, teamPosInPrevLeague = getLastPlayedLeagueIdx(playerIdx, ST_CLIENT)
    if prevLeagueIdx == 0:
        return getLastWrittenPlayerStateFromPlayerIdx(playerIdx, ST_CLIENT)
    else:
        return getAllStatesAtEndOfLeague(prevLeagueIdx, ST_CLIENT)

def getAllStatesAtEndOfLeague(leagueIdx, ST_CLIENT):
    return ST_CLIENT.leagues[leagueIdx].statesAtMatchday[-1]


def prepareDataToChallengeInitStates(leagueIdx, ST_CLIENT):
    thisLeague = duplicate(ST_CLIENT.leagues[leagueIdx])
    nTeams = len(thisLeague.usersInitData["teamIdxs"])
    dataToChallengeInitStates = [[None for player in range(NPLAYERS_PER_TEAM)] for team in range(nTeams)]
    # dimensions: [team, nPlayersInTeam]
    #   if that a given player is virtual, then it contains just its state
    #   if not, it contains all states of prev league's team
    for teamPos, teamIdx in enumerate(thisLeague.usersInitData["teamIdxs"]):
        for shirtNum, playerIdx in enumerate(ST_CLIENT.teams[teamIdx].playerIdxs):
            if playerIdx == 0:
                dataToChallengeInitStates[teamPos][shirtNum] = computeDataToChallengePlayerIdx(
                    getPlayerIdxFromTeamIdxAndShirt(teamIdx, shirtNum, ST_CLIENT),
                    ST_CLIENT
                )
            else:
                assert playerIdx == getPlayerIdxFromTeamIdxAndShirt(teamIdx, shirtNum, ST_CLIENT), "PlayerIdx should always coincide"
                dataToChallengeInitStates[teamPos][shirtNum] = computeDataToChallengePlayerIdx(playerIdx, ST_CLIENT)
    return dataToChallengeInitStates


def isCorrectStateForPlayerIdx(playerState, dataToChallengePlayerState, ST):
    # If player has never played a league, we can compute the playerState directly in the BC
    # It basically is equal to the birth skills, with ,potentially, a few team changes via sales.
    # If not, we can just compare the hash of the dataToChallengePlayerState with the stored hash in the prev league
    playerIdx = playerState.getPlayerIdx()
    prevLeagueIdx, teamPosInPrevLeague = getLastPlayedLeagueIdx(playerIdx, ST)
    if prevLeagueIdx == 0:
        return areEqualStructs(
            playerState,
            getPlayerStateBeforePlayingAnyLeague(playerIdx, ST)
        )
    else:
        assert isPlayerStateInsideDataToChallenge(playerState, dataToChallengePlayerState, teamPosInPrevLeague), \
            "The playerState provided is not part of the challengeData"
        return serialHash(dataToChallengePlayerState) == ST.leagues[prevLeagueIdx].statesAtMatchdayHashes[-1]

def isPlayerStateInsideDataToChallenge(playerState, dataToChallengePlayerState, teamPosInPrevLeague):
    return playerState in dataToChallengePlayerState[teamPosInPrevLeague]


def getPlayerStateFromChallengeData(playerIdx, dataToChallengePlayerState, ST):
    if type(dataToChallengePlayerState) == type([]):
        prevLeagueIdx, teamPosInPrevLeague = getLastPlayedLeagueIdx(playerIdx, ST)
        thisPlayerState = [s for s in dataToChallengePlayerState[teamPosInPrevLeague] if s.getPlayerIdx() == playerIdx]
        assert len(thisPlayerState) < 2, "This data contains more than once the required playerIdx"
        assert len(thisPlayerState) > 0, "This data does not contain the required playerIdx"
        return thisPlayerState[0]
    else:
        assert dataToChallengePlayerState.getPlayerIdx() == playerIdx, "This data does not contain the required playerIdx"
        return dataToChallengePlayerState


def createEmptyPlayerStatesForAllTeams(nTeams):
    return  [[None for playerPosInLeague in range(NPLAYERS_PER_TEAM)] for team in range(nTeams)]
