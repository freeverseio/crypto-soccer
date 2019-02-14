pragma solidity ^0.4.25;

import "./LeaguesComputer.sol";

contract Leagues is LeaguesComputer {
    constructor(address engine) public LeaguesComputer(engine) {
    }    
}