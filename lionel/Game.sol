contract Game {
    // assume nTeams = 10  => nGames = 90
    Stakers     stakers;
    
    address     firstUpdaterAddr;
    uint256     blockNumLastUpdate;

    bytes32     allInitStatesHash;      // hash of bytes32[10]
    bytes32[10] finalTeamStatesHash;        // 10 teams, for each: keccak(playerStatesForThatTeam)
    byte[90]    allGamePlayResults; // 1 result = 8bit, 10 teams => 90 total plays 
    bool[10]    wasTeamVerified;      //*// 1 bool for inithash, 1 for each finalStatesHash
    bool        wasInitStateVerified; //*// 1 bool for inithash, 1 for each finalStatesHash

    function updateForFirstTime(
        bytes32              _allInitStatesHash,
        bytes32[10] calldata _finalTeamStatesHash,
        byte[90] calldata    _allGamePlayResults
        uint _leagueIdx; //*//
    ) external {
        stakers.verify(msg.sender);     // valid actor
        require( leagues[leagueIdx].hasFinished()) ; //*// only do if league has finished
        require(blockNumLastUpdate==0x0); // not challanged yet

        updaterAddr   = msg.sender;
        blockNumLastUpdate = block.number;

        allInitStatesHash         = _allInitStatesHash;
        finalTeamStatesHash       = _finalTeamStatesHash;
        allGamePlayResultsHash    = _allGamePlayResultsHash;
        leagueIdx                 = _leagueIdx; //*//
    }
    
    function challengeUpdate(
        uint256              selectedTeam,
        byte[90]   calldata _allGamePlayResults,
        bytes32[10] calldata initStatesHash,  // not needed. You need the data, not the hashes (see initPlayersStates)
        bytes32[80] calldata initPlayersStates // 10 teams, 8 players per team
        // uint[110] calldata initPlayersStates // 10 teams, 11 players per team //*//
        uint leagueIdx; //*//
    ) external {
        stakers.verify(msg.sender);
        
        // require that _allGamePlayResults for already verified teams are OK
        for every team for which wasTeamVerified==1:
            check that allGamePlayResults[that team] =  _allGamePlayResults[that team]

        // require this team was not processed by BChain already //*//
        require(wasUpdatedByBLockchain[selectedTeam] == false);

        // verify allInitStatesHash  //*//
        bool isInitStateVerifiedByBC = wasInitStateVerified == 1;
        bool initStatesDiffer = keccak256(initPlayersStates) != allInitStatesHash;

        if (isInitStateVerifiedByBC && initStatesDiffer) { assert(false); }
        if (!isInitStateVerifiedByBC && initStatesDiffer) {
            // see who is lying:
            // checks that init states provided are OK for that league
            bool initStatesAreGood = verifyInitStates(initPlayersStates, leagueIdx);
            if(!initStatesAreGood) {
                stackers.slash(msg.sender);
                return;
            };
            // we proved that updater lied, and the BC verified that provided initStates are Good.
            stackers.slash(updaterAddr);
            allInitStatesHash = keccak256(initPlayersStates);
            updaterAddr   = msg.sender;
            blockNumLastUpdate = block.number;  //*// reopen for new challenges
            wasInitStateVerified = true;
            return;
        }

        // you get only if: initState was already verified, and coinciding with provided
        
        // verify allGamePlayResults
        bytes32 teamStateHash, byte[18] selectedTeamResults = playLeagueForTeam(selectedTeam, initPlayersStates); //*//
        // did previous updater lie in results?
        bool resultsDiffer = resultsDiffer(selectedTeamResults, selectedTeam, allGamePlayResults); // compares the 18 results

        // did previous updater lie in results?
        bool resultsDiffer = resultsDiffer(selectedTeamResults, selectedTeam, allGamePlayResults); // compares the 18 results

        // did previous updater lie in final states?
        bool finalHashDiffers = teamStateHash != _finalHash[selectedTeam];

        if (finalHashDiffers || resultsDiffer) {
            stackers.slash(updaterAddr);
            _finalTeamStateHash[selectedTeam] = teamStateHash; 
            allGamePlayResults = _allGamePlayResults; 
            blockNumLastUpdate = block.number;  //*// reopen for new challenges
            updaterAddr   = msg.sender;         //*// to be slashed by someone later
            wasUpdatedByBLockchain[selectedTeam] = true; //*// to avoid later dummy re-challenges
        }
    }
    

    function finishChallange() public {
        require(blockNumLastUpdate>0); // challange started
        require(block.number - blockNumLastUpdate > 80); // 80 blocks for challanging

        startNewGame();
    }
}