import React, { Component } from 'react';
import Header from './views/header';
import Main from './views/main';
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
    for (let team = 0; team < count; team++) {
      const name = await testingFacade.teamName(team);
      let players = [];
      for (let player = 0; player < 11; player++) {
        const skills = await testingFacade.playerSkills(team, player);
        const name = await testingFacade.playerName(team, player);
        players.push({ index: player, name, skills });
      }
      teams.push({ 
        index: team, 
        name: name,
        players: players
      });
    }

    return teams;
  }

  async componentDidMount() {
    const accounts = await this.web3.eth.getAccounts();
    const testingContract = await createTestingContract(this.web3);
    if (!testingContract)
      return;

    const testingFacade = new TestingFacade(testingContract, accounts[0]);
    const teams = await this.getTeams(testingFacade);

    testingContract.events.TeamCreation()
      .on('data', event => {
        this.getTeams(testingFacade)
          .then(teams => this.setState({ teams }));
      })
      .on('changed', reason => console.log("(WW): " + reason))
      .on('error', reason => console.log("(EE): " + reason));

    this.setState({ testingFacade, teams });
  }

  render() {
    const url = this.web3.currentProvider.connection.url;

    return (
      <div className="App">
        <Header url={url} />
        <Main {...this.state} />
      </div>
    );
  }
}

export default App;
