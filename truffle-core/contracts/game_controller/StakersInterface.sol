pragma solidity >=0.4.21 <0.6.0;

/// @dev bridge between Game and Stakers contract
interface StakersInterface {
  function initChallenge(address _staker) external;
  function lierChallenge(address _staker) external;
}

