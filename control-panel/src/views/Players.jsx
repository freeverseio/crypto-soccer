import React, { Component } from 'react';
import { Container, Form, Segment } from 'semantic-ui-react';

class SpecialPlayer extends Component {
    constructor(props) {
        super(props);

        this.handleSubmit = this.handleSubmit.bind(this);

        this.state = {
            defence: '',
            speed: ''
        }
    }
    handleSubmit(event) {
        console.log(this.state);
    }

    render() {
        return (
            <Container style={{ margin: 20 }} >
                <Segment>
                    <Form onSubmit={this.handleSubmit}>
                        <Form.Field>
                            <label>Defence</label>
                            <input placeholder='Defence' type='number' value={this.state.defence} onChange={event => this.setState({defence: event.target.value})}/>
                        </Form.Field>
                        <Form.Field>
                            <label>Speed</label>
                            <input placeholder='Speed' type='number' value={this.state.speed} onChange={event => this.setState({speed: event.target.value})}/>
                        </Form.Field>
                        <Form.Button type='submit'>Create</Form.Button>
                    </Form>
                </Segment>
            </Container>
        );
    }
};

export default SpecialPlayer;