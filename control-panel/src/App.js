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

const web3 = new Web3('http://localhost:8545');
const privileged = new web3.eth.Contract(privilegedJSON.abi, "0x72a2F9bfCD665Efadc58A05bCaf7Be380a8dE03B");

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
          <Main web3={web3} privileged={privileged}/>
        </div>
      </ApolloProvider>
    </Router>
  );
}

export default App;
