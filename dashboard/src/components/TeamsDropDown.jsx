import React from 'react';
import { Dropdown } from 'semantic-ui-react';
import gql from 'graphql-tag';
import { useQuery } from '@apollo/react-hooks';

const GET_TEAMS = gql`
query {
  allTeams {
      nodes {
        teamId
        name
      }
    }
}
`;

export default function TeamsDropDown(props) {
    const { onTeamIdChange } = props;

    const { loading, error, data } = useQuery(GET_TEAMS); if (loading) return 'Loading...';
    if (error) return `Error! ${error.message}`;

    const teams = data.allTeams.nodes;
    const teamOptions = teams.map(team => ({
        key: team.teamId,
        text: team.name,
        value: team.teamId,
    }))

    return (
        <Dropdown
            placeholder='Select Team'
            fluid
            selection
            search
            options={teamOptions}
            onChange={(_, props) => {
                onTeamIdChange && onTeamIdChange(props.value);
            }}
        />
    );
}