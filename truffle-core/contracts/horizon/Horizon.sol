pragma solidity ^0.4.24;

import "../ERC721/CryptoPlayers.sol";
import "../ERC721/CryptoTeams.sol";

contract Horizon {
    CryptoPlayers private _cryptoPlayers;
    CryptoTeams private _cryptoTeams;

    constructor(address cryptoPlayers, address cryptoTeams) public {
        _cryptoPlayers = CryptoPlayers(cryptoPlayers);
        _cryptoTeams = CryptoTeams(cryptoTeams);
    }
}