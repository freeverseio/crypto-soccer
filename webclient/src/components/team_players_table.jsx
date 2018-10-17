import React, { PureComponent } from 'react';
import { Icon, Label, Menu, Table } from 'semantic-ui-react'

class TeamPlayersTable extends PureComponent {
    render() {
        const { team } = this.props;

        console.log(team);

        return (
            <Table celled>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell>Team</Table.HeaderCell>
                        <Table.HeaderCell>Header</Table.HeaderCell>
                        <Table.HeaderCell>Header</Table.HeaderCell>
                    </Table.Row>
                </Table.Header>

                <Table.Body>
                    <Table.Row>
                        <Table.Cell>
                            <Label ribbon>First</Label>
                        </Table.Cell>
                        <Table.Cell>Cell</Table.Cell>
                        <Table.Cell>Cell</Table.Cell>
                    </Table.Row>
                    <Table.Row>
                        <Table.Cell>Cell</Table.Cell>
                        <Table.Cell>Cell</Table.Cell>
                        <Table.Cell>Cell</Table.Cell>
                    </Table.Row>
                    <Table.Row>
                        <Table.Cell>Cell</Table.Cell>
                        <Table.Cell>Cell</Table.Cell>
                        <Table.Cell>Cell</Table.Cell>
                    </Table.Row>
                </Table.Body>
            </Table>
        )
    }
}

export default TeamPlayersTable;