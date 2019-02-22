pragma solidity ^0.5.0;

contract Scores {
    uint256 constant public DIVIDER = uint(-1);

    function scoresConcat(uint256[2][] memory left, uint256[2][] memory right) public pure returns (uint256[2][] memory) {
        if(left.length == 0)
            return right;
        if(right.length == 0)
            return left;

        uint256[2][] memory result = new uint256[2][](left.length + right.length + 1);
        uint256 i;
        for (i = 0; i < left.length ; i++){
            result[i][0] = left[i][0];
            result[i][1] = left[i][1];
        }
        result[left.length][0] = DIVIDER;
        result[left.length][1] = DIVIDER;
        for (i = 0 ; i < right.length ; i++){
            result[left.length + 1 + i][0] = right[i][0];
            result[left.length + 1 + i][1] = right[i][1];
        }

        return result;        
    }

    function isValid(uint256[] memory scores) public pure returns (bool)
    {
        if (scores.length % 2 != 0)
            return false;
    }
}