pragma solidity >=0.4.21 <0.6.0;

import "./Engine.sol";
/**
 * @title Scheduling of leagues, and calls to Engine to resolve games.
 */

contract Cups {
    
    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 constant public LEAGUES_PER_CUP = 16;
    uint8 constant public TEAMS_PER_LEAGUE = 8;
    uint8 constant public MATCHDAYS = 14;
    uint8 constant public MATCHES_PER_DAY = 4;
    Engine private _engine;

    function setEngineAdress(address addr) public {
        _engine = Engine(addr);
    }

    function getEngineAddress() public view returns (address) {
        return address(_engine);
    }

    // groupIdx = 0,...,15
    // teamIdx  = 0,...,128
    function getTeamsInGroup(uint8 groupIdx) public pure returns(uint8[8] memory teamIdxs) {
        if (groupIdx % 2 == 0) {
            for (uint8 t = 0; t < 8; t++) {
                teamIdxs[t] = 8 * t + groupIdx / 2;
            }
        } else {
            for (uint8 t = 0; t < 8; t++) {
                teamIdxs[t] = 8 * t + (groupIdx - 1) / 2 + 64;
            }
        }
    }

}