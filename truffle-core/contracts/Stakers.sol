pragma solidity >= 0.6.3;

// TODO: leaving for later how the monthly grant is going to be computed/shared among L1 updaters

contract Owned {
  address public owner = address(0x0);

  constructor() public {
    owner = msg.sender;
  }
  function setOwner(address _address) external virtual onlyOwner {
    owner = _address;
  }
  modifier onlyOwner {
    require( msg.sender == owner, "Only owner can call this function.");
        _;
  }
}

contract Rewards is Owned{

  uint balanceToShare = 0;
  uint totalNumUpdates = 0;
  mapping (address => uint) pendingWithdrawals;
  mapping (address => uint) howManyUpdates;
  address [] whoToBeRewarded;

  function addReward() external payable {
    require (msg.value > 0, "failed to add reward of zero");
    balanceToShare += msg.value;
  }

  function execute() external onlyOwner {
    require (whoToBeRewarded.length > 0, "failed to execute rewards: empty array");
    require (balanceToShare >= whoToBeRewarded.length, "failed to execute rewards: Not enough balance to share");
    uint amount = balanceToShare / totalNumUpdates;
    for (uint i=0; i<whoToBeRewarded.length; i++) {
    }
    for (uint i=0; i<whoToBeRewarded.length; i++) {
      address who = whoToBeRewarded[i];
      pendingWithdrawals[who] = amount*howManyUpdates[who];
      howManyUpdates[who] = 0;
    }
    delete whoToBeRewarded;
    balanceToShare = 0;
    totalNumUpdates = 0;
  }

  function withdraw(address payable _addr) public onlyOwner {
    _addr.transfer(pendingWithdrawals[_addr]);
    pendingWithdrawals[_addr] = 0;
  }

  function push(address _addr) external onlyOwner {
    require (_addr != address(0x0));
    if (howManyUpdates[_addr] == 0)
    {
      whoToBeRewarded.push(_addr);
    }
    howManyUpdates[_addr] += 1;
    totalNumUpdates++;
  }
}

contract AddressStack is Owned{
  uint16 public constant capacity = 4;
  uint16 public length = 0;
  address[capacity] private array;

  /// @notice adds a new element. Reverts in case the element is found in array or array is full
  function push(address _address) external onlyOwner {
    require (length < capacity, "cannot push to a full AddressStack");
    require (!contains(_address), "cannot push, address is already in AddressStack");
    array[length++] = _address;
  }

  /// notice: removes the last element that was pushed. Reverts in case it is empty.
  /// returns the element that has been removed from the array
  function pop() external  onlyOwner returns (address _address) {
    require (length > 0, "cannot pop from an empty AddressStack");
    _address = array[--length];
  }

  function contains(address _address) public view returns (bool) {
    for (uint16 i=0; i<length; i++) {
      if (array[i] == _address) {
        return true;
      }
    }
    return false;
  }
}

contract AddressMapping is Owned{

  mapping(address => bool) public data;

  function add(address _address) external onlyOwner {
    require (!has(_address), "failed to add to AddressMapping, address already exists");
    data[_address] = true;
  }
  function has(address _address) public view returns (bool) {
    return data[_address] == true;
  }
}


