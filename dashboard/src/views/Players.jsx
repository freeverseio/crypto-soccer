import React, { useState } from 'react';
import { Container, Form, Segment, Label, Input, Card, Button, List } from 'semantic-ui-react';
import gql from 'graphql-tag';
import { useMutation, useQuery } from '@apollo/react-hooks';
import signPutAssetForSaleMTx from './marketUtils';
const uuidv1 = require('uuid/v1');

const ALL_PLAYER_IN_ACCADEMY = gql`
query {
    allPlayers(condition: { teamId: "1" }) {
        nodes {
          playerId
          name
          defence
          speed
          pass
          shoot
          endurance
        }
      }
}
`;

const CREATE_PLAYER = gql`
mutation CreateSpecialPlayer(
    $playerId: String!
    $name: String!
    $defence: Int!
    $speed: Int!
    $pass: Int!
    $shoot: Int!
    $endurance: Int!
    $preferredPosition: String!
    $potential: Int!
    $dayOfBirth: Int!
    ) {
        createSpecialPlayer(
            playerId: $playerId
            name: $name
            defence: $defence
            speed: $speed
            pass: $pass
            shoot: $shoot
            endurance: $endurance
            preferredPosition: $preferredPosition
            potential: $potential
            dayOfBirth: $dayOfBirth
        )
    }
`;

const DELETE_PLAYER = gql`
mutation DeleteAcademyPlayer(
    $playerId: String!
    ) {
        deleteSpecialPlayer(
            playerId: $playerId
        )
    }
`;

const CREATE_AUCTION = gql`
mutation CreateAuction(
  $uuid: UUID!
  $playerId: String!
  $currencyId: Int!
  $price: Int!
  $rnd: Int!
  $validUntil: String!
  $signature: String!
  $seller: String!
) {
  createAuction(
    input: {
      uuid: $uuid
      playerId: $playerId
      currencyId: $currencyId
      price: $price
      rnd: $rnd
      validUntil: $validUntil
      signature: $signature
      seller: $seller
    }
  )
}
`;

