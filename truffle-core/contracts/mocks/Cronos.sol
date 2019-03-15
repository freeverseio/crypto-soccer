pragma solidity ^0.5.0;

contract Cronos {
    uint256 private _counter;

    function wait() public {
        _counter++;
    }
}