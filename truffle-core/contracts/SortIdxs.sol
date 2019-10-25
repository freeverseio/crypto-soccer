pragma solidity ^0.5.0;

contract SortIdxs {
    
    function sort14(uint8[14] memory data) public pure returns(uint8[14] memory) {
       uint8[14] memory idxs;
       for (uint8 i = 0; i < 14; i++) idxs[i] = i;
       quickSort14(data, int(0), int(13), idxs);
       return idxs;
    }
    
    function quickSort14(uint8[14] memory arr, int left, int right, uint8[14] memory idxs) internal pure {
        int i = left;
        int j = right;
        if(i==j) return;
        uint pivot = arr[uint(left + (right - left) / 2)];
        while (i <= j) {
            while (arr[uint(i)] < pivot) i++;
            while (pivot < arr[uint(j)]) j--;
            if (i <= j) {
                (arr[uint(i)], arr[uint(j)]) = (arr[uint(j)], arr[uint(i)]);
                (idxs[uint(i)], idxs[uint(j)]) = (idxs[uint(j)], idxs[uint(i)]);
                i++;
                j--;
            }
        }
        if (left < j)
            quickSort14(arr, left, j, idxs);
        if (i < right)
            quickSort14(arr, i, right, idxs);
    }

}

