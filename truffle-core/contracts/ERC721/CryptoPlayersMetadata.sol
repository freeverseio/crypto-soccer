pragma solidity ^0.4.24;

import "./CryptoPlayersBase.sol";
import "../helpers.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";

contract CryptoPlayersMetadata is ERC721Metadata("CryptoSoccerPlayers", "CSP"), CryptoPlayersBase, HelperFunctions {
    string private _URI;

    function _setTokensURI(string uri) internal { // TODO add modifier
        _URI = uri;
    }

    function tokenURI(uint256 tokenId) external view returns (string) {
        require(_exists(tokenId), "unexistent token");
        uint256 state = getState(tokenId);
        string memory stateString = uint2str(state);
        return strConcat(_URI, "?state=", stateString);
    }
}