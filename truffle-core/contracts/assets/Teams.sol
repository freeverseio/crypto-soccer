pragma solidity >=0.4.21 <0.6.0;

import "./Players.sol";

contract Teams is Players {
    constructor(address playerState) public Players(playerState) {
    }
}