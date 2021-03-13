var startVotingData, selectedTaskId;

$('#start-voting').click(() => {

    socketInst.send(`StartVoting==${startVotingData}`)
    startVotingData = null;
});

$(document).on('click', '.task-room', (e) => {
    e.preventDefault()
    selectedTaskId = e.target.id.split('-')[1];
    startVotingData = `<Change><StartVoting taskId="${selectedTaskId}" isCurrentActive="1"/></Change>`
    
    console.log(selectedTaskId)
});

$('#finish-voting').click(() => {
    socketInst.send(`StopVoting==`)
    timerStarted = false;
    clearInterval(timerTask);
    $('.is-current-active-1').next().text("Completed");
    $(document).ready(function() {
        buttonStartVoting.prop("disabled", false); 
        buttonStartVoting.css({"background": "#40AA29"});                   
    });
});

$('#finish-planning').click(() => {
    socketInst.send(`FinishPlanning==`)
});

$(document).on('click', '.task-room', (e) => {
    $('.task-room').css({"font-size": "12px", "font-weight": "normal"})
    $('#'+e.target.id).css({"font-size": "14px", "font-weight": "bold"})
});