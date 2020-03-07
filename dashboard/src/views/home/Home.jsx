import React, {useState} from 'react';
import { Container } from 'semantic-ui-react';
import TeamsByRankingTable from '../../components/TeamsByRankingTable';
import TeamTable from '../../components/TeamTable';
import PlayerChart from '../../components/PlayerChart';

const data = [
  {
    "id": "japan",
    "color": "hsl(303, 70%, 50%)",
    "data": [
      {
        "x": "plane",
        "y": 261
      },
      {
        "x": "helicopter",
        "y": 88
      },
      {
        "x": "boat",
        "y": 272
      },
      {
        "x": "train",
        "y": 47
      },
      {
        "x": "subway",
        "y": 203
      },
      {
        "x": "bus",
        "y": 67
      },
      {
        "x": "car",
        "y": 192
      },
      {
        "x": "moto",
        "y": 160
      },
      {
        "x": "bicycle",
        "y": 104
      },
      {
        "x": "horse",
        "y": 230
      },
      {
        "x": "skateboard",
        "y": 154
      },
      {
        "x": "others",
        "y": 33
      }
    ]
  },
  {
    "id": "france",
    "color": "hsl(347, 70%, 50%)",
    "data": [
      {
        "x": "plane",
        "y": 205
      },
      {
        "x": "helicopter",
        "y": 38
      },
      {
        "x": "boat",
        "y": 64
      },
      {
        "x": "train",
        "y": 228
      },
      {
        "x": "subway",
        "y": 84
      },
      {
        "x": "bus",
        "y": 170
      },
      {
        "x": "car",
        "y": 288
      },
      {
        "x": "moto",
        "y": 234
      },
      {
        "x": "bicycle",
        "y": 142
      },
      {
        "x": "horse",
        "y": 189
      },
      {
        "x": "skateboard",
        "y": 190
      },
      {
        "x": "others",
        "y": 254
      }
    ]
  },
  {
    "id": "us",
    "color": "hsl(277, 70%, 50%)",
    "data": [
      {
        "x": "plane",
        "y": 152
      },
      {
        "x": "helicopter",
        "y": 66
      },
      {
        "x": "boat",
        "y": 263
      },
      {
        "x": "train",
        "y": 294
      },
      {
        "x": "subway",
        "y": 180
      },
      {
        "x": "bus",
        "y": 160
      },
      {
        "x": "car",
        "y": 61
      },
      {
        "x": "moto",
        "y": 267
      },
      {
        "x": "bicycle",
        "y": 204
      },
      {
        "x": "horse",
        "y": 62
      },
      {
        "x": "skateboard",
        "y": 53
      },
      {
        "x": "others",
        "y": 249
      }
    ]
  },
  {
    "id": "germany",
    "color": "hsl(203, 70%, 50%)",
    "data": [
      {
        "x": "plane",
        "y": 222
      },
      {
        "x": "helicopter",
        "y": 207
      },
      {
        "x": "boat",
        "y": 297
      },
      {
        "x": "train",
        "y": 83
      },
      {
        "x": "subway",
        "y": 241
      },
      {
        "x": "bus",
        "y": 184
      },
      {
        "x": "car",
        "y": 61
      },
      {
        "x": "moto",
        "y": 147
      },
      {
        "x": "bicycle",
        "y": 24
      },
      {
        "x": "horse",
        "y": 12
      },
      {
        "x": "skateboard",
        "y": 284
      },
      {
        "x": "others",
        "y": 166
      }
    ]
  },
  {
    "id": "norway",
    "color": "hsl(83, 70%, 50%)",
    "data": [
      {
        "x": "plane",
        "y": 267
      },
      {
        "x": "helicopter",
        "y": 120
      },
      {
        "x": "boat",
        "y": 80
      },
      {
        "x": "train",
        "y": 8
      },
      {
        "x": "subway",
        "y": 90
      },
      {
        "x": "bus",
        "y": 151
      },
      {
        "x": "car",
        "y": 77
      },
      {
        "x": "moto",
        "y": 81
      },
      {
        "x": "bicycle",
        "y": 280
      },
      {
        "x": "horse",
        "y": 28
      },
      {
        "x": "skateboard",
        "y": 24
      },
      {
        "x": "others",
        "y": 39
      }
    ]
  }
];
export default function Home(props) {
    const [teamId, setTeamId] = useState("");
    return (
        <Container>
            <TeamsByRankingTable onTeamIdChange={setTeamId} />
            <TeamTable teamId={teamId} />
            <Container style={{"height": "600px"}}>
                <PlayerChart data={data} />
            </Container>
        </Container>
    );
}
