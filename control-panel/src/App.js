import React from 'react';
// import logo from './logo.svg';
// import './App.css';
import 'semantic-ui-css/semantic.min.css'

import { HashRouter as Router } from 'react-router-dom';
import ApolloClient from 'apollo-boost';
import { ApolloProvider } from '@apollo/react-hooks';
import { gql } from "apollo-boost";
import Header from './views/Header';
import Main from './views/Main';

const url = 'http://165.22.66.118:4000/graphiql';

const client = new ApolloClient({
  uri: url,
});

const createPlayer = () => {
  client
    .query({
      query: gql`
      {
        allAuctions {
          nodes {
            uuid
            state
          }
        }
      }
    `
    })
    .then(result => console.log(result));
}

function App() {
  return (
    <Router>
      <ApolloProvider client={client}>
        <div className="App">
          <Header url={url} />
          <Main />
        </div>
      </ApolloProvider>
    </Router>
  );
}

export default App;
