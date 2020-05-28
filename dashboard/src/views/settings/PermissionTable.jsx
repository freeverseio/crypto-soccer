import React from 'react';
import { Table } from 'semantic-ui-react';
import COOCard from './COOCard';
import MarketCard from './MarketCard';
import RelayCard from './RelayCard';
import CryptoMarketCard from './CryptoMarketCard';
import SuperUserCard from './SuperUserCard';
import CompanyCard from './CompanyCard';

const proxyJSON = require("../../contracts/Proxy.json");
const assetsJSON = require("../../contracts/Assets.json");
const marketJSON = require("../../contracts/Market.json");

const PermissionTable = ({ web3, account, proxyAddress }) => {
    const proxyContract = new web3.eth.Contract(proxyJSON.abi, proxyAddress);
    const assetsContract = new web3.eth.Contract(assetsJSON.abi, proxyAddress);
    const marketContract = new web3.eth.Contract(marketJSON.abi, proxyAddress);

    return (
        <Table color='orange'>
            <Table.Header>
                <Table.Row>
                    <Table.HeaderCell width={1}></Table.HeaderCell>
                    <Table.HeaderCell width={1}></Table.HeaderCell>
                    <Table.HeaderCell width='six'></Table.HeaderCell>
                </Table.Row>
            </Table.Header>
            <Table.Body>
                <Table.Row>
                    <Table.Cell>proxy</Table.Cell>
                    <Table.Cell>{proxyAddress}</Table.Cell>
                </Table.Row>
                {proxyContract && <CompanyCard account={account} proxyContract={proxyContract} />}
                {proxyContract && <SuperUserCard account={account} proxyContract={proxyContract} />}
                {assetsContract && <COOCard account={account} assetsContract={assetsContract} />}
                {assetsContract && <RelayCard account={account} assetsContract={assetsContract} />}
                {assetsContract && <MarketCard account={account} assetsContract={assetsContract} />}
                {marketContract && <CryptoMarketCard account={account} marketContract={marketContract} />}
            </Table.Body>
        </Table>
    )
}

export default PermissionTable;