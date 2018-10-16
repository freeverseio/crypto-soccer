var Testing = artifacts.require("Testing");
const League = artifacts.require("League");

module.exports = function(deployer) { 
  deployer.deploy(Testing, "0x1234");
  deployer.deploy(League, "0x1234");
};

