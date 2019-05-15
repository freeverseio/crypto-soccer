pragma solidity ^ 0.5.0;

import "../stakers/StakersInterface.sol";

// TODOs:
// 1. Restricted window updates is still missing
// 2. Test restricted window update


contract GameController {

  string constant ERR_NO_STAKERS         = "err-no-stakers-contract-set";
  string constant ERR_WINDOW_NOT_STARTED = "err-window-not-started";
  string constant ERR_WINDOW_FINISHED    = "err-window-finished";
  string constant ERR_WINDOW_RESTRICTED  = "err-window-restricted";

  address constant kNullAddress = address(0x0);
  uint16 public constant kWindowBlocks = 100;
  uint16 public constant kWindowBlocksRestricted = 66;

  address owner;
  address stakersContractAddress;

  event UpdateEvent(uint256 id, address staker);
  event ChallengeEvent(uint256 id, address staker);

  // ----------------- modifiers -----------------------

  modifier onlyOwner {
    require(msg.sender == owner,
            "Only owner can call this function.");
    _;
  }

  modifier onlyIfStakersAddressValid {
    require(stakersContractAddress != kNullAddress, ERR_NO_STAKERS);
    _;
  }

  // ----------------- public functions -----------------------

  constructor() public {
    owner = msg.sender;
  }

  function setStakersContractAddress(address _address) public onlyOwner {
    require (_address != kNullAddress);
    stakersContractAddress = _address;
  }

  function getStakersContractAddress() public view returns (address) {
    return stakersContractAddress;
  }

  // ----------------- internal/protected functions -----------------------

  /// @dev called thru inheritance by game logic when league id is updated by a staker
  /// @param _id game identifier
  /// @param _windowStart should be the block number at which the league ended
  function updated(uint256 _id, uint256 _windowStart, address _updater) public onlyIfStakersAddressValid {
    checkUpdateWindow(_windowStart, _updater);
    StakersInterface(stakersContractAddress).initChallenge(_updater);
    emit UpdateEvent(_id, _updater);
  }

  /// @dev called thru inheritance by game logic when a challenger succesfully
  /// demonstrates that the updater was lying. Typically _windowStart should be
  /// the block number at which the league ended, however it should be reset when
  /// an updater lies.
  function challenged(uint256 _id, address _updater) public onlyIfStakersAddressValid {
    StakersInterface(stakersContractAddress).lierChallenge(_updater); // will revert if _updater was not in challengable state
    emit ChallengeEvent(_id, _updater);
  }

  // ----------------- private functions -----------------------

  /// @dev checks whether the requested update can proceed according to stakers contract and update window periods
  function checkUpdateWindow(uint256 _windowStart, address _staker) private view onlyIfStakersAddressValid returns (bool) {
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
