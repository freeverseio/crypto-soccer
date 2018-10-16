import React, { Component } from 'react';
import Connection from './components/connection';
import Main from './components/main';
import 'semantic-ui-css/semantic.min.css';

import EthLeagueManager from './eth_managers/eth_league_manager';

import leagueJSON from './contracts/League.json';
import testingJSON from './contracts/Testing.json';


const provider = 'http://localhost:8545';

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      ethLeagueManager: ''
    }
  }

  async componentDidMount() {
    const ethLeagueManager = await EthLeagueManager.createAsync(provider, testingJSON);
    this.setState({ ethLeagueManager });
  }

  render() {
    const { ethLeagueManager } = this.state;

    return (
      <div className="App">
        <Connection provider={provider} ethLeagueManager={ethLeagueManager} />
        <Main />
        <footer>This is the footer</footer>
      </div>
    );
  }
}

export default App;
