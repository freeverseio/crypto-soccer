import React, { PureComponent } from 'react';
import ApolloClient, { gql } from "apollo-boost";
import { ApolloProvider, Query } from "react-apollo";

import { Segment, Divider } from 'semantic-ui-react'
import TeamCreator from '../components/team_creator';
import TeamList from '../components/team_list';

class Teams extends PureComponent {
    render() {
        const { testingFacade, teams } = this.props;

        return (
            <Segment>
                <TeamCreator contract={testingFacade} />
                <Divider />
                <Query pollInterval={500}
                    query={gql`
           {
  allTeamsList {
    id
    name
  }
}
         `}
                >
                    {({ loading, error, data }) => {
                        if (loading) return <p>Loading...</p>;
                        if (error) return <p>Error :(</p>;

                        return <TeamList teams={data.allTeamsList} />
                    }}
                </Query>
            </Segment>
        )
    }
}

export default Teams;