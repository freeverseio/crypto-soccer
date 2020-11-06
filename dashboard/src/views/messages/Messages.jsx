import React from 'react';
import { useQuery } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import { Container, Form } from 'semantic-ui-react';
import Config from '../../Config';

const GET_AUCTIONS = gql`
{
  allAuctions {
    nodes {
      id
      playerId
      price
      validUntil
      seller
      state
      stateExtra
      playerByPlayerId {
        name
      }
    }
  }
}
`;

export default () => {
  return (
    <Container>
      <Form>
        <Form.Input fluid label="Type" />
        <Form.Input fluid label="Title" />
        <Form.TextArea label="Body" rows="10"/>
        <Form.Button>Submit</Form.Button>
      </Form>
    </Container>
  )
}