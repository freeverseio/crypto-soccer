pragma solidity ^0.4.23;

contract Oracle {
    uint public stackAmount;

    constructor(uint _stackAmount) public {
        stackAmount = _stackAmount;
    }

    function registerSolver() public payable returns (bool) {
        require(msg.value == stackAmount, "wrong stack amount");
        return true;
    }

    function unregisterSolver() public {

    }

    function setResult() public returns (bool) {
        return false;
    }

    function judge() public {

    }
}