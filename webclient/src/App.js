import React, { Component } from 'react';
import logo from './logo.svg';
import './App.css';

import EthLeagueManager from './eth_managers/eth_league_manager';

import leagueJSON from './contracts/League.json';
import testingJSON from './contracts/Testing.json';

const provider = 'http://localhost:8545';

class App extends Component {
  constructor(props){
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
        <header className="App-header">
          <img src={logo} className="App-logo" alt="logo" />
          <p>Ganache</p>
          <p>
            { ethLeagueManager ? "connected !" : "not connected !" }
          </p>
          <a
            className="App-link"
            href="https://reactjs.org"
            target="_blank"
            rel="noopener noreferrer"
          >
            Learn React
          </a>
        </header>
      </div>
    );
  }
}

export default App;
