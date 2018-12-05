// require('chai')
//     .use(require('chai-as-promised'))
//     .should();

// const CryptoPlayers = artifacts.require('CryptoPlayers');
// const CryptoTeams = artifacts.require('CryptoTeams');
// const TeamFactoryMock = artifacts.require("TeamFactoryMock");
// var k = require('../../jsCommons/constants.js');
// var f = require('../../jsCommons/functions.js');

// const skillNames = ["Age","Defense","Speed","Pass","Shoot","Endurance","Role"];

// contract('Teams', function (accounts) {
//   let instance;
//   let cryptoPlayers;
//   let cryptoTeams;

//   beforeEach(async () => {
//     cryptoPlayers = await CryptoPlayers.new().should.be.fulfilled;
//     cryptoTeams = await CryptoTeams.new().should.be.fulfilled;
//     instance = await TeamFactoryMock.new(cryptoTeams.address, cryptoPlayers.address).should.be.fulfilled;
//     cryptoTeams.addMinter(instance.address);
//   });

//   it("creates a single contract and computes the gas cost of deploying GameEngine", async () => {
//     var receipt = await web3.eth.getTransactionReceipt(instance.transactionHash);
//     console.log("GameEngine\n\tdeployment cost: ", receipt.gasUsed, "\n\tcontract address:", receipt.contractAddress)
//     assert.isTrue(receipt.gasUsed > 2000000);
//   });

//   it('get unexistent team', async () => {
//     await cryptoTeams.getName(0).should.be.rejected;
//     await cryptoTeams.getName(1).should.be.rejected;
//   })

//   it("creates an entire team, an checks that we have 11 players at the end", async () => {
//     let nCreatedPlayers = await instance.test_getNCreatedPlayers.call().should.be.fulfilled;
//     nCreatedPlayers.toNumber().should.be.equal(0);
//     teamName = "Mataro";
//     playerBasename = "Bogarde";
//     await f.createTeam(instance, teamName, playerBasename, k.MaxPlayersInTeam, f.createAlineacion(4, 3, 3), accounts[0]).should.be.fulfilled;
//     nCreatedPlayers = await instance.test_getNCreatedPlayers.call().should.be.fulfilled;
//     nCreatedPlayers.toNumber().should.be.equal(k.MaxPlayersInTeam)
//   });

//   it("create team", async () => {
//     const name = "Los Cojos";
//     await instance.createTeam(name).should.be.fulfilled;
//     const result = await cryptoTeams.getName(1).should.be.fulfilled;
//     result.should.be.equal(name)
//   });

//   it("checks that we cannot add 2 teams with same name", async () => {
//     const name = "Los Cojos";
//     await instance.createTeam(name).should.be.fulfilled;
//     await instance.createTeam(name).should.be.rejected;
//   });

//   it("creates a team via .call() instead of Tx and checks that you can create 2 teams with same name", async () => {
//     teamName = "test";
//     var newTeamIdx = await instance.totalSupply.call();
//     await instance.test_createTeam.call(teamName);
//     var newTeamIdx2 = await instance.totalSupply.call();
//     assert.equal(newTeamIdx.toNumber(), newTeamIdx2.toNumber()); // meaning that nothing has been stored in the blockchain
//     await instance.test_createTeam.call(teamName);
//   });
// });


