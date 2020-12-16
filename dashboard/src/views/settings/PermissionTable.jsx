import React, { useState, useEffect } from 'react';
import { Table } from 'semantic-ui-react';
import COOCard from './COOCard';
import MarketCard from './MarketCard';
import RelayCard from './RelayCard';
import CryptoMarketCard from './CryptoMarketCard';
import SuperUserCard from './SuperUserCard';
import CompanyCard from './CompanyCard';
import Config from '../../Config';

const proxyJSON = require("../../contracts/Proxy.json");
const assetsJSON = require("../../contracts/Assets.json");
const marketJSON = require("../../contracts/Market.json");
const multisigJSON = require("../../contracts/MultiSigWallet.json");

const PermissionTable = ({ web3, account, proxyAddress }) => {
    const [seconds, setSeconds] = useState(0);
    const proxyContract = new web3.eth.Contract(proxyJSON.abi, proxyAddress);
    const assetsContract = new web3.eth.Contract(assetsJSON.abi, proxyAddress);
    const marketContract = new web3.eth.Contract(marketJSON.abi, proxyAddress);
    const multisigContract = new web3.eth.Contract(multisigJSON.abi, Config.multiSigAddress);

    console.log(multisigContract)

    useEffect(() => {
        const interval = setInterval(() => {
            setSeconds(seconds => seconds + 1);
        }, 5000);
        return () => clearInterval(interval);
    },[seconds]);

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
                    <Table.Cell>Proxy Address</Table.Cell>
                    <Table.Cell>{proxyAddress}</Table.Cell>
                </Table.Row>
                {proxyContract && <CompanyCard account={account} proxyContract={proxyContract}  multisigContract={multisigContract}/>}
                {proxyContract && <SuperUserCard account={account} proxyContract={proxyContract} multisigContract={multisigContract}/>}
                {assetsContract && <COOCard account={account} assetsContract={assetsContract} />}
                {assetsContract && <RelayCard account={account} assetsContract={assetsContract} />}
                {assetsContract && <MarketCard account={account} assetsContract={assetsContract} />}
                {marketContract && <CryptoMarketCard account={account} marketContract={marketContract} />}
            </Table.Body>
        </Table>
    )
}

export default PermissionTable;