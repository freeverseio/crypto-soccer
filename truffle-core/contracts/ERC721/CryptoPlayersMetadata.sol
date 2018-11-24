pragma solidity ^0.4.24;

import "./CryptoPlayers.sol";

contract CryptoPlayersMetadata is CryptoPlayers {
    string private _tokenCID;

    function tokenURI(uint256 tokenId) external view returns (string) {
        require(_exists(tokenId), "unexistent token");
        uint256 state = getState(tokenId);
        string memory stateString = uint2str(state);
        return strConcat(_tokenCID, "/?state=", stateString);
    }

    function uint2str(uint i) internal pure returns (string){
        if (i == 0) return "0";
        uint j = i;
        uint length;
        while (j != 0){
            length++;
            j /= 10;
        }
        bytes memory bstr = new bytes(length);
        uint k = length - 1;
        uint tmp = i;
        while (tmp != 0){
            bstr[k--] = byte(48 + tmp % 10);
            tmp /= 10;
        }
        return string(bstr);
    }

    function strConcat(string _a, string _b, string _c) internal pure returns (string){
        bytes memory _ba = bytes(_a);
        bytes memory _bb = bytes(_b);
        bytes memory _bc = bytes(_c);
        string memory abcde = new string(_ba.length + _bb.length + _bc.length );
        bytes memory babcde = bytes(abcde);
        uint k = 0;
        for (uint i = 0; i < _ba.length; i++) babcde[k++] = _ba[i];
        for (i = 0; i < _bb.length; i++) babcde[k++] = _bb[i];
        for (i = 0; i < _bc.length; i++) babcde[k++] = _bc[i];
        return string(babcde);
    }
}