pragma solidity ^0.4.23;

contract Oracle {
    uint public deposit;
    mapping (address => uint) public solvers;

    constructor(uint _deposit) public {
        require(_deposit != 0, "deposit can't be 0");
        deposit = _deposit;
    }

    function registerSolver() public payable {
        require(msg.value == deposit, "wrong stack amount");
        require(solvers[msg.sender] == 0, "already registered");
        solvers[msg.sender] = msg.value;
    }

    function unregisterSolver() public {
        require(solvers[msg.sender] != 0, "not registered");
        msg.sender.transfer(solvers[msg.sender]);
        solvers[msg.sender] = 0;
    }

    function setResult() public returns (bool) {
        return false;
    }

    function judge() public {

    }
}