import React from 'react';
// import logo from './logo.svg';
import './App.css';
import 'semantic-ui-css/semantic.min.css'
import ApolloClient from 'apollo-boost';
import { gql } from "apollo-boost";

const client = new ApolloClient({
  uri: 'http://165.22.66.118:4000/graphiql',
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
    <div className="App">
      <header className="App-header">
        <div class="ui divided three column grid">
          <div class="row">
            <div class="ui right labeled input">
              <input type="text" placeholder="Enter weight..." />
              <div class="ui basic label label">defence</div>
            </div>
          </div>
        </div>
        <button class="ui button" onClick={() => createPlayer()}>Create Special Player</button>
      </header>
    </div>
  );
    }
    
    export default App;
