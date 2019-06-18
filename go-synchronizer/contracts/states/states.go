// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package states

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// StatesABI is the input ABI used to generate the binding from.
const StatesABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getSkills\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"state\",\"type\":\"uint256\"}],\"name\":\"isValidPlayerState\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getEndurance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"state\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPrevLeagueId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getSpeed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getDefence\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"defence\",\"type\":\"uint256\"},{\"name\":\"speed\",\"type\":\"uint256\"},{\"name\":\"pass\",\"type\":\"uint256\"},{\"name\":\"shoot\",\"type\":\"uint256\"},{\"name\":\"endurance\",\"type\":\"uint256\"},{\"name\":\"monthOfBirthInUnixTime\",\"type\":\"uint256\"},{\"name\":\"playerId\",\"type\":\"uint256\"},{\"name\":\"currentTeamId\",\"type\":\"uint256\"},{\"name\":\"currentShirtNum\",\"type\":\"uint256\"},{\"name\":\"prevLeagueId\",\"type\":\"uint256\"},{\"name\":\"prevTeamPosInLeague\",\"type\":\"uint256\"},{\"name\":\"prevShirtNumInLeague\",\"type\":\"uint256\"},{\"name\":\"lastSaleBlock\",\"type\":\"uint256\"}],\"name\":\"playerStateCreate\",\"outputs\":[{\"name\":\"state\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getPass\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getPrevTeamPosInLeague\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getShoot\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getPrevShirtNumInLeague\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"state\",\"type\":\"uint256\"},{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"setPrevTeamPosInLeague\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getMonthOfBirthInUnixTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"},{\"name\":\"delta\",\"type\":\"uint16\"}],\"name\":\"playerStateEvolve\",\"outputs\":[{\"name\":\"evolvedState\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"state\",\"type\":\"uint256\"},{\"name\":\"currentShirtNum\",\"type\":\"uint256\"}],\"name\":\"setCurrentShirtNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"state\",\"type\":\"uint256\"},{\"name\":\"lastSaleBlock\",\"type\":\"uint256\"}],\"name\":\"setLastSaleBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"},{\"name\":\"teamId\",\"type\":\"uint256\"}],\"name\":\"setCurrentTeamId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getLastSaleBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getSkillsVec\",\"outputs\":[{\"name\":\"skills\",\"type\":\"uint16[5]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getCurrentTeamId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getCurrentShirtNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getPlayerId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"getPrevLeagueId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"teamStateCreate\",\"outputs\":[{\"name\":\"state\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamState\",\"type\":\"uint256[]\"},{\"name\":\"playerState\",\"type\":\"uint256\"}],\"name\":\"teamStateAppend\",\"outputs\":[{\"name\":\"state\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamState\",\"type\":\"uint256[]\"}],\"name\":\"teamStateSize\",\"outputs\":[{\"name\":\"count\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamState\",\"type\":\"uint256[]\"},{\"name\":\"idx\",\"type\":\"uint256\"}],\"name\":\"teamStateAt\",\"outputs\":[{\"name\":\"playerState\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"state\",\"type\":\"uint256[]\"}],\"name\":\"isValidTeamState\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamState\",\"type\":\"uint256[]\"},{\"name\":\"delta\",\"type\":\"uint8\"}],\"name\":\"teamStateEvolve\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"teamState\",\"type\":\"uint256[]\"}],\"name\":\"computeTeamRating\",\"outputs\":[{\"name\":\"rating\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// StatesBin is the compiled bytecode used for deploying new contracts.
const StatesBin = `"0x608060405234801561001057600080fd5b506123e5806100206000396000f3fe608060405234801561001057600080fd5b50600436106101d95760003560e01c80637ee0aebc11610104578063c37b1c25116100a2578063e06da0bf11610071578063e06da0bf14610ca2578063eb78b7b714610db9578063f438591214610dfb578063f8bd3e6e14610e3d576101d9565b8063c37b1c2514610b68578063c566b5bc14610bb4578063cc1cc3d714610bf6578063cd2105e814610c60576101d9565b8063a95e858b116100de578063a95e858b14610934578063af26723a14610980578063af76cd0114610a50578063bd64d4fa14610a9c576101d9565b80637ee0aebc1461085657806385053566146108a25780638d216b52146108e4576101d9565b806351585b491161017c5780635c0107821161014b5780635c0107821461065957806365b4b47614610773578063666d0070146107b557806374833f73146107f7576101d9565b806351585b49146104d8578063530f63d61461051a57806355a6f86f146105d557806358a7a46a14610617576101d9565b806344328d1a116101b857806344328d1a146102a85780634444a2831461037e57806347f3d7161461044a5780634b93f75314610496576101d9565b806292bf78146101de57806319a4860c14610220578063258e5d9014610266575b600080fd5b61020a600480360360208110156101f457600080fd5b8101908080359060200190929190505050610e7f565b6040518082815260200191505060405180910390f35b61024c6004803603602081101561023657600080fd5b8101908080359060200190929190505050610f08565b604051808215151515815260200191505060405180910390f35b6102926004803603602081101561027c57600080fd5b8101908080359060200190929190505050610f1d565b6040518082815260200191505060405180910390f35b610368600480360360408110156102be57600080fd5b81019080803590602001906401000000008111156102db57600080fd5b8201836020820111156102ed57600080fd5b8035906020019184602083028401116401000000008311171561030f57600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f82011690508083019250505050505050919291929080359060200190929190505050610faa565b6040518082815260200191505060405180910390f35b6104346004803603602081101561039457600080fd5b81019080803590602001906401000000008111156103b157600080fd5b8201836020820111156103c357600080fd5b803590602001918460208302840111640100000000831117156103e557600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f8201169050808301925050505050505091929192905050506110b9565b6040518082815260200191505060405180910390f35b6104806004803603604081101561046057600080fd5b81019080803590602001909291908035906020019092919050505061113f565b6040518082815260200191505060405180910390f35b6104c2600480360360208110156104ac57600080fd5b81019080803590602001909291905050506111e0565b6040518082815260200191505060405180910390f35b610504600480360360208110156104ee57600080fd5b810190808035906020019092919050505061126d565b6040518082815260200191505060405180910390f35b6105bf60048036036101a081101561053157600080fd5b81019080803590602001909291908035906020019092919080359060200190929190803590602001909291908035906020019092919080359060200190929190803590602001909291908035906020019092919080359060200190929190803590602001909291908035906020019092919080359060200190929190803590602001909291905050506112f6565b6040518082815260200191505060405180910390f35b610601600480360360208110156105eb57600080fd5b8101908080359060200190929190505050611718565b6040518082815260200191505060405180910390f35b6106436004803603602081101561062d57600080fd5b81019080803590602001909291905050506117a5565b6040518082815260200191505060405180910390f35b61071c6004803603604081101561066f57600080fd5b810190808035906020019064010000000081111561068c57600080fd5b82018360208201111561069e57600080fd5b803590602001918460208302840111640100000000831117156106c057600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050509192919290803560ff169060200190929190505050611831565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561075f578082015181840152602081019050610744565b505050509050019250505060405180910390f35b61079f6004803603602081101561078957600080fd5b810190808035906020019092919050505061190c565b6040518082815260200191505060405180910390f35b6107e1600480360360208110156107cb57600080fd5b8101908080359060200190929190505050611999565b6040518082815260200191505060405180910390f35b6107ff611a25565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b83811015610842578082015181840152602081019050610827565b505050509050019250505060405180910390f35b61088c6004803603604081101561086c57600080fd5b810190808035906020019092919080359060200190929190505050611a2a565b6040518082815260200191505060405180910390f35b6108ce600480360360208110156108b857600080fd5b8101908080359060200190929190505050611ac6565b6040518082815260200191505060405180910390f35b61091e600480360360408110156108fa57600080fd5b8101908080359060200190929190803561ffff169060200190929190505050611b53565b6040518082815260200191505060405180910390f35b61096a6004803603604081101561094a57600080fd5b810190808035906020019092919080359060200190929190505050611bd9565b6040518082815260200191505060405180910390f35b610a366004803603602081101561099657600080fd5b81019080803590602001906401000000008111156109b357600080fd5b8201836020820111156109c557600080fd5b803590602001918460208302840111640100000000831117156109e757600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050509192919290505050611c78565b604051808215151515815260200191505060405180910390f35b610a8660048036036040811015610a6657600080fd5b810190808035906020019092919080359060200190929190505050611ccb565b6040518082815260200191505060405180910390f35b610b5260048036036020811015610ab257600080fd5b8101908080359060200190640100000000811115610acf57600080fd5b820183602082011115610ae157600080fd5b80359060200191846020830284011164010000000083111715610b0357600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050509192919290505050611d65565b6040518082815260200191505060405180910390f35b610b9e60048036036040811015610b7e57600080fd5b810190808035906020019092919080359060200190929190505050611e5e565b6040518082815260200191505060405180910390f35b610be060048036036020811015610bca57600080fd5b8101908080359060200190929190505050611f03565b6040518082815260200191505060405180910390f35b610c2260048036036020811015610c0c57600080fd5b8101908080359060200190929190505050611f93565b6040518082600560200280838360005b83811015610c4d578082015181840152602081019050610c32565b5050505090500191505060405180910390f35b610c8c60048036036020811015610c7657600080fd5b81019080803590602001909291905050506120e8565b6040518082815260200191505060405180910390f35b610d6260048036036040811015610cb857600080fd5b8101908080359060200190640100000000811115610cd557600080fd5b820183602082011115610ce757600080fd5b80359060200191846020830284011164010000000083111715610d0957600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f82011690508083019250505050505050919291929080359060200190929190505050612177565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b83811015610da5578082015181840152602081019050610d8a565b505050509050019250505060405180910390f35b610de560048036036020811015610dcf57600080fd5b810190808035906020019092919050505061221b565b6040518082815260200191505060405180910390f35b610e2760048036036020811015610e1157600080fd5b81019080803590602001909291905050506122a7565b6040518082815260200191505060405180910390f35b610e6960048036036020811015610e5357600080fd5b81019080803590602001909291905050506122bb565b6040518082815260200191505060405180910390f35b6000610e8a82610f08565b610efc576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f696e76616c696420706c6179657220737461746500000000000000000000000081525060200191505060405180910390fd5b60ba82901c9050919050565b600080610f14836122a7565b14159050919050565b6000610f2882610f08565b610f9a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f696e76616c696420706c6179657220737461746500000000000000000000000081525060200191505060405180910390fd5b613fff60ba83901c169050919050565b600082518210611022576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600c8152602001807f6f7574206f6620626f756e64000000000000000000000000000000000000000081525060200191505060405180910390fd5b61102b83611c78565b61109d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f696e76616c6964207465616d207374617465000000000000000000000000000081525060200191505060405180910390fd5b8282815181106110a957fe5b6020026020010151905092915050565b60006110c482611c78565b611136576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f696e76616c6964207465616d207374617465000000000000000000000000000081525060200191505060405180910390fd5b81519050919050565b6000630200000082106111ba576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f707265764c6561677565496478206f7574206f6620626f756e6400000000000081525060200191505060405180910390fd5b6dffffff80000000000000000000001983169250605782901b8317925082905092915050565b60006111eb82610f08565b61125d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f696e76616c696420706c6179657220737461746500000000000000000000000081525060200191505060405180910390fd5b613fff60e483901c169050919050565b600061127882610f08565b6112ea576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f696e76616c696420706c6179657220737461746500000000000000000000000081525060200191505060405180910390fd5b60f282901c9050919050565b60006140008e1061136f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f646566656e6365206f7574206f6620626f756e6400000000000000000000000081525060200191505060405180910390fd5b6140008d106113e6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f646566656e6365206f7574206f6620626f756e6400000000000000000000000081525060200191505060405180910390fd5b6140008c1061145d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f646566656e6365206f7574206f6620626f756e6400000000000000000000000081525060200191505060405180910390fd5b6140008b106114d4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f646566656e6365206f7574206f6620626f756e6400000000000000000000000081525060200191505060405180910390fd5b6140008a1061154b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f646566656e6365206f7574206f6620626f756e6400000000000000000000000081525060200191505060405180910390fd5b61400089106115a5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602381526020018061236d6023913960400191505060405180910390fd5b6000881180156115b85750631000000088105b61162a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f706c617965724964206f7574206f6620626f756e64000000000000000000000081525060200191505060405180910390fd5b60108310611683576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806123906021913960400191505060405180910390fd5b60f28e901b8117905060e48d901b8117905060d68c901b8117905060c88b901b8117905060ba8a901b8117905060ac89901b81179050609088901b811790506116cc8188611e5e565b90506116d88187611bd9565b90506116e4818661113f565b90506116f08185611a2a565b9050604b83901b811790506117058183611ccb565b90509d9c50505050505050505050505050565b600061172382610f08565b611795576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f696e76616c696420706c6179657220737461746500000000000000000000000081525060200191505060405180910390fd5b613fff60d683901c169050919050565b60006117b082610f08565b611822576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f696e76616c696420706c6179657220737461746500000000000000000000000081525060200191505060405180910390fd5b60ff604f83901c169050919050565b606061183c83611c78565b6118ae576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f696e76616c6964207465616d207374617465000000000000000000000000000081525060200191505060405180910390fd5b60008090505b8351811015611902576118dd8482815181106118cc57fe5b60200260200101518460ff16611b53565b8482815181106118e957fe5b60200260200101818152505080806001019150506118b4565b5082905092915050565b600061191782610f08565b611989576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f696e76616c696420706c6179657220737461746500000000000000000000000081525060200191505060405180910390fd5b613fff60c883901c169050919050565b60006119a482610f08565b611a16576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f696e76616c696420706c6179657220737461746500000000000000000000000081525060200191505060405180910390fd5b600f604b83901c169050919050565b606090565b60006101008210611aa3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260208152602001807f707265765465616d506f73496e4c6561677565206f7574206f6620626f756e6481525060200191505060405180910390fd5b6a7f800000000000000000001983169250604f82901b8317925082905092915050565b6000611ad182610f08565b611b43576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f696e76616c696420706c6179657220737461746500000000000000000000000081525060200191505060405180910390fd5b613fff60ac83901c169050919050565b6000611b5e83610f08565b611bd0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f696e76616c696420706c6179657220706c61796572537461746500000000000081525060200191505060405180910390fd5b82905092915050565b600060108210611c51576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601c8152602001807f63757272656e7453686972744e756d206f7574206f6620626f756e640000000081525060200191505060405180910390fd5b6e0f00000000000000000000000000001983169250607082901b8317925082905092915050565b600080600090505b8251811015611cc057611ca5838281518110611c9857fe5b6020026020010151610f08565b611cb3576000915050611cc6565b8080600101915050611c80565b50600190505b919050565b60006408000000008210611d47576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f6c61737453616c65426c6f636b206f7574206f6620626f756e6400000000000081525060200191505060405180910390fd5b654500000000001983169250602882901b8317925082905092915050565b6000611d7082611c78565b611de2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260128152602001807f696e76616c6964207465616d207374617465000000000000000000000000000081525060200191505060405180910390fd5b60008090505b8251811015611e58576000838281518110611dff57fe5b60200260200101519050611e128161126d565b83019250611e1f816111e0565b83019250611e2c81611718565b83019250611e398161190c565b83019250611e4681610f1d565b83019250508080600101915050611de8565b50919050565b600063100000008210611ed9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f63757272656e745465616d496478206f7574206f6620626f756e64000000000081525060200191505060405180910390fd5b71fffffff000000000000000000000000000001983169250607482901b8317925082905092915050565b6000611f0e82610f08565b611f80576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f696e76616c696420706c6179657220737461746500000000000000000000000081525060200191505060405180910390fd5b6407ffffffff602883901c169050919050565b611f9b61234a565b611fa482610f08565b612016576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f696e76616c696420706c6179657220737461746500000000000000000000000081525060200191505060405180910390fd5b61201f8261126d565b8160006005811061202c57fe5b602002019061ffff16908161ffff1681525050612048826111e0565b8160016005811061205557fe5b602002019061ffff16908161ffff168152505061207182611718565b8160026005811061207e57fe5b602002019061ffff16908161ffff168152505061209a8261190c565b816003600581106120a757fe5b602002019061ffff16908161ffff16815250506120c382610f1d565b816004600581106120d057fe5b602002019061ffff16908161ffff1681525050919050565b60006120f382610f08565b612165576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f696e76616c696420706c6179657220737461746500000000000000000000000081525060200191505060405180910390fd5b630fffffff607483901c169050919050565b606060018351016040519080825280602002602001820160405280156121ac5781602001602082028038833980820191505090505b50905060008090505b83518110156121f7578381815181106121ca57fe5b60200260200101518282815181106121de57fe5b60200260200101818152505080806001019150506121b5565b50818160018351038151811061220957fe5b60200260200101818152505092915050565b600061222682610f08565b612298576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f696e76616c696420706c6179657220737461746500000000000000000000000081525060200191505060405180910390fd5b600f607083901c169050919050565b6000630fffffff609083901c169050919050565b60006122c682610f08565b612338576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260148152602001807f696e76616c696420706c6179657220737461746500000000000000000000000081525060200191505060405180910390fd5b6301ffffff605783901c169050919050565b6040518060a0016040528060059060208202803883398082019150509050509056fe6d6f6e74684f664269727468496e556e697854696d65206f7574206f6620626f756e647072657653686972744e756d496e4c6561677565206f7574206f6620626f756e64a265627a7a7230582027f621ed9d78c9e46faf8e0b9dc623a2d4f7545383b252e70e1fe474a0d9a84d64736f6c63430005090032"`

// DeployStates deploys a new Ethereum contract, binding an instance of States to it.
func DeployStates(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *States, error) {
	parsed, err := abi.JSON(strings.NewReader(StatesABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StatesBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &States{StatesCaller: StatesCaller{contract: contract}, StatesTransactor: StatesTransactor{contract: contract}, StatesFilterer: StatesFilterer{contract: contract}}, nil
}

// States is an auto generated Go binding around an Ethereum contract.
type States struct {
	StatesCaller     // Read-only binding to the contract
	StatesTransactor // Write-only binding to the contract
	StatesFilterer   // Log filterer for contract events
}

// StatesCaller is an auto generated read-only Go binding around an Ethereum contract.
type StatesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StatesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StatesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StatesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StatesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StatesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StatesSession struct {
	Contract     *States           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StatesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StatesCallerSession struct {
	Contract *StatesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// StatesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StatesTransactorSession struct {
	Contract     *StatesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StatesRaw is an auto generated low-level Go binding around an Ethereum contract.
type StatesRaw struct {
	Contract *States // Generic contract binding to access the raw methods on
}

// StatesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StatesCallerRaw struct {
	Contract *StatesCaller // Generic read-only contract binding to access the raw methods on
}

// StatesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StatesTransactorRaw struct {
	Contract *StatesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStates creates a new instance of States, bound to a specific deployed contract.
func NewStates(address common.Address, backend bind.ContractBackend) (*States, error) {
	contract, err := bindStates(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &States{StatesCaller: StatesCaller{contract: contract}, StatesTransactor: StatesTransactor{contract: contract}, StatesFilterer: StatesFilterer{contract: contract}}, nil
}

// NewStatesCaller creates a new read-only instance of States, bound to a specific deployed contract.
func NewStatesCaller(address common.Address, caller bind.ContractCaller) (*StatesCaller, error) {
	contract, err := bindStates(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StatesCaller{contract: contract}, nil
}

// NewStatesTransactor creates a new write-only instance of States, bound to a specific deployed contract.
func NewStatesTransactor(address common.Address, transactor bind.ContractTransactor) (*StatesTransactor, error) {
	contract, err := bindStates(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StatesTransactor{contract: contract}, nil
}

// NewStatesFilterer creates a new log filterer instance of States, bound to a specific deployed contract.
func NewStatesFilterer(address common.Address, filterer bind.ContractFilterer) (*StatesFilterer, error) {
	contract, err := bindStates(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StatesFilterer{contract: contract}, nil
}

// bindStates binds a generic wrapper to an already deployed contract.
func bindStates(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StatesABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_States *StatesRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _States.Contract.StatesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_States *StatesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _States.Contract.StatesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_States *StatesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _States.Contract.StatesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_States *StatesCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _States.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_States *StatesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _States.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_States *StatesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _States.Contract.contract.Transact(opts, method, params...)
}

// ComputeTeamRating is a free data retrieval call binding the contract method 0xbd64d4fa.
//
// Solidity: function computeTeamRating(uint256[] teamState) constant returns(uint256 rating)
func (_States *StatesCaller) ComputeTeamRating(opts *bind.CallOpts, teamState []*big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "computeTeamRating", teamState)
	return *ret0, err
}

// ComputeTeamRating is a free data retrieval call binding the contract method 0xbd64d4fa.
//
// Solidity: function computeTeamRating(uint256[] teamState) constant returns(uint256 rating)
func (_States *StatesSession) ComputeTeamRating(teamState []*big.Int) (*big.Int, error) {
	return _States.Contract.ComputeTeamRating(&_States.CallOpts, teamState)
}

// ComputeTeamRating is a free data retrieval call binding the contract method 0xbd64d4fa.
//
// Solidity: function computeTeamRating(uint256[] teamState) constant returns(uint256 rating)
func (_States *StatesCallerSession) ComputeTeamRating(teamState []*big.Int) (*big.Int, error) {
	return _States.Contract.ComputeTeamRating(&_States.CallOpts, teamState)
}

// GetCurrentShirtNum is a free data retrieval call binding the contract method 0xeb78b7b7.
//
// Solidity: function getCurrentShirtNum(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetCurrentShirtNum(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getCurrentShirtNum", playerState)
	return *ret0, err
}

// GetCurrentShirtNum is a free data retrieval call binding the contract method 0xeb78b7b7.
//
// Solidity: function getCurrentShirtNum(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetCurrentShirtNum(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetCurrentShirtNum(&_States.CallOpts, playerState)
}

// GetCurrentShirtNum is a free data retrieval call binding the contract method 0xeb78b7b7.
//
// Solidity: function getCurrentShirtNum(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetCurrentShirtNum(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetCurrentShirtNum(&_States.CallOpts, playerState)
}

// GetCurrentTeamId is a free data retrieval call binding the contract method 0xcd2105e8.
//
// Solidity: function getCurrentTeamId(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetCurrentTeamId(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getCurrentTeamId", playerState)
	return *ret0, err
}

// GetCurrentTeamId is a free data retrieval call binding the contract method 0xcd2105e8.
//
// Solidity: function getCurrentTeamId(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetCurrentTeamId(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetCurrentTeamId(&_States.CallOpts, playerState)
}

// GetCurrentTeamId is a free data retrieval call binding the contract method 0xcd2105e8.
//
// Solidity: function getCurrentTeamId(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetCurrentTeamId(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetCurrentTeamId(&_States.CallOpts, playerState)
}

// GetDefence is a free data retrieval call binding the contract method 0x51585b49.
//
// Solidity: function getDefence(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetDefence(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getDefence", playerState)
	return *ret0, err
}

