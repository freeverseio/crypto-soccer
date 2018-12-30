pragma solidity ^0.4.24;

import "./CryptoPlayersStorage.sol";
import "openzeppelin-solidity/contracts/access/roles/MinterRole.sol";

/**
 * @title CryptoPlayersMintable
 * @dev CryptoPlayers minting logic
 */
contract CryptoPlayersMintable is CryptoPlayersStorage, MinterRole {
    function mintWithName(address to, string memory name) public onlyMinter {
        uint256 playerId = calculateId(name);
        _mint(to, playerId);
        _setName(playerId, name);
        _setState(playerId, 0);
    }

    function getPlayerId(string name) public view returns(uint256) {
        uint256 id = calculateId(name);
        require(_exists(id));
        return id;
    }

    function calculateId(string name) public pure returns (uint256) {
        bytes32 playerNameHash = keccak256(abi.encodePacked(name));
        uint256 id = uint256(playerNameHash);
        return id;
    }
}
