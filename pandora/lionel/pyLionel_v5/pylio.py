import numpy as np
import copy
import datetime
import sha3
from copy import deepcopy as duplicate
from constants import *
from structs import *
from merkle_tree import *


# ------------ Functions to take hashes, serialize structs ------------

#  serializes and converts to str in a complicated way
def serialize2str(object):
    return str(serialize(object).hex())

# Numpy accepts a max possible seed
def limitSeed(seed):
    return seed % MAX_RND_SEED_ALLOWED_BY_NUMPY

# Returns keccak of string in hex format
def hexHash(str):
    return sha3.keccak_256(str.encode('utf-8')).hexdigest()

# Returns kekkack of string in decimal format
def intHash(str):
    return int(hexHash(str), 16)

# serializes object and then taken intHash
def serialHash(obj):
    return intHash(serialize2str(obj))


# ------------ Functions to print structs ------------

def printPlayerFromSkills(ST_CLIENT, playerSkills):
    return printPlayer(ST_CLIENT.skillsToLastWrittenState(playerSkills))

# Simple player print
def printPlayer(playerState):
    toPrint =  "PlayerIdx: %s\n" % str(playerState.getPlayerIdx())
    toPrint += "Age      : %s\n" % str(playerState.getMonth()/12)
    toPrint += "Skills   : %s\n" % str(playerState.getSkills())
    toPrint += "TeamIdx  : %s\n" % str(playerState.getCurrentTeamIdx())
    toPrint += "ShirtNum : %s\n" % str(playerState.getCurrentShirtNum())
    toPrint += "SaleBlock: %s\n" % str(playerState.getLastSaleBlocknum())
    print("%s" % toPrint)
    return intHash(toPrint) % 1000

# Simple team print
def printTeam(teamIdx, ST_CLIENT):
    hash = 0
    print("Player for teamIdx %d, with teamName %s: " %(teamIdx, ST_CLIENT.teams[teamIdx].name))
    for shirtNum in range(PLAYERS_PER_TEAM_MAX):
        if ST_CLIENT.isShirtNumFree(teamIdx, shirtNum):
            continue
        playerIdx = ST_CLIENT.getPlayerIdxFromTeamIdxAndShirt(teamIdx, shirtNum)
        playerSkills = ST_CLIENT.getPlayerSkillsAtEndOfLastLeague(playerIdx)
        playerChallengeData = ST_CLIENT.computeDataToChallengePlayerSkills(playerSkills.getPlayerIdx())
        assert ST_CLIENT.areLatestSkills(playerChallengeData), "Player state not correctly in sync"
        hash += printPlayerFromSkills(ST_CLIENT, playerSkills)
    return hash

def isValidOrdering(playerOrders):
    # TODO: Currently not implemented. Check all nums are different and in [0, PLAYERS_PER_TEAM_INIT]
    return True

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

# mockup
def playMatch(initPlayerSkills1, initPlayerSkills2, tactics1, tactics2, teamOrders1, teamOrders2, MatchSeed):
    hash1 = intHash(str(MatchSeed)+serialize2str(initPlayerSkills1)+serialize2str(tactics1)+serialize2str(teamOrders1))
    hash2 = intHash(str(MatchSeed)+serialize2str(initPlayerSkills2)+serialize2str(tactics2)+serialize2str(teamOrders2))
    return hash1 % 4, hash2 % 4

# mockup: the rating of a team is just the sum of all skills
def computeTeamRating(playerSkills):
    return sum([sum(thisPlayerSkills.getSkills()) for thisPlayerSkills in playerSkills])


def addFixedPointsToAllPlayers(allPlayerSkills, points):
    for playerSkills in allPlayerSkills:
        playerSkills.setSkills(playerSkills.getSkills() + points)


# given the result, it computes the evolution points won per team, and applies them to their players
def updatePlayerSkillsAfterMatch(playerSkills1, playerSkills2, goals1, goals2):
    ps1 = duplicate(playerSkills1)
    ps2 = duplicate(playerSkills2)

    if goals1 == goals2:
        return ps1, ps2

    pointsWon = computePointsWon(ps1, ps2, goals1, goals2)
    if goals1 > goals2:
        addFixedPointsToAllPlayers(ps1, pointsWon)
    else:
        addFixedPointsToAllPlayers(ps2, pointsWon)

    return ps1, ps2


