pragma solidity ^0.5.0;

import "./LeaguesComputer.sol";

contract Leagues is LeaguesComputer {
    constructor(address engine) public LeaguesComputer(engine) {
    }    
}