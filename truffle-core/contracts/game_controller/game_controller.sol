pragma solidity ^ 0.5.0;

import "./GameControllerInterface.sol";
import "./Stakers.sol";

// TODOs:
// 1. Restricted window updates is still missing
// 2. Test restricted window update


contract GameController is GameControllerInterface, Stakers {

  string constant ERR_NO_STAKERS         = "err-no-stakers-contract-set";
  string constant ERR_WINDOW_NOT_STARTED = "err-window-not-started";
  string constant ERR_WINDOW_FINISHED    = "err-window-finished";
  string constant ERR_WINDOW_RESTRICTED  = "err-window-restricted";

  address constant kNullAddress = address(0x0);
  uint16 public constant kWindowBlocks = 100;
  uint16 public constant kWindowBlocksRestricted = 66;

  address owner;

  mapping (uint256 => address) public id2staker;

  event UpdateEvent(uint256 id, address staker);
  event ChallengeEvent(uint256 id, address staker);

  // ----------------- modifiers -----------------------

  modifier onlyOwner {
    require(msg.sender == owner,
            "Only owner can call this function.");
    _;
  }

  // ----------------- public functions -----------------------
  constructor(address _game) public Stakers(_game) {
    owner = msg.sender;
  }

  // ----------------- internal/protected functions -----------------------
  function updated(uint256 _id, uint256 _windowStart, address _updater) external {
    // checkUpdateWindow(_windowStart, _updater);
    initChallenge(_updater);
    id2staker[_id] = _updater;
    emit UpdateEvent(_id, _updater);
  }

  function challenged(uint256 _id) external {
    lierChallenge(id2staker[_id]); // will revert if _updater was not in challengable state
    emit ChallengeEvent(_id, id2staker[_id]);
  }

  // ----------------- private functions -----------------------

  /// @dev checks whether the requested update can proceed according to stakers contract and update window periods
  function checkUpdateWindow(uint256 _windowStart, address _staker) private view returns (bool) {
    uint256 windowEnd = _windowStart + kWindowBlocks;
    require (block.number > _windowStart, ERR_WINDOW_NOT_STARTED);
    require (block.number < windowEnd, ERR_WINDOW_FINISHED);

    uint256 windowEveryone = _windowStart + kWindowBlocksRestricted;
    if (block.number > windowEveryone) {
      return true;
    }
    return restrictedWindowAvailable(_staker);
  }

  /// @dev whether staker is allowed to participate in the restricted period
  function restrictedWindowAvailable(address /*_staker*/) private pure returns (bool) {
    // TODO
    return true;
  }
}
