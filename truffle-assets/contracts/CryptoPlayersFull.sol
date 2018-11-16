pragma solidity ^0.4.24;

import "./CryptoPlayersMetadata.sol";
import "./CryptoPlayersTeamed.sol";

contract CryptoPlayersFull is CryptoPlayersMetadata, CryptoPlayersTeamed {
    constructor(string CID) public 
    CryptoPlayersMetadata(CID)
    {
    }
}
