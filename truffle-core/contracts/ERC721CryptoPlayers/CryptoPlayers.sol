pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Full.sol";

contract CryptoPlayers is ERC721Full("CryptoSoccerPlayers", "CSP") {
    // Mapping from token ID to its state
    mapping (uint256 => uint256) private _tokenState;

    function getState(uint256 tokenId) public view returns (uint256) {
        require(_exists(tokenId), "unexistent token");
        return _tokenState[tokenId];
    }

    function _setState(uint256 tokenId, uint256 state) internal {
        require(_exists(tokenId), "unexistent token");
        _tokenState[tokenId] = state;
    }
}
