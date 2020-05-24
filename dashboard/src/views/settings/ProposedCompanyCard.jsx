import React, { useState } from 'react';
import { Table, Button } from 'semantic-ui-react';
import Config from '../../Config';

const proxyJSON = require("../../contracts/Proxy.json");

const ProposedCompanyWidget = (props) => {
    const { web3, proxyAddress, account } = props;
    const [company, setProposedCompany] = useState("");

    const proxyContract = new web3.eth.Contract(proxyJSON.abi, proxyAddress);
    proxyContract.methods.proposedCompany().call()
        .then(setProposedCompany)
        .catch(error => { setProposedCompany("") })

    const accept = () => {
        proxyContract.methods.acceptCompany().send({
            from: account,
            gasPrice: Config.gasPrice,
        })
            .on('transactionHash', hash => { console.log(hash) })
            .on('confirmation', (confirmationNumber, receipt) => {
                console.log("confiramtion");
            })
            .on('receipt', (receipt) => {
                console.log(receipt);
            })
            .on('error', (error, receipt) => { console.log(error) });
    };

    return (
        <Table.Row>
            <Table.Cell>ProposedCompany</Table.Cell>
            <Table.Cell>{company}</Table.Cell>
            <Table.Cell>
                <Button size='mini' color='red' onClick={accept} disabled={account ? false : true}>Accept</Button>
            </Table.Cell>
        </Table.Row>
    )
}

export default ProposedCompanyWidget;