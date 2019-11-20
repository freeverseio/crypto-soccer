pragma solidity >=0.4.21 <0.6.0;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

contract EncodingIDs {

    /**
     * @dev PlayerId and TeamId both serialize a total of 49 bits:
     *      timeZone        = 5 bits
     *      countryIdxInTZ  = 10 bits
     *      generation      = 6 bits
     *      val             = 28 bits (either  (playerIdxInCountry or teamIdxInCountry)
    **/
    function encodeTZCountryGenAndVal(uint8 timeZone, uint256 countryIdxInTZ, uint8 gen, uint256 val) public pure returns (uint256)
    {
        require(timeZone < 2**5, "defence out of bound");
        require(countryIdxInTZ < 2**10, "defence out of bound");
        require(val < 2**28, "defence out of bound");
        require(gen < 64, "generation out of bound");
        uint256 encoded = val;
        encoded |= uint256(gen) << 28;
        encoded |= countryIdxInTZ << 34;
        return (encoded | (uint256(timeZone) << 44));
    }

    function decodeTZCountryGenAndVal(uint256 encoded) public pure returns (uint8, uint256, uint8, uint256)
    {
        // 2**14 - 1 = 31;  2**10 - 1 = 1023; 2**28 - 1 = 268435455;
        return (
            uint8((encoded >> 44) & 31), 
            uint256((encoded >> 34) & 1023), 
            uint8((encoded >> 28) & 63),
            uint256(encoded & 268435455)
        );
    }

}