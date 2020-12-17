import React, { useState, useEffect } from 'react';
import { Table } from 'semantic-ui-react';

const CompanyWidget = ({proxyContract, account}) => {
    const [company, setCompany] = useState();

    useEffect(() => {
        proxyContract.methods.company().call()
            .then(setCompany)
            .catch(error => {
                console.error(error);
                setCompany("error");
            });
    }, [proxyContract]);

    return (
        <Table.Row>
            <Table.Cell singleLine>Company Role</Table.Cell>
            <Table.Cell>{company}</Table.Cell>
            <Table.Cell/>
        </Table.Row>
    )
}

export default CompanyWidget;