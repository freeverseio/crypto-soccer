const deployPair = async (sto, Contr, ContrView) => {
    contr = await Contr.at(sto.address).should.be.fulfilled;
    contrAsLib = await Contr.new().should.be.fulfilled;
    contrView = await ContrView.at(sto.address).should.be.fulfilled;
    contrViewAsLib = await ContrView.new().should.be.fulfilled;
    selectors = extractSelectorsFromAbi(Contr.abi);
    selectorsView = extractSelectorsFromAbi(ContrView.abi);
    selectors = removeDuplicatesFromFirstContract(selectors, selectorsView);
    return [contr, contrAsLib, contrViewAsLib, selectors, selectorsView]
};

const deployDelegate = async (StorageProxy, Assets, AssetsView, Market, MarketView) => {
    sto = await StorageProxy.new().should.be.fulfilled;
    // setting up StorageProxy delegate calls to Assets
    const {0: assets, 1: assetsAsLib, 2: assetsViewAsLib, 3: selectorsAssets, 4: selectorsAssetsView} = await deployPair(sto, Assets, AssetsView);
    const {0: market, 1: marketAsLib, 2: marketViewAsLib, 3: selectorsMarket, 4: selectorsMarketView} = await deployPair(sto, Market, MarketView);
    
    namesStr            = ['Assets',            'AssetsView',           'Market',           'MarketView'];
    contractsAsLib      = [assetsAsLib,         assetsViewAsLib,        marketAsLib,        marketViewAsLib];
    allSelectors        = [selectorsAssets,     selectorsAssetsView,    selectorsMarket,    selectorsMarketView];
    requiresPermission  = [true,                false,                  true,               false];

    addresses = [];                 
    names = [];
    addresses = [];
    contractIds = [];

    nContracts = requiresPermission.length;
    for (c = 0; c < nContracts; c++) {
        names.push(toBytes32(namesStr[c]));
        addresses.push(contractsAsLib[c].address);
        contractIds.push(c+1);
    }
    
    // Add all contracts to predicted Ids: [1, 2, ...]
    for (c = 0; c < nContracts; c++) {
        tx0 = await sto.addContract(contractIds[c], addresses[c], requiresPermission[c], allSelectors[c], names[c]).should.be.fulfilled;
    }

    // Activate all contracts atomically
    tx1 = await sto.deleteAndActivateContracts(deactivate = [], activate = contractIds).should.be.fulfilled;

    return [assets, market];

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

function removeDuplicatesFromFirstContract(selectors1, selectors2) {
    selectors1_filtered = [];
    for (s = 0; s < selectors1.length; s++) {
        if (!selectors2.includes(selectors1[s])) {
            selectors1_filtered.push(selectors1[s]);
        }
    }
    return selectors1_filtered;
}

function toBytes32(name) { return web3.utils.utf8ToHex(name); }


module.exports = {
    extractSelectorsFromAbi,
    removeDuplicatesFromFirstContract,
    deployDelegate
}

