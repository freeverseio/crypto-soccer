import React from 'react';
import { Select } from 'semantic-ui-react';

// [{ key: 'af', value: 'af', flag: 'af', text: 'Afghanistan' }, ...{}]

const TeamSelect = (params) => {
    console.log(params);
    return (
        <Select placeholder='Select your country' options={countryOptions} />
    );
}


export default TeamSelect;