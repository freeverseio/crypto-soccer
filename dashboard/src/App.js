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
const web3 = new Web3('http://167.172.161.183:8545');
const privileged = new web3.eth.Contract(privilegedJSON.abi, "0xd02E63F7130ABb255a4ccC7a74105D0E8dE46017");
const market = new web3.eth.Contract(marketJSON.abi, "0xC7Ef4b2CB85cc2e7764d4d9510EBe7F925163081");

const url = 'http://167.172.161.183:4000/graphiql';
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
