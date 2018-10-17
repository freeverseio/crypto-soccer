import React, { Component } from 'react';
import Connection from './components/connection';
import Main from './components/main';
import Web3 from 'web3';
import 'semantic-ui-css/semantic.min.css';

import EthLeagueManager from './eth_managers/eth_league_manager';

const provider = 'http://localhost:8545';

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      ethLeagueManager: ''
    }
  }

  async componentDidMount() {
    const web3 = new Web3(provider);
    const accounts = await web3.eth.getAccounts();
    const ethLeagueManager = await EthLeagueManager.createAsync(web3, accounts[0]);
    this.setState({ethLeagueManager});
  }

  render() {
    const { ethLeagueManager } = this.state;

    return (
      <div className="App">
        <Connection provider={provider} ethLeagueManager={ethLeagueManager} />
        <Main ethLeagueManager={ethLeagueManager}/>
      </div>
    );
  }
}

export default App;
