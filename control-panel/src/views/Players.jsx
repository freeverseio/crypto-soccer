import React, { Component } from 'react';
import { Container, Form, Segment } from 'semantic-ui-react';
import gql from 'graphql-tag';
import { useQuery } from '@apollo/react-hooks';

const GET_PLAYERS = gql`
{
    allAuctions {
        nodes {
            uuid
        }
    }
}
`;

class SpecialPlayer extends Component {
    constructor(props) {
        super(props);

        this.handleSubmit = this.handleSubmit.bind(this);

        this.state = {
            shoot: '2000',
            speed: '2000',
            pass: '2000',
            defence: '2000',
            endurance: '2000',
            potential: '5',
            forwardness: '3',
            leftishness: '3',
            aggressiveness: '2',
            age: '17',
            name: 'Johnnie Freeverse',
        }
    }

    GetPlayers() {
        const { loading, error, data } = useQuery(GET_PLAYERS, {
            pollInterval: 500,
        });

        console.log("here")
        if (loading) return null;
        if (error) return `Error! ${error}`;
        console.log(data);

        return (
            <div>
                {data.allAuctions.nodes.map(auction => <div key={auction.uuid}>{auction.uuid}</div>)}
            </div>
        );
    }

    async handleSubmit(event) {
        const { privileged } = this.props;
        const { 
            shoot, 
            speed, 
            pass, 
            defence, 
            endurance,
            potential,
            forwardness,
            leftishness,
            aggressiveness,
            age,
        } = this.state;

        const sk = [shoot, speed, pass, defence, endurance ];
        const traits = [potential, forwardness,  leftishness, aggressiveness];
        const secsInYear = 365 * 24 * 3600
        const internalId = Math.floor(Math.random()*1000000);
        
        const playerId = await privileged.methods.createSpecialPlayer(
            sk,
            age * secsInYear,
            traits,
            internalId
        ).call();

        console.log(playerId)
    }

    render() {
        return (
            <Container style={{ margin: 20 }} >
                <Segment>
                    <Form onSubmit={this.handleSubmit}>
                        <Form.Field>
                            <label>Name</label>
                            <input placeholder='Name' value={this.state.name} onChange={event => this.setState({ name: event.target.value })} />
                        </Form.Field>
                        <Form.Group>
                            <Form.Field>
                                <label>Shoot</label>
                                <input placeholder='Shoot' type='number' value={this.state.shoot} onChange={event => this.setState({ shoot: event.target.value })} />
                            </Form.Field>
                            <Form.Field>
                                <label>Speed</label>
                                <input placeholder='Speed' type='number' value={this.state.speed} onChange={event => this.setState({ speed: event.target.value })} />
                            </Form.Field>
                            <Form.Field>
                                <label>Pass</label>
                                <input placeholder='Pass' type='number' value={this.state.pass} onChange={event => this.setState({ pass: event.target.value })} />
                            </Form.Field>
                            <Form.Field>
                                <label>Defence</label>
                                <input placeholder='Defence' type='number' value={this.state.defence} onChange={event => this.setState({ defence: event.target.value })} />
                            </Form.Field>
                            <Form.Field>
                                <label>Endurance</label>
                                <input placeholder='Endurance' type='number' value={this.state.endurance} onChange={event => this.setState({ endurance: event.target.value })} />
                            </Form.Field>
                        </Form.Group>
                        <Form.Group>
                            <Form.Field>
                                <label>Potential</label>
                                <input placeholder='Potential' type='number' value={this.state.potential} onChange={event => this.setState({ potential: event.target.value })} />
                            </Form.Field>
                            <Form.Field>
                                <label>Forwardness</label>
                                <input placeholder='Forwardness' type='number' value={this.state.forwardness} onChange={event => this.setState({ forwardness: event.target.value })} />
                            </Form.Field>
                            <Form.Field>
                                <label>Leftishness</label>
                                <input placeholder='Leftishness' type='number' value={this.state.leftishness} onChange={event => this.setState({ leftishness: event.target.value })} />
                            </Form.Field>
                            <Form.Field>
                                <label>Aggressiveness</label>
                                <input placeholder='Aggressiveness' type='number' value={this.state.aggressiveness} onChange={event => this.setState({ aggressiveness: event.target.value })} />
                            </Form.Field>
                            <Form.Field>
                                <label>Age</label>
                                <input placeholder='Age' type='number' value={this.state.age} onChange={event => this.setState({ age: event.target.value })} />
                            </Form.Field>
                        </Form.Group>
                        <Form.Button type='submit'>Create</Form.Button>
                    </Form>
                </Segment>
                <this.GetPlayers/>
            </Container>
        );
    }
};

export default SpecialPlayer;