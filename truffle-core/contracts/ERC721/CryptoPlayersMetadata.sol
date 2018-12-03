pragma solidity ^0.4.24;

import "./CryptoPlayersBase.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";

contract CryptoPlayersMetadata is CryptoPlayersBase, ERC721Metadata("CryptoSoccerPlayers", "CSP") {
    string private _baseURI;

    function setBaseURI(string URI) public {
        _baseURI = URI;
    }

    function getBaseURI() external view returns (string) {
        return _baseURI;
    }

    // function tokenURI(uint256 tokenId) external view returns (string) {
    //     require(_exists(tokenId), "unexistent token");
    //     uint256 state = getState(tokenId);
    //     string memory stateString = uint2str(state);
    //     return strConcat(_baseURI, "/?state=", stateString);
    // }
}