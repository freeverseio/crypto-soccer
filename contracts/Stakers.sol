pragma solidity ^ 0.5.0;

// TODO: leaving for later how the monthly grant is going to be computed/shared among L1 updaters

contract Rewards {

  address public owner = address(0x0);
  address payable[] public toBeRewarded;

  modifier onlyOwner {
    require(msg.sender == owner,
            "Only owner can call this function.");
    _;
  }

  constructor() public {
    owner = msg.sender;
  }

  function() external payable { }

  function execute() external onlyOwner {
    require (toBeRewarded.length != 0, "failed to execute reward: empty array");
    uint amount = address(this).balance / toBeRewarded.length;
    for (uint i=0; i<toBeRewarded.length; i++) {
      toBeRewarded[i].transfer(amount);
    }
    delete toBeRewarded;
  }

  function push(address _addr) external onlyOwner {
    toBeRewarded.push(address(uint160(_addr)));
  }
}

contract AddressStack {
  uint16 public constant capacity = 4;
  uint16 public length = 0;
  address[capacity] private array;

  /// @notice adds a new element. Reverts in case the element is found in array or array is full
  function push(address _address) external {
    require (length < capacity, "cannot push to a full AddressStack");
    require (!contains(_address), "cannot push, address is already in AddressStack");
    array[length++] = _address;
  }

  /// @notice removes the last element that was pushed. Reverts in case it is empty.
  /// @return the element that has been removed from the array
  function pop() external returns (address _address) {
    require (length > 0, "cannot pop from an empty AddressStack");
    _address = array[--length];
  }

  //function clear() public {
  //  length = 0;
  //}

  function contains(address _address) public view returns (bool) {
    for (uint16 i=0; i<length; i++) {
      if (array[i] == _address) {
        return true;
      }
    }
    return false;
  }
}

contract Stakers {

  uint16 public constant kNumStakers = 32;
  uint public constant kRequiredStake = 4 ether;

  address public owner = address(0x0);
  address public game = address(0x0);
  AddressStack private updaters = new AddressStack();
  Rewards public rewards = new Rewards();
  address[kNumStakers] public stakers;
  address[] public slashed;
  address[] public trustedParties;


  // ----------------- modifiers -----------------------

  modifier onlyOwner {
    require(msg.sender == owner,
            "Only owner can call this function.");
    _;
  }
  modifier onlyGame {
    require(msg.sender == game && game != address(0x0),
            "Only game can call this function.");
    _;
  }

  // ----------------- public functions -----------------------

  constructor() public {
    owner = msg.sender;
  }

  /// @notice adds amount to rewards contract
  function addReward() external payable {
    address(rewards).transfer(msg.value);
  }

  /// @notice executes rewards
  function executeReward() external {
    rewards.execute();
  }

  /// @notice sets the address of the game that interacts with this contract
  function setGame(address _address) external onlyOwner {
    require (game == address(0x0),     "game is already set");
    require (_address != address(0x0), "invalid address 0x0");
    game = _address;
  }

  /// @notice adds address as trusted party
  function addTrustedParty(address _staker) external onlyOwner {
    require (!isTrustedParty(_staker), "failed to add trusted party");
    trustedParties.push(_staker);
  }

  /// @notice registers a new staker
  function enroll() external payable {
    require (msg.value == kRequiredStake, "failed to enroll: wrong stake amount");
    require (!isSlashed(msg.sender),      "failed to enroll: staker was slashed");
    require (isTrustedParty(msg.sender),  "failed to enroll: staker is not trusted party");
    require (addStaker(msg.sender),       "failed to enroll");

  }

  /// @notice unregisters a new staker
  function unEnroll() external {
    require (!updaters.contains(msg.sender), "failed to unenroll: staker currently updating");
    require (removeStaker(msg.sender),       "failed to unenroll");
    msg.sender.transfer(kRequiredStake);
  }

  /// @notice update to a new level
  /// @param _level to which update
  /// @param _staker address of the staker that reports this update
  /// @dev This function will also resolve previous updates when
  //       level is below current or level has reached the end
  function update(uint16 _level, address _staker) external onlyGame {
    require (_level <= level(),        "failed to update: wrong level");
    //require (_level <= maxNumLevels(), "failed to update: max level exceeded"); // already covered by previous require
    require (isStaker(_staker),        "failed to update: staker not registered");
    //require (!isSlashed(_staker),      "failed to update: staker was slashed"); // also covered by not being part of stakers, because slashing removes address from stakers

    if (_level < maxNumLevels()) {
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
      slash(updaters.pop());
      earnStake(_staker);
    }
  }

  /// @notice start new verse
  /// @dev current state will be resolved at this point.
  /// If called from level 1, then staker is rewardeds.
  /// When called from any other level, means that everys
  /// other staker told the truth but the one in between
  /// told a lie.
  function start() external onlyGame {
    require (level() > 0, "failed to start: wrong level");
    while (level() > 1) {
      resolve();
    }
    if (level() == 1) {
      rewards.push(updaters.pop());
    }
    require (level() == 0, "failed to start: no updaters should have been left");
  }

  /// @notice get the current level
  function level() public view returns (uint) {
    return updaters.length();
  }

  /// @notice get the maximum level
  function maxNumLevels() public view returns (uint) {
    return updaters.capacity();
  }

  // ----------------- private functions -----------------------

  function contains(address[] storage _array, address _value) private view returns (bool) {
    for (uint i=0; i<_array.length; i++) {
      if (_array[i] == _value) {
        return true;
      }
    }
    return false;
  }

  function isTrustedParty(address _addr) private view returns (bool) {
    return contains(trustedParties, _addr);
  }

  function isSlashed(address _addr) private view returns (bool) {
    return contains(slashed, _addr);
  }

  function isStaker(address _addr) private view returns (bool) {
    for (uint16 i=0; i<kNumStakers; i++) {
      if (stakers[i] == _addr) {
        return true;
      }
    }
    return false;
  }

  function addStaker(address _staker) private returns (bool) {
    for (uint16 i = 0; i<kNumStakers; i++){
      if (stakers[i] == _staker) {
        // staker already registered
        return false;
      }
      if (stakers[i] == address(0x0)) {
        stakers[i] = _staker;
        return true;
      }
    }
    return false;
  }

  function removeStaker(address _staker) private returns (bool){
    // find index of staker
    uint16 stakerIndex = 0;
    while (stakerIndex < kNumStakers) {
      if (stakers[stakerIndex] == _staker) {
        break;
      }
      ++stakerIndex;
    }

    if (stakerIndex < kNumStakers) {
      // remove gaps
      for (uint16 i = stakerIndex; i<kNumStakers-1; i++){
       stakers[i] = stakers[i+1];
      }
      stakers[kNumStakers-1] = address(0x0);
      return true;
    }
    return false;
  }

  function resolve() private {
    earnStake(updaters.pop());
    slash(updaters.pop());
  }

  function slash(address _staker) private {
    require (removeStaker(_staker), "failed to slash: staker not found");
    slashed.push(_staker);
  }

  function earnStake(address _addr) private {
    address(uint160(_addr)).transfer(kRequiredStake);
    // TODO: alternatively it has been proposed to burn stake, and reward true tellers with the monthly pool.
    // The idea behind it, is not to promote interest in stealing someone else's stake
    // address(0x0).transfer(kRequiredStake); // burn stake
  }
}
