import React, { useState, useEffect } from 'react';
import { Container, Table } from 'semantic-ui-react';
import Config from '../../Config';
import COOCard from './COOCard';
import MarketCard from './MarketCard';
import RelayCard from './RelayCard';
import CryptoMarketCard from './CryptoMarketCard';
import SuperUserCard from './SuperUserCard';
import CompanyCard from './CompanyCard';

const directoryJSON = require("../../contracts/Directory.json");
const proxyJSON = require("../../contracts/Proxy.json");
const assetsJSON = require("../../contracts/Assets.json");
const marketJSON = require("../../contracts/Market.json");

const Settings = ( { web3, account }) => {
    const [proxyAddress, setProxyAddress] = useState();
    const [proxyContract, setProxyContract] = useState();
    const [assetsContract, setAssetsContract] = useState();
    const [marketContract, setMarketContract] = useState();

    useEffect(() => {
        const directoryContract = new web3.eth.Contract(directoryJSON.abi, Config.directory_address);
        const proxyKey = web3.utils.utf8ToHex('PROXY');
        directoryContract.methods.getAddress(proxyKey).call()
            .then(proxyAddress => {
                setProxyAddress(proxyAddress);
                const proxyContract = new web3.eth.Contract(proxyJSON.abi, proxyAddress);
                setProxyContract(proxyContract);
                const assetsContract = new web3.eth.Contract(assetsJSON.abi, proxyAddress);
                setAssetsContract(assetsContract);
                const marketContract = new web3.eth.Contract(marketJSON.abi, proxyAddress);
                setMarketContract(marketContract);
            })
            .catch(error => {
                console.error(error);
                setProxyAddress("error");
            });
    }, [web3]);

    return (
        <Container>
            <Table columns={2} color='blue'>
                <Table.Body>
                    {
                        Object.entries(Config).map((entry, i) => (
                            <Table.Row key={entry[0]}>
                                <Table.Cell>{entry[0]}</Table.Cell>
                                <Table.Cell>{entry[1]}</Table.Cell>
                            </Table.Row>
                        ))
                    }
                </Table.Body>
            </Table>

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
        </Container>
    )
}

export default Settings;