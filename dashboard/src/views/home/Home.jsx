import React, {useState} from 'react';
import { Container } from 'semantic-ui-react';
import TeamsByRankingTable from '../../components/TeamsByRankingTable';
import TeamTable from '../../components/TeamTable';
import PlayerChart from '../../components/PlayerChart';

export default function Home(props) {
    const [teamId, setTeamId] = useState("");
    const [playerId, setPlayerId] = useState("");

    return (
        <Container>
            <TeamsByRankingTable onTeamIdChange={setTeamId} />
            <TeamTable teamId={teamId} onPlayerIdChange={setPlayerId}/>
            <Container style={{"height": "300px"}}>
                <PlayerChart  playerId={playerId} />
            </Container>
        </Container>
    );
}
