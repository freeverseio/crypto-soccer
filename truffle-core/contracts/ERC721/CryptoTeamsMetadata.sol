pragma solidity ^0.4.24;

import "./CryptoTeamsBase.sol";
import "openzeppelin-solidity/contracts/introspection/ERC165.sol";
import "openzeppelin-solidity/contracts/token/ERC721/IERC721Metadata.sol";

contract CryptoTeamsMetadata is ERC165, CryptoTeamsBase, IERC721Metadata {
    // Token name
    string constant private _name = "CryptoSoccerTeams";

    // Token symbol
    string constant private _symbol = "CST";

    string private _teamsURI;

    bytes4 private constant InterfaceId_ERC721Metadata = 0x5b5e139f;
    /**
     * 0x5b5e139f ===
     *     bytes4(keccak256('name()')) ^
     *     bytes4(keccak256('symbol()')) ^
     *     bytes4(keccak256('tokenURI(uint256)'))
     */

    /**
     * @dev Constructor function
     */
    constructor () public {
        // register the supported interfaces to conform to ERC721 via ERC165
        _registerInterface(InterfaceId_ERC721Metadata);
    }

    /**
     * @dev Gets the token name
     * @return string representing the token name
     */
    function name() external view returns (string) {
        return _name;
    }

    /**
     * @dev Gets the token symbol
     * @return string representing the token symbol
     */
    function symbol() external view returns (string) {
        return _symbol;
    }

    /**
     * @dev Returns an URI for a given token ID
     * Throws if the token ID does not exist. May return an empty string.
     * @param tokenId uint256 ID of the token to query
     */
    function tokenURI(uint256 tokenId) external view returns (string) {
        require(_exists(tokenId));
        return _teamsURI;
    }

    /**
     * @dev Internal function to set the token URI for all token
     * @param uri string URI to assign
     */
    function setTokensURI(string uri) public {
        _teamsURI = uri;
    }
}

