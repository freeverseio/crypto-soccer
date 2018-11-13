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
        return _tokenCID;
    }

    function defense(uint256 tokenId) public view returns (uint256) {
        require(_exists(tokenId), "unexistent token");
        return _tokenState[tokenId];
    }
}