contract Stakers is Owned {

  uint16 public constant NUM_STAKERS = 32;

  uint public requiredStake;
  address public gameOwner = address(0x0);
  AddressStack private updaters = new AddressStack();
  Rewards public rewards = new Rewards();
  address[NUM_STAKERS] public stakers;
  AddressMapping public slashed = new AddressMapping();
  AddressMapping public trustedParties = new AddressMapping();

  mapping (address => uint) stakes;


  // ----------------- modifiers -----------------------

  modifier onlyGame {
    require(msg.sender == gameOwner && gameOwner != address(0x0),
            "Only gameOwner can call this function.");
    _;
  }

  // ----------------- public functions -----------------------

  constructor(uint stake) public {
    requiredStake =  stake;
  }

  function setOwner(address _address) external override (Owned) onlyOwner {
    owner = _address;
    updaters.setOwner(_address);
    rewards.setOwner(_address);
    slashed.setOwner(_address);
    trustedParties.setOwner(_address);
  }

  /// @notice adds amount to rewards contract
  function addReward() external payable {
    rewards.addReward.value(msg.value)();
  }

  /// @notice executes rewards
  function executeReward() external {
    rewards.execute();
  }

  /// @notice transfers earnings to the calling staker
  function withdraw() external {
    require (isStaker(msg.sender), "failed to withdraw: staker not registered");
    rewards.withdraw(msg.sender);
    if (stakes[msg.sender] > requiredStake)
    {
      uint amount = stakes[msg.sender] - requiredStake;
      stakes[msg.sender] = requiredStake;
      msg.sender.transfer(amount);
    }
  }

  /// @notice sets the address of the gameOwner that interacts with this contract
  function setGameOwner(address _address) external onlyOwner {
    require (gameOwner == address(0x0),     "gameOwner is already set");
    require (_address != address(0x0), "invalid address 0x0");
    gameOwner = _address;
  }

  /// @notice adds address as trusted party
  function addTrustedParty(address _staker) external onlyOwner {
    trustedParties.add(_staker);
  }

  /// @notice registers a new staker
  function enroll() external payable {
    require (msg.value == requiredStake, "failed to enroll: wrong stake amount");
    require (!isSlashed(msg.sender),      "failed to enroll: staker was slashed");
    require (isTrustedParty(msg.sender),  "failed to enroll: staker is not trusted party");
    require (addStaker(msg.sender),       "failed to enroll");
    stakes[msg.sender] = msg.value;

  }

  /// @notice unregisters a new staker and transfers all earnings
  function unEnroll() external {
    require (!updaters.contains(msg.sender), "failed to unenroll: staker currently updating");
    require (removeStaker(msg.sender),       "failed to unenroll");
    rewards.withdraw(msg.sender);
    uint amount = stakes[msg.sender];
    stakes[msg.sender]  = 0;
    msg.sender.transfer(amount);
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
      address badStaker = updaters.pop();
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
      rewards.push(updaters.pop());
    }
    require (level() == 0, "failed to finalize: no updaters should have been left");
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

  // WARNING: careful with this list ever-growing and not being able to check in one TX.
  // TODO: create a mapping isSlashed(address => bool) instead of an array.
  function contains(address[] storage _array, address _value) private view returns (bool) {
    for (uint i=0; i<_array.length; i++) {
      if (_array[i] == _value) {
        return true;
      }
    }
    return false;
  }

  function isTrustedParty(address _addr) private view returns (bool) {
    return trustedParties.has(_addr);
  }

  function isSlashed(address _addr) private view returns (bool) {
    return slashed.has(_addr);
  }

  function isStaker(address _addr) private view returns (bool) {
    for (uint16 i=0; i<NUM_STAKERS; i++) {
      if (stakers[i] == _addr) {
        return true;
      }
    }
    return false;
  }

  function addStaker(address _staker) private returns (bool) {
    for (uint16 i = 0; i<NUM_STAKERS; i++){
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
    while (stakerIndex < NUM_STAKERS) {
      if (stakers[stakerIndex] == _staker) {
        break;
      }
      ++stakerIndex;
    }

    if (stakerIndex < NUM_STAKERS) {
      // remove gaps
      for (uint16 i = stakerIndex; i<NUM_STAKERS-1; i++){
       stakers[i] = stakers[i+1];
      }
      stakers[NUM_STAKERS-1] = address(0x0);
      return true;
    }
    return false;
  }

  function resolve() private {
    address goodStaker = updaters.pop();
    address badStaker = updaters.pop();
    earnStake(goodStaker,badStaker);
    slash(badStaker);
  }

  function slash(address _staker) private {
    require (removeStaker(_staker), "failed to slash: staker not found");
    slashed.add(_staker);
  }

  function earnStake(address _goodStaker, address _badStaker) private {
    require (stakes[_badStaker] > 0);
    uint amount = stakes[_badStaker];
    stakes[_badStaker] = 0;
    stakes[_goodStaker] += amount;
    // TODO: alternatively it has been proposed to burn stake, and reward true tellers with the monthly pool.
    // The idea behind it, is not to promote interest in stealing someone else's stake
    // address(0x0).transfer(requiredStake); // burn stake
  }
}
