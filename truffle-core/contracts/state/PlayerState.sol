pragma solidity ^0.5.0;

/// @title the state of a player
contract PlayerState {
    uint256 constant public TEAMSTATEEND = 0;

    function isValidPlayerState(uint256 state) public pure returns (bool) {
        return state != TEAMSTATEEND;
    }

    /**
     * @dev encoding:
     *        defence:   0xff00000000
     *        speed:     0x00ff000000
     *        pass:      0x0000ff0000
     *        shoot:     0x000000ff00
     *        endurance: 0x00000000ff
     */
    function playerStateCreate(
        uint8 defence,
        uint8 speed,
        uint8 pass,
        uint8 shoot,
        uint8 endurance 
    )
        public 
        pure
        returns (uint256 state)
    {
        state |= uint256(defence) << 8 * 4;
        state |= uint256(speed) << 8 * 3;
        state |= uint256(pass) << 8 * 2;
        state |= uint256(shoot) << 8;
        state |= endurance;
    }

    function getDefence(uint256 playerState) public pure returns (uint8) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint8(playerState >> 8 * 4);
    }
    
    function getSpeed(uint256 playerState) public pure returns (uint8) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint8(playerState >> 8 * 3);
    }

    function getPass(uint256 playerState) public pure returns (uint8) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint8(playerState >> 8 * 2);
    }

    function getShoot(uint256 playerState) public pure returns (uint8) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint8(playerState >> 8);
    }

    function getEndurance(uint256 playerState) public pure returns (uint8) {
        require(isValidPlayerState(playerState), "invalid player state");
        return uint8(playerState);
    }

    /// increase the skills of delta
    function playerStateEvolve(uint256 playerState, uint8 delta) public pure returns (uint256) {
        require(isValidPlayerState(playerState), "invalid player playerState");
        uint8 defence = getDefence(playerState) + delta;
        uint8 speed = getSpeed(playerState) + delta;
        uint8 pass = getPass(playerState) + delta;
        uint8 shoot = getShoot(playerState) + delta;
        uint8 endurance = getEndurance(playerState) + delta;
        return playerStateCreate(defence, speed, pass, shoot, endurance);
    }
}