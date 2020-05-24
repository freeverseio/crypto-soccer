import React, { useState } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';
import RoleCard from './RoleCard';

const proxyJSON = require("../../contracts/Proxy.json");

const SuperUserWidget = ({web3, proxyAddress, account}) => {
    const [superUser, setSuperUser] = useState("");

    const proxyContract = new web3.eth.Contract(proxyJSON.abi, proxyAddress);
    proxyContract.methods.superUser().call()
        .then(setSuperUser)
        .catch(error => { setSuperUser("") })

    const setAddress = (address) => {
        proxyContract.methods.setSuperUser(address).send({ from: account, gasPrice: Config.gasPrice })
            .on('error', (error, receipt) => { console.error(error) });
    }

    return (
        <Table.Row>
            <Table.Cell>SuperUser</Table.Cell>
            <Table.Cell>{superUser}</Table.Cell>
            <Table.Cell>
                <RoleCard onChange={setAddress} />
            </Table.Cell>
        </Table.Row>
    )
}

export default SuperUserWidget;