# simple mockup of what the evolution points could look like.
def computePointsWon(playerState1, playerState2, goals1, goals2):
    ratingDiff              = computeTeamRating(playerState1) - computeTeamRating(playerState2)
    winnerWasBetter         = (ratingDiff > 0 and goals1>goals2) or (ratingDiff < 0 and goals1<goals2)

    if ratingDiff == 0:
        return 5
    else:
        return (2 if winnerWasBetter else 10)


# plays all games in a given matchday, using the provided input for how the teams
# were right at the beginning of that matchday
def computeStatesAtMatchday(matchday, prevSkills, tactics, teamOrders, matchdaySeed):
    nTeams = len(prevSkills)
    nMatchesPerMatchday = nTeams//2
    scores = np.zeros([nMatchesPerMatchday, 2], int)
    skillsAtMatchday = createEmptyPlayerStatesForAllTeams(nTeams)

    for match in range(nMatchesPerMatchday):
        team1, team2 = getTeamsInMatch(matchday, match, nTeams)
        matchSeed = intHash(str(matchdaySeed + match))
        goals1, goals2 = playMatch(
            prevSkills[team1],
            prevSkills[team2],
            tactics[team1],
            tactics[team2],
            teamOrders[team1],
            teamOrders[team2],
            matchSeed
        )
        scores[match] = [goals1, goals2]
        skillsAtMatchday[team1], skillsAtMatchday[team2] = \
            updatePlayerSkillsAfterMatch(
                    prevSkills[team1],
                    prevSkills[team2],
                    goals1,
                    goals2
                )
    return skillsAtMatchday, scores

# checks if 2 structs are equal by comparing the hash of their serialization
def areEqualStructs(st1, st2):
    return serialHash(st1) == serialHash(st2)


def createEmptyPlayerStatesForAllTeams(nTeams):
    return arrayDims(PLAYERS_PER_TEAM_MAX, nTeams)


# ---------------- FUNCTIONS TO ADVANCE BLOCKS IN THE BC AND CLIENT ----------------
# advances both BC and CLIENT, and commits the userActions if it goes through a verse
def advanceToBlock(n, ST, ST_CLIENT):
    nBlocksToAdvance = n - ST.currentBlock
    assert nBlocksToAdvance > 0, "cannot advance less than 1 block"
    for block in range(nBlocksToAdvance):
        assert ST.isCrossingVerse() == ST_CLIENT.isCrossingVerse(), "CLIENT and BC not synced in verse crossing"
        if ST.isCrossingVerse():
            ST_CLIENT.syncTimeZoneCommits(ST)
        ST.incrementBlock()
        ST_CLIENT.incrementBlock()

def advanceNBlocks(deltaN, ST, ST_CLIENT):
    advanceToBlock(
        ST.currentBlock + deltaN,
        ST,
        ST_CLIENT
    )

def advanceNVerses(nVerses, ST, ST_CLIENT):
    nBlocks = nVerses*ST.blocksBetweenVerses
    if nBlocks == 0:
        return
    advanceNBlocks(nBlocks, ST, ST_CLIENT)

def advanceToEndOfLeague(leagueIdx, ST, ST_CLIENT):
    verseFinal = ST.leagues[leagueIdx].verseFinal()
    while ST.currentVerse < verseFinal:
        advanceNBlocks(1, ST, ST_CLIENT)

# ------------------------------------------------


# selects a pseudo-random element of an array (obtained from hashing the seed)
def getRandomElement(arr, seed):
    nElems = len(arr)
    return arr[intHash(serialize2str(seed)) % nElems]


# returns an 1D-array from a 2D-array
def flatten(statesPerTeam):
    flatStates = []
    for statesTeam in statesPerTeam:
        for statePlayer in statesTeam:
            flatStates.append(PlayerSkills(statePlayer)) # select only skills and playerIdx
    return flatStates

def challengeLevel1(verse, addr, ST, ST_CLIENT, lie):
    superRoots, leagueRoots = ST_CLIENT.computeLeagueHashesForVerse(verse)
    if lie == 0:
        ST.challengeLevel1(verse, superRoots, addr)
    else:
        superRootsLie, leagueRootsLie = createLieSuperRoot(superRoots, leagueRoots, lie)
        ST.challengeLevel1(verse, superRootsLie, addr)