// GetDefence is a free data retrieval call binding the contract method 0x51585b49.
//
// Solidity: function getDefence(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetDefence(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetDefence(&_States.CallOpts, playerState)
}

// GetDefence is a free data retrieval call binding the contract method 0x51585b49.
//
// Solidity: function getDefence(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetDefence(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetDefence(&_States.CallOpts, playerState)
}

// GetEndurance is a free data retrieval call binding the contract method 0x258e5d90.
//
// Solidity: function getEndurance(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetEndurance(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getEndurance", playerState)
	return *ret0, err
}

// GetEndurance is a free data retrieval call binding the contract method 0x258e5d90.
//
// Solidity: function getEndurance(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetEndurance(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetEndurance(&_States.CallOpts, playerState)
}

// GetEndurance is a free data retrieval call binding the contract method 0x258e5d90.
//
// Solidity: function getEndurance(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetEndurance(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetEndurance(&_States.CallOpts, playerState)
}

// GetLastSaleBlock is a free data retrieval call binding the contract method 0xc566b5bc.
//
// Solidity: function getLastSaleBlock(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetLastSaleBlock(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getLastSaleBlock", playerState)
	return *ret0, err
}

// GetLastSaleBlock is a free data retrieval call binding the contract method 0xc566b5bc.
//
// Solidity: function getLastSaleBlock(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetLastSaleBlock(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetLastSaleBlock(&_States.CallOpts, playerState)
}