export default function SpecialPlayer(props) {
    const [shoot, setShoot] = useState(50);
    const [speed, setSpeed] = useState(50);
    const [pass, setPass] = useState(50);
    const [defence, setDefence] = useState(50);
    const [endurance, setEndurance] = useState(50);
    const [potential, setPotential] = useState(5);
    const [forwardness, setForwardness] = useState(3);
    const [leftishness, setLeftishness] = useState(3);
    const [aggressiveness, setAggressiveness] = useState(3);
    const [age, setAge] = useState(19);
    const [name, setName] = useState('Johnnie Freeverse');
    const [price, setPrice] = useState(50);
    const [timeout, setTimeout] = useState(3600);
    const [createAuction] = useMutation(CREATE_AUCTION);
    const [createAcademyPlayer] = useMutation(CREATE_PLAYER);
    const [deleteAcademyPlayer] = useMutation(DELETE_PLAYER);

    async function generatePlayerId() {
        const { privileged } = props;

        const sk = [shoot, speed, pass, defence, endurance];
        const traits = [potential, forwardness, leftishness, aggressiveness];
        const secsInYear = 365 * 24 * 3600
        const internalId = Math.floor(Math.random() * 1000000);

        console.log(sk, traits, secsInYear, internalId);

        const playerId = await privileged.methods.createSpecialPlayer(
            sk,
            age * secsInYear,
            traits,
            internalId
        ).call();

console.log("here")
        return playerId;
    }

    function AccademyPlayers() {
        const { loading, error, data } = useQuery(ALL_PLAYER_IN_ACCADEMY, {
            pollInterval: 2000,
        });

        if (loading) return null;
        if (error) return `Error! ${error}`;

        const players = data.allPlayers.nodes;
        return (
            <Card.Group>
                {
                    players.map((player, key) => {
                        return (
                            <Card key={key}>
                                <Card.Content>
                                    <Card.Header>{player.name}</Card.Header>
                                    <Card.Meta>id: {player.playerId}</Card.Meta>
                                    <Card.Description>
                                        <List>
                                            <List.Item>
                                                <List.Icon name='users' />
                                                <List.Content>{player.shoot}</List.Content>
                                            </List.Item>
                                        </List>
                                    </Card.Description>
                                </Card.Content>
                                <Card.Content extra>
                                    <Form>
                                        <Form.Field>
                                            <Input labelPosition='right' type='number' placeholder='Amount' value={price} onChange={event => setPrice(event.target.value)}>
                                                <Label basic>Price</Label>
                                                <input />
                                                <Label>€</Label>
                                            </Input>
                                        </Form.Field>
                                        <Form.Field>
                                            <Input labelPosition='right' type='number' value={timeout} onChange={event => setTimeout(event.target.value)}>
                                                <Label basic>Timeout</Label>
                                                <input />
                                                <Label>sec</Label>
                                            </Input>
                                        </Form.Field>
                                    </Form>
                                    <div className='ui two buttons'>
                                        <Button basic color='green' onClick={async () => {
                                            const { web3, market } = props;
                                            const rnd = Math.floor(Math.random() * 1000000);
                                            const now = new Date();
                                            const validUntil = (Math.round(now.getTime() / 1000) + timeout).toString();
                                            const sellerAccount = await web3.eth.accounts.create("iamaseller");
                                            console.log("Becoming the owner of the Academy...");
                                            await market.methods.setAcademyAddr(sellerAccount.address);
                                            console.log("Becoming the owner of the Academy...done");
                                            const currencyId = 1;
                                            const signature = await signPutAssetForSaleMTx(web3, currencyId, price, rnd, validUntil, player.playerId, sellerAccount);
                                            const seller = sellerAccount.address;
                                            createAuction({
                                                variables: {
                                                    uuid: uuidv1(),
                                                    playerId: player.playerId,
                                                    currencyId: currencyId,
                                                    price: Number(price),
                                                    rnd: Number(rnd),
                                                    validUntil: validUntil,
                                                    signature: signature.signature,
                                                    seller: seller,
                                                }
                                            });
                                        }}>
                                            Sell
                                        </Button>
                                        <Button basic color='red' onClick={() => {
                                            deleteAcademyPlayer({
                                                variables: {
                                                    playerId: player.playerId,
                                                }
                                            })
                                        }
                                        }>
                                            Kill
                                        </Button>
                                    </div>
                                </Card.Content>
                            </Card>
                        );
                    })
                }
            </Card.Group>
        )
    }

    async function handleSubmit(e) {
        e.preventDefault();

        const playerId = await generatePlayerId();
        const now = new Date();

        console.log("Creating player ", playerId);

        createAcademyPlayer({ // use the block chain to retrive all the values from the playerId
            variables: {
                playerId: playerId,
                name: name,
                defence: Number(defence),
                speed: Number(speed),
                pass: Number(pass),
                shoot: Number(shoot),
                endurance: Number(endurance),
                preferredPosition: "TODO",
                potential: Number(potential),
                dayOfBirth: 16950,
            }
        });
        console.log("dewewewew")
    }

    return (
        <Container style={{ margin: 20 }} >
            <Segment>
                <Form onSubmit={handleSubmit}>
                    <Form.Field>
                        <Input labelPosition='left' type='text' value={name} onChange={event => setName(event.target.value)}>
                            <Label basic>Name</Label>
                            <input />
                        </Input>
                    </Form.Field>
                    <Form.Group>
                        <Form.Field>
                            <label>Shoot</label>
                            <input placeholder='Shoot' type='number' value={shoot} onChange={event => setShoot(event.target.value)} />
                        </Form.Field>
                        <Form.Field>
                            <label>Speed</label>
                            <input placeholder='Speed' type='number' value={speed} onChange={event => setSpeed(event.target.value)} />
                        </Form.Field>
                        <Form.Field>
                            <label>Pass</label>
                            <input placeholder='Pass' type='number' value={pass} onChange={event => setPass(event.target.value)} />
                        </Form.Field>
                        <Form.Field>
                            <label>Defence</label>
                            <input placeholder='Defence' type='number' value={defence} onChange={event => setDefence(event.target.value)} />
                        </Form.Field>
                        <Form.Field>
                            <label>Endurance</label>
                            <input placeholder='Endurance' type='number' value={endurance} onChange={event => setEndurance(event.target.value)} />
                        </Form.Field>
                    </Form.Group>
                    <Form.Group>
                        <Form.Field>
                            <label>Potential</label>
                            <input placeholder='Potential' type='number' value={potential} onChange={event => setPotential(event.target.value)} />
                        </Form.Field>
                        <Form.Field>
                            <label>Forwardness</label>
                            <input placeholder='Forwardness' type='number' value={forwardness} onChange={event => setForwardness(event.target.value)} />
                        </Form.Field>
                        <Form.Field>
                            <label>Leftishness</label>
                            <input placeholder='Leftishness' type='number' value={leftishness} onChange={event => setLeftishness(event.target.value)} />
                        </Form.Field>
                        <Form.Field>
                            <label>Aggressiveness</label>
                            <input placeholder='Aggressiveness' type='number' value={aggressiveness} onChange={event => setAggressiveness(event.target.value)} />
                        </Form.Field>
                        <Form.Field>
                            <label>Age</label>
                            <input placeholder='Age' type='number' value={age} onChange={event => setAge(event.target.value)} />
                        </Form.Field>
                    </Form.Group>
                    <Form.Button type='submit'>Create</Form.Button>
                </Form>
            </Segment>
            <AccademyPlayers />
        </Container>
    );
};
