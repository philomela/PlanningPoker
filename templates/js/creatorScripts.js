$('#start-voting').click(() => {
    socketInst.send("StartVoting")
});