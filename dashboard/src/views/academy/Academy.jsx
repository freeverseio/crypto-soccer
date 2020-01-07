import React, { useState } from 'react';
import { Container, Form, Segment, Card, Grid, Header, Button, Divider } from 'semantic-ui-react';
import gql from 'graphql-tag';
import { useMutation, useQuery } from '@apollo/react-hooks';
import PlayerCard from '../../components/PlayerCard';

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
      potential
      dayOfBirth
      auctionsByPlayerId(orderBy: VALID_UNTIL_DESC, first: 1) {
        nodes {
          validUntil
          bidsByAuction {
            totalCount
          }
        }
      }
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

    const [createAcademyPlayer] = useMutation(CREATE_PLAYER);

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
        const { web3 } = props;
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
                            <PlayerCard key={key} player={player} web3={web3} />
                        );
                    })

                }
            </Card.Group>
        )
    }

    async function handleSubmit(e) {
        e.preventDefault();

        const playerId = await generatePlayerId();

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
    }

    return (
        <Container style={{ margin: 20 }} >
            <Grid textAlign='center' verticalAlign='middle'>
                <Grid.Column style={{ maxWidth: 650 }}>
                    <Header as='h2' color='teal' textAlign='center'>Academy</Header>
                    <Form size='large' onSubmit={handleSubmit}>
                        <Segment stacked>
                        <Form.Input fluid labelPosition='left' type='text' value={name} onChange={event => setName(event.target.value)}/>
                            <Form.Group  widths='equal'>
                                <Form.Input fluid label='Shoot' placeholder='Shoot' type='number' value={shoot} onChange={event => setShoot(event.target.value)} />
                                <Form.Input fluid label='Speed' placeholder='Speed' type='number' value={speed} onChange={event => setSpeed(event.target.value)} />
                                <Form.Input fluid label='Pass' placeholder='Pass' type='number' value={pass} onChange={event => setPass(event.target.value)} />
                                <Form.Input fluid label='Defence' placeholder='Defence' type='number' value={defence} onChange={event => setDefence(event.target.value)} />
                                <Form.Input fluid label='Endurance' placeholder='Endurance' type='number' value={endurance} onChange={event => setEndurance(event.target.value)} />
                                <Form.Input fluid label='Speed' placeholder='Speed' type='number' value={speed} onChange={event => setSpeed(event.target.value)} />
                            </Form.Group>
                            <Form.Group  widths='equal'>
                                <Form.Input fluid label='Potential' placeholder='Potential' type='number' value={potential} onChange={event => setPotential(event.target.value)} />
                                <Form.Input fluid label='Forwardness' placeholder='Forwardness' type='number' value={forwardness} onChange={event => setForwardness(event.target.value)} />
                                <Form.Input fluid label='Leftishness' placeholder='Leftishness' type='number' value={leftishness} onChange={event => setLeftishness(event.target.value)} />
                                <Form.Input fluid label='Aggressiveness' placeholder='Aggressiveness' type='number' value={aggressiveness} onChange={event => setAggressiveness(event.target.value)} />
                                <Form.Input fluid label='Age' placeholder='Age' type='number' value={age} onChange={event => setAge(event.target.value)} />
                            </Form.Group>
                            <Button type='submit' color='teal' fluid size='large'>Create</Button>
                        </Segment>
                    </Form>
                </Grid.Column>
            </Grid>
            <Divider />
            <AccademyPlayers />
        </Container>
    );
};
