import React, { useState, useEffect } from 'react';
import { Table, Button, Divider } from 'semantic-ui-react';
import RoleCard from './RoleCard';
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

    const proposeCompany = (address) => {
        proxyContract.methods.proposeCompany(address).send({
            from: account,
            gasPrice: Config.gasPrice,
        })
            .catch(console.error);
    }

    const validAddress = proposedCompany !== '0x0000000000000000000000000000000000000000' && Web3.utils.isAddress(proposedCompany);

    return (
        <Table.Row>
            <Table.Cell singleLine>Company Role</Table.Cell>
            <Table.Cell>{company}</Table.Cell>
            <Table.Cell/>
        </Table.Row>
    )
}

export default CompanyWidget;