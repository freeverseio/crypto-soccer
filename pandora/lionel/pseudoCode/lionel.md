# Mappings
    - mapping(bytes32 => address) public teamToOwnerAddr; /// from team hash(name) to the owner's address.
    - mapping(uint => state) public playerIdxToState; 
    - mapping(uint => address) public playerIdxToOwnerAddr; // for non-team players



# Numerology

## Bits to store block numbers? 

About 1 block is produced every 10s. 
    - about 3M blocks per year: we would need 22 bit
    - about 8K blocks per day: we would need 13 bit
    - if we want the game to last 100 years => 29bit
    - if we want the time between games to be 1 week => 16 bit


## Bits to store monthOfBirthAfterUnixEpoch (i.e. after 1970)

14bit = 1364 years  (ENOUGH)

## Bits to store player skills (they start with an avg of 50 points per skill)

14bit = 32K values per skill!
5 skills => 70bit


## Bits to store playerIdx (i.e. what's the max num of players that can be created?)

28 = 270M players ==> 9 players in a uint256. 


## Bits to store laegueIdx (i.e. what's the max num of leagues that can be created?)

25bit = 33M leagues.


## Bits to store teamIdx (i.e. what's the max num of teams that can be created?)

28bit = 268 teams (allows 9 per uint, 28x9=252)


## Leagues
- Games played by a team in a league: nGamesPerTeam = 2 (nTeams - 1)
- Total games that make a league: nTeams (nTeams-1)    (a square minus de diagonal)

references:
10 teams => 90 games overall, 18 for one team
11 teams => 110 games, 20 for one team
16 teams => 240 games, 30 for one team



## Max games that can be computed in one block

Max gas per block = 8M.
Max gas for a reasonable large Tx = 4M.
Currently, 1 game = 260K. Let's say 300K.

18 games = 5.4M  => so a league of 10 seems like the upper limit if we want to allow challenging in one atomic Tx.

LETS LIMIT TO LEAGUES OF 10 for the time being.


## Bits per game result

if we limit max goals per team per game to 15 => use 4bit per team per game

if nTeams = 10
4 x 2 x 10 x (10-1) = 7.5 * 256 => 3 uints

## Bits for shirtNum

this is basically the "dorsal", which specifies the ordering of players in a team.
nBits = 4 => max 16 player in a team 

## Tactics

When a team joins a league, he commits to an order (sorting) of its players, so his join TX ensures (playerIdx, shirtNums), where shirtNums = 14 x 4 = 56b

Tactics are a choice among 8 (442, 433,...) plus a choice of 'pressing bool' => 4 bit total => 8 bit per game.

Since there are 90 games => 720 bits = 3 uint or uint8[90]


# Main Structs/Data

## Player state

    - uint state = serialization of:
        - skills: 5 x 14 = 70bit
        - monthOfBirthAfterUnixEpoch: 14bit
        - currentTeamIdx: 28bit
        - playerIdx: 28bit
        - lastSaleBlocknum: 29bit

    Total = 169b



## Teams
Team is a struct that has:
    - string tname = unique
    - uint256 playersIdxA = serialization of 9 playerIdx
    - uint256 playersIdxB = serialization of up to 7 playerIdx (7x28=196) 
                                + previousLeagueIdx ( +28 = 224)
                                + currentLeagueIdx ( +28 = 252)



## League

Updaters have the state of all teams at league join, built from genesis.

A challenge in a TX can process at most one team through the league. 

For that, they need to know the state of all players, and all teams, at start of league.
These are not written anywhere, so they need to be provided as input in the 
challenger's TX, and the contract will check that their hashes are contained
in the previous leagues final hashes.

struct:

  Provided by a user creating the league (and others that join):

    - uint256 init: containing a serialization of (n0, nStep, nTeams, teamIdx_A for 7 teams)
    - uint256 teamIdx_B (space for 9 more teams)
    - bytes32 usersDataHash
    - tactics uint8[90]

usersDataHash is built from the order at which inscriptions took place.
    - each inscription is a uint: I1 = (teamIdx1, shirtNums1[14]),...
    - userDataHash = H(I10 + H(I9 + ... H(I0)...))



  Provided by UPDATER (in Lionel v3)

    - initHash = hash(  hash(stateTeam1) + ... + hash(stateTeamN) )
    - finalTeamStatesHash[10] = [hash(stateTeam1), ..., hash(stateTeamN)]
    - uint8 scores[90]   
    - address updaterAddr = who wrote the update for this league
    - blockNumLastUpdate (28b)  

Nomenclature:  
D = 'Data' = { initHash, finalStatesHash[10], scores[90] }

