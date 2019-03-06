pragma solidity ^0.5.0;

import "./LeaguesComputer.sol";

contract LeagueUpdater is LeaguesComputer {
    constructor(address engine, address leagueState) public LeaguesComputer(engine, leagueState) {
    } 
}