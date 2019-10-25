pragma solidity ^0.5.0;

contract SortIdxs {
    
    uint256 constant private N_IDXS = 8;
    
    function sortIdxs(uint8[N_IDXS] memory data) public pure returns(uint8[N_IDXS] memory) {
       uint8[N_IDXS] memory idxs;
       for (uint8 i = 0; i < N_IDXS; i++) idxs[i] = i;
       quickSort(data, idxs, int(0), int(N_IDXS - 1));
       return idxs;
    }
    
    function quickSort(uint8[N_IDXS] memory arr, uint8[N_IDXS] memory idxs, int left, int right) internal pure {
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
            quickSort(arr, idxs, left, j);
        if (i < right)
            quickSort(arr, idxs, i, right);
    }

}

