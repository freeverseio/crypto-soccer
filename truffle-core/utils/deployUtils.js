var Web3 = require('web3');
var web3 = new Web3(Web3.givenProvider);

const deployPair = async (proxyAddress, Contr) => {
    if (Contr == "") return ["", "", []];
    selectors = extractSelectorsFromAbi(Contr.abi);
    contr = await Contr.at(proxyAddress).should.be.fulfilled;
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

const addContracts = async (proxy, addresses, allSelectors, names, firstNewContractId) => {
    // Add all contracts to ids = [firstNewContractId, firstNewContractId+1,...]
    nContracts = namesStr.length;
    newContractIds = [];
    for (c = 0; c < nContracts; c++) {
        if (allSelectors[c].length > 0) {
            newContractIds.push(firstNewContractId + c);
            tx0 = await proxy.addContract(newContractIds[c], addresses[c], allSelectors[c], names[c]).should.be.fulfilled;
        }
    }
    return newContractIds;
}

const assertActiveStatusIs = async (contractIds, status, proxy) => {
    for (c = 0; c < contractIds.length; c++) {
        var {0: addr, 1: nom, 2: sels, 3: isActive} = await proxy.getContractInfo(contractIds[c]).should.be.fulfilled;
        if (isActive != status) throw new Error("unexpected contract state");
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
  
function informNoCollisions(Assets, Market, Updates) {
    allSelectors = [];
    allSelectors = allSelectors.concat(extractSelectorsFromAbi(Assets.abi));
    allSelectors = allSelectors.concat(extractSelectorsFromAbi(Market.abi));
    allSelectors = allSelectors.concat(extractSelectorsFromAbi(Updates.abi));
    duplicates = findDuplicates(allSelectors);
    if (duplicates.length != 0) {
        console.log("There are collisions between the contracts to delegate. No panic. It is normal when they inherit common libs")
        console.log("The important thing is that there are no collisions with the proxy.")
    }
}

function assertNoCollisionsWithProxy(Proxy, Assets, Market, Updates, Challenges) {
    proxySelectors = extractSelectorsFromAbi(Proxy.abi);

    duplicates = findDuplicates(proxySelectors.concat(extractSelectorsFromAbi(Assets.abi)));
    if (duplicates.length != 0) throw new Error("duplicates found proxy-Assets!!!");

    duplicates = findDuplicates(proxySelectors.concat(extractSelectorsFromAbi(Market.abi)));
    if (duplicates.length != 0) throw new Error("duplicates found proxy-Market!!!");

    duplicates = findDuplicates(proxySelectors.concat(extractSelectorsFromAbi(Updates.abi)));
    if (duplicates.length != 0) throw new Error("duplicates found proxy-Updates!!!");

    duplicates = findDuplicates(proxySelectors.concat(extractSelectorsFromAbi(Challenges.abi)));
    if (duplicates.length != 0) throw new Error("duplicates found proxy-Challenges!!!");

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
// - proxyAddress needs only be specified if versionNumber > 0.
const deploy = async (versionNumber, Proxy, proxyAddress = "0x0", Assets, Market, Updates, Challenges) => {
    // Inform about possible collisions between contracts to delegate (among them, and with proxy)
    // informNoCollisions(Proxy, Assets, Market, Updates);
    assertNoCollisionsWithProxy(Proxy, Assets, Market, Updates, Challenges);
    
    // Next: proxy is built either by deploy, or by assignement to already deployed address
    var proxy;
    if (versionNumber == 0) {
        const proxySelectors = extractSelectorsFromAbi(Proxy.abi);
        proxy = await Proxy.new(proxySelectors).should.be.fulfilled;
    } else {
        proxy = await Proxy.at(proxyAddress).should.be.fulfilled;
    }

    // Check that the number of contracts already declared in Proxy is as expected.
    //  - contactId = 0 is null, so the first available contract on a clean deploy is 1, and every version adds 3 contracts
    const nContractsToProxy = 4;
    const firstNewContractId = 1 + versionNumber * nContractsToProxy;
    const nContractsNum = await proxy.countContracts().should.be.fulfilled;
    if (firstNewContractId != nContractsNum.toNumber()) throw new Error("mismatch between firstNewContractId and nContractsNum");

    // The following line does:
    //  - deploy new contracts (not proxy) to delegate to, and return their addresses
    //  - build interfaces to those contracts which point to the proxy address
    const {0: assets, 1: market, 2: updates, 3: challenges, 4: addresses, 5: allSelectors, 6: names} = 
        await deployContractsToDelegateTo(proxy.address, Assets, Market, Updates, Challenges);
        
    const versionedNames = appendVersionNumberToNames(names, versionNumber);

    // Adds new contracts to proxy
    const newContractIds = await addContracts(proxy, addresses, allSelectors, versionedNames, firstNewContractId);

    // Build list of contracts to deactivate
    //  - example: when deploying v1, we have activated already [0,1,2,3]
    //  - so newId = 4, and we need to deactivate [1,2,3]
    var deactivateContractIds;
    if (versionNumber == 0) {
        deactivateContractIds = [];
    } else {
        deactivateContractIds = Array.from(new Array(nContractsToProxy), (x,i) => firstNewContractId - nContractsToProxy + i);
    }

    // await assertActiveStatusIs(deactivateContractIds, true, proxy);
    // Deactivate and Activate all contracts atomically
    tx1 = await proxy.deactivateAndActivateContracts(deactivateContractIds, newContractIds).should.be.fulfilled;

    // await assertActiveStatusIs(deactivateContractIds, false, proxy);
    // await assertActiveStatusIs(newContractIds, true, proxy);
    return [proxy, assets, market, updates, challenges];
}

async function deployAndConfigureStakers(Stakers, owner, parties, updates) {
    stakers  = await Stakers.new({from:owner});
    stakers.setGameOwner(updates.address, {from:owner}).should.be.fulfilled;
    stake = await stakers.kRequiredStake();
    await addTrustedParties(stakers, owner, parties);
    await enroll(stakers, stake, parties);
    await updates.setStakersAddress(stakers.address).should.be.fulfilled;
    return stakers;
}

async function addTrustedParties(contract, owner, addresses) {
    await asyncForEach(addresses, async (address) => {
        contract.addTrustedParty(address, {from:owner})
    });
}
async function enroll(contract, stake, addresses) {
    await asyncForEach(addresses, async (address) => {
        await contract.enroll({from:address, value: stake})
    });
}

async function asyncForEach(array, callback) {
    for (let index = 0; index < array.length; index++) {
        await callback(array[index], index, array);
    }
}

async function unenroll(contract, addresses) {
    await asyncForEach(addresses, async (address) => {
        await contract.unEnroll({from:address})
    });
}


module.exports = {
    extractSelectorsFromAbi,
    informNoCollisions,
    assertNoCollisionsWithProxy,
    deploy,
    deployAndConfigureStakers,
    addTrustedParties,
    enroll,
}

