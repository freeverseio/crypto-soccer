pragma solidity ^0.5.0;

contract Sort {
    
    function sort11(uint8[11] memory data) public pure returns(uint8[11] memory) {
       quickSort11(data, int(0), int(10));
       return data;
    }
    
    function quickSort11(uint8[11] memory arr, int left, int right) internal pure {
        int i = left;
        int j = right;
        if(i==j) return;
        uint pivot = arr[uint(left + (right - left) / 2)];
        while (i <= j) {
            while (arr[uint(i)] < pivot) i++;
            while (pivot < arr[uint(j)]) j--;
            if (i <= j) {
                (arr[uint(i)], arr[uint(j)]) = (arr[uint(j)], arr[uint(i)]);
                i++;
                j--;
            }
        }
        if (left < j)
            quickSort11(arr, left, j);
        if (i < right)
            quickSort11(arr, i, right);
    }

}

