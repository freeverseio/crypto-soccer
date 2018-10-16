import React, { Component } from 'react';
import { Form, Container, Segment, Divider } from 'semantic-ui-react'
import Teams from './teams';

class Main extends Component {
    constructor(props){
        super(props);
        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);

        this.state = {
            name: ''
        }
    }

    handleSubmit(event) {
        const { ethLeagueManager } = this.props;
        ethLeagueManager.createTeam(this.state.name);
        this.setState({name: ''});
        event.preventDefault();
    }

    handleChange(event){
        this.setState({name: event.target.value});
    }
    
    render() {
        const { ethLeagueManager } = this.props;

        return (
            <Container>
                <Segment>
                    <Form onSubmit={this.handleSubmit}>
                        <Form.Group widths='equal'>
                            <Form.Input placeholder='Name' name='name' value={this.state.name} onChange={this.handleChange} />
                            <Form.Button content='Create Team' />
                        </Form.Group>
                    </Form>
                    <Divider/>
                    <Teams ethLeagueManager={ethLeagueManager}/>
                </Segment>
            </Container>
        )
    }
}

export default Main;