import React, { useState, useEffect } from 'react';
import { Table, Button, Divider } from 'semantic-ui-react';
import RoleCard from './RoleCard';
import Config from '../../Config';
const Web3 = require('web3');

const CompanyWidget = ({proxyContract, multisigContract, account}) => {
    const [company, setCompany] = useState();
    const [owners, setOwners] = useState([]);
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

    useEffect(() => {
        multisigContract.methods.getOwners().call()
            .then(setOwners)
            .catch(error => {
                console.error(error);
                setOwners("error");
            });
    }, [multisigContract]);

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
        <React.Fragment>
            <Table.Row>
                <Table.Cell singleLine>Company Role</Table.Cell>
                <Table.Cell>{company}</Table.Cell>
            </Table.Row>
            <Table.Row>
                <Table.Cell singleLine>Company owners</Table.Cell>
                <Table.Cell>
                    {owners.map(owner => owner + " ")}
                </Table.Cell>
            </Table.Row>
        </React.Fragment>
    )
}

export default CompanyWidget;