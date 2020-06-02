var Web3 = require('web3');
var assert = require('assert')
var web3 = new Web3(Web3.givenProvider);
const NULL_ADDR = web3.utils.toHex("0");

const deployPair = async (proxyAddress, Contr) => {
    if (Contr == "") return ["", "", []];
    selectors = extractSelectorsFromAbi(Contr.abi);
    contr = await Contr.at(proxyAddress).should.be.fulfilled;
    let contrAsLib
    contrAsLib = await Contr.new().should.be.fulfilled;
    return [contr, contrAsLib, selectors];
};

const deployContractsToDelegateTo = async (proxyAddress, Assets, Market, Updates, Challenges) => {
    // setting up StorageProxy delegate calls to Assets
    const {0: assets, 1: assetsAsLib, 2: selectorsAssets} = await deployPair(proxyAddress, Assets);
    const {0: market, 1: marketAsLib, 2: selectorsMarket} = await deployPair(proxyAddress, Market);
    const {0: updates, 1: updatesAsLib, 2: selectorsUpdates} = await deployPair(proxyAddress, Updates);
    const {0: challenges, 1: challengesAsLib, 2: selectorsChallenges} = await deployPair(proxyAddress, Challenges);
    
    namesStr            = ['Assets', 'Market', 'Updates', 'Challenges'];
    contractsAsLib      = [assetsAsLib, marketAsLib, updatesAsLib, challengesAsLib];
    allSelectors        = [selectorsAssets, selectorsMarket, selectorsUpdates, selectorsChallenges];

    addresses = [];                 
    names = [];
    addresses = [];
    contractIds = [];

    nContracts = namesStr.length;
    for (c = 0; c < nContracts; c++) {
        if (allSelectors[c].length > 0) {
            names.push(toBytes32(namesStr[c]));
            addresses.push(contractsAsLib[c].address);
            contractIds.push(c+1);
        }
    }
    return [assets, market, updates, challenges, addresses, allSelectors, names];
}

const addContracts = async (superuser, proxy, addresses, allSelectors, names, firstNewContractId) => {
    // Add all contracts to ids = [firstNewContractId, firstNewContractId+1,...]
    nContracts = namesStr.length;
    newContractIds = [];
    concatSelectors = [];
    nSelectorsPerContract = [];
    for (c = 0; c < nContracts; c++) {
        if (allSelectors[c].length > 0) {
            newContractIds.push(firstNewContractId + c);
        }
        nSelectorsPerContract.push(allSelectors[c].length);
        concatSelectors = concatSelectors.concat(allSelectors[c]);
    }
    tx0 = await proxy.addContracts(newContractIds, addresses, nSelectorsPerContract, concatSelectors, names, {from: superuser}).should.be.fulfilled;
    return newContractIds;
}

const assertActiveStatusIs = async (contractIds, status, proxy) => {
    for (c = 0; c < contractIds.length; c++) {
        var {0: addr, 1: nom, 2: sels, 3: isActive} = await proxy.getContractInfo(contractIds[c]).should.be.fulfilled;
        assert.equal(isActive, status, "unexpected contract state");
    }
}

function extractSelectorsFromAbi(abi) {
    functions = [];
    for (i = 0; i < abi.length; i++) { 
        if (abi[i].type == "function") {
            functions.push(web3.eth.abi.encodeFunctionSignature(abi[i]));
        }
    }    
    return functions;
}

function toBytes32(name) { return web3.utils.utf8ToHex(name); }
function fromBytes32(name) { return web3.utils.hexToUtf8(name); }


function findDuplicates(data) {
    let result = [];
    for (i = 0; i < data.length; i++) {
        thisEntry = data[i]
        for (j = 0; j < i; j++) {
            if (thisEntry == data[j]) {
                result.push(thisEntry);
            }
        }
    }
    return result;
}
  
function informNoCollisions(contractsArray) {
    allSelectors = [];
    for (contract of contractsArray){ 
        allSelectors = allSelectors.concat(extractSelectorsFromAbi(contract.abi));
    }
    duplicates = findDuplicates(allSelectors);
    if (duplicates.length != 0) {
        console.log("There are collisions between the contracts to delegate. No panic. It is normal when they inherit common libs")
        console.log("The important thing is that there are no collisions with the proxy.")
    }
}

function assertNoCollisionsWithProxy(Proxy, Assets, Market, Updates, Challenges) {
    proxySelectors = extractSelectorsFromAbi(Proxy.abi);

    duplicates = findDuplicates(proxySelectors.concat(extractSelectorsFromAbi(Assets.abi)));
    assert.equal(duplicates.length, 0, "duplicates found proxy-Assets!!!");

    duplicates = findDuplicates(proxySelectors.concat(extractSelectorsFromAbi(Market.abi)));
    assert.equal(duplicates.length, 0, "duplicates found proxy-Market!!!");

    duplicates = findDuplicates(proxySelectors.concat(extractSelectorsFromAbi(Updates.abi)));
    assert.equal(duplicates.length, 0, "duplicates found proxy-Updates!!!");

    duplicates = findDuplicates(proxySelectors.concat(extractSelectorsFromAbi(Challenges.abi)));
    assert.equal(duplicates.length, 0, "duplicates found proxy-Challenges!!!");

    // console.log("No collisions were found with the proxy.")
}

