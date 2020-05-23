import React from 'react';
import { Container, Table } from 'semantic-ui-react';
import Config from '../../Config';

export default function settings(params) {
    // const [playerid, setplayerid] = usestate("");
    const entries = Object.entries(Config) ;

    console.log(entries)

    return (
        <Container>
            <Table columns={2}>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell>Property</Table.HeaderCell>
                        <Table.HeaderCell>Value</Table.HeaderCell>
                    </Table.Row>
                </Table.Header>

                <Table.Body>
                    {
                        entries.map((entry, i) => (
                            <Table.Row key={entry[0]}>
                                <Table.Cell>{entry[0]}</Table.Cell>
                                <Table.Cell>{entry[1]}</Table.Cell>
                            </Table.Row>
                        ))
                    }
                </Table.Body>

            </Table>
        </Container>
    )
}