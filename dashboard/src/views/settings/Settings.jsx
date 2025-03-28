import React from 'react';
import { Container, Table } from 'semantic-ui-react';
import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import Config from '../../Config';
import PermissionTable from './PermissionTable';
import MultisigTable from './MultisigTable';

const GET_PROXY_ADDRESS = gql`
    {
        paramByName(name: "PROXY"){
            value
        }
    }
`;

const Settings = ({ web3, account }) => {
    const { loading, error, data } = useQuery(GET_PROXY_ADDRESS);

    if (loading) return 'Loading...';
    if (error) return `Error! ${error.message}`;

    const proxyAddress = data.paramByName.value;

    return (
        <Container>
            <Table columns={2} color='blue'>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell colSpan='3' textAlign='center'>options</Table.HeaderCell>
                    </Table.Row>
                </Table.Header>
                <Table.Body>
                    {
                        Object.entries(Config).map((entry, i) => (
                            <Table.Row key={entry[0]}>
                                <Table.Cell>{entry[0]}</Table.Cell>
                                <Table.Cell>{entry[1]}</Table.Cell>
                            </Table.Row>
                        ))
                    }
                </Table.Body>
            </Table>
            <MultisigTable web3={web3} account={account} proxyAddress={proxyAddress} />
            <PermissionTable web3={web3} account={account} proxyAddress={proxyAddress} />
        </Container>
    )
}

export default Settings;