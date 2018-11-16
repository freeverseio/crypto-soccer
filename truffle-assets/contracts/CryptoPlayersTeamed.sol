pragma solidity ^0.4.24;

import "./CryptoPlayersMintable.sol";

contract CryptoPlayersTeamed is CryptoPlayersMintable {
    // Mapping from team ID to its name
    mapping (uint256 => string) private _teamName;


}