// GetLastSaleBlock is a free data retrieval call binding the contract method 0xc566b5bc.
//
// Solidity: function getLastSaleBlock(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetLastSaleBlock(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetLastSaleBlock(&_States.CallOpts, playerState)
}

// GetMonthOfBirthInUnixTime is a free data retrieval call binding the contract method 0x85053566.
//
// Solidity: function getMonthOfBirthInUnixTime(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetMonthOfBirthInUnixTime(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getMonthOfBirthInUnixTime", playerState)
	return *ret0, err
}

// GetMonthOfBirthInUnixTime is a free data retrieval call binding the contract method 0x85053566.
//
// Solidity: function getMonthOfBirthInUnixTime(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetMonthOfBirthInUnixTime(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetMonthOfBirthInUnixTime(&_States.CallOpts, playerState)
}

// GetMonthOfBirthInUnixTime is a free data retrieval call binding the contract method 0x85053566.
//
// Solidity: function getMonthOfBirthInUnixTime(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetMonthOfBirthInUnixTime(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetMonthOfBirthInUnixTime(&_States.CallOpts, playerState)
}

// GetPass is a free data retrieval call binding the contract method 0x55a6f86f.
//
// Solidity: function getPass(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetPass(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getPass", playerState)
	return *ret0, err
}

