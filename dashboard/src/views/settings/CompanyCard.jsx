import React, { useState } from 'react';
import { Table } from 'semantic-ui-react';
import ProposedCompanyCard from './ProposedCompanyCard';

const proxyJSON = require("../../contracts/Proxy.json");

const CompanyWidget = ({web3, proxyAddress, account}) => {
    const [company, setCompany] = useState("");
    const [proposedCompany, setProposedCompany] = useState("");

    const proxyContract = new web3.eth.Contract(proxyJSON.abi, proxyAddress);
    proxyContract.methods.company().call()
        .then(setCompany)
        .catch(error => { setCompany("") })

    proxyContract.methods.proposedCompany().call()
        .then(setProposedCompany)
        .catch(error => { setProposedCompany("") })

    return (
        <Table.Row>
            <Table.Cell>Company</Table.Cell>
            <Table.Cell>{company}</Table.Cell>
            <Table.Cell>
                <ProposedCompanyCard
                    web3={web3}
                    account={account}
                    proxyAddress={proxyAddress}
                    proposedCompany={proposedCompany}
                />
            </Table.Cell>
        </Table.Row>
    )
}

export default CompanyWidget;