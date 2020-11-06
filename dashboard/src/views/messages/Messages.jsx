import React from 'react';
import { useMutation } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import { Container, Form } from 'semantic-ui-react';
import { useState } from 'react';

const SEND_MESSADE = gql`
  mutation SetBroadcastMessage(
    $category: String!,
    $title: String!,
    $text: String!
  )
  {
    setBroadcastMessage(input: {
    category: $category,
    auctionId:"",
    title: $title,
    text: $text,
    customImageUrl:"",
    metadata:"{}"})
  }
`;

export default () => {
  const [sendMessage] = useMutation(SEND_MESSADE);
  const [category, setCategory] = useState("");
  const [title, setTitle] = useState("");
  const [text, setText] = useState("");

  const handleSubmit = () => {
    sendMessage({
      variables: {
        category: category,
        title: title,
        text: text,
      }
    })
    .catch(console.error);
  }

  return (
    <Container>
      <Form onSubmit={handleSubmit}>
        <Form.Input fluid label="Category" value={category} onChange={e => setCategory(e.value)} />
        <Form.Input fluid label="Title" value={title} onChange={e => setTitle(e.value)} />
        <Form.TextArea label="Body" value={text} onChange={e => setText(e.value)} rows="10" />
        <Form.Button>Submit</Form.Button>
      </Form>
    </Container>
  )
}