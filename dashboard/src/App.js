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
const web3 = new Web3('https://prod.gorengine.com/pippolo');
const privileged = new web3.eth.Contract(privilegedJSON.abi, "0x615668099Cc46D035b3c34aCdf01204Ac4A4F446");
const market = new web3.eth.Contract(marketJSON.abi, "0xFB1436D488726D64a0441081D508b238fF756802");

function App() {
  const url = process.env.REACT_APP_API;

  const client = new ApolloClient({
    request: (operation) => {
      operation.setContext({
        headers: {
          authorization: 'Bearer joshua'
        }
      })
    },
    uri: url,
  });

  return (
    <Router>
      <ApolloProvider client={client}>
        <div className="App">
          <Header url={url} />
          <Main web3={web3} privileged={privileged} market={market} />
        </div>
      </ApolloProvider>
    </Router>
  );
}

export default App;
