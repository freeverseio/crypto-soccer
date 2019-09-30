pragma solidity >=0.4.21 <0.6.0;
/**
 * @title Library of functions to serialize values into uints, and deserialize back
 */

contract EncodingIDs {

    /**
     * @dev PlayerId and TeamId both serialize a total of 43 bits:
     *      timeZone        = 5 bits
     *      countryIdxInTZ  = 10 bits
     *      val             = 28 bits (either  (playerIdxInCountry or teamIdxInCountry)
    **/
    function encodeTZCountryAndVal(uint8 timeZone, uint256 countryIdxInTZ, uint256 val) public pure returns (uint256)
    {
        require(timeZone < 2**5, "defence out of bound");
        require(countryIdxInTZ < 2**10, "defence out of bound");
        require(val < 2**28, "defence out of bound");
        uint256 encoded  = uint256(timeZone) << 38;        // 43 - 5
        encoded         |= countryIdxInTZ << 28;  // 38 - 10
        return (encoded | val);            // 28 - 28
    }

    function decodeTZCountryAndVal(uint256 encoded) public pure returns (uint8, uint256, uint256)
    {
        // 2**14 - 1 = 31;  2**10 - 1 = 1023; 2**28 - 1 = 268435455;
        return (uint8(encoded >> 38 & 31), uint256(encoded >> 28 & 1023), uint256(encoded & 268435455));
    }

}