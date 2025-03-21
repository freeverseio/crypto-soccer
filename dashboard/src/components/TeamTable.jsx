import React from 'react';
import { Table } from 'semantic-ui-react'
import gql from 'graphql-tag';
import { useQuery } from '@apollo/react-hooks';
import Config from '../Config';
import { useState } from 'react';

const GET_TEAM = gql`
query teamByTeamId($teamId: String!){
  teamByTeamId(teamId: $teamId) {
    name
    playersByTeamId(orderBy: SHIRT_NUMBER_ASC) {
      nodes {
        playerId
        shirtNumber
        name
        defence
        pass
        speed
        endurance
        dayOfBirth
      }
    }
  }
}
`;

export default function TeamTable({ teamId, onPlayerIdChange }) {
    const [playerId, setPlayerId] = useState("");
    const { loading, error, data } = useQuery(GET_TEAM, {
        variables: { teamId },
        pollInterval: Config.polling_ms,
    });

    if (teamId === '') {
        return <div/>
    }

    if (loading) return 'Loading...';
    if (error) return `Error! ${error.message}`;
    
    const team = data.teamByTeamId
    if (!team) return <div/>;

    const players = team.playersByTeamId.nodes;

    return (
        <Table color='grey' inverted selectable >
            <Table.Header>
                <Table.Row>
                    <Table.HeaderCell>Shirt</Table.HeaderCell>
                    <Table.HeaderCell>Name</Table.HeaderCell>
                    <Table.HeaderCell>Defence</Table.HeaderCell>
                    <Table.HeaderCell>Pass</Table.HeaderCell>
                    <Table.HeaderCell>Speed</Table.HeaderCell>
                    <Table.HeaderCell>Endurance</Table.HeaderCell>
                </Table.Row>
            </Table.Header>
            <Table.Body >
                {
                    players.map(player => (
                        <Table.Row 
                            key={player.playerId}
                            active={playerId === player.playerId}
                            onClick={() => {
                                setPlayerId(player.playerId);
                                onPlayerIdChange && onPlayerIdChange(player.playerId);
                            }}>
                            <Table.Cell>{player.shirtNumber}</Table.Cell>
                            <Table.Cell>{player.name}</Table.Cell>
                            <Table.Cell>{player.defence}</Table.Cell>
                            <Table.Cell>{player.pass}</Table.Cell>
                            <Table.Cell>{player.speed}</Table.Cell>
                            <Table.Cell>{player.endurance}</Table.Cell>
                        </Table.Row>
                    ))
                }
            </Table.Body>
        </Table>
    );
}