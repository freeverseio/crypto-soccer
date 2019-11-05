pragma solidity >=0.4.21 <0.6.0;

contract EncodingMatchLogPart4 {
    function getInGameSubsHappened(uint256 log, uint8 posInHalf, bool is2ndHalf) public pure returns (uint8) {
        uint8 offset = 173 + 2 * (posInHalf + (is2ndHalf ? 3 : 0));
        return uint8((log >> offset) & 3);
    }
}