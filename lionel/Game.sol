contract Game {
    // assume nTeams = 10  => nGames = 90
    Stakers     stakers;
    
    address     firstUpdaterAddr;
    uint256     blockNumLastUpdate;

    bytes32     allInitStatesHash;      // hash of bytes32[10]
    bytes32[10] finalStatesHash;        // 10 teams
    bytes32     allGamePlayResultsHash; // hash of byte[90], 10 teams => 90 total plays 

    function challengeUpdateForFirstTime(
        bytes32              _allInitStatesHash,
        bytes32[10] calldata _finalStatesHash,
        byte[90] calldata    _allGamePlayResultsHash
    ) external {
        stakers.verify(msg.sender);     // valid actor
        require(blockNumLastUpdate==0x0); // not challanged yet

        firstUpdaterAddr   = msg.sender;
        blockNumLastUpdate = block.number;

        allInitStatesHash         = _allInitStatesHash;
        finalStatesHash           = _finalStatesHash;
        allGamePlayResultsHash    = _allGamePlayResultsHash;
    }
    
    function challengeUpdate(
        uint256              selectedTeam,
        bytes[90]   calldata allGamePlayResults,
        bytes32[10] calldata initStatesHash,
        bytes32[80] calldata initPlayersStates // 10 teams, 8 players per team
    ) external {
        stakers.verify(msg.sender);
        
        // verify allGamePlayResultsHash
        require(keccak256(allGamePlayResults)==allGamePlayResultsHash);

        // verify allInitStatesHash
        require(keccak256(initStatesHash)==allInitStatesHash);

        // verify allGamePlayResults
        bytes32 finalHashForTeam = playGame(selectedTeam,initStatesHash,initPlayersStates);
        if (finalHashForTeam != _finalHash[selectedTeam]) {
            stackers.slash(firstUpdaterAddr);
            blockNumLastUpdate = 0x0;
        }
    }
    
    function finishChallange() public {
        require(blockNumLastUpdate>0); // challange started
        require(block.number - blockNumLastUpdate > 80); // 80 blocks for challanging

        startNewGame();
    }
}