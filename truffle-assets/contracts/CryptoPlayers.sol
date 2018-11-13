pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Full.sol";

contract CryptoPlayers is ERC721Full {
    // Mapping from token ID to its state
    mapping (uint256 => uint256) private _tokenState;
    string private _tokenCID;

    constructor(string name, string symbol, string CID) public 
    ERC721Full(name, symbol)
    {
        _tokenCID = CID;
    }

    function tokenURI(uint256 tokenId) external view returns (string) {
        require(_exists(tokenId), "unexistent token");
        string memory stateString = uint2str(_tokenState[tokenId]);
        return strConcat(_tokenCID, "/?state=", stateString);
    }

    function state(uint256 tokenId) public view returns (uint256) {
        require(_exists(tokenId), "unexistent token");
        return _tokenState[tokenId];
    }

    function _setState(uint256 tokenId, uint256 state) internal {
        require(_exists(tokenId), "unexistent token");
        _tokenState[tokenId] = state;
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

    function strConcat(string _a, string _b, string _c) internal returns (string){
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
