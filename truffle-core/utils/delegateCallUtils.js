var Web3 = require('web3');
var web3 = new Web3(Web3.givenProvider);

const deployPair = async (proxy, Contr) => {
    if (Contr == "") return ["", "", []];
    selectors = extractSelectorsFromAbi(Contr.abi);
    contr = await Contr.at(proxy.address).should.be.fulfilled;
    contrAsLib = await Contr.new().should.be.fulfilled;
    return [contr, contrAsLib, selectors];
};

const deployDelegate = async (proxy, Assets, Market = "", Updates = "") => {
    // setting up StorageProxy delegate calls to Assets
    const {0: assets, 1: assetsAsLib, 2: selectorsAssets} = await deployPair(proxy, Assets);
    const {0: market, 1: marketAsLib, 2: selectorsMarket} = await deployPair(proxy, Market);
    const {0: updates, 1: updatesAsLib, 2: selectorsUpdates} = await deployPair(proxy, Updates);
    
    namesStr            = ['Assets', 'Market', 'Updates'];
    contractsAsLib      = [assetsAsLib, marketAsLib, updatesAsLib];
    allSelectors        = [selectorsAssets, selectorsMarket, selectorsUpdates];

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
    
    // Add all contracts to predicted Ids: [1, 2, ...]
    for (c = 0; c < nContracts; c++) {
        if (allSelectors[c].length > 0) {
            tx0 = await proxy.addContract(contractIds[c], addresses[c], allSelectors[c], names[c]).should.be.fulfilled;
        }
    }

    // Activate all contracts atomically
    tx1 = await proxy.deactivateAndActivateContracts(deactivate = [], activate = contractIds).should.be.fulfilled;

    return [assets, market, updates, allSelectors];

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

function assertNoCollisionsWithProxy(Proxy, Assets, Market, Updates) {
    proxySelectors = extractSelectorsFromAbi(Proxy.abi);

    duplicates = findDuplicates(proxySelectors.concat(extractSelectorsFromAbi(Assets.abi)));
    if (duplicates.length != 0) throw new Error("duplicates found proxy-Assets!!!");

    duplicates = findDuplicates(proxySelectors.concat(extractSelectorsFromAbi(Market.abi)));
    if (duplicates.length != 0) throw new Error("duplicates found proxy-Market!!!");

    duplicates = findDuplicates(proxySelectors.concat(extractSelectorsFromAbi(Updates.abi)));
    if (duplicates.length != 0) throw new Error("duplicates found proxy-Updates!!!");

    console.log("No collisions were found with the proxy.")
}

module.exports = {
    extractSelectorsFromAbi,
    deployDelegate,
    informNoCollisions,
    assertNoCollisionsWithProxy
}

