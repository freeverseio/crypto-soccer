pragma solidity ^0.5.0;

import "./LeagueChallengeable.sol";
import "./LeaguesComputer.sol";

contract Leagues is LeagueChallengeable {
    constructor(address engine, address state) 
    public 
    LeagueChallengeable(engine, state)
    {
    }    
}