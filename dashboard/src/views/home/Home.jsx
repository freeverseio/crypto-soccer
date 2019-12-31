import React, { useState, useEffect } from 'react';
import { Table, Container } from 'semantic-ui-react'
import gql from 'graphql-tag';
import { useQuery } from '@apollo/react-hooks';

const GET_10_TEAMS_BY_RANKING = gql`
query {
  allTeams (first: 1, orderBy: RANKING_POINTS_DESC){
    nodes {
      teamId
      rankingPoints
    }
  }
}
`;

export default function TableExampleInvertedColors() {
    return (
        <Container>
            <Table color='teal' inverted>
                <Table.Header>
                    <Table.Row>
                        <Table.HeaderCell>Team</Table.HeaderCell>
                        <Table.HeaderCell>Country</Table.HeaderCell>
                        <Table.HeaderCell>League</Table.HeaderCell>
                        <Table.HeaderCell>Ranking</Table.HeaderCell>
                    </Table.Row>
                </Table.Header>
                <Table.Body>
                    <Table.Row>
                        <Table.Cell>Apples</Table.Cell>
                        <Table.Cell>200</Table.Cell>
                        <Table.Cell>0g</Table.Cell>
                    </Table.Row>
                    <Table.Row>
                        <Table.Cell>Orange</Table.Cell>
                        <Table.Cell>310</Table.Cell>
                        <Table.Cell>0g</Table.Cell>
                    </Table.Row>
                </Table.Body>
            </Table>
        </Container>
    );
}
