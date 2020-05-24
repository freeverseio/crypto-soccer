import React, { useState } from 'react';
import { Table } from 'semantic-ui-react';

const assetsJSON = require("../../contracts/Assets.json");

const COOWidget = (props) => {
    const {web3, assetsAddress} = props;
    const [COO, setCOO] = useState("");

    const assetsContract = new web3.eth.Contract(assetsJSON.abi, assetsAddress);
    assetsContract.methods.COO().call()
        .then(setCOO)
        .catch(error => { setCOO("") })

    return (
        <Table.Row>
            <Table.Cell>COO</Table.Cell>
            <Table.Cell>{COO}</Table.Cell>
        </Table.Row>
    )
}

export default COOWidget;