def challengeLevel2(verse, subVerse, addr, ST, ST_CLIENT, lie):
    superRoots, leagueRoots = ST_CLIENT.computeLeagueHashesForVerse(verse)
    if lie == 0:
        ST.challengeLevel2(verse, subVerse, leagueRoots[subVerse], addr)
    else:
        superRootsLie, leagueRootsLie = createLieSuperRoot(superRoots, leagueRoots, lie)
        ST.challengeLevel2(verse, subVerse, leagueRootsLie[subVerse], addr)

def challengeLevel3(verse, leagueIdx, addr, ST, ST_CLIENT, lie):
    # lie = 0   => tell the truth
    # lie = 1   => lie in initSkillsHash
    # lie = 2   => lie in first dataAtMatchdayHashes[0]
    # lie = 3   => lie in first dataAtMatchdayHashes[1]
    # ...
    dataToChallengeLeague = ST_CLIENT.leagues[leagueIdx].dataToChallengeLeague
    posInSubverse = ST.getPosInSubverse(verse, leagueIdx)
    if lie == 0:
        ST.challengeLevel3(verse, posInSubverse, dataToChallengeLeague, addr)
    else:
        dataToChallengeLeagueLie = pylio.duplicate(dataToChallengeLeague)
        if lie == 1:
            dataToChallengeLeagueLie.initSkillsHash += 1
        else:
            dataToChallengeLeagueLie.dataAtMatchdayHashes[lie-2] += 1
        ST.challengeLevel3(verse, posInSubverse, dataToChallengeLeagueLie, addr)


def challengeLevel4(selectedMatchday, verse, ST, ST_CLIENT):
    ST.assertCanChallengeStatus(verse, UPDT_LEVEL4)
    posInSubVerse = ST.verseToLeagueCommits[verse].posInSubVerse
    leagueRoot = ST.verseToLeagueCommits[verse].leagueRoots[posInSubVerse]
    assert leagueRoot != 0, "You cannot challenge a league that is not part of the verse commit"
    leagueIdx = ST.getLeagueIdxFromPosInSubverse(verse, posInSubVerse)

    if selectedMatchday == LEAGUE_INIT_SKILLS_ID:
        ST.challengeLevel4InitSkills(
            verse,
            ST_CLIENT.leagues[leagueIdx].usersInitData,
            duplicate(ST_CLIENT.leagues[leagueIdx].dataToChallengeInitSkills)
        )
    else:
        challengeLevel4MatchDay(selectedMatchday, verse, leagueIdx, ST, ST_CLIENT)


# It uses the CLIENT data to submit a challenge to the BC
def challengeLevel4MatchDay(selectedMatchday, verse, leagueIdx, ST, ST_CLIENT):
    # ...first, it selects a matchday, and gathers the data at that matchday (states, tactics, teamOrders)
    dataAtPrevMatchday = ST_CLIENT.getPrevMatchdayData(leagueIdx, selectedMatchday)
    # ...next, it builds the Merkle proof for the actions commited on the corresponding verse, for that league
    merkleProofDataForMatchday = ST_CLIENT.getActionsMerkleProofForMatchday(leagueIdx, selectedMatchday)

    assert pylio.areEqualStructs(
        ST_CLIENT.leagues[leagueIdx].actionsPerMatchday[selectedMatchday],
        merkleProofDataForMatchday.leaf[1]
    ), "The Merkle Proof does not contain the correct pre-hash actions for that day"

    # ...finally, it does the challenge. If successful, it will reset() the leauge update
    ST.challengeMatchdayStates(
        verse,
        selectedMatchday,
        dataAtPrevMatchday,
        duplicate(ST_CLIENT.leagues[leagueIdx].usersInitData),
        merkleProofDataForMatchday
    )


def verifyMerkleProof(root, merkleProof, hashFunction):
    # the current library requires the leaf & leafIdx to be formatted inside 'values' as follows:
    values = {merkleProof.leafIdx: merkleProof.leaf}
    return verify(root, merkleProof.depth, values, merkleProof.neededHashes, hashFunction)

# ------------

