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
    timerStarted = false;
    clearInterval(timerTask);
    $('.is-current-active-1').next().text("Completed");
});

$('#finish-planning').click(() => {
    socketInst.send(`FinishPlanning==`)
});

$(document).on('click', '.task-room', (e) => {
    console.log(e)
    console.log($(this))
    console.log(e.target.id)
    //$(this).css("font-size", "16px")
    //$(this).css({"font-size": "16px"})
    //$(e.target.id).css("font-size", "16px")
    $('.task-room').css({"font-size": "12px", "font-weight": "normal"})
    $('#'+e.target.id).css({"font-size": "14px", "font-weight": "bold"})
    
    console.log('#'+e.target.id)
});