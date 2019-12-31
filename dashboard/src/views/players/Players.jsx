import React, { useState } from 'react';
import { Container } from 'semantic-ui-react'
import PlayersDropDown from '../../components/PlayersDropDown';

export default function Teams(params) {
    const [playerId, setPlayerId] = useState("");

    return (
        <Container>
            <PlayersDropDown onPlayerIdChange={setPlayerId} />
        </Container>
    )
}