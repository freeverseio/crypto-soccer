pragma solidity >=0.5.12 <=0.6.3;

 /**
 * @title Entry point to submit user actions, and timeZone root updates, which makes time evolve.
 */

contract Merkle {
    
    // This function will revert if nLevels = 0.
    // nLeafs = 2**nLevels
    //  nLevels = 1 => 3 leafs = lev1 - root // lev0 - 2 leafs
    //  nLevels = 2 => 7 leafs = lev2 - root // lev1 - 2 leafs // lev0 - 4 leafs
    function merkleRoot(bytes32[] memory array, uint256 nLevels) public pure returns(bytes32) {
        uint256 nLeafs = 2**nLevels;
        require(array.length == nLeafs, "number of leafs is not = pow(2,nLevels)");
        for (uint8 level = 0; level < nLevels; level++) {
            nLeafs /= 2;
            for (uint256 pos = 0; pos < nLeafs; pos++) {
                array[pos] = hash_node(array[2 * pos], array[2 * pos + 1]);      
            }
        }
        return array[0];
    }
  
    function hash_node(bytes32 left, bytes32 right) public pure returns (bytes32 hash) {
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