// GetPass is a free data retrieval call binding the contract method 0x55a6f86f.
//
// Solidity: function getPass(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetPass(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPass(&_States.CallOpts, playerState)
}

// GetPass is a free data retrieval call binding the contract method 0x55a6f86f.
//
// Solidity: function getPass(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetPass(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPass(&_States.CallOpts, playerState)
}

// GetPlayerId is a free data retrieval call binding the contract method 0xf4385912.
//
// Solidity: function getPlayerId(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetPlayerId(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getPlayerId", playerState)
	return *ret0, err
}

// GetPlayerId is a free data retrieval call binding the contract method 0xf4385912.
//
// Solidity: function getPlayerId(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetPlayerId(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPlayerId(&_States.CallOpts, playerState)
}

// GetPlayerId is a free data retrieval call binding the contract method 0xf4385912.
//
// Solidity: function getPlayerId(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetPlayerId(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPlayerId(&_States.CallOpts, playerState)
}

// GetPrevLeagueId is a free data retrieval call binding the contract method 0xf8bd3e6e.
//
// Solidity: function getPrevLeagueId(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetPrevLeagueId(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getPrevLeagueId", playerState)
	return *ret0, err
}

// GetPrevLeagueId is a free data retrieval call binding the contract method 0xf8bd3e6e.
//
// Solidity: function getPrevLeagueId(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetPrevLeagueId(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPrevLeagueId(&_States.CallOpts, playerState)
}

