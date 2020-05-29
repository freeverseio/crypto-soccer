import React, { useState, useEffect } from 'react';
import { Table, Button } from 'semantic-ui-react';
import Config from '../../Config';
const Web3 = require('web3');

const CompanyWidget = ({proxyContract, account}) => {
    const [company, setCompany] = useState();
    const [proposedCompany, setProposedCompany] = useState();

    useEffect(() => {
        proxyContract.methods.company().call()
            .then(setCompany)
            .catch(error => {
                console.error(error);
                setCompany("error");
            });

        proxyContract.methods.proposedCompany().call()
            .then(setProposedCompany)
            .catch(error => {
                console.error(error);
                setProposedCompany("error");
            });
    }, [proxyContract]);

    const accept = () => {
        proxyContract.methods.acceptCompany().send({
            from: account,
            gasPrice: Config.gasPrice,
        })
            .catch(console.error);
    };

    const validAddress = proposedCompany !== '0x0000000000000000000000000000000000000000' && Web3.utils.isAddress(proposedCompany);

    return (
        <Table.Row>
            <Table.Cell singleLine>Company Role</Table.Cell>
            <Table.Cell>{company}</Table.Cell>
            <Table.Cell >
                <Button.Group size='mini' fluid>
                    <Button color='grey' onClick={accept} disabled={true}>{proposedCompany}</Button>
                    <Button color='red' onClick={accept} disabled={!validAddress || !account}>Accept</Button>
                </Button.Group>
            </Table.Cell>
        </Table.Row>
    )
}

export default CompanyWidget;