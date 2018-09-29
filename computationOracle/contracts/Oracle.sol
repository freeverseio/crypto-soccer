pragma solidity ^0.4.23;

contract Oracle {
    uint public stackAmount;
    mapping (address => uint) public solvers;

    constructor(uint _stackAmount) public {
        require(_stackAmount != 0, "deposit can't be 0");
        stackAmount = _stackAmount;
    }

    function registerSolver() public payable {
        require(msg.value == stackAmount, "wrong stack amount");
        require(solvers[msg.sender] == 0, "already registered");
        solvers[msg.sender] = msg.value;
    }

    function unregisterSolver() public {
        require(solvers[msg.sender] != 0, "not registered");
        solvers[msg.sender] = 0;
    }

    function setResult() public returns (bool) {
        return false;
    }

    function judge() public {

    }
}