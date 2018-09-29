pragma solidity ^0.4.23;

contract Oracle {
    uint public stackAmount;
    mapping (address => uint) public solvers;

    constructor(uint _stackAmount) public {
        stackAmount = _stackAmount;
    }

    function registerSolver() public payable {
        require(msg.value == stackAmount, "wrong stack amount");
        require(solvers[msg.sender] == 0, "already registered");
        solvers[msg.sender] = msg.value;
    }

    function unregisterSolver() public {

    }

    function setResult() public returns (bool) {
        return false;
    }

    function judge() public {

    }
}