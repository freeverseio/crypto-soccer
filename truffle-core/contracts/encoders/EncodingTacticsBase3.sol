pragma solidity >= 0.6.3;

/**
 @title Subset library to serialize/deserialize match tactics decided by users
 @author Freeverse.io, www.freeverse.io
*/ 
 
contract EncodingTacticsBase3 {
    function getTacticsId(uint256 tactics) public pure returns(uint8) {
        return uint8(tactics & 63);
    }

    function getExtraAttack(uint256 tactics, uint8 p) public pure returns(bool) {
        return (((tactics >> (6 + p)) & 1) == 1 ? true : false); /// 2^1 - 1
    }
}
