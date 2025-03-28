import React, { useState } from 'react';
import gql from 'graphql-tag';
import { useQuery, useMutation } from '@apollo/react-hooks';
import { Button, Form, Grid, Header, Segment, Container, Divider, Image } from 'semantic-ui-react';
import ShopItemCard from '../../components/ShopItemCard';
import uuidv1 from 'uuid/v1';
import Config from '../../Config';

const ALL_SHOPS_ITEMS = gql`
query {
      allShopItems {
        nodes {
        uuid
        name
        url
        type
        price
        quantity
        }
    }
}`;

const CREATE_SHOP_ITEM = gql`
mutation CreateShopItem(
    $uuid: UUID!
    $name: String!
    $url: String!
    $type: Int!
    $price: Int!
    $quantity: Int!
    ){
  createShopItem(
    input: { 
        uuid: $uuid
        name: $name
        url: $url
        type: $type
        price: $price
        quantity: $quantity
    }
  )
}`;

const boostOptions = [
    {
        key: 'Nike GT1000T',
        text: 'Nike GT1000T',
        value: 0,
        image: '/nike_shoes.png',
    },
    {
        key: 'RedBull',
        text: 'RedBull',
        value: 1,
        image: '/redbull.png',
    },
    {
        key: 'Adidas Shirt',
        text: 'Adidas Shirt',
        value: 2,
        image: '/adidas_shirt.png',
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
            pollInterval: Config.polling_ms,
        });

        if (loading) return null;
        if (error) return `Error! ${error}`;

        const items = data.allShopItems.nodes;
        return (
            <Grid columns={5}>
                {
                    items.map((item, key) => {
                        return (
                            <Grid.Column key={key}>
                                <ShopItemCard item={item} options={boostOptions[item.type]} />
                            </Grid.Column>
                        );
                    })

                }
            </Grid>
        )
    }

    function createItem(e) {
        e.preventDefault();
        createShopItem({
            variables: {
                uuid: uuidv1(),
                name: name,
                url: url,
                quantity: Number(quantity),
                type: Number(type),
                price: Number(price),
            }
        })
        .catch(console.error);
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
                                    <Divider />
                                    <Form.Dropdown fluid selection options={boostOptions} placeholder='Type' value={type} onChange={(_, { value }) => setType(value)} />
                                </Grid.Column>
                                <Grid.Column>
                                    <Form.Input required label='name' fluid placeholder='Name' value={name} onChange={event => setName(event.target.value)} />
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