import React, { useState, useEffect } from 'react';
import { Grid, Button, Input, GridRow } from 'semantic-ui-react';
const Web3 = require('web3');

const ProposedCompanyWidget = ({ disabled, onChange }) => {
    const [address, setAddress] = useState("");

    const validAddress = Web3.utils.isAddress(address);

    return (
        <Input fluid
            size='mini'
            error={!validAddress}
            icon='ethereum'
            iconPosition='left'
            placeholder='Address ...'
            value={address}
            onChange={event => setAddress(event.target.value)}
            action={
                <Button size='mini' color='red' onClick={() => onChange(address)} disabled={disabled || !validAddress}>Set</Button>
            }
        />
    )
}

export default ProposedCompanyWidget;