import React, { useState, useEffect } from 'react';
import { Table, Input, Button } from 'semantic-ui-react';
import Config from '../../Config';
const Web3 = require('web3');

const CompanyWidget = ({web3, proxyContract, account}) => {
    const [company, setCompany] = useState();
    const [proposedCompany, setProposedCompany] = useState();
    const [address, setAddress] = useState("");

    useEffect(() => {
        proxyContract.methods.company().call()
            .then(setCompany)
            .catch(error => { setCompany("") })

        proxyContract.methods.proposedCompany().call()
            .then(setProposedCompany)
            .catch(error => { setProposedCompany("") })
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
                    disabled
                    icon='ethereum'
                    iconPosition='left'
                    placeholder={proposedCompany}
                    value={address}
                    onChange={event => setAddress(event.target.value)}
                    action={
                        <Button size='mini' color='red' onClick={accept} disabled={!validAddress}>Accept</Button>
                    }
                />
            </Table.Cell>
        </Table.Row>
    )
}

export default CompanyWidget;