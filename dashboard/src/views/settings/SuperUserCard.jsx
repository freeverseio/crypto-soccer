import React, { useState } from 'react';
import { Table } from 'semantic-ui-react';

const proxyJSON = require("../../contracts/Proxy.json");

const SuperUserWidget = (props) => {
    const {web3, proxyAddress} = props;
    const [superUser, setSuperUser] = useState("");

    const proxyContract = new web3.eth.Contract(proxyJSON.abi, proxyAddress);
    proxyContract.methods.superUser().call()
        .then(setSuperUser)
        .catch(error => { setSuperUser("") })

    return (
        <Table.Row>
            <Table.Cell>SuperUser</Table.Cell>
            <Table.Cell>{superUser}</Table.Cell>
        </Table.Row>
    )
}

export default SuperUserWidget;