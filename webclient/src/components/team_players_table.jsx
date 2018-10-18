import React, { PureComponent } from 'react';
import { Label, Table } from 'semantic-ui-react'

class TeamPlayersTable extends PureComponent {
    render() {
        const { team } = this.props;

        const players = () => {
            if (!team)
                return <Table.Cell />

            return (
                team.players.map(player => (
                    <Table.Row key={player.index}>
                        <Table.Cell>
                            <Label ribbon>{player.index}</Label>
                        </Table.Cell>
                        <Table.Cell>{player.skills[0]}</Table.Cell>
                        <Table.Cell>{player.skills[1]}</Table.Cell>
                        <Table.Cell>{player.skills[2]}</Table.Cell>
                        <Table.Cell>{player.skills[3]}</Table.Cell>
                        <Table.Cell>{player.skills[4]}</Table.Cell>
                        <Table.Cell>{player.skills[5]}</Table.Cell>
                        <Table.Cell>{player.skills[6]}</Table.Cell>
                    </Table.Row>
                ))
            )
        }

        return (
            <Table celled>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell>{team && team.name}</Table.HeaderCell>
                        <Table.HeaderCell>Skill 0</Table.HeaderCell>
                        <Table.HeaderCell>Skill 1</Table.HeaderCell>
                        <Table.HeaderCell>Skill 2</Table.HeaderCell>
                        <Table.HeaderCell>Skill 3</Table.HeaderCell>
                        <Table.HeaderCell>Skill 4</Table.HeaderCell>
                        <Table.HeaderCell>Skill 5</Table.HeaderCell>
                        <Table.HeaderCell>Skill 6</Table.HeaderCell>
                    </Table.Row>
                </Table.Header>

                <Table.Body>
                    {players()}
               </Table.Body>
            </Table>
        )
    }
}

export default TeamPlayersTable;