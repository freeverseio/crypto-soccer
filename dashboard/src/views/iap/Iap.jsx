import React, {useState} from 'react';
import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';

const GET_PROXY_ADDRESS = gql`
    {
        paramByName(name: "PROXY"){
            value
        }
    }
`;

export default () => {
    return (
        <div/>
    )
}