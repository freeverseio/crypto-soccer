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
        uint256 state = _computeState(playerId, 0);
        _mint(to, playerId);
        _setName(playerId, name);
        // _setGenome(playerId, state);
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

    /// @dev Main interface to create a player by users. We receive a random number,
    /// @dev computed elsewhere (e.g. from hash(name+userChoice+dorsal)) and create 
    /// @dev a balanced player whose skills add up to 250.
    function _computeState(uint256 rndSeed, uint8 playerRole) internal view returns(uint256)
    {
        /// @dev Get random numbers between 0 and 9999 and assign them to states, where:
        /// @dev state[0] -> age, state[6] -> role
        /// @dev state[1]...state[5] -> skills
        uint16[] memory states = decode(kNumStates, rndSeed, kBitsPerState);

        /// @dev Last number is role, as provided from outside. Just store it.
        states[kStatRole] = playerRole;

        /// @dev Ensure that age, in years at moment of creation, can vary between 16 and 35.
        states[kStatBirth] = 16 + (states[0] % 20);

        /// @dev Convert age to monthOfBirthAfterUnixEpoch.
        /// @dev TODO: We can optimize by not declaring these as variables, and putting the exact numbers. 
        /// @dev I leave it this way for clarity, for the time being.
        uint years2secs = 365 * 24 * 3600;
        uint month2secs = 30 * 24 * 3600;
        states[kStatBirth] = uint16((block.timestamp - states[0] * years2secs) / month2secs);

        /// @dev The next 5 are states skills. Adjust them to so that they add up to, maximum, 5*50 = 250.
        uint16 excess;
        for (uint8 sk = kStatDef; sk <= kStatEndur; sk++) {
            states[sk] = states[sk] % 50;
            excess += states[sk];
        }
        /// @dev At this point, at most, they add up to 5*49=245. Share the excess to reach 250:
        excess = (250 - excess)/kNumSkills;
        for (sk = kStatDef; sk <= kStatEndur; sk++) {
            states[sk] = states[sk] + excess;
        }

        return serialize(kNumStates, states, kBitsPerState);
    }
}
