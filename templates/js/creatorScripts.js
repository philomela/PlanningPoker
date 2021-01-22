var startVotingData;
$('#start-voting').click(() => {

    socketInst.send(`StartVoting==${startVotingData}`)
    console.log(startVotingData)
});

$(document).on('click', '.task-room', (e) => {
    e.preventDefault()
    startVotingData = `<Change><StartVoting taskId="${e.target.id.split('-')[1]}" isCurrentActive="1"/></Change>`
    console.log(startVotingData)
});

$('#finish-voting').click(() => {
    socketInst.send(`StopVoting==`)
    timerStarted = 0;
})