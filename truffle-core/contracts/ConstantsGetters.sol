pragma solidity >=0.5.12 <=0.6.3;

import "./Constants.sol";

/**
 * @title Constants used in the project
 */

contract ConstantsGetters is Constants {

    function get_PLAYERS_PER_TEAM_MAX() external pure returns(uint256) { return PLAYERS_PER_TEAM_MAX;}
    function get_PLAYERS_PER_TEAM_INIT() external pure returns(uint256) { return PLAYERS_PER_TEAM_INIT;}
    function get_NULL_ADDR() external pure returns(address) { return NULL_ADDR;}
    function get_LEAGUES_PER_DIV() external pure returns(uint8) { return LEAGUES_PER_DIV;}
    function get_TEAMS_PER_LEAGUE() external pure returns(uint8) { return TEAMS_PER_LEAGUE;}
    function get_FREE_PLAYER_ID() external pure returns(uint256) { return FREE_PLAYER_ID;}
    function get_AUCTION_TIME() external pure returns(uint256) { return AUCTION_TIME;}
    function get_POST_AUCTION_TIME() external pure returns(uint256) { return POST_AUCTION_TIME;}
    function get_NULL_TIMEZONE() external pure returns(uint8) { return NULL_TIMEZONE;}

}