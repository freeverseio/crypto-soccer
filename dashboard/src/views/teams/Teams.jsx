import React, { useState } from 'react';
import { Container } from 'semantic-ui-react'
import TeamsDropDown from '../../components/TeamsDropDown';
import TeamTable from '../../components/TeamTable';

export default function Teams(params) {
    const [teamId, setTeamId] = useState("");

    return (
        <Container>
            <TeamsDropDown onTeamIdChange={setTeamId} />
            <TeamTable teamId={teamId} />
        </Container>
    )
}