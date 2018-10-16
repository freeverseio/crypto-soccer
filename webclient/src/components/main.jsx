import React, { Component } from 'react';
import { Form, Container } from 'semantic-ui-react'
import Teams from './teams';

class Main extends Component {
    render() {
        return (
            <Container>
                <Form onSubmit={this.handleSubmit}>
                    <Form.Group>
                        <Form.Input placeholder='Name' name='name' value="ciao" onChange={this.handleChange} />
                        <Form.Button content='Create Team' />
                    </Form.Group>
                </Form>
                <Teams />
            </Container>
        )
    }
}

export default Main;