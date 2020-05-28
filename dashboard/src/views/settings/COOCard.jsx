import React, { useState, useEffect } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';
import RoleCard from './RoleCard';

const COOWidget = ({ assetsContract, account }) => {
    const [COO, setCOO] = useState();

    useEffect(() => {
        assetsContract.methods.COO().call()
            .then(setCOO)
            .catch(error => {
                console.error(error);
                setCOO("error");
            });
    }, [assetsContract]);

    const setAddress = (address) => {
        assetsContract.methods.setCOO(address).send({ from: account, gasPrice: Config.gasPrice })
            .catch(console.error);
    }

    return (
        <Table.Row>
            <Table.Cell singleLine>COO Role</Table.Cell>
            <Table.Cell>{COO}</Table.Cell>
            <Table.Cell>
                <RoleCard onChange={setAddress} account={account}/>
            </Table.Cell>
        </Table.Row>
    )
}

export default COOWidget;