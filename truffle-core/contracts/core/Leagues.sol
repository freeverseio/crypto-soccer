pragma solidity ^0.5.0;

import "./LeagueUpdater.sol";
import "./LeaguesComputer.sol";

contract Leagues is LeagueUpdater, LeaguesComputer {
    constructor(address engine, address state) 
    public 
    LeaguesComputer(engine, state)
    LeagueUpdater(state) 
    {
    }    
}