import React, {useState} from 'react';
import { useQuery } from '@apollo/react-hooks';
import { Button, Form, Grid, Header, Image, Message, Segment } from 'semantic-ui-react'

const boostOptions = [
    {
        key: 'Speed Boost',
        text: 'Speed Boost',
        value: 'Speed Boost',
        icon: 'fast forward',
    },
    {
        key: 'Shoot Boost',
        text: 'Shoot Boost',
        value: 'Shoot Boost',
        icon: 'fire',
    },
    {
        key: 'Happy Boost',
        text: 'Happy Boost',
        value: 'Happy Boost',
        icon: 'thumbs up',
    }
]

export default function Shop(props) {
    const [type, setType] = useState(boostOptions[0].value);
    const [name, setName] = useState("");
    const [price, setPrice] = useState(50);

    const createItem = () => {
        console.log("Submitted");
    }

    return (
        <Grid textAlign='center' verticalAlign='middle'>
            <Grid.Column style={{ maxWidth: 450 }}>
                <Header as='h2' color='teal' textAlign='center'>Shop</Header>
                <Form size='large' onSubmit={createItem}>
                    <Segment stacked>
                        <Form.Dropdown fluid selection options={boostOptions} placeholder='Type' value={type} onChange={(_, {value}) => setType(value)}/>
                        <Form.Input fluid type='number' icon='dollar' iconPosition='left' placeholder='Price' value={price} onChange={event => setPrice(event.target.value)}/>
                        <Form.Input fluid icon='user' iconPosition='left' placeholder='Name' value={name} onChange={event => setName(event.target.value)}/>
                        <Button type='submit' color='teal' fluid size='large'>Create</Button>
                    </Segment>
                </Form>
            </Grid.Column>
        </Grid>
    )
}