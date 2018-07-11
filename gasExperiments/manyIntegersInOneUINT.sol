pragma solidity ^0.4.24;

// an illustration of how to store all the things we need about a player in one single int
// I've been using this to compare the cost of reading numbers from a long number via hardcore
// explicit writing of the outputs, or via a for loop. 
// The latter is a bit more costly. For 5 numbers, 2403 versus 1144...
// Also, using uint fac = fac() is cheaper than calling fac() all the time.
// Finally, using fac() vs storage uint factor when deploying but not when executing.

contract manyIntsInOne {
    // assume we want to store many numbers, n1, n2, n3... each within 0 and 9999 in one long uint
    // we will store them sequentially, starting from the end.
    // each number is 4 digits long.
    uint factor = 10000;
    function fac() internal pure returns (uint) { return 10000; }

    function divide(uint numerator, uint denominator) internal pure returns (uint quotient, uint16 remainder)  {
        quotient  = numerator / denominator;
        remainder = uint16( numerator - denominator * quotient ) ;
    }
    
/*  
    function store(uint16 n1, uint16 n2, uint16 n3, uint16 n4, uint16 n5) public {
      state = n1 + n2 * factor + n3 * factor*factor 
                     + n4 * factor*factor*factor + n5 * factor*factor*factor*factor;
    }
*/
    
    function readNumbers(uint state) public returns (uint16 n1, uint16 n2, uint16 n3, uint16 n4, uint16 n5) {
        uint quotient;
        uint fact = fac();
        (quotient,n1)   = divide(state,fact); 
        (quotient,n2)   = divide(quotient,fact); 
        (quotient,n3)   = divide(quotient,fact); 
        (quotient,n4)   = divide(quotient,fact); 
        (quotient,n5)   = divide(quotient,fact); 
    }
    
    function readNumbersArray(uint8 nNumbers, uint state) public returns (uint16[] memory result) {
        uint quotient;
        uint fact = fac();
        result = new uint16[](nNumbers);
        (quotient,result[0]) = divide(state,fact);
        for (uint8 n=1;n<nNumbers;n++){ 
            (quotient,result[n]) = divide(quotient,fact); 
        }
    }
    
    function generateRndNumbersArray(uint8 nNumbers, bytes32 rndSeed) public returns (uint16[] memory result) {
        result = new uint16[](nNumbers);
        result[0] = uint16(keccak256(rndSeed));
        for (uint8 n=1;n<nNumbers;n++){ 
            result[n] = uint16(keccak256(rndSeed)); 
        }
    }
}


