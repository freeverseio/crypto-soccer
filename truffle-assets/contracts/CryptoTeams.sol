pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Enumerable.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";
import "openzeppelin-solidity/contracts/access/roles/MinterRole.sol";

contract CryptoTeams is ERC721Enumerable, ERC721Metadata, MinterRole {
    // Mapping from team ID to its name
    mapping (uint256 => string) private _teamName;

    constructor(string name, string symbol) ERC721Metadata(name, symbol)
        public
    {
    }

    function getName(uint256 teamId) public view returns (string) {
        require(_exists(teamId), "unexistent team");
        return _teamName[teamId];
    }

    function _setName(uint256 teamId, string name) public view returns (string) {
        require(_exists(teamId), "unexistent team");
        return _teamName[teamId] = name;
    }

    /**
    * @dev Function to mint tokens
    * @param to The address that will receive the minted tokens.
    * @param tokenId The token id to mint.
    * @param name The token name
    * @return A boolean that indicates if the operation was successful.
    */
    function mint(
        address to,
        uint256 tokenId,
        string name
    )
    public
    onlyMinter
    returns (bool)
    {
        _mint(to, tokenId);
        _setName(tokenId, name);
        return true;
    }
}