// GetPrevLeagueId is a free data retrieval call binding the contract method 0xf8bd3e6e.
//
// Solidity: function getPrevLeagueId(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetPrevLeagueId(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPrevLeagueId(&_States.CallOpts, playerState)
}

// GetPrevShirtNumInLeague is a free data retrieval call binding the contract method 0x666d0070.
//
// Solidity: function getPrevShirtNumInLeague(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetPrevShirtNumInLeague(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getPrevShirtNumInLeague", playerState)
	return *ret0, err
}

// GetPrevShirtNumInLeague is a free data retrieval call binding the contract method 0x666d0070.
//
// Solidity: function getPrevShirtNumInLeague(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetPrevShirtNumInLeague(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPrevShirtNumInLeague(&_States.CallOpts, playerState)
}

// GetPrevShirtNumInLeague is a free data retrieval call binding the contract method 0x666d0070.
//
// Solidity: function getPrevShirtNumInLeague(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetPrevShirtNumInLeague(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPrevShirtNumInLeague(&_States.CallOpts, playerState)
}

// GetPrevTeamPosInLeague is a free data retrieval call binding the contract method 0x58a7a46a.
//
// Solidity: function getPrevTeamPosInLeague(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetPrevTeamPosInLeague(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getPrevTeamPosInLeague", playerState)
	return *ret0, err
}

// GetPrevTeamPosInLeague is a free data retrieval call binding the contract method 0x58a7a46a.
//
// Solidity: function getPrevTeamPosInLeague(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetPrevTeamPosInLeague(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPrevTeamPosInLeague(&_States.CallOpts, playerState)
}

// GetPrevTeamPosInLeague is a free data retrieval call binding the contract method 0x58a7a46a.
//
// Solidity: function getPrevTeamPosInLeague(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetPrevTeamPosInLeague(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetPrevTeamPosInLeague(&_States.CallOpts, playerState)
}

// GetShoot is a free data retrieval call binding the contract method 0x65b4b476.
//
// Solidity: function getShoot(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetShoot(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getShoot", playerState)
	return *ret0, err
}

// GetShoot is a free data retrieval call binding the contract method 0x65b4b476.
//
// Solidity: function getShoot(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetShoot(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetShoot(&_States.CallOpts, playerState)
}

// GetShoot is a free data retrieval call binding the contract method 0x65b4b476.
//
// Solidity: function getShoot(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetShoot(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetShoot(&_States.CallOpts, playerState)
}

// GetSkills is a free data retrieval call binding the contract method 0x0092bf78.
//
// Solidity: function getSkills(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetSkills(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getSkills", playerState)
	return *ret0, err
}

// GetSkills is a free data retrieval call binding the contract method 0x0092bf78.
//
// Solidity: function getSkills(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetSkills(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetSkills(&_States.CallOpts, playerState)
}

// GetSkills is a free data retrieval call binding the contract method 0x0092bf78.
//
// Solidity: function getSkills(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetSkills(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetSkills(&_States.CallOpts, playerState)
}

// GetSkillsVec is a free data retrieval call binding the contract method 0xcc1cc3d7.
//
// Solidity: function getSkillsVec(uint256 playerState) constant returns(uint16[5] skills)
func (_States *StatesCaller) GetSkillsVec(opts *bind.CallOpts, playerState *big.Int) ([5]uint16, error) {
	var (
		ret0 = new([5]uint16)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getSkillsVec", playerState)
	return *ret0, err
}

// GetSkillsVec is a free data retrieval call binding the contract method 0xcc1cc3d7.
//
// Solidity: function getSkillsVec(uint256 playerState) constant returns(uint16[5] skills)
func (_States *StatesSession) GetSkillsVec(playerState *big.Int) ([5]uint16, error) {
	return _States.Contract.GetSkillsVec(&_States.CallOpts, playerState)
}

// GetSkillsVec is a free data retrieval call binding the contract method 0xcc1cc3d7.
//
// Solidity: function getSkillsVec(uint256 playerState) constant returns(uint16[5] skills)
func (_States *StatesCallerSession) GetSkillsVec(playerState *big.Int) ([5]uint16, error) {
	return _States.Contract.GetSkillsVec(&_States.CallOpts, playerState)
}

// GetSpeed is a free data retrieval call binding the contract method 0x4b93f753.
//
// Solidity: function getSpeed(uint256 playerState) constant returns(uint256)
func (_States *StatesCaller) GetSpeed(opts *bind.CallOpts, playerState *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "getSpeed", playerState)
	return *ret0, err
}

// GetSpeed is a free data retrieval call binding the contract method 0x4b93f753.
//
// Solidity: function getSpeed(uint256 playerState) constant returns(uint256)
func (_States *StatesSession) GetSpeed(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetSpeed(&_States.CallOpts, playerState)
}

