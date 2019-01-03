pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";
import "./CryptoPlayersTeam.sol";
import "./URIerRole.sol";

contract CryptoPlayersMetadata is ERC721Metadata("CryptoSoccerPlayers", "CSP"), CryptoPlayersTeam, URIerRole {
    string private _URI;

    function _setTokensURI(string uri) internal onlyURIer { 
        _URI = uri;
    }

    function tokenURI(uint256 tokenId) external view returns (string) {
        require(_exists(tokenId), "unexistent token");
        uint88 genome = getGenome(tokenId);
        string memory genomeString = uint2str(genome);
        return strConcat(_URI, "?genome=", genomeString);
    }
}