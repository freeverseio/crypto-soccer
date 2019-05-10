pragma solidity ^ 0.5.0;

import "./game_controller.sol";

contract TestableGameController is GameController {

  function test_updated(uint256 id, uint256 leagueStartBlockNumber, address staker) public {
    updated(id, leagueStartBlockNumber, staker);
  }
  function test_challenged(uint256 id, address staker) public {
    challenged(id, staker);
  }

  function computeKeccak256(string calldata s, uint n1) external pure returns(uint) {
      return uint(keccak256(abi.encodePacked(s, n1)));
  }
}