function appendVersionNumberToNames(names, versionNumber) {
    newNames = [...names];
    for (n = 0; n < newNames.length; n++) {
        newNames[n] = toBytes32( fromBytes32(newNames[n]) + versionNumber.toString() );
    }
    return newNames;
}

// - versionNumber = 0 for first deploy
const deploy = async (owners, Proxy, Assets, Market, Updates, Challenges) => {
    assertNoCollisionsWithProxy(Proxy, Assets, Market, Updates, Challenges);

    // Optionally inform about duplicates between the contracts themselves:
    // informNoCollisions(Assets, Market, Updates, Challenges);
    
    // Next: proxy is built either by deploy, or by assignement to already deployed address
    const proxySelectors = extractSelectorsFromAbi(Proxy.abi);
    const proxy = await Proxy.new(owners.company, owners.superuser, proxySelectors).should.be.fulfilled;

    // Check that the number of contracts already declared in Proxy is as expected.
    //  - contactId = 0 is null, so the first available contract on a clean deploy is 1, and every version adds 3 contracts
    const nContractsToProxy = 4;
    const firstNewContractId = 1
    const nContractsNum = await proxy.countContracts().should.be.fulfilled;
    assert.equal(firstNewContractId, nContractsNum.toNumber(), "mismatch between firstNewContractId and nContractsNum");

    // The following line does:
    //  - deploy new contracts (not proxy) to delegate to, and return their addresses
    //  - build interfaces to those contracts which point to the proxy address
    const {0: assets, 1: market, 2: updates, 3: challenges, 4: addresses, 5: allSelectors, 6: names} = 
        await deployContractsToDelegateTo(proxy.address, Assets, Market, Updates, Challenges);
        
    const versionedNames = appendVersionNumberToNames(names, versionNumber = 0);

    // Adds new contracts to proxy
    const newContractIds = await addContracts(owners.superuser, proxy, addresses, allSelectors, versionedNames, firstNewContractId);

    // await assertActiveStatusIs(deactivateContractIds, true, proxy);
    // Deactivate and Activate all contracts atomically
    tx1 = await proxy.activateContracts(newContractIds, {from: owners.superuser}).should.be.fulfilled;

    // await assertActiveStatusIs(deactivateContractIds, false, proxy);
    // await assertActiveStatusIs(newContractIds, true, proxy);
    return [proxy, assets, market, updates, challenges];
}

// - versionNumber = 0 for first deploy
// - proxyAddress needs only be specified for upgrades
// Step 1: deploy all contracts except for proxy (permisionless)
//  - deploy all non-proxy, non-directory contracts
//  - deploy directory pointing to these new-deployed contracts
// Step 2: "proxy.addContracts" (superUser)
//  - build input parameters, such as: "selectors"...
// Step 3: "proxy.upgrade", (superUser)
//  - atomic: deactivateOld + activateNew + setNewDirectoryAddress    
const upgrade = async (versionNumber, owners, Proxy, proxyAddress, Assets, Market, Updates, Challenges, Directory, namesAndAddresses) => {
    assert.notEqual(versionNumber, 0, "version number must be larger than 0 for upgrades");
    assert.notEqual(proxyAddress, "0x0", "proxyAddress must different from 0x0 for upgrades");
    
    assertNoCollisionsWithProxy(Proxy, Assets, Market, Updates, Challenges);

    // Optionally inform about duplicates between the contracts themselves
    // informNoCollisions(Assets, Market, Updates, Challenges);
    
    const proxy = await Proxy.at(proxyAddress).should.be.fulfilled;

    // Check that the number of contracts already declared in Proxy is as expected.
    //  - contactId = 0 is null, so the first available contract on a clean deploy is 1, and every version adds 3 contracts
    const nContractsToProxy = 4;
    const firstNewContractId = 1 + versionNumber * nContractsToProxy;
    const nContractsNum = await proxy.countContracts().should.be.fulfilled;
    assert.equal(firstNewContractId, nContractsNum.toNumber(), "mismatch between firstNewContractId and nContractsNum");

    // The following line does:
    //  - deploy new contracts (not proxy) to delegate to, and return their addresses
    //  - build interfaces to those contracts which point to the proxy address
    const {0: assets, 1: market, 2: updates, 3: challenges, 4: addresses, 5: allSelectors, 6: names} = 
        await deployContractsToDelegateTo(proxy.address, Assets, Market, Updates, Challenges);
        
    const versionedNames = appendVersionNumberToNames(names, versionNumber);

    // Adds new contracts to proxy in one single TX signed by owners.superuser
    const newContractIds = await addContracts(owners.superuser, proxy, addresses, allSelectors, versionedNames, firstNewContractId);

    // Stores new addresses in Directory contract
    const {0: dummy, 1: nonProxyNames, 2: nonProxyAddresses} = splitNamesAndAdresses(namesAndAddresses);
    directory = await Directory.new(nonProxyNames, nonProxyAddresses).should.be.fulfilled;
    
    // Build list of contracts to deactivate
    //  - example: when deploying v1, we have activated already [0,1,2,3]
    //  - so newId = 4, and we need to deactivate [1,2,3]
    const deactivateContractIds = Array.from(new Array(nContractsToProxy), (x,i) => firstNewContractId - nContractsToProxy + i);

    // await assertActiveStatusIs(deactivateContractIds, true, proxy);
    // Deactivate and Activate all contracts atomically
    // And point to the new directory
    await proxy.upgrade(deactivateContractIds, newContractIds, directory.address, {from: owners.superuser}).should.be.fulfilled;

    return [proxy, assets, market, updates, challenges];
}

