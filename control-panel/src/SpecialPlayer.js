import React, { Component } from 'react';
import { Container, Form, Segment } from 'semantic-ui-react';

class SpecialPlayer extends Component {
    handleSubmit(event) {
        console.log(event);
    }

    render() {
        return (
            <Container style={{ margin: 20 }} >
                <Segment>
                    <Form onSubmit={this.handleSubmit}>
                        <Form.Field>
                            <label>Defence</label>
                            <input placeholder='Defence' onChange={this.setState} />
                        </Form.Field>
                        <Form.Field>
                            <label>Speed</label>
                            <input placeholder='Speed' />
                        </Form.Field>
                        <Form.Button type='submit'>Create</Form.Button>
                    </Form>
                </Segment>
            </Container>
        );
    }
};

export default SpecialPlayer;