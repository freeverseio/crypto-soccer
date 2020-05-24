import React, { useState, useEffect } from 'react';
import { Table, Input, Button } from 'semantic-ui-react';
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

    const validAddress = Web3.utils.isAddress(proposedCompany);

    return (
        <Table.Row>
            <Table.Cell>Company</Table.Cell>
            <Table.Cell>{company}</Table.Cell>
            <Table.Cell>
                <Input fluid
                    size='mini'
                    error={!validAddress}
                    icon='ethereum'
                    iconPosition='left'
                    value={proposedCompany}
                    action={
                        <Button size='mini' color='red' onClick={accept} disabled={!validAddress || !account}>Accept</Button>
                    }
                />
            </Table.Cell>
        </Table.Row>
    )
}

export default CompanyWidget;