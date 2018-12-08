pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";
import "./CryptoPlayersLink.sol";
import "./URIerRole.sol";
import "../helpers.sol";

contract CryptoPlayersMetadata is ERC721Metadata("CryptoSoccerPlayers", "CSP"), CryptoPlayersLink, URIerRole, HelperFunctions {
    string private _URI;

    function _setTokensURI(string uri) internal onlyURIer { 
        _URI = uri;
    }

    function tokenURI(uint256 tokenId) external view returns (string) {
        require(_exists(tokenId), "unexistent token");
        uint256 state = getState(tokenId);
        string memory stateString = uint2str(state);
        return strConcat(_URI, "?state=", stateString);
    }
}