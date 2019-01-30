pragma solidity ^ 0.4.24;

contract Leagues {
    uint256 private _init;
    uint256 private _final;

    function getInit() public view returns (uint256) {
        return _init;
    }

    function getFinal() public view returns (uint256) {
        return _final;
    }
}