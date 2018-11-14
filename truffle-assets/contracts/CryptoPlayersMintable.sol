pragma solidity ^0.4.24;

import "./CryptoPlayers.sol";
import "openzeppelin-solidity/contracts/access/roles/MinterRole.sol";

contract CryptoPlayersMintable is CryptoPlayers, MinterRole {
    /// @dev Event fired whenever a new player is created
    event PlayerCreation(string playerName, uint playerIdx, uint playerState);

    constructor(string name, string symbol, string CID) public 
    CryptoPlayers(name, symbol, CID)
    {
    }

  /**
   * @dev Function to mint tokens
   * @param to The address that will receive the minted tokens.
   * @param tokenId The token id to mint.
   * @return A boolean that indicates if the operation was successful.
   */
    function mint(
        address to,
        uint256 tokenId
    )
        public
        onlyMinter
        returns (bool)
    {
        uint state = createBalancedState();
        _mint(to, tokenId);
        _setState(tokenId, state);
        return true;
    }

    function createBalancedState() internal pure returns (uint) {
        uint _monthOfBirthAfterUnixEpoch = 4;
        uint _defense = 50;
        uint _speed = 50;
        uint _pass = 50;
        uint _shoot = 50;
        uint _endurance = 50;
        uint _role = 50;
        uint kBitsPerState = 14;
        uint bits = kBitsPerState;
        uint state = _monthOfBirthAfterUnixEpoch;
        state += (_defense << bits);
        state += (_speed << (bits*2));
        state += (_pass << (bits*3));
        state += (_shoot << (bits*4));
        state += (_endurance << (bits*5));
        state += (_role << (bits*6));
        return state;
    }


}
