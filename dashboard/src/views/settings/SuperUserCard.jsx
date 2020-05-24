import React, { useState, useEffect } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';
import RoleCard from './RoleCard';

const SuperUserWidget = ({ proxyContract, account }) => {
    const [superUser, setSuperUser] = useState();

    useEffect(() => {
        proxyContract.methods.superUser().call()
            .then(setSuperUser)
            .catch(error => {
                console.error(error);
                setSuperUser("error");
            });
    }, [proxyContract]);

    const setAddress = (address) => {
        proxyContract.methods.setSuperUser(address).send({ from: account, gasPrice: Config.gasPrice })
        .catch(console.error);
    }

    return (
        <Table.Row>
            <Table.Cell>SuperUser</Table.Cell>
            <Table.Cell>{superUser}</Table.Cell>
            <Table.Cell>
                <RoleCard onChange={setAddress} account={account}/>
            </Table.Cell>
        </Table.Row>
    )
}

export default SuperUserWidget;