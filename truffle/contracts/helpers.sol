pragma solidity ^ 0.4.24;

// Contract containing reusable generic functions
contract HelperFunctions {

    /// @dev Serializes an array of nElem numbers into a single uint with specific bits used for each num.
    function serialize(uint8 nElem, uint16[] nums, uint bits) internal pure returns(uint256 result) {
        require(bits*nElem <= 256, "Not enough space in a uint256 to serialize");
        require(bits <= 16, "Not enough bits to encode each number, since they are read as uint16");
        result = 0;
        uint usedBits = 0;
        for (uint8 i = 0; i < nElem; i++) {
            result += (uint(nums[i]) << usedBits);
            usedBits += bits;
        }
        return result;
    }

    /// @dev Decodes a serialized uint256 into an array of nums with specific bits used for each num
    function decode(uint8 nElem, uint serialized, uint bits) internal pure returns(uint16[] decoded) {
        require (bits <= 16, "Not enough bits to encode each number, since they are read as uint16");
        uint mask = (1 << bits)-1; // (2**bits)-1
        decoded = new uint16[](nElem);
        for (uint8 i=0; i<nElem; i++) {
            decoded[i] = uint16(serialized & mask);
            serialized >>= bits;
        }
    }

    /// @dev Returns value at a given position (index) from a serialized uint256
    function getNumAtIndex(uint serialized, uint8 index, uint bits) public pure returns(uint) {
        return (serialized >> (bits*index))&((1 << bits)-1);
    }

    /// @dev Sets the number at a given index in a serialized uint256
    function setNumAtIndex(uint value, uint serialized, uint8 index, uint bits) 
        internal 
        pure 
        returns(uint) 
    {
        uint maxnum = 1<<bits; // 2**bits
        require(value < maxnum, "Value too large to fit in available space");
        uint b = bits*index;
        uint mask = (1 << bits)-1; // (2**bits)-1
        serialized &= ~(mask << b); // clear all bits at index
        return serialized + (value << b);
    }

    /// @dev Returns the hash of a uint using the hash function used in this game.
    /// @dev Only used for testing since web3.eth.solidityUtils not yet available
    function computeKeccak256ForNumber(uint n)
        internal
        pure
        returns(uint)
    {
        return uint(keccak256(abi.encodePacked(n)));
    }

    /// @dev Returns the hash of concat(string,uint,uint) using the hash function used in this game.
    /// @dev Only used for testing since web3.eth.solidityUtils not yet available
    function computeKeccak256(string s, uint n1, uint n2)
        internal
        pure
        returns(uint)
    {
        return uint(keccak256(abi.encodePacked(s, n1, n2)));
    }

    /// @dev Throws a dice that returns 0 with probability weight1/(weight1+weight2), and 1 otherwise.
    /// @dev So, returning 0 has semantics: "the responsible for weight1 is selected".
    /// @dev We return a uint8, not bool, to allow the return to be used as an idx in an array by the callee.
    /// @dev The formula is derived as follows. Throw a random number R in the range [0,maxR].
    /// @dev Then, w1 wins if (w1+w2)*(R/maxR) < w1, and w2 wins otherise. 
    /// @dev maxRndNum controls the resolution or fine-graining of the algorithm.
    function throwDice(uint weight1, uint weight2, uint rndNum, uint maxRndNum)
        internal
        pure
        returns(uint8)
    {
        if( ( (weight1 + weight2) * rndNum ) < ( weight1 * (maxRndNum-1) ) ) {
            return 0;
        } else {
            return 1;
        }
    }

    /// @dev Generalization of the previous to any number of input weights
    /// @dev It therefore throws any number of dice and returns the winner's idx.
    function throwDiceArray(uint[] memory weights, uint rndNum, uint maxRndNum)
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
            if( uniformRndInSumOfWeights < ( cumSum * (maxRndNum-1) )) {
                return w;
            }
        }
        return w;
    }

    /// @dev A function needed for game scheduling:  
    ///  P(t) = { t if t < T; t-(N-1) otherwise }
    function circulate(uint8 t, uint8 nTeams) internal pure returns(uint8) {
        if (t < nTeams) { return t; }
        else { return t-(nTeams-1); }
    }

    /// @dev For a given round in a league, and a given game number
    ///  in that round, it returns the teams that play that game,
    ///  in order (first plays at home), according to formula:  
    ///  game(n,r) = ( P(N-n+r),  P(n+1+r) )   (except at game 0)
    function teamsInGame(uint8 round, uint8 game, uint8 nTeams)
        internal
        pure
        returns(uint8 team1, uint8 team2)
    {
        if (game > 0) {
            return (circulate(nTeams-game+round, nTeams), circulate(nTeams+1+round, nTeams));
        } else {
            return (0, circulate(nTeams+1+round, nTeams));
        }
    }

}
