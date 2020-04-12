{
  "cluster": {
    "secret": "09968797372bbfe9c4b7bb70b738e61db814fa44cb33b4faea520496607c63ed",
    "listen_multiaddress": [
      "/ip4/0.0.0.0/tcp/9096",
      "/ip4/0.0.0.0/udp/9096/quic"
    ],
    "connection_manager": {
      "high_water": 400,
      "low_water": 100,
      "grace_period": "2m"
    },
    "state_sync_interval": "2h",
    "ipfs_sync_interval": "2m",
    "monitor_ping_interval": "15s",
    "peer_watch_interval": "5s",
    "mdns_interval": "0",
    "disable_repinning": true,
    "follower_mode": true,
    "peer_addresses": [
      "/ip4/178.128.141.147/tcp/9096/p2p/12D3KooW9qyokQDee1H7Z3ym1RWXdxzuZpQzmJrp78vnnEKW49F9",
      "/ip4/178.128.141.147/tcp/9096/p2p/12D3KooWBQ2jpp812nNp1PfzWXaoqnWVm3faMsT4xgrP3EUmPKN4"
    ]
  },
  "consensus": {
    "crdt": {
      "cluster_name": "ipfs-cluster",
      "trusted_peers": [
          "12D3KooW9qyokQDee1H7Z3ym1RWXdxzuZpQzmJrp78vnnEKW49F9",
	  "12D3KooWBQ2jpp812nNp1PfzWXaoqnWVm3faMsT4xgrP3EUmPKN4"
      ],
      "rebroadcast_interval": "10s",
      "peerset_metric": "ping"
    }
  },
  "api": {},
  "ipfs_connector": {
    "ipfshttp": {
      "node_multiaddress": "/ip4/127.0.0.1/tcp/5001",
      "connect_swarms_delay": "30s",
      "ipfs_request_timeout": "5m0s",
      "pin_timeout": "10m",
      "unpin_timeout": "30m",
      "unpin_disable": false
    }
  },
  "pin_tracker": {
    "maptracker": {
      "max_pin_queue_size": 20000,
      "concurrent_pins": 15
    },
    "stateless": {
      "max_pin_queue_size": 1000000,
      "concurrent_pins": 10
    }
  },
  "monitor": {
    "pubsubmon": {
      "check_interval": "15s",
      "failure_threshold": 3
    }
  },
  "informer": {
    "disk": {
      "metric_ttl": "1m",
      "metric_type": "freespace"
    },
    "numpin": {
      "metric_ttl": "5m"
    }
  },
  "observations": {},
  "datastore": {
    "badger": {
      "badger_options": {
        "dir": "",
        "value_dir": "",
        "sync_writes": true,
        "table_loading_mode": 0,
        "value_log_loading_mode": 0,
        "num_versions_to_keep": 1,
        "max_table_size": 67108864,
        "level_size_multiplier": 10,
        "max_levels": 7,
        "value_threshold": 32,
        "num_memtables": 5,
        "num_level_zero_tables": 5,
        "num_level_zero_tables_stall": 10,
        "level_one_size": 268435456,
        "value_log_file_size": 1073741823,
        "value_log_max_entries": 1000000,
        "num_compactors": 2,
        "compact_l_0_on_close": true,
        "read_only": false,
        "truncate": false
      }
    }
  }
}
