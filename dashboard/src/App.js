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

const web3 = new Web3(Web3.givenProvider || 'http://165.22.66.118:8545');
const privileged = new web3.eth.Contract(privilegedJSON.abi, "0x7058c864D752bC8c12C001611fc417b92B370A2C");
const market = new web3.eth.Contract(marketJSON.abi, "0x7Ff2f2A42191226FF7C1Aa7f5ca4Fb70dF8DE469");

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
