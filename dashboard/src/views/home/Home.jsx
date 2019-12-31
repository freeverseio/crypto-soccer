import React from 'react';
import { Container } from 'semantic-ui-react'
import TeamsByRankingTable from '../../components/TeamsByRankingTable'
import TeamTable from '../../components/TeamTable'

export default function Home(props) {
    return (
        <Container>
            <TeamsByRankingTable  onTeamIdChange={console.log}/>
            <TeamTable />
        </Container>
    );
}
