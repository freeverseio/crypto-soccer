pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";
import "./URIerRole.sol";

contract ERC721MetadataBaseURI is ERC721Metadata, URIerRole {
    string private _URI;

    constructor (string name, string symbol) ERC721Metadata(name, symbol) public {}

    function getBaseTokenURI() external view returns (string) {
        return _URI;
    }

    function setBaseTokenURI(string uri) public onlyURIer { 
        _URI = uri;
    }

    function tokenURI(uint256 tokenId) external view returns (string) {
        require(_exists(tokenId), "unexistent token");
        return strConcat(
            _URI, 
            uint2str(tokenId),
            ""
            );
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