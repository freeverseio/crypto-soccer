const UNDEF = undefined;
    
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

const deployDelegate = async (StorageProxy, Assets, AssetsView, Market, MarketView) => {
    sto = await StorageProxy.new().should.be.fulfilled;
    // setting up StorageProxy delegate calls to Assets
    assets = await Assets.at(sto.address).should.be.fulfilled;
    assetsAsLib = await Assets.new().should.be.fulfilled;
    assetsView = await AssetsView.at(sto.address).should.be.fulfilled;
    assetsViewAsLib = await AssetsView.new().should.be.fulfilled;
    market = await Market.at(sto.address).should.be.fulfilled;
    marketAsLib = await Market.new().should.be.fulfilled;
    marketView = await MarketView.at(sto.address).should.be.fulfilled;
    marketViewAsLib = await MarketView.new().should.be.fulfilled;
    
    selectorsAssets = extractSelectorsFromAbi(Assets.abi);
    selectorsAssetsView = extractSelectorsFromAbi(AssetsView.abi);
    selectorsAssets = removeDuplicatesFromFirstContract(selectorsAssets, selectorsAssetsView);
    
    selectorsMarket = extractSelectorsFromAbi(Market.abi);
    selectorsMarketView = extractSelectorsFromAbi(MarketView.abi);
    selectorsMarket = removeDuplicatesFromFirstContract(selectorsMarket, selectorsMarketView);
    
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


module.exports = {
    extractSelectorsFromAbi,
    removeDuplicatesFromFirstContract,
    deployDelegate
}

