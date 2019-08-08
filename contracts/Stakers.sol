pragma solidity ^ 0.5.0;

// TODO: leaving for later how the monthly grant is going to be computed/shared among L1 updaters

contract AddressStack {
  uint16 public constant capacity = 4;
  uint16 public length = 0;
  address[capacity] private array;

  /// @notice adds a new element. Reverts in case the element is found in array or array is full
  function push(address _address) public {
    require (length < capacity, "cannot push to a full AddressStack");
    require (!contains(_address), "cannot push, address is already in AddressStack");
    array[length++] = _address;
  }

  /// @notice removes the last element that was pushed. Reverts in case it is empty.
  /// @return the element that has been removed from the array
  function pop() public returns (address _address) {
    require (length > 0, "cannot pop from an empty AddressStack");
    _address = array[--length];
  }

  function clear() public {
    length = 0;
  }

  function contains(address _address) private view returns (bool) {
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
  address[kNumStakers] public stakers;
  address[] public slashed;


  // ----------------- modifiers -----------------------

  modifier onlyOwner {
    require(msg.sender == owner,
            "Only owner can call this function.");
    _;
  }
  modifier onlyGame {
    require(msg.sender == game,
            "Only game can call this function.");
    _;
  }

  // ----------------- public functions -----------------------

  constructor() public {
    owner = msg.sender;
  }

  /// @notice sets the address of the game that interacts with this contract
  function setGame(address _address) public onlyOwner {
    require (game == address(0x0), "game is already set");
    require (_address == address(0x0), "invalid address 0x0");
    game = _address;
  }

  /// @notice registers a new staker
  /// @param _staker address that will be registered
  function enroll(address payable _staker) public payable onlyOwner onlyGame {
    require (msg.value == kRequiredStake, "failed to enroll: not enough stake");
    require (addStaker(_staker), "failed to enroll");
  }

  /// @notice unregisters a new staker
  /// @param _staker address that will be unregistered
  function unEnroll(address payable _staker) public onlyOwner onlyGame {
    require (removeStaker(_staker), "failed to unenroll");
    _staker.transfer(kRequiredStake);
  }

  /// @notice update to a new level
  /// @param _level to which update
  /// @param _staker address of the staker that reports this update
  /// @dev if some state from previous updates can be resolved, it will be done at this point. That means previous updates could be slashed
  function update(uint16 _level, address _staker) public onlyGame {
    require (_level == level() + 1, "cannout update: unexpected update level");
    require (_level < maxNumLevels() + 1, "cannot update: level too large");
    require (!isSlashed(_staker), "cannot update: staker was slashed");
    // TODO: add logic of the stakers game. For now just simply push
    updaters.push(_staker);
  }

  /// @notice start new verse
  /// @dev current state will be resolved at this point
  function start() public onlyGame {
    require (updaters.length() > 0 && updaters.length() < 3, "cannot start new verse from current level");
    // TODO: add logic to resolve the previous stakers game. For now just simply clear state
    updaters.clear();
  }

  /// @notice get the current level
  function level() public view returns (uint) {
    return updaters.length();
  }

  /// @notice get the current level
  function maxNumLevels() public view returns (uint) {
    return updaters.capacity();
  }

  // ----------------- private functions -----------------------

  function isSlashed(address _addr) private view returns (bool) {
    for (uint i=0; i<slashed.length; i++) {
      if (slashed[i] == _addr) {
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

    // staker not found
    return false;
  }
}
