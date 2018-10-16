import React, { Component } from 'react';
import { Button, Card, Image } from 'semantic-ui-react'

const TeamCard = () => (
    <Card>
        <Card.Content>
            <Image floated='right' size='mini' src='https://react.semantic-ui.com/images/avatar/large/steve.jpg' />
            <Card.Header>Steve Sanders</Card.Header>
            <Card.Meta>Friends of Elliot</Card.Meta>
            <Card.Description>
                Steve wants to add you to the group <strong>best friends</strong>
            </Card.Description>
        </Card.Content>
        <Card.Content extra>
            <Button basic color='red'>
                Delete
                </Button>
        </Card.Content>
    </Card>
);

export default TeamCard;