
// start a Web3 server
const web3rpc = "ws://127.0.0.1:8545"
const web3 = new Web3(web3rpc);

// path to ABI file for interfacing with our blockchain functions
const ABIpath = '../build/contracts/Testing.json'

// TXs data
const gasPrice = 4000000000
const maxGas = 4000000

// useful variables
const maxNumPlayers = 11;
const nSkills = 7;
const bitPerState = 14;
const KEEPER = 0;
const DEF = 1;
const MID = 2;
const ATT = 3;
const tactica = createAlineacion(4,3,3);
const skillNames = ["Age","Defense","Speed","Pass","Shoot","Endurance","Role"];
const eth2wei = 1000000000000000000;

//Global variables
var addresses
var contractABI = ''
var contract
var nTeams = 0

var jsonTeamName
var jsonSelectedAddr



function retrieveAddr() {
    web3.eth.getAccounts(function (error, result) {
        if (error) {
            console.log(error)
            alert("Cannot retrieve addresses")
        } else {
            addresses = result
            jsonSelectedAddr = addresses[0]
            getABI()
        }    
    })
}

function deploy() {

    try {
        var newContract = new web3.eth.Contract(contractABI.abi)
        newContract.deploy({
            data: contractABI.bytecode
        })
            .send({
                from: jsonSelectedAddr,
                gas: maxGas,
                gasPrice: gasPrice
            }, function (error, transactionHash) { console.error(error) })
            .on('error', function (error) { console.error(error) })
            .on('transactionHash', function (transactionHash) { console.log(transactionHash) })
            .on('receipt', function (receipt) {
                console.log(receipt.contractAddress) // contains the new contract address
            })
            .on('confirmation', function (confirmationNumber, receipt) { console.log('Confirmed deploy') })
            .then(function (newContractInstance) {
                console.log(newContractInstance.options.address)
                //contractAddress = newContractInstance.options.address // instance with the new contract address
                contract = newContractInstance
                getDeployedContractAddress()
            });


    } catch (e) {
        console.error('Problem compiling:' + e);
    }
}

function getABI() {

    $.getJSON(ABIpath, function (result) {
        contractABI = result
        if (contractABI) {
            deploy()
        } else {
            log.error("Problem getting compiled contract, wrong build path?")
            alert("Problem getting compiled contract, wrong build path?")
        }
    });
}

window.onload = function () {
    retrieveAddr()
};

function getGanacheAddresses() {

    web3.eth.getAccounts(function (error, result) {
        if (error) {
            console.log(error)
            alert("getAccounts error.");
        } else {
            var code = ''
            for (var key in result) {
                code += '<a class="dropdown-item" href="#" onclick="loadAddressView(\'' + result[key] + '\'); return false;">' + result[key] + '</a>';
            }
            $('#listOfAccounts').html(code)
        }
    })
}

function createTeam() {
 
    if ($('#teamName').val() == '') {
        jsonTeamName = Math.random().toString(36).substring(2, 8)
        $('#teamName').html(jsonTeamName)
    } else {
        jsonTeamName = $('#teamName').val()
    }

    contract.methods.test_createTeam(jsonTeamName).send({ from: jsonSelectedAddr,
        gasPrice: 20000000000, gas: 6721975 }, function (error, transactionHash) {
        if (error) {
            console.log(error)
            alert("createTeam error.");
        } else {
            console.log("Create team successful");
            alert("Create team successful.");
            //Create the players for the created team
            createBalancedPlayer(nTeams, "messi",1, 0)
            nTeams++
        }
    })
    getSelectedAddressInfo()
}

//Creates team players
function createBalancedPlayer(teamId, playerName, userChoice, playerNumberInTeam) {

    var thisName = Math.random().toString(36).substring(2, 8)

    console.log("Creating player --> name:" + thisName + " id:" + teamId + " userCh:" + userChoice + " playerid:" + playerNumberInTeam + " role:" + tactica[playerNumberInTeam])
    contract.methods.test_createBalancedPlayer(thisName, teamId, userChoice, playerNumberInTeam, tactica[playerNumberInTeam])
        .send({ from: addresses[0], gasPrice: gasPrice, gas: maxGas }, function (error, transactionHash) {
            if (playerNumberInTeam < maxNumPlayers -1) {
                createBalancedPlayer(teamId, playerName, userChoice, playerNumberInTeam + 1)
            } else {
                getCreatedTeams()
            }
        })
}

function getCreatedTeams() {

    contract.methods.test_getNCreatedTeams().call({ from: jsonSelectedAddr }, function (error, result) {
        if (error) {
            console.log(error)
            alert("getNCreatedTeams error.");
        } else {
            console.log("accessing getCreatedTeams " + result);
            $('#createdTeams').html('')
            for (i = 0; i < result; i++) {
                console.log("team info: ");
                $('#createdTeams').append('<div onclick="getTeamInfo(' + i + ')"style="cursor: pointer; display: inline-block; background:#FFFF00;">' + i + '</div>');
            }
        }
    })
}


