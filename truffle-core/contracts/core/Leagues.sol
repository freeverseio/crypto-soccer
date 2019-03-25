pragma solidity ^0.5.0;

import "./LeagueChallengeable.sol";
import "./LeaguesComputer.sol";

contract Leagues is LeagueChallengeable, LeaguesComputer {
    constructor(address engine, address state) 
    public 
    LeaguesComputer(engine, state)
    {
    }    
}