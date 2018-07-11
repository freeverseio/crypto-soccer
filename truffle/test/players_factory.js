const cryptoSoccer = artifacts.require("GameEngine");

contract('Players', function(accounts) {

  var instance;
  var nTotalPlayers=0;
  var sourceBalance;
  console.log('Funds in the source account:');
  console.log(web3.eth.getBalance(web3.eth.accounts[0]).toNumber()/web3.toWei(1, "ether"));

  //// deploy the contract before each test
  //beforeEach(async () => {
  //  instance = await cryptoSoccer.deployed();
  //});

  it("tests if contract is deployed correctly", async () => {
    instance = await cryptoSoccer.deployed();
    assert.isTrue(true);
  });

  it("creates an empty team, checks that nTeams moves from 0 to 1", async () =>{
    var nTeams = await  instance.getNCreatedTeams.call();
    assert.isTrue(nTeams==0);
    await instance.createTeam("Barca");
    var nTeams = await  instance.getNCreatedTeams.call();
    assert.isTrue(nTeams==1);
  });

  it("adds a player to the previously created empty team, and checks nPlayers goes from 0 to 1", async () =>{
    var nTeams = await  instance.getNCreatedTeams.call();
    assert.equal(nTeams,1);
    const playerName="Messi";
    const teamName="Barca";
    const teamIdx = 0; // todo: get this idx from the team name, using the contract mappings
    const userChoice=1;
    const playerNumberInTeam=2;
    const playerRole=3; 
    nTotalPlayers = await instance.getNCreatedPlayers.call();
    assert.equal(nTotalPlayers,0);
    await instance.createRandomPlayer(playerName,teamIdx,userChoice,playerNumberInTeam,playerRole);
    nTotalPlayers = await instance.getNCreatedPlayers.call();
    assert.equal(nTotalPlayers,1);
    playerState = await instance.getPlayerState.call(0);
    var states = await instance.readNumbersFromUint.call(7,playerState,10000);
  });

  it("tries to add a player with the same name, and checks that it fails", async () =>{
    nTotalPlayers = await instance.getNCreatedPlayers.call();
    hasFailed = false;
    try{ 
      // If you create a player with an existing name, it won't let you, no matter what the rest of stuff is
      await instance.createRandomPlayer("Messi",12,43,22,33);
    } catch (err) {
      // Great, the transaction failed
      hasFailed = true;
    }
    assert.isTrue(hasFailed);
  });

})


