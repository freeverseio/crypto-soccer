import React, { useState } from 'react';
import gql from 'graphql-tag';
import { useQuery } from '@apollo/react-hooks';
import { Button, Form, Grid, Header, Card, Segment, Container, Divider } from 'semantic-ui-react';
import ShopItem from '../../components/ShopItem';

const ALL_SHOPS_ITEMS = gql`
query {
      allShopItems {
        nodes {
        uuid
        type
        price
        }
    }
}
`;

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
];

export default function Shop(props) {
    const [type, setType] = useState(boostOptions[0].value);
    const [name, setName] = useState("");
    const [price, setPrice] = useState(50);

    const createItem = () => {
        console.log("Submitted");
    }

    const Shop = () => {
        const { loading, error, data } = useQuery(ALL_SHOPS_ITEMS, {
            pollInterval: 2000,
        });

        if (loading) return null;
        if (error) return `Error! ${error}`;

        const items = data.allShopItems.nodes;
        return (
            <Card.Group>
                {
                    items.map((item, key) => {
                        return (
                            <ShopItem key={key} item={item} />
                        );
                    })

                }
            </Card.Group>
        )
    }

    return (
        <Container>
            <Grid textAlign='center' verticalAlign='middle'>
                <Grid.Column style={{ maxWidth: 450 }}>
                    <Header as='h2' color='teal' textAlign='center'>Shop</Header>
                    <Form size='large' onSubmit={createItem}>
                        <Segment stacked>
                            <Form.Dropdown fluid selection options={boostOptions} placeholder='Type' value={type} onChange={(_, { value }) => setType(value)} />
                            <Form.Input fluid type='number' icon='dollar' iconPosition='left' placeholder='Price' value={price} onChange={event => setPrice(event.target.value)} />
                            <Form.Input fluid icon='user' iconPosition='left' placeholder='Name' value={name} onChange={event => setName(event.target.value)} />
                            <Button type='submit' color='teal' fluid size='large'>Create</Button>
                        </Segment>
                    </Form>
                </Grid.Column>
            </Grid>
            <Divider />
            <Shop />
        </Container>
    )
}