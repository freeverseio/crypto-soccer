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

  address public owner = address(0x0);
  address public game = address(0x0);
  AddressStack private updaters = new AddressStack();
  address[] public stakers;
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
    game = _address;
  }

  /// @notice registers a new staker
  /// @param _staker address that will be registered
  function register(address _staker) public onlyOwner onlyGame {
    require (stakers.length < kNumStakers, "");
    stakers.push(_staker);
  }

  /// @notice update to a new level
  /// @param _level to which update
  /// @param _staker address of the staker that reports this update
  /// @dev if some state from previous updates can be resolved, it will be done at this point. That means previous updates could be slashed
  function update(uint16 _level, address _staker) public onlyGame {
    require (_level == level() + 1);
    require (_level < maxNumLevels() + 1);
    // TODO: add logic of the stakers game. For now just simply push
    updaters.push(_staker);
  }

  /// @notice start new verse
  /// @dev current state will be resolved at this point
  function start() public onlyGame {
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
}
