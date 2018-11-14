pragma solidity ^0.4.24;

import "./CryptoPlayersMetadata.sol";
import "./CryptoPlayersMintable.sol";

contract CryptoPlayersFull is CryptoPlayersMetadata, CryptoPlayersMintable {
    constructor(string CID) public 
    CryptoPlayersMetadata(CID)
    {
    }
}
