pragma solidity >=0.5.12 <0.6.2;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

contract EncodingTacticsPart2 {

    uint8 constant private PLAYERS_PER_TEAM_MAX = 25;

    function setStaminaRecovery(uint256 tactics, uint8[PLAYERS_PER_TEAM_MAX] memory vals) public pure returns(uint256) {

        for (uint8 p = 0; p < PLAYERS_PER_TEAM_MAX; p++) {
            require(vals[p] < 4, "staminaRecovery must be < 4");
            tactics = (tactics & ~(uint256(3) << (110 + 2*p))) | (uint256(vals[p]) << (110 + 2*p));
        }
        return tactics;
    }

    function setItemId(uint256 tactics, uint16 val) public pure returns(uint256) {
        require(val < 2**13, "staminaRecovery must be < 2**13");
        return (tactics & ~(uint256(8191) << 160)) | (uint256(val) << 160);
    }

    function setItemBoost(uint256 tactics, uint32 val) public pure returns(uint256) {
        return (tactics & ~(uint256(4294967295) << 173)) | (uint256(val) << 173);
    }
    
    function getItemsData(uint256 tactics) public pure returns(uint8[PLAYERS_PER_TEAM_MAX] memory staminas, uint16, uint32) {
        for (uint8 p = 0; p < PLAYERS_PER_TEAM_MAX; p++) {
            staminas[p] = uint8((tactics >> (110 + 2*p)) & 3);
        }
        return (
            staminas, 
            uint16((tactics >> 160) & 8191), 
            uint32((tactics >> 173) & 4294967295)
        );
    }

}
