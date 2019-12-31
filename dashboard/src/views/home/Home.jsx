import React, {useState} from 'react';
import { Container } from 'semantic-ui-react'
import TeamsByRankingTable from '../../components/TeamsByRankingTable'
import TeamTable from '../../components/TeamTable'

export default function Home(props) {
    const [teamId, setTeamId] = useState("");
    return (
        <Container>
            <TeamsByRankingTable  onTeamIdChange={setTeamId}/>
            <TeamTable teamId={teamId}/>
        </Container>
    );
}
