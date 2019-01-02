pragma solidity ^0.4.24;

import "./CryptoPlayersStorage.sol";
import "../CryptoSoccer.sol";
import "../HelperFunctions.sol";
import "openzeppelin-solidity/contracts/access/roles/MinterRole.sol";

/**
 * @title CryptoPlayersMintable
 * @dev CryptoPlayers minting logic
 */
contract CryptoPlayersMintable is CryptoPlayersStorage, CryptoSoccer, HelperFunctions, MinterRole {
    function mintWithName(address to, string memory name) public onlyMinter {
        uint256 playerId = _computeId(name);
        _mint(to, playerId);
        _setName(playerId, name);
        uint16 birth = uint16(block.number);  // TODO: reformulate
        uint16[5] memory skills = _computeSkills(name);
        _setGenome(
            playerId, 
            birth,
            skills[0],
            skills[1],
            skills[2],
            skills[3],
            skills[4]
            );
    }

    function getPlayerId(string name) public view returns(uint256) {
        uint256 id = _computeId(name);
        require(_exists(id));
        return id;
    }

    function _computeId(string name) internal pure returns (uint256) {
        bytes32 playerNameHash = keccak256(abi.encodePacked(name));
        uint256 id = uint256(playerNameHash);
        return id;
    }

    function _computeSkills(string name) internal pure returns (uint16[5]) {
        bytes32 playerNameHash = keccak256(abi.encodePacked(name));
        uint256 seed = uint256(playerNameHash);
        uint16[5] memory states = decodeHere(seed);
        return states; 
    }

    function decodeHere(uint256 serialized) internal pure returns (uint16[5]) {
        uint256 copy = serialized;
        uint16[5] memory states;
        for (uint8 i = 0; i<5; i++) {
            states[i] = uint16(copy & 0x3fff);
            copy >>= 14;
        }
        /// @dev The next 5 are states skills. Adjust them to so that they add up to, maximum, 5*50 = 250.
        uint16 excess;
        for (uint8 sk = 0; sk < 5; sk++) {
            states[sk] = states[sk] % 50;
            excess += states[sk];
        }
        /// @dev At this point, at most, they add up to 5*49=245. Share the excess to reach 250:
        excess = (250 - excess)/5;
        for (sk = 0; sk < 5; sk++) {
            states[sk] = states[sk] + excess;
        }
        return states;
    }
}
