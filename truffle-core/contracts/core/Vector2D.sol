pragma solidity ^0.4.25;

library Vector2D {
    uint256 constant public DIVIDER = 0;

    function concat(uint256[] memory left, uint256[] memory right) public pure returns (uint256[] memory) {
        if(left.length == 0)
            return right;
        if(right.length == 0)
            return left;

        uint256[] memory result = new uint256[](left.length + right.length + 1);
        uint256 i;
        for (i = 0; i < left.length ; i++)
            result[i] = left[i];
        result[left.length] = DIVIDER;
        for (i = 0 ; i < right.length ; i++)
            result[left.length + 1 + i] = right[i];

        return result;        
    }
}