import React, { Component } from 'react';
import Connection from './components/connection';
import Main from './components/main';
import Web3 from 'web3';
import 'semantic-ui-css/semantic.min.css';

import { createTestingContract, TestingFacade } from './ethereum/testing_contract';

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      testingFacade: '',
      teams: []
    }

    this.web3Provider = new Web3.providers.WebsocketProvider('ws://localhost:8545');
    this.web3 = new Web3(this.web3Provider);
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
    const accounts = await this.web3.eth.getAccounts();
    const testingContract = await createTestingContract(this.web3);
    const testingFacade = new TestingFacade(testingContract, accounts[0]);
    const teams = await this.getTeams(testingFacade);

    testingContract.events.TeamCreation()
      .on('data', event => {
          this.getTeams(testingFacade)
          .then(teams => this.setState({teams}));
      })
      .on('changed', reason => console.log("(WW): " + reason))
      .on('error', reason => console.log("(EE): " + reason));

    this.setState({ testingFacade, teams });
  }

  render() {
    return (
      <div className="App">
        <Connection provider={this.web3Provider} {...this.state} />
        <Main {...this.state} />
      </div>
    );
  }
}

export default App;
