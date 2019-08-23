pragma solidity >=0.4.21 <0.6.0;

// FUNCTIONS FOR TESTING ONLY
contract AbiTesting {

    // Raw ABI ENCONDING 
    function getAbiAddress(address x) external pure returns (bytes memory) {
        return abi.encode(x);
    }
    function getAbiUint8(uint8 x) external pure returns (bytes memory) {
        return abi.encode(x);
    }
    function getAbiUint256(uint8 x) external pure returns (bytes memory) {
        return abi.encode(x);
    }
    function getAbiBytes32(bytes32 x) external pure returns (bytes memory) {
        return abi.encode(x);
    }
    function getAbiBytes32AndBytes32(bytes32 x, bytes32 y) external pure returns (bytes memory) {
        return abi.encode(x, y);
    }
    
    // HASHING after encoding 
    function getHashOfAddress(address x) external pure returns (bytes32) {
        return keccak256(abi.encode(x));
    }
    function getHashOfUint8(uint8 x) external pure returns (bytes32) {
        return keccak256(abi.encode(x));
    }
    function getHashOfUint256(uint256 x) external pure returns (bytes32) {
        return keccak256(abi.encode(x));
    }
    function getHashOfBytes32(bytes32 x) external pure returns (bytes32) {
        return keccak256(abi.encode(x));
    }

    bytes32 stor = "0x23";

    function repeatAbi(bytes32 x1, bytes32 x2, bytes32 x3, bytes32 x4) public returns(bytes32) {
        stor = keccak256(abi.encode(x1, x2, x3, x4));
    }

    function repeatHash(bytes32 x1, bytes32 x2, bytes32 x3, bytes32 x4) public returns(bytes32) {
        stor = keccak256(abi.encode(keccak256(abi.encode(keccak256(abi.encode(x1, x2)),x3)),x4));
    }
}

