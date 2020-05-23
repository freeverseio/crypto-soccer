import React from 'react';
import { Container, Table } from 'semantic-ui-react';
import Config from '../../Config';

const directoryJSON = require("../../contracts/Directory.json");

export default function settings(params) {
    const { web3 } = params;
    // const directory = new web3.eth.Contract(directoryJSON.abi, Config.directory_address);
    // const d = directory.methods.getDirectory().call().then(console.log);
    // console.log(d)
    // const [playerid, setplayerid] = usestate("");
    const entries = Object.entries(Config);

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