function splitNamesAndAdresses(namesAndAddresses) {    
    names = [];
    namesBytes32 = [];
    addresses = [];
    for (c = 0; c < namesAndAddresses.length; c++) {
        names.push(namesAndAddresses[c][0]);
        namesBytes32.push(web3.utils.utf8ToHex(namesAndAddresses[c][0]));
        addresses.push(namesAndAddresses[c][1]);
    }
    return [names, namesBytes32, addresses];
}

async function addTrustedParties(contract, owner, addresses) {
    await asyncForEach(addresses, async (address) => {
        await contract.addTrustedParty(address, {from:owner}).should.be.fulfilled;
    });
}
async function enrol(contract, stake, addresses) {
    await asyncForEach(addresses, async (address) => {
        await contract.enrol({from:address, value: stake}).should.be.fulfilled;
    });
}

async function asyncForEach(array, callback) {
    for (let index = 0; index < array.length; index++) {
        await callback(array[index], index, array);
    }
}

async function unenroll(contract, addresses) {
    await asyncForEach(addresses, async (address) => {
        await contract.unEnroll({from:address}).should.be.fulfilled;
    });
}

function getDefaultSetup(accounts) {
    return {
      singleTimezone: -1,
      owners: {
        company:  accounts[0],
        superuser:  accounts[1],
        COO:  accounts[2],
        market:  accounts[3],
        relay:  accounts[4],
        trustedParties: [accounts[5]]
      },
      requiredStake: 1000000000000,
    }
  }
  
  function getAccount0Owner(account0) {
    return {
          company:  account0,
          superuser:  account0,
          COO:  account0,
          market:  account0,
          relay:  account0,
          trustedParties: [account0]
      }
  }
  
  function getExplicitOrDefaultSetup(networkParams, accounts) {
    const { singleTimezone, owners, requiredStake } = networkParams;
    // Safety check: either ALL or NONE of the networkParams must be defined (otherwise, expect having forgotten to assign some)
    numDefined = (singleTimezone ? 1 : 0) +  (owners ? 1 : 0) + (requiredStake ? 1 : 0);
    isValidSetup = (numDefined == 3) || (numDefined == 0);
    assert.equal(isValidSetup, true, "only some of the setup parameters are assigned in deployer.networks");
    // Set up default values only if needed:
    needsDefaultValues = (numDefined == 0);
    return needsDefaultValues ? getDefaultSetup(accounts) : networkParams;
  }
  
async function setProxyContractOwners(proxy, assets, owners, prevCompany) {
    // Order matters. First, company is established:
    await proxy.proposeCompany(owners.company, {from: prevCompany}).should.be.fulfilled;
    await proxy.acceptCompany({from: owners.company}).should.be.fulfilled;
    // company authorizes superUser to do everything else:
    await proxy.setSuperUser(owners.superuser, {from: owners.company}).should.be.fulfilled;
    // finally, superUser sets the rest of the roles:
    await assets.setCOO(owners.COO, {from: owners.superuser}).should.be.fulfilled;
    await assets.setMarket(owners.market, {from: owners.superuser}).should.be.fulfilled;
    await assets.setRelay(owners.relay, {from: owners.superuser}).should.be.fulfilled;
  }


module.exports = {
    extractSelectorsFromAbi,
    informNoCollisions,
    assertNoCollisionsWithProxy,
    deploy,
    addTrustedParties,
    enrol,
    unenroll,
    getExplicitOrDefaultSetup,
    getDefaultSetup,
    setProxyContractOwners,
    getAccount0Owner,
    upgrade,
    splitNamesAndAdresses
}

