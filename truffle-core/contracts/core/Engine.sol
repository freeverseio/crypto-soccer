pragma solidity ^0.5.0;

import "./Leagues.sol";
import "../state/PlayerState.sol";

contract Engine is PlayerState {
    // @dev Max num of players allowed in a team
    uint8 constant kMaxPlayersInTeam = 11;
    /**
     * @dev playMatch returns the result of a match
     * @param seed the pseudo-random number to use as a seed for the match
     * @param state0 a vector with the state of the players of team 0
     * @param state1 a vector with the state of the players of team 1
     * @param tactic0 a vector[3] with the tactic (ex. [4,4,2]) of team 0 
     * @param tactic1 a vector[3] with the tactic (ex. [4,4,2]) of team 1
     * @return the score of the match
     */
    function playMatch(
        bytes32 seed,
        uint256[] memory state0,
        uint256[] memory state1, 
        uint8[3] memory tactic0, 
        uint8[3] memory tactic1
    ) 
        public 
        pure 
        returns (uint8 home, uint8 visitor) 
    {
        require(state0.length == 11, "Team 0 needs 11 players");
        require(state1.length == 11, "Team 1 needs 11 players");
        require(tactic0[0] + tactic0[1] + tactic0[2] == 10, "wrong tactic for team 0");
        require(tactic1[0] + tactic1[1] + tactic1[2] == 10, "wrong tactic for team 1");
        bytes32 hash0 = keccak256(abi.encode(uint256(seed) + state0[0] + tactic0[0]));
        bytes32 hash1 = keccak256(abi.encode(uint256(seed) + state1[0] + tactic1[0]));
        return (uint8(uint256(hash0) % 4), uint8(uint256(hash1) % 4));
    }


    /// @dev Computes basic data, including globalSkills, needed during the game.
    /// @dev Basically implements the formulas:
    // move2attack =    defence(defenders + 2*midfields + attackers) +
    //                  speed(defenders + 2*midfields) +
    //                  pass(defenders + 3*midfields)
    // createShoot =    speed(attackers) + pass(attackers)
    // defendShoot =    speed(defenders) + defence(defenders);
    // blockShoot  =    shoot(keeper);
    function getTeamGlobSkills(uint256[] memory teamState, uint8[3] memory tactic)
        public
        pure
        returns (
            uint[5] memory globSkills,
            uint[] memory attackersSpeed, 
            uint[] memory attackersShoot
        )
    {
        attackersSpeed = new uint[](tactic[2]); 
        attackersShoot = new uint[](tactic[2]); 

        uint8 p = 1;
        uint8 i = 0;
/*
        for (;p <= tactic[0]; p++) {
        }
        for (;p <= tactic[1]; p++) {
        }
        */
        for (;p <= tactic[2]; p++) {
            attackersSpeed[i] = getSpeed(teamState[p]); 
            attackersShoot[i] = getShoot(teamState[p]); 
            i++;
        }


/*
        uint move2attack;
        uint createShoot;
        uint defendShoot;
        uint blockShoot;
        uint endurance;

        nAttackers = 0;

        for (uint8 p = 0; p < kMaxPlayersInTeam; p++) {
            uint16[] memory skills = decode(kNumStates, getStatePlayerInTeam(p, _teamIdx), kBitsPerState);
            endurance += skills[kStatEndur];
            if (skills[kStatRole] == kRoleKeeper) {
                blockShoot = skills[kStatShoot];
            }
            else if (skills[kStatRole] == kRoleDef) {
                move2attack = move2attack + skills[kStatDef] + skills[kStatSpeed] + skills[kStatPass];
                defendShoot = defendShoot + skills[kStatSpeed] + skills[kStatDef];
            }
            else if (skills[kStatRole] == kRoleMid) {
                move2attack = move2attack + 2 * skills[kStatDef] + 2 * skills[kStatSpeed] + 3 * skills[kStatPass];
            }
            else if (skills[kStatRole] == kRoleAtt) {
                move2attack = move2attack + skills[kStatDef];
                createShoot = createShoot + skills[kStatSpeed] + skills[kStatPass];
                attackersSpeed[nAttackers] = skills[kStatSpeed];
                attackersShoot[nAttackers] = skills[kStatShoot];
                nAttackers++;
            }
        }
        /// @dev endurance is converted to a percentage, 
        /// @dev used to multiply (and hence decrease) the start endurance.
        /// @dev 100 is super-endurant (1500), 70 is bad, for an avg starting team (550).
        if (endurance < 500) {
            endurance = 70;
        } else if (endurance < 1400) {
            endurance = 100 - (1400-endurance)/30;
        } else {
            endurance = 100;
        }

        return (
            [move2attack, createShoot, defendShoot, blockShoot, endurance],
            nAttackers,
            attackersSpeed,
            attackersShoot
        );
*/

    }
}

