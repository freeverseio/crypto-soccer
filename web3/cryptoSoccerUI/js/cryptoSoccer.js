
var url = "http://localhost:3000"

function tmp() {
    $('#players').text = 'Here will appear the players'
}


function getGanacheAddresses() {

    $.get(url + '/ganache', function (data) {

        if (!data || jQuery.isEmptyObject(data)) {
            $('#ganacheAddresses').html('ERROR: Ganache is not running')
        } else {
            var html = ''
            for (var key in data) {
                html += data[key] + '<br>'
            }
            $('#ganacheAddresses').html(html)
        }

    })
}

function createTeam() {

    $.get(url + '/createTeam', function (data) {

        if (!data || jQuery.isEmptyObject(data)) {
            $('#createdTeams').append('ERROR: Problem creating team')
        } else {

            $('#createdTeams').html('')
            for (i = 0; i < data; i++) {
                $('#createdTeams').append(i + '<br>')
            }
        }
    })
}

function playMatch() {

    var local = $('#localTeamInput').val()
    var visitor = $('#visitorTeamInput').val()

    var queryUrl = url + '/playMatch?localTeam=' + local + '&visitorTeam=' + visitor

    $.get(queryUrl, function (data) {

        if (!data || jQuery.isEmptyObject(data)) {
            $('#resultMatch').append('ERROR: Problem playing the match')
        } else {

            $('#resultMatch').html('')

            $('#createdTeams').append(data)

        }
    })
}

function getCreatedTeams() {

    $.get(url + '/createdTeams', function (data) {

        if (!data || jQuery.isEmptyObject(data)) {
            $('#createdTeams').html('ERROR: Problem getting created team')
        } else {
            if (data > 0) {
                $('#createdTeams').html('')
                for (i = 0; i < data; i++) {
                    $('#createdTeams').append(i + '<br>')
                }
            }
        }
    })

}

function getDeployedContractAddress() {

    $.get(url + '/deployedContractAddress', function (data) {

        if (!data || jQuery.isEmptyObject(data)) {
            $('#contractAddress').html('ERROR: Problem getting contract address')
        } else {
            $('#contractAddress').html(data)
        }
    })
}

function getPlayersOfTeam() {

    var teamId = $('#teamIdInput').val()
    $.get(url + '/playersOfTeam?teamId=' + teamId, function (data) {

        if (!data || jQuery.isEmptyObject(data)) {
            $('#playersList').html('ERROR: Problem getting players list for team ' + teamId)
        } else {
            $('#playersList').html(data)
        }
    })


}