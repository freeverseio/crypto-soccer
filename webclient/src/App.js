import React, { Component } from 'react';
import Connection from './components/connection';
import Main from './components/main';
import Web3 from 'web3';
import 'semantic-ui-css/semantic.min.css';

import { createTestingContract, TestingFacade } from './ethereum/testing_contract';

const provider = 'http://localhost:8545';

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      testingFacade: '',
      teams: []
    }
  }

  async getTeams(testingFacade) {
    const count = await testingFacade.countTeams();

    let teams = [];
    for (let i = 0; i < count; i++) {
      const name = await testingFacade.teamName(i);
      teams.push({index: i, name: name});
    }

    return teams;
  }

  async componentDidMount() {
    const web3 = new Web3(provider);
    const accounts = await web3.eth.getAccounts();
    const testingContract = await createTestingContract(web3);
    const testingFacade = new TestingFacade(testingContract, accounts[0]);
    const teams = await this.getTeams(testingFacade);

    this.setState({ testingFacade, teams });
  }

  render() {
    return (
      <div className="App">
        <Connection provider={provider} {...this.state} />
        <Main {...this.state} />
      </div>
    );
  }
}

export default App;
