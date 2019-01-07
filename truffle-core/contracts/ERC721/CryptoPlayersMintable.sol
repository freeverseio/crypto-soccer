pragma solidity ^0.4.24;

import "./CryptoPlayersProps.sol";
import "../CryptoSoccer.sol";
import "../HelperFunctions.sol";
import "openzeppelin-solidity/contracts/access/roles/MinterRole.sol";

/**
 * @title CryptoPlayersMintable
 * @dev CryptoPlayers minting logic
 */
contract CryptoPlayersMintable is CryptoPlayersProps, CryptoSoccer, HelperFunctions, MinterRole {
    /**
     * @dev Function to mint players
     * @param to The address that will receive the minted players.
     * @param name The player name to mint.
     * @return A boolean that indicates if the operation was successful.
     */
    function mint(address to, string memory name) public onlyMinter {
        uint256 playerId = _computeId(name);
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
    }

    /**
     * @return player ID from player name
     */
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

    /**
     * @dev Compute the pseudorandom skills, sum of the skills is 250
     * @param seed to generate the skills
     * @return 5 skills
     */
    function _computeSkills(uint256 seed) internal view returns (uint16[5]) {
        uint256 rand = uint256(keccak256(abi.encodePacked(blockhash(block.number-1))));
        uint256 rna = rand + seed;

        uint16[5] memory skills;
        for (uint8 i = 0; i<5; i++) {
            skills[i] = uint16(rna & 0x3fff);
            rna >>= 14;
        }

        /// @dev The next 5 are skills skills. Adjust them to so that they add up to, maximum, 5*50 = 250.
        uint16 excess;
        for (i = 0; i < 5; i++) {
            skills[i] = skills[i] % 50;
            excess += skills[i];
        }

        /// @dev At this point, at most, they add up to 5*49=245. Share the excess to reach 250:
        uint16 delta = (250 - excess) / 5;
        for (i = 0; i < 5; i++) 
            skills[i] = skills[i] + delta;

        uint16 remainder = (250 - excess) % 5;
        for (i = 0 ; i < remainder ; i++)
            skills[i]++;

        return skills;
    }
}
