import React, { useState } from 'react';
import gql from 'graphql-tag';
import { useQuery, useMutation } from '@apollo/react-hooks';
import { Button, Form, Grid, Header, Card, Segment, Container, Divider } from 'semantic-ui-react';
import ShopItemCard from '../../components/ShopItemCard';
import uuidv1 from 'uuid/v1';

const ALL_SHOPS_ITEMS = gql`
query {
      allShopItems {
        nodes {
        uuid
        type
        price
        }
    }
}`;

const CREATE_SHOP_ITEM = gql`
mutation CreateShopItem(
    $uuid: UUID!
    $type: Int!
    $price: Int!
    ){
  createShopItem(
    input: { 
        uuid: $uuid
        type: $type
        price: $price
    }
  )
}`;

const boostOptions = [
    {
        key: 'Speed Boost',
        text: 'Speed Boost',
        value: 0,
        icon: 'fast forward',
    },
    {
        key: 'Shoot Boost',
        text: 'Shoot Boost',
        value: 1,
        icon: 'fire',
    },
    {
        key: 'Happy Boost',
        text: 'Happy Boost',
        value: 2,
        icon: 'thumbs up',
    }
];

export default function Shop(props) {
    const [type, setType] = useState(boostOptions[0].value);
    const [name, setName] = useState("");
    const [price, setPrice] = useState(50);
    const [createShopItem] = useMutation(CREATE_SHOP_ITEM);

    const Shop = () => {
        const { loading, error, data } = useQuery(ALL_SHOPS_ITEMS, {
            pollInterval: 2000,
        });

        if (loading) return null;
        if (error) return `Error! ${error}`;

        const items = data.allShopItems.nodes;
        return (
            <Card.Group itemsPerRow={5}>
                {
                    items.map((item, key) => {
                        return (
                            <ShopItemCard key={key} item={item} />
                        );
                    })

                }
            </Card.Group>
        )
    }

    function createItem(e) {
        e.preventDefault();
        console.log(type)
        createShopItem({
            variables: {
                uuid: uuidv1(),
                type: Number(type),
                price: Number(price),
            }
        });
    };

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