// GetSpeed is a free data retrieval call binding the contract method 0x4b93f753.
//
// Solidity: function getSpeed(uint256 playerState) constant returns(uint256)
func (_States *StatesCallerSession) GetSpeed(playerState *big.Int) (*big.Int, error) {
	return _States.Contract.GetSpeed(&_States.CallOpts, playerState)
}

// IsValidPlayerState is a free data retrieval call binding the contract method 0x19a4860c.
//
// Solidity: function isValidPlayerState(uint256 state) constant returns(bool)
func (_States *StatesCaller) IsValidPlayerState(opts *bind.CallOpts, state *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "isValidPlayerState", state)
	return *ret0, err
}

// IsValidPlayerState is a free data retrieval call binding the contract method 0x19a4860c.
//
// Solidity: function isValidPlayerState(uint256 state) constant returns(bool)
func (_States *StatesSession) IsValidPlayerState(state *big.Int) (bool, error) {
	return _States.Contract.IsValidPlayerState(&_States.CallOpts, state)
}

// IsValidPlayerState is a free data retrieval call binding the contract method 0x19a4860c.
//
// Solidity: function isValidPlayerState(uint256 state) constant returns(bool)
func (_States *StatesCallerSession) IsValidPlayerState(state *big.Int) (bool, error) {
	return _States.Contract.IsValidPlayerState(&_States.CallOpts, state)
}

// IsValidTeamState is a free data retrieval call binding the contract method 0xaf26723a.
//
// Solidity: function isValidTeamState(uint256[] state) constant returns(bool)
func (_States *StatesCaller) IsValidTeamState(opts *bind.CallOpts, state []*big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "isValidTeamState", state)
	return *ret0, err
}

// IsValidTeamState is a free data retrieval call binding the contract method 0xaf26723a.
//
// Solidity: function isValidTeamState(uint256[] state) constant returns(bool)
func (_States *StatesSession) IsValidTeamState(state []*big.Int) (bool, error) {
	return _States.Contract.IsValidTeamState(&_States.CallOpts, state)
}

// IsValidTeamState is a free data retrieval call binding the contract method 0xaf26723a.
//
// Solidity: function isValidTeamState(uint256[] state) constant returns(bool)
func (_States *StatesCallerSession) IsValidTeamState(state []*big.Int) (bool, error) {
	return _States.Contract.IsValidTeamState(&_States.CallOpts, state)
}

// PlayerStateCreate is a free data retrieval call binding the contract method 0x530f63d6.
//
// Solidity: function playerStateCreate(uint256 defence, uint256 speed, uint256 pass, uint256 shoot, uint256 endurance, uint256 monthOfBirthInUnixTime, uint256 playerId, uint256 currentTeamId, uint256 currentShirtNum, uint256 prevLeagueId, uint256 prevTeamPosInLeague, uint256 prevShirtNumInLeague, uint256 lastSaleBlock) constant returns(uint256 state)
func (_States *StatesCaller) PlayerStateCreate(opts *bind.CallOpts, defence *big.Int, speed *big.Int, pass *big.Int, shoot *big.Int, endurance *big.Int, monthOfBirthInUnixTime *big.Int, playerId *big.Int, currentTeamId *big.Int, currentShirtNum *big.Int, prevLeagueId *big.Int, prevTeamPosInLeague *big.Int, prevShirtNumInLeague *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "playerStateCreate", defence, speed, pass, shoot, endurance, monthOfBirthInUnixTime, playerId, currentTeamId, currentShirtNum, prevLeagueId, prevTeamPosInLeague, prevShirtNumInLeague, lastSaleBlock)
	return *ret0, err
}

// PlayerStateCreate is a free data retrieval call binding the contract method 0x530f63d6.
//
// Solidity: function playerStateCreate(uint256 defence, uint256 speed, uint256 pass, uint256 shoot, uint256 endurance, uint256 monthOfBirthInUnixTime, uint256 playerId, uint256 currentTeamId, uint256 currentShirtNum, uint256 prevLeagueId, uint256 prevTeamPosInLeague, uint256 prevShirtNumInLeague, uint256 lastSaleBlock) constant returns(uint256 state)
func (_States *StatesSession) PlayerStateCreate(defence *big.Int, speed *big.Int, pass *big.Int, shoot *big.Int, endurance *big.Int, monthOfBirthInUnixTime *big.Int, playerId *big.Int, currentTeamId *big.Int, currentShirtNum *big.Int, prevLeagueId *big.Int, prevTeamPosInLeague *big.Int, prevShirtNumInLeague *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _States.Contract.PlayerStateCreate(&_States.CallOpts, defence, speed, pass, shoot, endurance, monthOfBirthInUnixTime, playerId, currentTeamId, currentShirtNum, prevLeagueId, prevTeamPosInLeague, prevShirtNumInLeague, lastSaleBlock)
}

// PlayerStateCreate is a free data retrieval call binding the contract method 0x530f63d6.
//
// Solidity: function playerStateCreate(uint256 defence, uint256 speed, uint256 pass, uint256 shoot, uint256 endurance, uint256 monthOfBirthInUnixTime, uint256 playerId, uint256 currentTeamId, uint256 currentShirtNum, uint256 prevLeagueId, uint256 prevTeamPosInLeague, uint256 prevShirtNumInLeague, uint256 lastSaleBlock) constant returns(uint256 state)
func (_States *StatesCallerSession) PlayerStateCreate(defence *big.Int, speed *big.Int, pass *big.Int, shoot *big.Int, endurance *big.Int, monthOfBirthInUnixTime *big.Int, playerId *big.Int, currentTeamId *big.Int, currentShirtNum *big.Int, prevLeagueId *big.Int, prevTeamPosInLeague *big.Int, prevShirtNumInLeague *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _States.Contract.PlayerStateCreate(&_States.CallOpts, defence, speed, pass, shoot, endurance, monthOfBirthInUnixTime, playerId, currentTeamId, currentShirtNum, prevLeagueId, prevTeamPosInLeague, prevShirtNumInLeague, lastSaleBlock)
}

// PlayerStateEvolve is a free data retrieval call binding the contract method 0x8d216b52.
//
// Solidity: function playerStateEvolve(uint256 playerState, uint16 delta) constant returns(uint256 evolvedState)
func (_States *StatesCaller) PlayerStateEvolve(opts *bind.CallOpts, playerState *big.Int, delta uint16) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "playerStateEvolve", playerState, delta)
	return *ret0, err
}

