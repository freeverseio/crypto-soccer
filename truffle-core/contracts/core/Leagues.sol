pragma solidity ^0.5.0;

import "./LeagueUpdater.sol";

contract Leagues is LeagueUpdater {
    constructor(address engine, address leagueState) public LeagueUpdater(engine, leagueState) {
    }    
}