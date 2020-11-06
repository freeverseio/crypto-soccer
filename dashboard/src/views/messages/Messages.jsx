import React from 'react';
import { useMutation } from '@apollo/react-hooks';
import gql from 'graphql-tag';
import { Container, Form } from 'semantic-ui-react';
import { useState } from 'react';

const SEND_MESSAGE = gql`
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
  const [sendMessage, { error }] = useMutation(SEND_MESSAGE);
  const [category, setCategory] = useState("");
  const [title, setTitle] = useState("");
  const [text, setText] = useState("");

  function handleSubmit() {
    sendMessage({
      variables: {
        category: category,
        title: title,
        text: text,
      }
    }).catch(console.error);
  }

  if (error) return `Error! ${error.message}`;

  return (
    <Container>
      <Form onSubmit={handleSubmit}>
        <Form.Input fluid label="Category" value={category} onChange={e => setCategory(e.target.value)} error={category === ""} />
        <Form.Input fluid label="Title" value={title} onChange={e => setTitle(e.target.value)} error={title === ""} />
        <Form.TextArea label="Body" value={text} onChange={e => setText(e.target.value)} rows="10" />
        <Form.Button active={category !== "" && title !== ""}>Submit</Form.Button>
      </Form>
    </Container>
  )
}