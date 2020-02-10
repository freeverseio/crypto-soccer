pragma solidity >=0.5.12 <0.6.2;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

contract EncodingTacticsPart2 {

    function setStaminaRecovery(uint256 tactics, uint8 val) public pure returns(uint256) {
        require(val < 4, "staminaRecovery must be < 4");
        return (tactics & ~(uint256(3) << 110)) | (uint256(val) << 110);
    }

    function setItemId(uint256 tactics, uint16 val) public pure returns(uint256) {
        require(val < 2**13, "staminaRecovery must be < 2**13");
        return (tactics & ~(uint256(8191) << 112)) | (uint256(val) << 112);
    }

    function setItemBoost(uint256 tactics, uint32 val) public pure returns(uint256) {
        return (tactics & ~(uint256(4294967295) << 125)) | (uint256(val) << 125);
    }
    
    function getItemsData(uint256 tactics) public pure returns(uint8, uint16, uint32) {
        return (
            uint8((tactics >> 110) & 3), 
            uint16((tactics >> 112) & 8191), 
            uint32((tactics >> 125) & 4294967295)
        );
    }

}
