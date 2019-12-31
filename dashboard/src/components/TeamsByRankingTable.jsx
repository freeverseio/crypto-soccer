import React, { useState, useEffect } from 'react';
import { Table, Container } from 'semantic-ui-react'
import gql from 'graphql-tag';
import { useQuery } from '@apollo/react-hooks';

const GET_TEAMS_BY_RANKING = gql`
query {
  allTeams (first: 10, orderBy: RANKING_POINTS_DESC){
    nodes {
      teamId
      name
      rankingPoints
      timezoneIdx
      countryIdx
      leagueIdx
      teamIdxInLeague
    }
  }
}
`;

export default function TeamsByRanking(props) {
    const { onTeamIdChange } = props;
    const { loading, error, data } = useQuery(GET_TEAMS_BY_RANKING, {
        pollInterval: 1000,
    });

    if (loading) return 'Loading...';
    if (error) return `Error! ${error.message}`;

    const teams = data.allTeams.nodes;

    return (
        <Table color='teal' inverted selectable >
            <Table.Header>
                <Table.Row>
                    <Table.HeaderCell>Team</Table.HeaderCell>
                    <Table.HeaderCell>Ranking</Table.HeaderCell>
                    <Table.HeaderCell>Timezone</Table.HeaderCell>
                    <Table.HeaderCell>Country</Table.HeaderCell>
                    <Table.HeaderCell>League</Table.HeaderCell>
                </Table.Row>
            </Table.Header>
            <Table.Body >
                {
                    teams.map(team => (
                        <Table.Row key={team.teamId} onClick={() => { onTeamIdChange && onTeamIdChange(team.teamId)}}>
                            <Table.Cell>{team.name}</Table.Cell>
                            <Table.Cell>{team.rankingPoints}</Table.Cell>
                            <Table.Cell>{team.timezoneIdx}</Table.Cell>
                            <Table.Cell>{team.countryIdx}</Table.Cell>
                            <Table.Cell>{team.leagueIdx}</Table.Cell>
                        </Table.Row>
                    ))
                }
            </Table.Body>
        </Table>
    );
}