function getDeployedContractAddress() {

    $('#contractAddress').html(contract.options.address)
    getContractAddressBalance()
    getGanacheAddresses()
}

function playMatch() {

    var local = parseInt($('#localTeamInput').val())
    var visitor = parseInt($('#visitorTeamInput').val())
    nRounds=18
    var rnd1 = getRandArray(nRounds);
    var rnd2 = getRandArray(nRounds);
    var rnd3 = getRandArray(nRounds);
    var rnd4 = getRandArray(nRounds);
    //toni

    contract.methods.test_playGame(local, visitor, rnd1, rnd2, rnd3, rnd4).call({ from: jsonSelectedAddr, gasPrice: gasPrice, gas: maxGas }, function (error, result) {
        if (error) {
            console.log(error)
            alert("playGame error.");
        } else {
            $('#resultMatch').append("<p>SCORE: " + result[0] + " - " + result[1] + " (team " + local + " vs team " + visitor + ")</p>")
            $('#score1').html(result[0])
            $('#score2').html(result[1])
        }    
    })
}

function getRandArray(nElems) {
    randArray = [];
    for (var elem = 0; elem < nElems; elem++)
    {
        randArray.push(getRandomInt(1,1000));
    }
    return randArray;
}


async function getPlayersOfTeam(teamId) {

    if (teamId != 0 && !teamId) {
        teamId = $('#teamIdInput').val()
    }

    var html = ''

    html += "<b>Players in team " + teamId + ":<br>";
    for (key in skillNames) {
        html += " " + skillNames[key] + " "
    }
    html += "<br>";
    $('#playersList').html(html)


    for (var p=0;p<maxNumPlayers;p++) {
        getPlayerHTML(teamId,p);
    }
    console.log("HTML: "+html);

    $('#playersList').html(html)
}


async function getPlayerHTML(teamId, p) {
    contract.methods.test_getSkill(teamId,p).call({ from: jsonSelectedAddr, gasPrice: gasPrice, gas: maxGas }, function (error, result) {
        if (error) {
            console.log(error)
            alert("test_getSkill error.");
        } else {
            var serializedSkills = result;
            contract.methods.test_decode(nSkills,serializedSkills,bitPerState).call({ from: jsonSelectedAddr, gasPrice: gasPrice, gas: maxGas }, function (error, result) {
                if (error) {//toni
                    console.log(error)
                    alert("test_decode error.");
                } else {
                    state = result;
                    html = "<b>Player " + p + "</b>";
                    for (var sk = 0; sk < nSkills; sk++) {
                        if (sk==0) state[sk] = unixMonthToAge(state[sk]);
                        html += " " + state[sk] + " "
                    }
                    html += "<br>";
                    console.log("html: " + html);
                    $('#playersList').append(html);
                }
            });
        }
    });
}


function unixMonthToAge(unixMonthOfBirth) {
    // in July 2018, we are at month 582 after 1970.
    age = (582 - unixMonthOfBirth)/12;
    return parseInt(age*10)/10;
  }

  /**
 * Returns a random integer between min (inclusive) and max (inclusive)
 * Using Math.round() will give you a non-uniform distribution!
 */
function getRandomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

function loadAddressView(_selectedAddress) {
    $('#selectedAddress').html(_selectedAddress)
    jsonSelectedAddr = _selectedAddress; 
    getSelectedAddressInfo()
}

function getSelectedAddressInfo() {
    web3.eth.getBalance(jsonSelectedAddr, function (error, result) {
        if (error) {
            console.log(error)
            alert("getSelectedAddressInfo error.");
        } else {       
            $('#selectedAddressBalance').html(result/eth2wei)
        }
    })
}

function getContractAddressBalance() {
    web3.eth.getBalance(contract.options.address, function (error, result) {
         if (error) {
             console.log(error)
             alert("getContractAddressBalance error.");
         } else {
            $('#contractAddressBalance').html(result/eth2wei)
         }
    })
}

function getTeamInfo (teamId) {
    getPlayersOfTeam(teamId)
    getSelectedTeamName(teamId)
}

function getSelectedTeamName(teamId){
    contract.methods.test_getTeamName(teamId).call({ from: jsonSelectedAddr, gasPrice: gasPrice, gas: maxGas }, function (error, result) {
        if (error) {
            console.log(error)
            alert("getSelectedTeamName error.");
        } else {
            $('#selectedTeamName').html(result)
        }
    })
}




function createAlineacion(nDef,nMid,nAtt) {
    alineacion = [0];
    for (var p = 0; p<nDef; p++) {
        alineacion.push(DEF);
    }
    for (var p = 0; p<nMid; p++) {
        alineacion.push(MID);
    }
    for (var p = 0; p<nAtt; p++) {
        alineacion.push(ATT);
    }
    return alineacion;
}


