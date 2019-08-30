import React, { Component } from 'react';
import { Form } from 'semantic-ui-react'

class TeamCreator extends Component {
    constructor(props){
        super(props);
        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);

        this.state = {
            name: ''
        }
    }

    handleChange(event) {
        this.setState({ name: event.target.value });
    }

    handleSubmit(event) {
        const { contract } = this.props;
        contract.createTeam(this.state.name);
        this.setState({name: ''});
        event.preventDefault();
    }

    render() {
        return (
            <Form onSubmit={this.handleSubmit}>
                <Form.Group widths='equal'>
                    <Form.Input placeholder='Name' name='name' value={this.state.name} onChange={this.handleChange} />
                    <Form.Button content='Create Team' />
                </Form.Group>
            </Form>
        )
    }
}

export default TeamCreator;
