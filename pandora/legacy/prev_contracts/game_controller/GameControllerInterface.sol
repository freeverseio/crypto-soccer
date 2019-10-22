pragma solidity >=0.4.21 <0.6.0;

/// @dev bridge between Game and Stakers contract
interface GameControllerInterface {
  /// @dev called by game logic when league id is updated by a staker
  /// @param _id game identifier
  /// @param _windowStart should be the block number at which the league ended
  function updated(uint256 _id, uint256 _windowStart, address _updater) external;

  /// @dev called by game logic when a challenger succesfully
  /// demonstrates that the updater was lying. Typically _windowStart should be
  /// the block number at which the league ended, however it should be reset when
  /// an updater lies.
  function challenged(uint256 _id) external;
}