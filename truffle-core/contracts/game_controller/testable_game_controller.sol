pragma solidity ^ 0.5.0;

import "./game_controller.sol";

contract TestableGameController is GameController {
  function computeKeccak256(string calldata s, uint n1) external pure returns(uint) {
      return uint(keccak256(abi.encodePacked(s, n1)));
  }
}
