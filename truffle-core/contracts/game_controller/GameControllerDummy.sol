pragma solidity ^ 0.5.0;

import "./GameControllerInterface.sol";

contract GameControllerDummy is GameControllerInterface {
  event StakersUpdated(uint256 id, uint256 windowStart, address updater);
  event StakersChallenged(uint256 id);

  function updated(uint256 _id, uint256 _windowStart, address _updater) external {
    emit StakersUpdated(_id, _windowStart, _updater);
  }

  function challenged(uint256 _id) external {
    emit StakersChallenged(_id);
  }
}