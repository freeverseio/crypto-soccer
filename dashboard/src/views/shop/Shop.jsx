import React, { useState } from 'react';
import gql from 'graphql-tag';
import { useQuery, useMutation } from '@apollo/react-hooks';
import { Button, Form, Grid, Header, Card, Segment, Container, Divider, Image } from 'semantic-ui-react';
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
        image: '/speed.png',
    },
    {
        key: 'Shoot Boost',
        text: 'Shoot Boost',
        value: 1,
        image: '/kick.png',
    },
    {
        key: 'Happy Boost',
        text: 'Happy Boost',
        value: 2,
        image: '/pony.png',
    }
];

export default function Shop(props) {
    const [type, setType] = useState(boostOptions[0].value);
    const [name, setName] = useState("");
    const [quantity, setQuantity] = useState(1);
    const [url, setUrl] = useState("");
    const [createShopItem] = useMutation(CREATE_SHOP_ITEM);
    const [shoot, setShoot] = useState(200);
    const [speed, setSpeed] = useState(200);
    const [pass, setPass] = useState(200);
    const [defence, setDefence] = useState(200);
    const [endurance, setEndurance] = useState(200);

    const exponent = (Number(shoot)+Number(speed)+Number(pass)+Number(defence)+Number(endurance))/500;
    const price = Math.floor(Number(quantity) * Math.pow(10,exponent));

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
                            <ShopItemCard key={key} item={item} options={boostOptions[item.type]} />
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
                <Grid.Column style={{ maxWidth: 650 }}>
                    <Header as='h2' color='teal' textAlign='center'>Sponsorship</Header>
                    <Form size='large' onSubmit={createItem}>
                        <Segment stacked>
                            <Grid columns={2}>
                                <Grid.Column>
                                    <Image src={boostOptions[type].image} size='small' centered />
                                    <Form.Dropdown fluid selection options={boostOptions} placeholder='Type' value={type} onChange={(_, { value }) => setType(value)} />
                                </Grid.Column>
                                <Grid.Column>
                                    <Form.Input label='name' fluid placeholder='Name' value={name} onChange={event => setName(event.target.value)} />
                                    <Form.Input type='url' label='url' fluid placeholder='url' value={url} onChange={event => setUrl(event.target.value)} />
                                </Grid.Column>
                            </Grid>
                            <Grid columns={2}>
                                <Grid.Column>
                                    <Form.Input label='quantity' type='number' fluid placeholder='quantity' value={quantity} onChange={event => setQuantity(event.target.value)} />
                                </Grid.Column>
                                <Grid.Column>
                                    <Form.Input icon='euro' iconPosition='left' type='number' label='price' fluid placeholder='price' value={price} disabled />
                                </Grid.Column>
                            </Grid>
                            <Form.Group widths='equal'>
                                <Form.Input fluid label='+Shoot' placeholder='Shoot' type='number' min='0' step='100' max='10000' value={shoot} onChange={event => setShoot(event.target.value)} />
                                <Form.Input fluid label='+Speed' placeholder='Speed' type='number' min='0' step='100' max='10000' value={speed} onChange={event => setSpeed(event.target.value)} />
                                <Form.Input fluid label='+Pass' placeholder='Pass' type='number' min='0' step='100' max='10000' value={pass} onChange={event => setPass(event.target.value)} />
                                <Form.Input fluid label='+Defence' placeholder='Defence' type='number' min='0' step='100' max='10000' value={defence} onChange={event => setDefence(event.target.value)} />
                                <Form.Input fluid label='+Endurance' placeholder='Endurance' type='number' min='0' step='100' max='10000' value={endurance} onChange={event => setEndurance(event.target.value)} />
                            </Form.Group>
                            <Divider />
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