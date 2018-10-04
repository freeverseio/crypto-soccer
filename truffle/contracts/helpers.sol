pragma solidity ^ 0.4.24;

// contract containing reusable generic functions
contract HelperFunctions {

    /// @dev encodes an array of nums into a single uint with specific bits
    function encode(uint8 nElem, uint16[] nums, uint bits) internal pure returns(uint result) {
        require(bits <= 16);
        result = 0;
        uint b = 0;
        uint maxnum = 1<<bits; // 2**bits
        for (uint8 i=0; i<nElem; i++) {
            require(nums[i] < maxnum);
            result += (uint(nums[i]) << b);
            b += bits;
        }
        return result;
    }

    /// @dev decodes a uint256 into an array of nums with specific bits
    function decode(uint8 nNumbers, uint longState, uint bits) internal pure returns(uint16[] result) {
        require (bits <= 16);
        uint mask = (1 << bits)-1; // (2**bits)-1
        result = new uint16[](nNumbers);
        for (uint8 i=0; i<nNumbers; i++) {
            result[i] = uint16(longState & mask);
            longState >>= bits;
        }
    }
    
    /// obtains value at index from big uint256 longState
    function getNumAtIndex(uint longState, uint8 index, uint bits) internal pure returns(uint) {
        return (longState >> (bits*index))&((1 << bits)-1);
    }

    /// encodes value at specific index into longState
    function setNumAtIndex(uint value, uint longState, uint8 index, uint bits) internal pure returns(uint) {
        uint maxnum = 1<<bits; // 2**bits
        require(value < maxnum);
        uint b = bits*index;
        uint mask = (1 << bits)-1; // (2**bits)-1
        longState &= ~(mask << b); // clear all bits at index
        return longState + (value << b);
    }

    // only used for testing since web3.eth.solidityUtils not yet available
    function computeKeccak256ForNumber(uint n)
    internal
    pure
    returns(uint)
    {
        return uint(keccak256(abi.encodePacked(n)));
    }
    // only used for testing since web3.eth.solidityUtils not yet available
    function computeKeccak256(string s, uint n1, uint n2)
    internal
    pure
    returns(uint)
    {
        return uint(keccak256(abi.encodePacked(s, n1, n2)));
    }

    // throws a dice that returns 0 with probability weight1/(weight1+weight2), and 1 otherwise.
    // In other words, the responsible for weight1 is selected if return = 0.
    // We return a uint, not bool, to allow the return to be used as an idx in an array.
    // The formula is derived as follows. Throw a random number R in the range [0,M]
    // Then, w1 wins if (w1+w2)*(R/M) < w1, and w2 wins otherise. Clear, this is a weighted dice.
    function throwDice(uint weight1, uint weight2, uint rndNum, uint factor)
        internal
        pure
        returns(uint8)
    {
        if( ((weight1+weight2)*rndNum)<(weight1 * (factor-1)) ) {
            return 0;
        } else {
            return 1;
        }
    }

    // Generalization of the previous to any number of weights
    function throwDiceArray(uint[] memory weights, uint rndNum, uint factor)
        internal
        pure
        returns(uint8)
    {
        uint uniformRndInSumOfWeights;
        for (uint8 w = 0; w<weights.length; w++) {
            uniformRndInSumOfWeights += weights[w];
        }
        uniformRndInSumOfWeights *= rndNum;

        uint cumSum = 0;
        for (w = 0; w<weights.length-1; w++) {
            cumSum += weights[w];
            if( uniformRndInSumOfWeights < ( cumSum * (factor-1) )) {
                return w;
            }
        }
        return w;
    }
}
