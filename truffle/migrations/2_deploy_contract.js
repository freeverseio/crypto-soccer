var GameEngine = artifacts.require("GameEngine");
var HelperFunctions = artifacts.require("HelperFunctions");

module.exports = function(deployer) { 
  deployer.deploy(GameEngine, "0x1234");
  deployer.deploy(HelperFunctions, "0x1235");
};

