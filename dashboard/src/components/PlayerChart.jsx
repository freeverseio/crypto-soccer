import React from 'react';
import { ResponsiveLine } from '@nivo/line'
import gql from 'graphql-tag';
import { useQuery } from '@apollo/react-hooks';

// make sure parent container have a defined height when using
// responsive component, otherwise height will be 0 and
// no chart will be rendered.
// website examples showcase many properties,
// you'll often use just a few of them.

const GET_TEAMS_BY_RANKING = gql`
query allPlayersHistories($playerId: String!){
    allPlayersHistories(condition: { playerId: $playerId }) {
        nodes {
        blockNumber
        defence
        speed
        endurance
        pass
        shoot
        }
    }
}
`;



const MyResponsiveLine = ({ playerId }) => {
    const { loading, error, data } = useQuery(GET_TEAMS_BY_RANKING, {
        variables: {playerId},
    });

    if (playerId === ''){
        return <div/>
    }

    if (loading) return 'Loading...';
    if (error) return `Error! ${error.message}`

    const history = data.allPlayersHistories.nodes ? data.allPlayersHistories.nodes : [];

    const chartData = [
        {
            "id": "defence",
            "data": [],
        },
        {
            "id": "speed",
            "data": [],
        },
        {
            "id": "endurance",
            "data": [],
        },
        {
            "id": "pass",
            "data": [],
        },
        {
            "id": "shoot",
            "data": [],
        },
    ];

    history.forEach(state => {
        chartData[0].data.push({ "x": state.blockNumber, "y": state.defence });
        chartData[1].data.push({ "x": state.blockNumber, "y": state.speed });
        chartData[2].data.push({ "x": state.blockNumber, "y": state.endurance });
        chartData[3].data.push({ "x": state.blockNumber, "y": state.pass });
        chartData[4].data.push({ "x": state.blockNumber, "y": state.shoot });
    });

    return (
        <ResponsiveLine
            data={chartData}
            margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
            xScale={{ type: 'point' }}
            yScale={{ type: 'linear', min: 'auto', max: 'auto', stacked: true, reverse: false }}
            axisTop={null}
            axisRight={null}
            axisBottom={{
                orient: 'bottom',
                tickSize: 5,
                tickPadding: 5,
                tickRotation: 0,
                legend: 'blocknumber',
                legendOffset: 36,
                legendPosition: 'middle'
            }}
            axisLeft={{
                orient: 'left',
                tickSize: 5,
                tickPadding: 5,
                tickRotation: 0,
                // legend: 'count',
                legendOffset: -40,
                legendPosition: 'middle'
            }}
            colors={{ scheme: 'nivo' }}
            pointSize={10}
            pointColor={{ theme: 'background' }}
            pointBorderWidth={2}
            pointBorderColor={{ from: 'serieColor' }}
            pointLabel="y"
            pointLabelYOffset={-12}
            useMesh={true}
            legends={[
                {
                    anchor: 'bottom-right',
                    direction: 'column',
                    justify: false,
                    translateX: 100,
                    translateY: 0,
                    itemsSpacing: 0,
                    itemDirection: 'left-to-right',
                    itemWidth: 80,
                    itemHeight: 20,
                    itemOpacity: 0.75,
                    symbolSize: 12,
                    symbolShape: 'circle',
                    symbolBorderColor: 'rgba(0, 0, 0, .5)',
                    effects: [
                        {
                            on: 'hover',
                            style: {
                                itemBackground: 'rgba(0, 0, 0, .03)',
                                itemOpacity: 1
                            }
                        }
                    ]
                }
            ]}
        />
    );
}

export default MyResponsiveLine;