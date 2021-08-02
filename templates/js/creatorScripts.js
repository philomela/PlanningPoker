var startVotingData, selectedTaskId;

$('#start-voting').click(() => {

    socketInst.send(`StartVoting==${startVotingData}`) //вынести в константу
    startVotingData = null;
    buttonFinishVoting = $('#finish-voting').prop("disabled", false);
    buttonFinishVoting.css({ "background": "#ff0000" });
});

$(document).on('click', '.task-room', (e) => {
    e.preventDefault()
    selectedTaskId = e.target.id.split('-')[1];
    startVotingData = `<Change><StartVoting taskId="${selectedTaskId}" isCurrentActive="1"/></Change>` //вынести в константу
});

$('#finish-voting').click(() => {
    socketInst.send(`StopVoting==`) //вынести в константу
    timerStarted = false;
    clearInterval(timerTask);
    $('.is-current-active-1').next().text("Completed");
    buttonFinishVoting = $('#finish-voting').prop("disabled", true);
    buttonFinishVoting.css({ "background": "#b3b3b3" });
});

$('#finish-planning').click(() => {
    socketInst.send(`FinishPlanning==`)
});

$(document).on('click', '.task-room', (e) => {
    buttonStartVoting = $(`#start-voting`).prop("disabled", false);
    buttonStartVoting.css({ "background": "#40AA29" });
    $('.task-room').css({ "font-size": "12px", "font-weight": "normal" })
    $('#' + e.target.id).css({ "font-size": "14px", "font-weight": "bold" })

});