pragma solidity ^0.5.0;

import "./PlayersProps.sol";
import "openzeppelin-solidity/contracts/access/roles/MinterRole.sol";

/**
 * @title PlayersMintable
 * @dev Players minting logic
 */
contract PlayersMintable is PlayersProps, MinterRole {
    mapping(bytes32 => uint256) private _nameHashToId;

    /**
     * @dev Function to mint players
     * @param to The address that will receive the minted players.
     * @param name The player name to mint.
     * @return A boolean that indicates if the operation was successful.
     */
    function mint(address to, string memory name) public onlyMinter {
        bytes32 playerNameHash = keccak256(abi.encodePacked(name));
        require(_nameHashToId[playerNameHash] == 0, "player already exists");
        uint256 playerId = totalSupply() + 1;
        _nameHashToId[playerNameHash] = playerId;
        uint16 birth = uint16(block.number);  // TODO: reformulate
        uint16[5] memory skills = _computeSkills(playerId);
        _mint(to, playerId);
        _setName(playerId, name);
        _setGenome(
            playerId, 
            birth,
            skills[0],
            skills[1],
            skills[2],
            skills[3],
            skills[4]
        );
        _setCurrentHistory(
            playerId, 
            0,
            0,
            0,
            0,
            0,
            0
        );    
    }

    /**
     * @return player ID from player name
     */
    function getPlayerId(string memory name) public view returns(uint256) {
        bytes32 playerNameHash = keccak256(abi.encodePacked(name));
        uint256 id = _nameHashToId[playerNameHash];
        require(id != 0, "unexistent player");
        return _nameHashToId[playerNameHash];
    }

    function _computeId(string memory name) internal pure returns (uint256) {
        bytes32 playerNameHash = keccak256(abi.encodePacked(name));
        uint256 id = uint256(playerNameHash);
        return id;
    }

    /**
     * @dev Compute the pseudorandom skills, sum of the skills is 250
     * @param seed to generate the skills
     * @return 5 skills
     */
    function _computeSkills(uint256 seed) internal view returns (uint16[5] memory) {
        uint256 rand = uint256(keccak256(abi.encodePacked(blockhash(block.number-1))));
        uint256 rna = rand + seed;

        uint16[5] memory skills;
        for (uint8 i = 0; i<5; i++) {
            skills[i] = uint16(rna & 0x3fff);
            rna >>= 14;
        }

        /// @dev The next 5 are skills skills. Adjust them to so that they add up to, maximum, 5*50 = 250.
        uint16 excess;
        for (uint8 i = 0; i < 5; i++) {
            skills[i] = skills[i] % 50;
            excess += skills[i];
        }

        /// @dev At this point, at most, they add up to 5*49=245. Share the excess to reach 250:
        uint16 delta = (250 - excess) / 5;
        for (uint8 i = 0; i < 5; i++) 
            skills[i] = skills[i] + delta;

        uint16 remainder = (250 - excess) % 5;
        for (uint8 i = 0 ; i < remainder ; i++)
            skills[i]++;

        return skills;
    }
}
