pragma solidity ^0.5.0;

import "./LeagueUpdatable.sol";
import "./LeaguesComputer.sol";

contract Leagues is LeagueUpdatable, LeaguesComputer {
    constructor(address engine, address state) 
    public 
    LeaguesComputer(engine, state)
    {
    }    
}