import React, { useState } from 'react';
import { Container, Form, Segment, Label, Input } from 'semantic-ui-react';
import gql from 'graphql-tag';
import { useMutation } from '@apollo/react-hooks';
import { concatHash, getMessageHash, signPutAssetForSaleMTx, signAgreeToBuyPlayerMTx, signAgreeToBuyTeamMTx } from './marketUtils'
const uuidv1 = require('uuid/v1');

const ALL_PLAYER_IN_ACCADEMY = gql`
query {
    allPlayers(condition: { teamId: "1" }) {
        nodes {
          playerId
          name
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
    const [shoot, setShoot] = useState(2000);
    const [speed, setSpeed] = useState(2000);
    const [pass, setPass] = useState(2000);
    const [defence, setDefence] = useState(2000);
    const [endurance, setEndurance] = useState(2000);
    const [potential, setPotential] = useState(2000);
    const [forwardness, setForwardness] = useState(2000);
    const [leftishness, setLeftishness] = useState(2000);
    const [aggressiveness, setAggressiveness] = useState(2000);
    const [age, setAge] = useState(2000);
    const [name, setName] = useState('Johnnie Freeverse');
    const [price, setPrice] = useState(50);
    const [timeout, setTimeout] = useState(3600);
    const [createAuction] = useMutation(CREATE_AUCTION);

    async function generatePlayerId() {
        const { privileged } = props;

        const sk = [shoot, speed, pass, defence, endurance];
        const traits = [potential, forwardness, leftishness, aggressiveness];
        const secsInYear = 365 * 24 * 3600
        const internalId = Math.floor(Math.random() * 1000000);

        const playerId = await privileged.methods.createSpecialPlayer(
            sk,
            age * secsInYear,
            traits,
            internalId
        ).call();

        return playerId;
    }

    async function handleSubmit(e) {
        const { web3, market } = props;
        e.preventDefault();

        const playerId = await generatePlayerId();
        const rnd = Math.floor(Math.random() * 1000000);
        const now = new Date();
        const validUntil = (Math.round(now.getTime() / 1000) + timeout).toString();
        const sellerAccount = await web3.eth.accounts.create("iamaseller");
        console.log("Becoming the owner of the Academy...");
        await market.methods.setAcademyAddr(sellerAccount.address);
        console.log("Becoming the owner of the Academy...done");
      
        const currencyId = 1;
        const signature = await signPutAssetForSaleMTx(web3, currencyId, price, rnd, validUntil, playerId, sellerAccount);
        const seller = sellerAccount.address;
        createAuction({variables: {
            uuid: uuidv1(),
            playerId: playerId,
            currencyId: currencyId,
            price: Number(price),
            rnd: Number(rnd),
            validUntil: validUntil,
            signature: signature.signature,
            seller: seller,
        }});
        console.log("Correctly sent new Academy player creation!");
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
                    <Form.Button type='submit'>Create</Form.Button>
                </Form>
            </Segment>
        </Container>
    );
};