// PlayerStateEvolve is a free data retrieval call binding the contract method 0x8d216b52.
//
// Solidity: function playerStateEvolve(uint256 playerState, uint16 delta) constant returns(uint256 evolvedState)
func (_States *StatesSession) PlayerStateEvolve(playerState *big.Int, delta uint16) (*big.Int, error) {
	return _States.Contract.PlayerStateEvolve(&_States.CallOpts, playerState, delta)
}

// PlayerStateEvolve is a free data retrieval call binding the contract method 0x8d216b52.
//
// Solidity: function playerStateEvolve(uint256 playerState, uint16 delta) constant returns(uint256 evolvedState)
func (_States *StatesCallerSession) PlayerStateEvolve(playerState *big.Int, delta uint16) (*big.Int, error) {
	return _States.Contract.PlayerStateEvolve(&_States.CallOpts, playerState, delta)
}

// SetCurrentShirtNum is a free data retrieval call binding the contract method 0xa95e858b.
//
// Solidity: function setCurrentShirtNum(uint256 state, uint256 currentShirtNum) constant returns(uint256)
func (_States *StatesCaller) SetCurrentShirtNum(opts *bind.CallOpts, state *big.Int, currentShirtNum *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "setCurrentShirtNum", state, currentShirtNum)
	return *ret0, err
}

// SetCurrentShirtNum is a free data retrieval call binding the contract method 0xa95e858b.
//
// Solidity: function setCurrentShirtNum(uint256 state, uint256 currentShirtNum) constant returns(uint256)
func (_States *StatesSession) SetCurrentShirtNum(state *big.Int, currentShirtNum *big.Int) (*big.Int, error) {
	return _States.Contract.SetCurrentShirtNum(&_States.CallOpts, state, currentShirtNum)
}

// SetCurrentShirtNum is a free data retrieval call binding the contract method 0xa95e858b.
//
// Solidity: function setCurrentShirtNum(uint256 state, uint256 currentShirtNum) constant returns(uint256)
func (_States *StatesCallerSession) SetCurrentShirtNum(state *big.Int, currentShirtNum *big.Int) (*big.Int, error) {
	return _States.Contract.SetCurrentShirtNum(&_States.CallOpts, state, currentShirtNum)
}

// SetCurrentTeamId is a free data retrieval call binding the contract method 0xc37b1c25.
//
// Solidity: function setCurrentTeamId(uint256 playerState, uint256 teamId) constant returns(uint256)
func (_States *StatesCaller) SetCurrentTeamId(opts *bind.CallOpts, playerState *big.Int, teamId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "setCurrentTeamId", playerState, teamId)
	return *ret0, err
}

// SetCurrentTeamId is a free data retrieval call binding the contract method 0xc37b1c25.
//
// Solidity: function setCurrentTeamId(uint256 playerState, uint256 teamId) constant returns(uint256)
func (_States *StatesSession) SetCurrentTeamId(playerState *big.Int, teamId *big.Int) (*big.Int, error) {
	return _States.Contract.SetCurrentTeamId(&_States.CallOpts, playerState, teamId)
}

// SetCurrentTeamId is a free data retrieval call binding the contract method 0xc37b1c25.
//
// Solidity: function setCurrentTeamId(uint256 playerState, uint256 teamId) constant returns(uint256)
func (_States *StatesCallerSession) SetCurrentTeamId(playerState *big.Int, teamId *big.Int) (*big.Int, error) {
	return _States.Contract.SetCurrentTeamId(&_States.CallOpts, playerState, teamId)
}

// SetLastSaleBlock is a free data retrieval call binding the contract method 0xaf76cd01.
//
// Solidity: function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) constant returns(uint256)
func (_States *StatesCaller) SetLastSaleBlock(opts *bind.CallOpts, state *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "setLastSaleBlock", state, lastSaleBlock)
	return *ret0, err
}

// SetLastSaleBlock is a free data retrieval call binding the contract method 0xaf76cd01.
//
// Solidity: function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) constant returns(uint256)
func (_States *StatesSession) SetLastSaleBlock(state *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _States.Contract.SetLastSaleBlock(&_States.CallOpts, state, lastSaleBlock)
}

// SetLastSaleBlock is a free data retrieval call binding the contract method 0xaf76cd01.
//
// Solidity: function setLastSaleBlock(uint256 state, uint256 lastSaleBlock) constant returns(uint256)
func (_States *StatesCallerSession) SetLastSaleBlock(state *big.Int, lastSaleBlock *big.Int) (*big.Int, error) {
	return _States.Contract.SetLastSaleBlock(&_States.CallOpts, state, lastSaleBlock)
}

// SetPrevLeagueId is a free data retrieval call binding the contract method 0x47f3d716.
//
// Solidity: function setPrevLeagueId(uint256 state, uint256 value) constant returns(uint256)
func (_States *StatesCaller) SetPrevLeagueId(opts *bind.CallOpts, state *big.Int, value *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "setPrevLeagueId", state, value)
	return *ret0, err
}

// SetPrevLeagueId is a free data retrieval call binding the contract method 0x47f3d716.
//
// Solidity: function setPrevLeagueId(uint256 state, uint256 value) constant returns(uint256)
func (_States *StatesSession) SetPrevLeagueId(state *big.Int, value *big.Int) (*big.Int, error) {
	return _States.Contract.SetPrevLeagueId(&_States.CallOpts, state, value)
}

// SetPrevLeagueId is a free data retrieval call binding the contract method 0x47f3d716.
//
// Solidity: function setPrevLeagueId(uint256 state, uint256 value) constant returns(uint256)
func (_States *StatesCallerSession) SetPrevLeagueId(state *big.Int, value *big.Int) (*big.Int, error) {
	return _States.Contract.SetPrevLeagueId(&_States.CallOpts, state, value)
}

// SetPrevTeamPosInLeague is a free data retrieval call binding the contract method 0x7ee0aebc.
//
// Solidity: function setPrevTeamPosInLeague(uint256 state, uint256 value) constant returns(uint256)
func (_States *StatesCaller) SetPrevTeamPosInLeague(opts *bind.CallOpts, state *big.Int, value *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "setPrevTeamPosInLeague", state, value)
	return *ret0, err
}

