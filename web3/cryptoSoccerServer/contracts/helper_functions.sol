pragma solidity ^ 0.4.24;

import "./SafeMath.sol";
import "./Math.sol";

// contract containing reusable generic functions
contract HelperFunctions {

    using SafeMath for uint256;
    using Math for uint256;

    function divideUint(uint numerator, uint denominator) internal pure returns(uint quotient, uint16 remainder) {
        quotient = numerator.div(denominator);
        remainder = uint16(numerator.sub(denominator.mul(quotient)));
    }

    function divideUint8(uint8 numerator, uint8 denominator) internal pure returns(uint8 quotient, uint8 remainder) {
        quotient = numerator / denominator;
        remainder = numerator - denominator * quotient;
    }

    function min(uint16 a, uint16 b) pure private returns(uint16) {
        return a < b ? a : b;
    }

    function encodeIntoLongIntArray(uint8 nElem, uint16[] rnds, uint factor) internal pure returns(uint) {
        uint state = rnds[0];
        for (uint8 n = 1; n < nElem-1; n++) {
            state = state + uint(rnds[n]) * factor;
            factor *= factor;
        }
        return (state + rnds[nElem-1] * factor);
    }

    // reads an arbitrary number of numbers from a long one
    function readNumbersFromUint(uint8 nNumbers, uint longState, uint factor)  
        public
        pure 
        returns(uint16[] memory result) 
    {
        uint quotient;
        result = new uint16[](nNumbers);
        (quotient, result[0]) = divideUint(longState, factor);
        for (uint8 n = 1; n < nNumbers; n++) {
            (quotient, result[n]) = divideUint(quotient, factor);
        }
    }

    function power(uint x, uint8 exponent) internal pure returns(uint) {
        uint result = 1;
        for (uint8 e=0; e<exponent;e++) {
            result *= x;
        }
        return result;
    }

    function getNumAtPos(uint longState, uint8 pos, uint factor)  
        internal
        pure 
        returns(uint) 
    {
        return (longState/power(factor,pos)) % factor;
    }   


    function setNumAtPos(uint num, uint longState, uint8 pos, uint factor)  
        public
        pure 
        returns(uint) 
    {
        return longState + (num - getNumAtPos(longState, pos, factor))*power(factor, pos);
    }   

    // uses the previous function, feeding it with a uint256 hash, only for test use
    function readNumbersFromHash(uint8 nNumbers, uint seed, uint factor)  
        public
        pure 
        returns(uint16[] memory result) 
    {
        uint longState = uint(keccak256(abi.encodePacked(seed)));
        return readNumbersFromUint(nNumbers, longState, factor);
    }
    

    // throws a dice that returns 0 with probability weight1/(weight1+weight2), and 1 otherwise.
    // In other words, the responsible for weight1 is selected if return = 0.
    // We return a uint, not bool, to allow the return to be used as an idx in an array.
    // The formula is derived as follows. Throw a random number R in the range [0,M]
    // Then, w1 wins if (w1+w2)*(R/M) < w1, and w2 wins otherise. Clear, this is a weighted dice.
    function throwDice(uint weight1, uint weight2, uint rndNum, uint factor)   
        public
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
    function throwDiceArray(uint16[] memory weights, uint rndNum, uint factor)   
        public
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