# The CLIENT retrieves its last full player state, and sends a TX to the BC to make sure it can be certified
# To do so, that TX includes the merkleProof that the player skills were part of a previous league last matchday hash.
def assertPlayerStateInClientIsCertifiable(playerIdx, ST, ST_CLIENT):
    playerState = ST_CLIENT.getCurrentPlayerState(playerIdx)
    dataToChallengePlayerSkills = ST_CLIENT.computeDataToChallengePlayerSkills(playerIdx)
    assert ST.certifyPlayerState(playerState, dataToChallengePlayerSkills),\
        "Current player state in CLIENT is not certifiable by BC.."

def arrayDims(dim1, dim2):
    return [[None for d1 in range(dim1)] for d2 in range(dim2)]

def shouldFail(f, msg):
    itFailed = False
    try:
        f(0)
    except AssertionError as error:
        print("Expected fail:")
        print("..." + msg + "...with internal error: " + str(error))
        itFailed = True
    assert itFailed, "We should have failed, but did not"

def createLieSuperRoot(superRoots, leagueRoots, factor):
    superRootsLie = duplicate(superRoots)
    leagueRootsLie = duplicate(leagueRoots)

    for s, supers in enumerate(superRootsLie):
        for l, leagues in enumerate(leagueRootsLie[s]):
            leagueRootsLie[s][l] *= factor
        tree = MerkleTree(leagueRootsLie[s])
        superRootsLie[s] = tree.root

    return superRootsLie, leagueRootsLie

def createTeam(teamName, addr, ST, ST_CLIENT):
    teamIdx1 = ST.createTeam(teamName, addr)
    teamIdx2 = ST_CLIENT.createTeam(teamName, addr)
    assert teamIdx1 == teamIdx2, "ST and ST_CLIENT not in sync"
    return teamIdx1

def exchangePlayers(playerIdx1, playerIdx2, ST, ST_CLIENT):
    addr1 = ST.getOwnerAddrFromPlayerIdx(playerIdx1)
    addr2 = ST.getOwnerAddrFromPlayerIdx(playerIdx2)
    assert addr1 == ST_CLIENT.getOwnerAddrFromPlayerIdx(playerIdx1), "ST and ST_CLIENT not in sync"
    assert addr2 == ST_CLIENT.getOwnerAddrFromPlayerIdx(playerIdx2), "ST and ST_CLIENT not in sync"

    ST.exchangePlayers(
        playerIdx1, addr1,
        playerIdx2, addr2
    )
    ST_CLIENT.exchangePlayers(
        playerIdx1, addr1,
        playerIdx2, addr2
    )

def exchangePlayers(playerIdx1, playerIdx2, ST, ST_CLIENT):
    addr1 = ST.getOwnerAddrFromPlayerIdx(playerIdx1)
    addr2 = ST.getOwnerAddrFromPlayerIdx(playerIdx2)
    assert addr1 == ST_CLIENT.getOwnerAddrFromPlayerIdx(playerIdx1), "ST and ST_CLIENT not in sync"
    assert addr2 == ST_CLIENT.getOwnerAddrFromPlayerIdx(playerIdx2), "ST and ST_CLIENT not in sync"

    ST.exchangePlayers(
        playerIdx1, addr1,
        playerIdx2, addr2
    )
    ST_CLIENT.exchangePlayers(
        playerIdx1, addr1,
        playerIdx2, addr2
    )


def movePlayerToTeam(playerIdx, teamIdx, ST, ST_CLIENT):
    ST.movePlayerToTeam(playerIdx, teamIdx)
    ST_CLIENT.movePlayerToTeam(playerIdx, teamIdx)


def createCountry(timeZone, ST, ST_CLIENT):
    countryIdx = ST.createCountry(timeZone)
    countryIdx_client = ST_CLIENT.createCountry(timeZone)
    assert countryIdx == countryIdx_client, "ST/ST_CLIENT not in sync"
    return countryIdx

def addDivision(countryIdx, ST, ST_CLIENT):
    divisionIdx = ST.addDivision(countryIdx)
    divisionIdx_client = ST_CLIENT.addDivision(countryIdx)
    assert divisionIdx == divisionIdx_client, "ST/ST_CLIENT not in sync"
    return divisionIdx


# TODO: all this can be precompiled and remove calls to cycleIdx
def cycleIdx(day, turnInDay):
    return (day - 1) * 4 + turnInDay

