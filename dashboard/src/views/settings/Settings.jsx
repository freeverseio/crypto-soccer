import React, { useState } from 'react';
import { Container, Table } from 'semantic-ui-react';
import Config from '../../Config';

const directoryJSON = require("../../contracts/Directory.json");
const assetsJSON = require("../../contracts/Assets.json");
const proxyJSON = require("../../contracts/Assets.json");

const notAvailable = '0x0000000000000000000000000000000000000000';

const Settings = (params) => {
    const { web3 } = params;
    const [ proxyAddress, setProxyAddress ] = useState(notAvailable);
    const [ superUserAddress, setSuperUserAddress ] = useState(notAvailable);

    const directoryContract = new web3.eth.Contract(directoryJSON.abi, Config.directory_address);
    const proxyContract = new web3.eth.Contract(proxyJSON.abi, proxyAddress);
    const assetsContract = new web3.eth.Contract(assetsJSON.abi, proxyAddress);

    const proxyKey = web3.utils.utf8ToHex('PROXY');
    directoryContract.methods.getAddress(proxyKey).call()
    .then(setProxyAddress)
    .catch(error => {setProxyAddress(notAvailable)})

    return (
        <Container>
            <Table columns={2}>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell>Property</Table.HeaderCell>
                        <Table.HeaderCell>Value</Table.HeaderCell>
                    </Table.Row>
                </Table.Header>
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
                        (proxyAddress != notAvailable) &&
                        <Table.Row>
                            <Table.Cell>superUser</Table.Cell>
                            <Table.Cell>{superUserAddress}</Table.Cell>
                        </Table.Row>
                    }
                </Table.Body>
            </Table>
        </Container>
    )
}

export default Settings;