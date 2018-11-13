pragma solidity ^0.4.24;

import "./CryptoPlayers.sol";
import "openzeppelin-solidity/contracts/access/roles/MinterRole.sol";

contract CryptoPlayersMintable is CryptoPlayers, MinterRole {
    constructor(string name, string symbol, string CID) public 
    CryptoPlayers(name, symbol, CID)
    {
    }

  /**
   * @dev Function to mint tokens
   * @param to The address that will receive the minted tokens.
   * @param tokenId The token id to mint.
   * @return A boolean that indicates if the operation was successful.
   */
    function mint(
        address to,
        uint256 tokenId
    )
        public
        onlyMinter
        returns (bool)
    {
       _mint(to, tokenId);
        return true;
    }
}
