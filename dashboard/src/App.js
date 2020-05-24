import React, { useState } from 'react';
// import logo from './logo.svg';
// import './App.css';
import 'semantic-ui-css/semantic.min.css'

import { HashRouter as Router } from 'react-router-dom';
import ApolloClient from 'apollo-boost';
import { InMemoryCache } from 'apollo-cache-inmemory';
import { ApolloProvider } from '@apollo/react-hooks';
import Web3 from 'web3';
import Header from './views/Header';
import Main from './views/Main';
import Config from './Config';
const privilegedJSON = require("./contracts/Privileged.json");
const marketJSON = require("./contracts/Market.json");

const cache = new InMemoryCache({
  dataIdFromObject: object => object.nodeId || null,
});

const client = new ApolloClient({
  uri: Config.url,
  cache,
});

function App() {
  const [account, setAccount] = useState(window.ethereum ? window.ethereum.selectedAddress : undefined);

  let web3 = null;
  let privileged = null;
  let market = null;

  if (window.ethereum) {
    window.ethereum.on('accountsChanged', function (accounts) {
      setAccount(accounts[0]);
    });
    web3 = new Web3(window.ethereum);

    privileged = new web3.eth.Contract(privilegedJSON.abi, "0x615668099Cc46D035b3c34aCdf01204Ac4A4F446");
    market = new web3.eth.Contract(marketJSON.abi, "0xFB1436D488726D64a0441081D508b238fF756802");
  }

  return (
    <Router web3={web3}>
      <ApolloProvider client={client}>
        <div className="App">
          <Header account={account} />
          <Main web3={web3}  account={account} privileged={privileged} market={market} />
        </div>
      </ApolloProvider>
    </Router>
  );
}

export default App;
