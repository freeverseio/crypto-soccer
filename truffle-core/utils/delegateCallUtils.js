var Web3 = require('web3');
var web3 = new Web3(Web3.givenProvider);

const deployPair = async (proxy, Contr) => {
    if (Contr == "") return ["", "", []];
    contr = await Contr.at(proxy.address).should.be.fulfilled;
    contrAsLib = await Contr.new().should.be.fulfilled;
    selectors = extractSelectorsFromAbi(Contr.abi);
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
            tx0 = await proxy.addContract_magicx(contractIds[c], addresses[c], allSelectors[c], names[c]).should.be.fulfilled;
        }
    }

    // Activate all contracts atomically
    tx1 = await proxy.deactivateAndActivateContracts_magicx(deactivate = [], activate = contractIds).should.be.fulfilled;

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


module.exports = {
    extractSelectorsFromAbi,
    deployDelegate
}

