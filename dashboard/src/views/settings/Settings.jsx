import React, { useState, useEffect } from 'react';
import { Container, Table } from 'semantic-ui-react';
import Config from '../../Config';
import COOCard from './COOCard';
import MarketCard from './MarketCard';
import RelayCard from './RelayCard';
import CryptoMarketCard from './CryptoMarketCard';
import SuperUserCard from './SuperUserCard';
import CompanyCard from './CompanyCard';
import ProposedCompanyCard from './ProposedCompanyCard';

const directoryJSON = require("../../contracts/Directory.json");

const Settings = (params) => {
    const notAvailable = 'n/a';
    const { web3, account } = params;
    const [proxyAddress, setProxyAddress] = useState(notAvailable);

    useEffect(() => {
        const directoryContract = new web3.eth.Contract(directoryJSON.abi, Config.directory_address);
        const proxyKey = web3.utils.utf8ToHex('PROXY');
        directoryContract.methods.getAddress(proxyKey).call()
            .then(setProxyAddress)
            .catch(error => { setProxyAddress(notAvailable) })
    }, [web3.eth.Contract, web3.utils]);

    return (
        <Container>
            <Table columns={2}>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell>Property</Table.HeaderCell>
                        <Table.HeaderCell>Value</Table.HeaderCell>
                    </Table.Row>
                </Table.Header>
            </Table>

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

            <Table columns={2} color='orange'>
                <Table.Body>
                    <Table.Row>
                        <Table.Cell>proxy</Table.Cell>
                        <Table.Cell>{proxyAddress}</Table.Cell>
                    </Table.Row>
                    {
                        (proxyAddress !== notAvailable) &&
                        <React.Fragment>
                            <CompanyCard web3={web3} proxyAddress={proxyAddress} />
                            <SuperUserCard web3={web3} proxyAddress={proxyAddress} />
                            <COOCard web3={web3} assetsAddress={proxyAddress} />
                            <RelayCard web3={web3} assetsAddress={proxyAddress} />
                            <MarketCard web3={web3} assetsAddress={proxyAddress} />
                            <CryptoMarketCard web3={web3} assetsAddress={proxyAddress} />
                        </React.Fragment>
                    }
                </Table.Body>
            </Table>

            {
                (proxyAddress !== notAvailable) &&
                <Table color='red' >
                    <Table.Body>
                        <ProposedCompanyCard web3={web3} account={account} proxyAddress={proxyAddress} />
                    </Table.Body>
                </Table>
            }
        </Container>
    )
}

export default Settings;