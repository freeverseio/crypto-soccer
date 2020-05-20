pragma solidity >= 0.6.3;

import "@openzeppelin/contracts/math/SafeMath.sol";

// TODO: leaving for later how the monthly grant is going to be computed/shared among L1 updaters

contract Stakers {
  using SafeMath for uint256;
  
  address constant internal NULL_ADDR = address(0x0);
  uint16 public constant updatersCapacity = 4;

  address public owner;
  address public gameOwner;

  mapping (address => bool) public isStaker;
  mapping (address => bool) public isSlashed;
  mapping (address => bool) public isTrustedParty;

  mapping (address => uint256) public stakes;
  mapping (address => uint256) public pendingWithdrawals;
  mapping (address => uint256) public howManyUpdates;
  
  uint256 public nStakers;  
  uint256 public requiredStake;
  uint256 public potBalance;
  uint256 public totalNumUpdates;
  address [] public toBeRewarded;
  address[] public updaters;

  modifier onlyOwner {
    require( msg.sender == owner, "Only owner can call this function.");
        _;
  }

  modifier onlyGame {
    require(msg.sender == gameOwner && gameOwner != NULL_ADDR,
            "Only gameOwner can call this function.");
    _;
  }

  constructor(uint256 stake) public {
    requiredStake = stake;
    owner = msg.sender;
  }
  
  function setOwner(address _address) external onlyOwner {
    owner = _address;
  }

  /// @notice sets the address of the external contract that interacts with this contract
  function setGameOwner(address _address) external onlyOwner {
    require (_address != NULL_ADDR, "invalid address 0x0");
    gameOwner = _address;
  }
  

  /// @notice executes rewards
  function executeReward() external onlyOwner {
    require (toBeRewarded.length > 0, "failed to execute rewards: empty array");
    require (potBalance >= toBeRewarded.length, "failed to execute rewards: Not enough balance to share");
    for (uint256 i = 0; i < toBeRewarded.length; i++) {
      address who = toBeRewarded[i];
      // better to multiply, and then divide, each time, to minimize rounding errors.
      pendingWithdrawals[who] += (potBalance * howManyUpdates[who]) / totalNumUpdates;
      howManyUpdates[who] = 0;
    }
    delete toBeRewarded;
    potBalance = 0; // there could be a negligible loss of funds in the Pot.
    totalNumUpdates = 0;
  }  

  /// @notice transfers pendingWithdrawals to the calling staker; the stake remains until unenrol is called
  function withdraw() external {
    // no need to require (isStaker[msg.sender], "failed to withdraw: staker not registered");
    uint256 amount = pendingWithdrawals[msg.sender];
    require(amount > 0, "nothing to withdraw by this msg.sender");
    pendingWithdrawals[msg.sender] = 0;
    msg.sender.transfer(amount);
  }

  function assertGoodCandidate(address _addr) public view {
    require(_addr != NULL_ADDR, "candidate is null addr");
    require(!isSlashed[_addr], "candidate was slashed previously");
    require(stakes[_addr] == 0, "candidate already has a stake");
  }

  /// @notice adds address as trusted party
  function addTrustedParty(address _staker) external onlyOwner {
    assertGoodCandidate(msg.sender);
    require(!isTrustedParty[_staker], "trying to add a trusted party that is already trusted");
    isTrustedParty[_staker] = true;
  }

  /// @notice registers a new staker
  function enroll() external payable {
    assertGoodCandidate(msg.sender);
    require (msg.value == requiredStake, "failed to enroll: wrong stake amount");
    require (isTrustedParty[msg.sender], "failed to enroll: staker is not trusted party");
    require (addStaker(msg.sender), "failed to enroll: cannot add staker");
    stakes[msg.sender] = msg.value;
  }

  /// @notice unregisters a new staker and transfers all earnings, and pot
  function unEnroll() external {
    require (!alreadyDidUpdate(msg.sender), "failed to unenroll: staker currently updating");
    require (removeStaker(msg.sender), "failed to unenroll");
    uint256 amount = pendingWithdrawals[msg.sender] + stakes[msg.sender];
    pendingWithdrawals[msg.sender] = 0;
    stakes[msg.sender]  = 0;
    if (amount > 0) { msg.sender.transfer(amount); }
  }

  /// @notice update to a new level
  /// @param _level to which update
  /// @param _staker address of the staker that reports this update
  /// @dev This function will also resolve previous updates when
  //       level is below current or level has reached the end
  function update(uint16 _level, address _staker) external onlyGame {
    require (_level <= level(),        "failed to update: wrong level");
    //require (_level <= updatersCapacity, "failed to update: max level exceeded"); // already covered by previous require
    require (isStaker[_staker],        "failed to update: staker not registered");
    //require (!isSlashed(_staker),      "failed to update: staker was slashed"); // also covered by not being part of stakers, because slashing removes address from stakers

    if (_level < updatersCapacity) {
      if (_level < level()) {
        // If level is below current, it means the challenge
        // period has passed, so last updater told the truth.
        // The last updater should be rewarded, the one before
        // last should be slashed and level moves back two positions
        require (_level > 0 && _level == level() - 2, "failed to update: resolving wrong level");
        resolve();
      }
      updaters.push(_staker);
    }
    else {
      // The very last possible update: the challenge.
      // It resolves immediately by slashing the last
      // updater
      address badStaker = popUpdaters();
      slash(badStaker);
      earnStake(_staker, badStaker);
    }
  }

  /// @notice finalize current game, get ready for next one.
  /// @dev current state will be resolved at this point.
  /// If called from level 1, then staker is rewarded.
  /// When called from any other level, means that every
  /// other staker told the truth but the one in between
  /// lied.
  function finalize() external onlyGame {
    require (level() > 0, "failed to finalize: wrong level");
    while (level() > 1) {
      resolve();
    }
    if (level() == 1) {
      addRewardToUpdater(popUpdaters());
    }
    require (level() == 0, "failed to finalize: no updaters should have been left");
  }

  /// @notice get the current level
  function level() public view returns (uint256) {
    return updaters.length;
  }

  function addStaker(address _staker) private returns (bool) {
    if (_staker == NULL_ADDR) return false; // prevent null addr
    if (isStaker[_staker]) return false; // staker already registered
    isStaker[_staker] = true;
    nStakers++;
    return true;
  }

  function removeStaker(address _staker) private returns (bool){
    if (_staker == NULL_ADDR) return false; // prevent null addr
    if (!isStaker[_staker]) return false; // staker not registered
    isStaker[_staker] = false;
    nStakers--;
    return true;
  }

  function resolve() private {
    address goodStaker = popUpdaters();
    address badStaker = popUpdaters();
    earnStake(goodStaker, badStaker);
    slash(badStaker);
  }

  function slash(address _staker) private {
    require (removeStaker(_staker), "failed to slash: staker not found");
    isSlashed[_staker] = true;
  }

  // the slashed stake goes into the "pendingWithdrawals" of the good staker,
  // not to his "stake". This way, he can cash it without unenrolling.
  function earnStake(address _goodStaker, address _badStaker) private {
    uint256 amount = stakes[_badStaker];
    stakes[_badStaker] = 0;
    pendingWithdrawals[_goodStaker] += amount;
    // TODO: alternatively it has been proposed to burn stake, and reward true tellers with the monthly pool.
    // The idea behind it, is not to promote interest in stealing someone else's stake
    // NULL_ADDR.transfer(requiredStake); // burn stake
  }
  
  function addRewardToPot() external payable {
    require (msg.value > 0, "failed to add reward of zero");
    potBalance += msg.value;
  }

  function addRewardToUpdater(address _addr) private {
    if (howManyUpdates[_addr] == 0) {
      toBeRewarded.push(_addr);
    }
    howManyUpdates[_addr] += 1;
    totalNumUpdates++;
  }

  function popUpdaters() private returns (address _address) {
    uint256 updatersLength = updaters.length;
    require (updatersLength > 0, "cannot pop from an empty AddressStack");
    _address = updaters[updatersLength - 1];
    updaters.pop();
  }

  // this function iterates over a storage array, but of max length 4.
  function alreadyDidUpdate(address _address) public view returns (bool) {
    for (uint256 i = 0; i < updaters.length; i++) {
      if (updaters[i] == _address) {
        return true;
      }
    }
    return false;
  }
}

