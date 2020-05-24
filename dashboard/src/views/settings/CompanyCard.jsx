import React, { useState } from 'react';
import { Table } from 'semantic-ui-react';

const proxyJSON = require("../../contracts/Proxy.json");

const CompanyWidget = (props) => {
    const {web3, proxyAddress} = props;
    const [company, setCompany] = useState("");

    const proxyContract = new web3.eth.Contract(proxyJSON.abi, proxyAddress);
    proxyContract.methods.company().call()
        .then(setCompany)
        .catch(error => { setCompany("") })

    return (
        <Table.Row>
            <Table.Cell>Company</Table.Cell>
            <Table.Cell>{company}</Table.Cell>
        </Table.Row>
    )
}

export default CompanyWidget;