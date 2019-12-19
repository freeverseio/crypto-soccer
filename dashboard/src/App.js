import React from 'react';
// import logo from './logo.svg';
// import './App.css';
import 'semantic-ui-css/semantic.min.css'

import { HashRouter as Router } from 'react-router-dom';
import ApolloClient from 'apollo-boost';
import { ApolloProvider } from '@apollo/react-hooks';
import Web3 from 'web3';
import Header from './views/Header';
import Main from './views/Main';
const privilegedJSON = require("./contracts/Privileged.json");
const marketJSON = require("./contracts/Market.json");

// const web3 = new Web3(Web3.givenProvider || 'http://165.22.66.118:8545');
const web3 = new Web3('http://165.22.66.118:8545');
const privileged = new web3.eth.Contract(privilegedJSON.abi, "0x169dCd95cCD6384D58D51549a33904EDD28A9e7D");
const market = new web3.eth.Contract(marketJSON.abi, "0xF51a67e33A5534DaA8bbBaEb4334d849Dc574C05");

const url = 'http://165.22.66.118:4000/graphiql';
const client = new ApolloClient({
  uri: url,
});

function App() {
  return (
    <Router>
      <ApolloProvider client={client}>
        <div className="App">
          <Header url={url} />
          <Main web3={web3} privileged={privileged} market={market}/>
        </div>
      </ApolloProvider>
    </Router>
  );
}

export default App;
