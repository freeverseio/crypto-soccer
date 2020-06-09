import React, { useState } from 'react';
import { Button, Input } from 'semantic-ui-react';
const Web3 = require('web3');

function isValidIntArray(x) {
    var res = x.split(",");
    if (!Array.isArray(res)) return false;
    for (var elem of res) {
        if (!parseInt(elem)) return false;
    }
    return true;
}

const IntArrayCard = ({ onChange }) => {
    const [intArray, setIntArray] = useState("");

    const validIntArray = isValidIntArray(intArray);
    
    return (
        <Input fluid
            size='mini'
            error={!validIntArray}
            // icon='ethereum'
            iconPosition='left'
            placeholder='Comma-separated uints...'
            value={intArray}
            onChange={event => setIntArray(event.target.value)}
            // action={
            //     <Button size='mini' color='red' onClick={() => onChange(intArray)} disabled={!validIntArray}>Set</Button>
            // }
        />
    )
}

export default IntArrayCard;