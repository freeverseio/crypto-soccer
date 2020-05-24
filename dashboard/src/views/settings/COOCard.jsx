import React, { useState } from 'react';
import { Table } from 'semantic-ui-react';
import Config from '../../Config';
import RoleCard from './RoleCard';

const assetsJSON = require("../../contracts/Assets.json");

const COOWidget = ({web3, assetsAddress, account}) => {
    const [COO, setCOO] = useState("");

    const assetsContract = new web3.eth.Contract(assetsJSON.abi, assetsAddress);
    assetsContract.methods.COO().call()
        .then(setCOO)
        .catch(error => { setCOO("") })

    const setAddress = (address) => {
        assetsContract.methods.setCOO(address).send({ from: account, gasPrice: Config.gasPrice })
            .on('error', (error, receipt) => { console.error(error) });
    }

    return (
        <Table.Row>
            <Table.Cell>COO</Table.Cell>
            <Table.Cell>{COO}</Table.Cell>
            <Table.Cell>
                <RoleCard onChange={setAddress} disabled={account ? false : true} />
            </Table.Cell>
        </Table.Row>
    )
}

export default COOWidget;