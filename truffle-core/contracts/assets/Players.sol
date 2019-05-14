pragma solidity >=0.4.21 <0.6.0;

import "./Storage.sol";

contract Players is Storage {
    constructor(address playerState) public Storage(playerState) {
    }
  }