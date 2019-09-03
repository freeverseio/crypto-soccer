pragma solidity ^0.5.0;

import "../core/Engine.sol";

contract EngineMock is Engine {
    function teamsGetTired(uint[5] memory skillsTeamA, uint[5]  memory skillsTeamB)
        public
        pure
        returns (uint[5] memory, uint[5] memory)
    {
        return _teamsGetTired(skillsTeamA, skillsTeamB);
    }

    function getNRandsFromSeed(uint16 nRands, uint256 seed) public pure returns (uint16[] memory rnds) {
        return _getNRandsFromSeed(nRands, seed);
    }

    function throwDice(uint weight1, uint weight2, uint rndNum) public pure returns(uint8) {
        return _throwDice(weight1, weight2, rndNum);
    }

    function throwDiceArray(uint[] memory weights, uint rndNum) public pure returns(uint8 w) {
        return _throwDiceArray(weights, rndNum);
    }

    function managesToShoot(uint8 teamThatAttacks, uint[5][2] memory globSkills, uint rndNum)
        public
        pure
        returns (bool)
    {
        return _managesToShoot(teamThatAttacks, globSkills, rndNum);
    }

    function managesToScore(
        uint8 nAttackers,
        uint[] memory attackersSpeed,
        uint[] memory attackersShoot,
        uint blockShoot,
        uint rndNum1,
        uint rndNum2
    )
        public
        pure
        returns (bool)
    {
        return _managesToScore(
            nAttackers,
            attackersSpeed,
            attackersShoot,
            blockShoot,
            rndNum1,
            rndNum2
        );
    }

    function getTeamGlobSkills(uint256[] memory teamState, uint8[3] memory tactic)
        public
        pure
        returns (
            uint[5] memory globSkills,
            uint[] memory attackersSpeed, 
            uint[] memory attackersShoot
        )
    {
        return _getTeamGlobSkills(teamState, tactic);
    }
}

