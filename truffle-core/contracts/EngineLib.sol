pragma solidity ^0.5.0;

import "./EncodingSkills.sol";

contract EngineLib is EncodingSkills {
    uint8 private constant BITS_PER_RND     = 36;   // Number of bits allowed for random numbers inside match decisisons
    uint256 public constant MAX_RND         = 68719476735; // Max random number allowed inside match decisions: 2^36-1
    // // Idxs for vector of globSkills: [0=move2attack, 1=globSkills[IDX_CREATE_SHOOT], 2=globSkills[IDX_DEFEND_SHOOT], 3=blockShoot, 4=currentEndurance]
    uint256 private constant SECS_IN_DAY    = 86400; // 24 * 3600 

    /// @dev Throws a dice that returns 0 with probability weight0/(weight0+weight1), and 1 otherwise.
    /// @dev So, returning 0 has semantics: "the responsible for weight0 is selected".
    /// @dev We return a uint8, not bool, to allow the return to be used as an idx in an array by the callee.
    /// @dev The formula is derived as follows. Throw a random number R in the range [0,maxR].
    /// @dev Then, w0 wins if (w0+w1)*(R/maxR) < w0, and w1 wins otherise. 
    /// @dev MAX_RND controls the resolution or fine-graining of the algorithm.
    function throwDice(uint256 weight0, uint256 weight1, uint256 rndNum) public pure returns(uint8) {
        if( ( (weight0 + weight1) * rndNum ) < ( weight0 * (MAX_RND-1) ) ) {
            return 0;
        } else {
            return 1;
        }
    }

    /// @dev Generalization of the previous to any number of input weights
    /// @dev It therefore throws any number of dice and returns the winner's idx.
    function throwDiceArray(uint256[] memory weights, uint256 rndNum) public pure returns(uint8 w) {
        uint256 uniformRndInSumOfWeights;
        for (w = 0; w < weights.length; w++) {
            uniformRndInSumOfWeights += weights[w];
        }
        uniformRndInSumOfWeights *= rndNum;
        uint256 cumSum = 0;
        for (w = 0; w < weights.length-1; w++) {
            cumSum += weights[w];
            if( uniformRndInSumOfWeights < ( cumSum * (MAX_RND-1) )) {
                return w;
            }
        }
        return w;
    }

    function getNRandsFromSeed(uint256 seed, uint8 nRnds) public pure returns (uint64[] memory) {
        uint256 currentBigRnd = uint256(keccak256(abi.encode(seed)));
        uint8 remainingBits = 255;
        uint64[] memory rnds = new uint64[](nRnds);
        for (uint8 n = 0; n < nRnds; n++) {
            if (remainingBits < BITS_PER_RND) {
                currentBigRnd = uint256(keccak256(abi.encode(seed, n)));
                remainingBits = 255;
            }
            rnds[n] = uint64(currentBigRnd & MAX_RND);
            currentBigRnd >>= BITS_PER_RND;
            remainingBits -= BITS_PER_RND;
        }
        return rnds;
    }

    // no penalty at all => return 1M,  max penalty => return 0
    // for each day that passes over 31 years (=11315 days), we subtract 0,0274%, so that you get 10.001% less per year
    // on a max of 1M, this is 274 per day.
    // so, 3649 days after 31 (ten years), he will reach penalty 0. He'll be useless when reaching 41.
    function penaltyPerAge(uint256 playerSkills, uint256 matchStartTime) public pure returns (uint256) {
        uint256 ageDays = (7 * matchStartTime)/SECS_IN_DAY - 7 * getBirthDay(playerSkills);
        if (ageDays > 14964) return 0; // 3649 + 11315 (41 years)
        return ageDays < 11316 ? 1000000 : 1000000 - 274 * (ageDays - 11315);
    }

    function getNDefenders(uint8[9] memory playersPerZone) public pure returns (uint8) {
        return 2 * playersPerZone[0] + playersPerZone[1];
    }

    function getNMidfielders(uint8[9] memory playersPerZone) public pure returns (uint8) {
        return 2 * playersPerZone[3] + playersPerZone[4];
    }

    function getNAttackers(uint8[9] memory playersPerZone) public pure returns (uint8) {
        return 2 * playersPerZone[6] + playersPerZone[7];
    }
    
    // TODO: can this be expressed as
    // translates from a high level tacticsId (e.g. 442) to a format that describes how many
    // players play in each of the 9 zones in the field (Def, Mid, Forw) x (L, C, R), 
    // We impose left-right symmetry: DR = DL, MR = ML, FR = FL.
    // So we only manage 6 numbers: [DL, DM, ML, MM, FL, FM], and force 
    function getPlayersPerZone(uint8 tacticsId) public pure returns (uint8[9] memory) {
        require(tacticsId < 4, "we currently support only 4 different tactics");
        if (tacticsId == 0) return [1,2,1,1,2,1,0,2,0];  // 0 = 442
        if (tacticsId == 1) return [1,3,1,1,2,1,0,1,0];  // 0 = 541
        if (tacticsId == 2) return [1,2,1,1,1,1,1,1,1];  // 0 = 433
        if (tacticsId == 3) return [1,2,1,1,3,1,0,1,0];  // 0 = 451
    }

}

