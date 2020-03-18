pragma solidity >=0.5.12 <=0.6.3;

 /**
 * @title Entry point to submit user actions, and timeZone root updates, which makes time evolve.
 */

contract Merkle {
    
    // This function will revert if nLevels = 0, unless the array is empty.
    // nLeafs = 2**nLevels
    //  nLevels = 0 => 1 leafs = lev0 - root (nothing to do)
    //  nLevels = 1 => 3 leafs = lev1 - root // lev0 - 2 leafs
    //  nLevels = 2 => 7 leafs = lev2 - root // lev1 - 2 leafs // lev0 - 4 leafs
    function merkleRoot(bytes32[] memory array, uint256 nLevels) public pure {
        require(array.length == 2**nLevels, "number of leafs is not = pow(2,nLevels)");
        for (uint8 level = 0; level < nLevels - 1; level++) {
            for (uint32 pos = 0; pos < 2**(nLevels - level); pos++) {
                array[pos] = hash_node(array[2 * pos], array[2 * pos + 1]);      
            }
        }
    }
  
    function hash_node(bytes32 left, bytes32 right) internal pure returns (bytes32 hash) {
        assembly {
            mstore(0x00, left)
            mstore(0x20, right)
            hash := keccak256(0x00, 0x40)
        }
        return hash;
    }

    // if nLevels = 1, we need 1 element in the proof
    // if nLevels = 2, we need 2 elements...
    //        .
    //     ..   ..
    //   .. .. .. ..
    //   01 23 45 67
    
    function verify(bytes32 root, bytes32[] memory proof, bytes32 leafHash, uint256 leafPos) public pure returns(bool) {
        for (uint32 pos = 0; pos < proof.length; pos++) {
            if ((leafPos % 2) == 0) {
                leafHash = hash_node(leafHash, proof[pos]);
            } else {
                leafHash = hash_node(proof[pos], leafHash);
            }
            leafPos /= 2;
        }     
        return root == leafHash;   
    }
    
}
