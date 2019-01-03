pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";
import "./CryptoTeamsPlayers.sol";
import "./URIerRole.sol";
import "../HelperFunctions.sol";

contract CryptoTeamsMetadata is ERC721Metadata("CryptoSoccerTeams", "CST"), CryptoTeamsPlayers, URIerRole, HelperFunctions  {
    string private _teamsURI;

    /**
     * @dev Returns an URI for a given token ID
     * Throws if the token ID does not exist. May return an empty string.
     * @param tokenId uint256 ID of the token to query
     */
    function tokenURI(uint256 tokenId) external view returns (string) {
        require(_exists(tokenId));
        uint256 playersId = 0; // TODO getPlayersIds(tokenId);
        string memory playersIdString = uint2str(playersId);
        string memory uri = strConcat(_teamsURI, "?playersId=", playersIdString);
        return uri;
    }

    /**
     * @dev Internal function to set the token URI for all token
     * @param uri string URI to assign
     */
    function _setTokensURI(string uri) internal onlyURIer { 
        _teamsURI = uri;
    }
}

