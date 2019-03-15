pragma solidity ^0.5.0;

import "./LeagueUpdatable.sol";
import "./LeaguesComputer.sol";
import "./LeagueUsersAlongData.sol";

contract Leagues is LeagueUpdatable, LeaguesComputer, LeagueUsersAlongData {
    constructor(address engine, address state) 
    public 
    LeaguesComputer(engine, state)
    {
    }    
}