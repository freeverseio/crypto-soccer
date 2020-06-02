pragma solidity >= 0.6.3;

/**
 @title Subset library to serialize/deserialize match tactics decided by users
 @author Freeverse.io, www.freeverse.io
*/ 
 
contract EncodingTacticsBase3 {
    function getTacticsId(uint256 tactics) public pure returns(uint8) {
        return uint8(tactics & 63);
    }
}
