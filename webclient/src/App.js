import React, { Component } from 'react';
import Connection from './components/connection';
import Main from './components/main';
import Web3 from 'web3';
import 'semantic-ui-css/semantic.min.css';

import { createTestingContract, EthLeagueManager } from './eth_managers/eth_league_manager';

const provider = 'http://localhost:8545';

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      ethLeagueManager: '',
      teams: []
    }
  }

  async getTeams(ethLeagueManager) {
    const count = await ethLeagueManager.countTeams();

    let teams = [];
    for (let i = 0; i < count; i++) {
      const name = await ethLeagueManager.teamName(i);
      teams.push({index: i, name: name});
    }

    return teams;
  }

  async componentDidMount() {
    const web3 = new Web3(provider);
    const accounts = await web3.eth.getAccounts();
    const testingContract = await createTestingContract(web3);
    const ethLeagueManager = new EthLeagueManager(testingContract, accounts[0]);
    const teams = await this.getTeams(ethLeagueManager);

    this.setState({ ethLeagueManager, teams });
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
