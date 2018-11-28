require('chai')
  .use(require('chai-as-promised'))
  .should();

const CryptoPlayers = artifacts.require('CryptoPlayers');
const CryptoTeams = artifacts.require('CryptoTeams');
const cryptoSoccer = artifacts.require("TeamFactoryMock");
var k = require('../../jsCommons/constants.js');
var f = require('../../jsCommons/functions.js');

contract('Players', function(accounts) {

  var instance;
  var nTotalPlayers=0;

  const playerName="Messi";
  const teamName="Barca";
  const teamIdx = 1; 
  const userChoice=1;
  const playerNumberInTeam=2;
  const playerRole=3; 

  console.log('Funds in the source account:');
  console.log(web3.eth.getBalance(web3.eth.accounts[0]).toNumber()/web3.toWei(1, "ether"));

  it("tests if contract is deployed correctly", async () => {
    const cryptoPlayers = await CryptoPlayers.new().should.be.fulfilled;
    const cryptoTeams = await CryptoTeams.new().should.be.fulfilled;
    instance = await cryptoSoccer.new(cryptoTeams.address, cryptoPlayers.address);
  });

  it('team name of unexistent player', async () => {
    await instance.teamNameByPlayer("unexistent").should.be.rejected;
  });

  it('team name by player', async () => {
    const cryptoPlayers = await CryptoPlayers.new().should.be.fulfilled;
    const cryptoTeams = await CryptoTeams.new().should.be.fulfilled;
    const contract = await cryptoSoccer.new(cryptoTeams.address, cryptoPlayers.address);
    const team = "team";
    const player = "player";
    const playerState = 44535;
    await cryptoTeams.addTeam(team, accounts[0]);
    await cryptoPlayers.addPlayer(player, playerState, 1);
    const index = await cryptoPlayers.getTeamIndexByPlayer(player).should.be.fulfilled;
    index.toNumber().should.be.equal(1);
    const name = await contract.teamNameByPlayer(player).should.be.fulfilled;
    name.should.be.equal(team);
  });

  it("creates an empty team, checks that nTeams moves from 0 to 1", async () =>{
    var nTeams = await  instance.test_getNCreatedTeams.call();
    assert.isTrue(nTeams==0);
    await instance.test_createTeam(teamName);
    var nTeams = await  instance.test_getNCreatedTeams.call();
    assert.isTrue(nTeams==1);
  });

  it("adds a player to the previously created empty team, and checks nPlayers goes from 0 to 1", async () =>{
    var nTeams = await  instance.test_getNCreatedTeams.call().should.be.fulfilled;
    assert.equal(nTeams,1);
    nTotalPlayers = await instance.test_getNCreatedPlayers.call().should.be.fulfilled;
    assert.equal(nTotalPlayers,0); 
    await instance.test_createBalancedPlayer(playerName,teamIdx,userChoice,playerNumberInTeam,playerRole).should.be.fulfilled;
    nTotalPlayers = await instance.test_getNCreatedPlayers.call();
    assert.equal(nTotalPlayers,1);
  });

  it("reads skills of a player and check it is as expected", async () =>{
    playerIdx = 0;
    var state = await instance.test_getPlayerState(playerIdx);
    var decoded = await instance.test_decode(k.NumStates, state, k.BitsPerState);
    console.log(decoded);
    expected = [193, 47, 61, 46, 34, 62, 3]; 
    info = "Player " + playerIdx+ " skills: "
    for (var st=0; st<k.NumStates; st++) {
        thisState = decoded[st].toNumber();
        assert(thisState == expected[st]);
        info += " " + thisState;
    }
    console.log(info);
  });

  it("tries to add a player with the same name, and checks that it fails", async () =>{
    var nTeams = await  instance.test_getNCreatedTeams.call();
    var nPlayers = await  instance.test_getNCreatedPlayers.call();
    var lastPlayerName = await instance.test_getPlayerName(0);
    console.log("Teams created so far " + nTeams + " team, and nPlayers = " + nPlayers);
    console.log("lastPlayerName = " + lastPlayerName);
    hasFailed = false;
    try{ 
      // If you create a player with an existing name, it won't let you, no matter what the rest of stuff is
     await instance.test_createBalancedPlayer(lastPlayerName,teamIdx,userChoice+1,playerNumberInTeam+1,playerRole);
    } catch (err) {
      // Great, the transaction failed
      hasFailed = true;
    }
    assert.isTrue(hasFailed);
  });


})