Some data on bits required:

- bits left for teamidx_A = 203 = 28x7 + 7 => space for 7 teams
- bits for n0: 31 (max of 400 years of game)
- bits for nStep: 17 (max of 2 weaks between games)
- bits for nTeams: 5 (max of 32 teams)
- bits per teamIdx: 28 (max 266M teams)

- for results: 4bit (15 goals) x 2 x nTeams x (nTeams-1) 
                            (for nTeams = 10)               = 3 uints

- on initHash: just concatenate all playerStates of a teams, and hash. Concatenate all hashes, and hash. The reason why we include hashes instead of concatenating all playerStates of all teams, is that each hash has to be done anyway to prove that it's in the previous-league Merkle tree, and it sounds cheaper.




# Main Functions to control Updaters/Challengers


Updaters keep the entire state. We will disregard for the time being possible 
evolution of players between leagues. 


## isPlayer or isTeam or isLeague in challenging state, or safely updated?

- isPlayerBusyPlayingLeague?
Player.currentTeam -> team.currentLeague -> n + nGames x deltaN < now ?

- isPlayerStateUpdated?
Player.currentTeam -> team.currentLeague -> blockNumLastUpdate != 0 ?

- isPlayerInChallengingPeriod?
Player.currentTeam -> team.currentLeague -> blockNumLastUpdate + challengePeriod < now ?

- isPlayerReadyForNewStuff = isPlayerStateUpdated && !isPlayerInChallengingPeriod

Same for teams or leagues


finalStates are written using sort!

## When selling a player

Updater must provide:
    - state before sale (basically, skills) 
    - teamNumInLeagueBeforeSale
    - shirtNumInLeagueBeforeSale
    - rest of playerStates of the team at the end of that league

The TX:
    - finds out the latest leagueIdx truly played by comparing:
        if currentTeam.currentLeague.n0 > lastSaleBlocknum:
            - use currentLeagueIdx
            - checks coincidence with hash(teamPlayerStates), for which you need:
                - teamNumInLeagueBeforeSale
                - shirtNumInLeagueBeforeSale
        else:
            it was updated before last sale, no need to check nor modify skills

    - rewrites 1 uint:
        - lastSaleBlocknum
        - new skills (if needed)
        - new currentTeamIdx



## UpdateLeague(leagueIdx, D)

He provides the collection of data D, as defined above. The contract:

0. checks if league is in challenging period, otherwise halts.
1. Writes initHash, finalTeamStatesHash[10], results[3], challengerUpdaterAddr, blockNumLastUpdate


## ChallengeInitHash(leagueIdx, allLeaguePlayersStates[140], I[10], MerkleProofs[140][3])

recall: I[n] = {teamIDx_n, sortedPlayer_n}

0. require league is in challenging period.
1. require initHash(allLeaguePlayersStates) != initHash (make sure you claim he lied)
2. require leauge.usersDataHash == H(I+H(...))
2. for (team, shirtNums) in I:
     for (playerIdx, shirtNum) in (team.playerIdxs, shirtNums):
        # make sure that the provided playerState refers to this playerIdx
        # the position to look for in allLeaguePlayersStates[140] is:
        #   allLeaguePlayersStates[teamNumInLeague x 14+ shirtNum]
        - require that:
            player idx in allLeaguePlayersStates[teamNumInLeague x 14+ shirtNum] = playerIdx
        - find previous league for that player by getting the newest of:
                if currentTeam.previousLeagueIdx.n0 > lastSaleBlocknum:
                    - use previousLeagueIdx
                    - compute hash(teamFinalState) by using: 
                     player state, player.previousLeagueShirtNum, MerkleProofs[player][3]
                    - require it coincides with previous league 
                        allLeaguePlayersStates[player.teamNumInLeagueBeforeSale]
                else:
                    player state was updated after last sale, so just
                    - require it coincides with player.state


3. if you get here, the new initHash is correct. Slash updater. Reset.


## ChallengeFinalTeamStatesHash(leagueIdx, D', allLeaguePlayersStates, selectedTeam)

0. require league is in challenging period.

1. require initHash(allLeaguePlayersStates) == initHash (make sure we agree with init)

2. plays all league (many playGame) for selectTeam, based on the provided allLeaguePlayersStates in this TX call. Before each game, read from tactics[90]. 
At the end the TX will have computed:
    - finalState for selectedTeam
    - scores[18] for selectedTeam

3. require that at least one of this is true:
    - hash(finalState) !=  finalTeamStatesHash(selectedTeam)
    - scores[18] != corresponding 18 entries in scores[90]

4. slash updater. Reset.



