pragma solidity ^0.5.0;

import "./Evolution.sol";

contract PlayAndEvolve {

    Evolution private _evo;

    function setEvolutionAddress(address addr) public {
        _evo = Evolution(addr);
    }

    

}

