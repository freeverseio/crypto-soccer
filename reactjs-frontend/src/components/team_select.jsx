import React from 'react';
import { Select } from 'semantic-ui-react';

// [{ key: 'af', value: 'af', flag: 'af', text: 'Afghanistan' }, ...{}]

const TeamSelect = (params) => {
    const teams = params.teams;
    const options = teams.map(team => ({ key: team.index, value: team.index, text: team.name }))
    return (
        <Select  {...params} options={options} />
    );
}

export default TeamSelect;