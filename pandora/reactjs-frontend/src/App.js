import React from 'react';
import { HashRouter } from 'react-router-dom';
import ApolloClient from "apollo-boost";
import { ApolloProvider } from "react-apollo";
import Header from './views/header';
import Main from './views/main';
// import logo from './logo.svg';
import 'semantic-ui-css/semantic.min.css';
import './App.css';

const url = "http://localhost:5000/graphql"

const client = new ApolloClient({
  uri: url
});

const App = () => (
  <HashRouter>
    <ApolloProvider client={client}>
      <div className="App">
        <Header url={url} />
        <Main />
      </div>
    </ApolloProvider>
  </HashRouter>
    // {/* <div>
    //   <h2>My first Apollo app</h2>
    //   <Query
    //     query={gql`
    //       {
    //         allTeams {
    //           id
    //           name
    //         }
    //       }
    //     `}
    //   >
    //     {({ loading, error, data }) => {
    //       if (loading) return <p>Loading...</p>;
    //       if (error) return <p>Error :(</p>;

    //       return data.allTeams.map(({ id, name }) => (
    //         <div key={id}>
    //           <p>{id}: {name}</p>
    //         </div>
    //       ));
    //     }}
    //   </Query>
    //   </div> */}
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
