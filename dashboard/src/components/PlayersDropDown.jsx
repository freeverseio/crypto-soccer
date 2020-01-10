import React from 'react';
import { Dropdown } from 'semantic-ui-react';
import gql from 'graphql-tag';
import { useQuery } from '@apollo/react-hooks';

const GET_PLAYERS = gql`
query {
  allPlayers {
      nodes {
        playerId
        name
      }
    }
}
`;

export default function PlayersDropDown(props) {
    const { onPlayerIdChange } = props;

    const { loading, error, data } = useQuery(GET_PLAYERS); if (loading) return 'Loading...';
    if (error) return `Error! ${error.message}`;

    const players = data.allPlayers.nodes;
    const playerOptions = players.map(player => ({
        key: player.playerId,
        text: player.name,
        value: player.playerId,
    }))

    return (
        <Dropdown
            placeholder='Select Player'
            fluid
            selection
            search
            options={playerOptions}
            onChange={(_, props) => {
                onPlayerIdChange && onPlayerIdChange(props.value);
            }}
        />
    );
}