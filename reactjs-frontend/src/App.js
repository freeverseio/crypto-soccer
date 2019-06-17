import React from 'react';
import ApolloClient, { gql } from "apollo-boost";
import { ApolloProvider, Query } from "react-apollo";
// import logo from './logo.svg';
import './App.css';

const client = new ApolloClient({
  uri: "http://192.168.10.231:4000/"
});

const App = () => (
  <ApolloProvider client={client}>
    <div>
      <h2>My first Apollo app</h2>
      <Query
        query={gql`
          {
            allTeams {
              id
              name
            }
          }
        `}
      >
        {({ loading, error, data }) => {
          if (loading) return <p>Loading...</p>;
          if (error) return <p>Error :(</p>;

          return data.allTeams.map(({ id, name }) => (
            <div key={id}>
              <p>{id}: {name}</p>
            </div>
          ));
        }}
      </Query>
    </div>
  </ApolloProvider>
);

// function App() {
//   return (
//     <div className="App">
//       <header className="App-header">
//         <img src={logo} className="App-logo" alt="logo" />
//         <p>
//           Edit <code>src/App.js</code> and save to reload.
//         </p>
//         <a
//           className="App-link"
//           href="https://reactjs.org"
//           target="_blank"
//           rel="noopener noreferrer"
//         >
//           Learn React
//         </a>
//       </header>
//     </div>
//   );
// }

export default App;