// SetPrevTeamPosInLeague is a free data retrieval call binding the contract method 0x7ee0aebc.
//
// Solidity: function setPrevTeamPosInLeague(uint256 state, uint256 value) constant returns(uint256)
func (_States *StatesSession) SetPrevTeamPosInLeague(state *big.Int, value *big.Int) (*big.Int, error) {
	return _States.Contract.SetPrevTeamPosInLeague(&_States.CallOpts, state, value)
}

// SetPrevTeamPosInLeague is a free data retrieval call binding the contract method 0x7ee0aebc.
//
// Solidity: function setPrevTeamPosInLeague(uint256 state, uint256 value) constant returns(uint256)
func (_States *StatesCallerSession) SetPrevTeamPosInLeague(state *big.Int, value *big.Int) (*big.Int, error) {
	return _States.Contract.SetPrevTeamPosInLeague(&_States.CallOpts, state, value)
}

// TeamStateAppend is a free data retrieval call binding the contract method 0xe06da0bf.
//
// Solidity: function teamStateAppend(uint256[] teamState, uint256 playerState) constant returns(uint256[] state)
func (_States *StatesCaller) TeamStateAppend(opts *bind.CallOpts, teamState []*big.Int, playerState *big.Int) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "teamStateAppend", teamState, playerState)
	return *ret0, err
}

// TeamStateAppend is a free data retrieval call binding the contract method 0xe06da0bf.
//
// Solidity: function teamStateAppend(uint256[] teamState, uint256 playerState) constant returns(uint256[] state)
func (_States *StatesSession) TeamStateAppend(teamState []*big.Int, playerState *big.Int) ([]*big.Int, error) {
	return _States.Contract.TeamStateAppend(&_States.CallOpts, teamState, playerState)
}

// TeamStateAppend is a free data retrieval call binding the contract method 0xe06da0bf.
//
// Solidity: function teamStateAppend(uint256[] teamState, uint256 playerState) constant returns(uint256[] state)
func (_States *StatesCallerSession) TeamStateAppend(teamState []*big.Int, playerState *big.Int) ([]*big.Int, error) {
	return _States.Contract.TeamStateAppend(&_States.CallOpts, teamState, playerState)
}

// TeamStateAt is a free data retrieval call binding the contract method 0x44328d1a.
//
// Solidity: function teamStateAt(uint256[] teamState, uint256 idx) constant returns(uint256 playerState)
func (_States *StatesCaller) TeamStateAt(opts *bind.CallOpts, teamState []*big.Int, idx *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "teamStateAt", teamState, idx)
	return *ret0, err
}

// TeamStateAt is a free data retrieval call binding the contract method 0x44328d1a.
//
// Solidity: function teamStateAt(uint256[] teamState, uint256 idx) constant returns(uint256 playerState)
func (_States *StatesSession) TeamStateAt(teamState []*big.Int, idx *big.Int) (*big.Int, error) {
	return _States.Contract.TeamStateAt(&_States.CallOpts, teamState, idx)
}

// TeamStateAt is a free data retrieval call binding the contract method 0x44328d1a.
//
// Solidity: function teamStateAt(uint256[] teamState, uint256 idx) constant returns(uint256 playerState)
func (_States *StatesCallerSession) TeamStateAt(teamState []*big.Int, idx *big.Int) (*big.Int, error) {
	return _States.Contract.TeamStateAt(&_States.CallOpts, teamState, idx)
}

// TeamStateCreate is a free data retrieval call binding the contract method 0x74833f73.
//
// Solidity: function teamStateCreate() constant returns(uint256[] state)
func (_States *StatesCaller) TeamStateCreate(opts *bind.CallOpts) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "teamStateCreate")
	return *ret0, err
}

// TeamStateCreate is a free data retrieval call binding the contract method 0x74833f73.
//
// Solidity: function teamStateCreate() constant returns(uint256[] state)
func (_States *StatesSession) TeamStateCreate() ([]*big.Int, error) {
	return _States.Contract.TeamStateCreate(&_States.CallOpts)
}

// TeamStateCreate is a free data retrieval call binding the contract method 0x74833f73.
//
// Solidity: function teamStateCreate() constant returns(uint256[] state)
func (_States *StatesCallerSession) TeamStateCreate() ([]*big.Int, error) {
	return _States.Contract.TeamStateCreate(&_States.CallOpts)
}

// TeamStateEvolve is a free data retrieval call binding the contract method 0x5c010782.
//
// Solidity: function teamStateEvolve(uint256[] teamState, uint8 delta) constant returns(uint256[])
func (_States *StatesCaller) TeamStateEvolve(opts *bind.CallOpts, teamState []*big.Int, delta uint8) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "teamStateEvolve", teamState, delta)
	return *ret0, err
}

// TeamStateEvolve is a free data retrieval call binding the contract method 0x5c010782.
//
// Solidity: function teamStateEvolve(uint256[] teamState, uint8 delta) constant returns(uint256[])
func (_States *StatesSession) TeamStateEvolve(teamState []*big.Int, delta uint8) ([]*big.Int, error) {
	return _States.Contract.TeamStateEvolve(&_States.CallOpts, teamState, delta)
}

// TeamStateEvolve is a free data retrieval call binding the contract method 0x5c010782.
//
// Solidity: function teamStateEvolve(uint256[] teamState, uint8 delta) constant returns(uint256[])
func (_States *StatesCallerSession) TeamStateEvolve(teamState []*big.Int, delta uint8) ([]*big.Int, error) {
	return _States.Contract.TeamStateEvolve(&_States.CallOpts, teamState, delta)
}

// TeamStateSize is a free data retrieval call binding the contract method 0x4444a283.
//
// Solidity: function teamStateSize(uint256[] teamState) constant returns(uint256 count)
func (_States *StatesCaller) TeamStateSize(opts *bind.CallOpts, teamState []*big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _States.contract.Call(opts, out, "teamStateSize", teamState)
	return *ret0, err
}

// TeamStateSize is a free data retrieval call binding the contract method 0x4444a283.
//
// Solidity: function teamStateSize(uint256[] teamState) constant returns(uint256 count)
func (_States *StatesSession) TeamStateSize(teamState []*big.Int) (*big.Int, error) {
	return _States.Contract.TeamStateSize(&_States.CallOpts, teamState)
}

// TeamStateSize is a free data retrieval call binding the contract method 0x4444a283.
//
// Solidity: function teamStateSize(uint256[] teamState) constant returns(uint256 count)
func (_States *StatesCallerSession) TeamStateSize(teamState []*big.Int) (*big.Int, error) {
	return _States.Contract.TeamStateSize(&_States.CallOpts, teamState)
}
