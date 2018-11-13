pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Full.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Mintable.sol";

contract CryptoPlayers is ERC721Full, ERC721Mintable {
    struct Props {
        uint8 defense;
        uint8 attack;
    }

    // Mapping from token ID to its props
    mapping (uint256 => Props) private _tokenProps;
    string private _tokenCID;

    constructor(string name, string symbol, string CID) public 
    ERC721Full(name, symbol)
    {
        _tokenCID = CID;
    }

    function tokenURI(uint256 tokenId) external view returns (string) {
        require(_exists(tokenId), "unexistent token");
        return _tokenCID;
    }

    function defense(uint256 tokenId) public view returns (uint) {
        require(_exists(tokenId), "unexistent token");
        return _tokenProps[tokenId].defense;
    }
}
