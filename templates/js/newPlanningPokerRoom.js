let xmlTasks, nameRoom, tasks = new Set();
function createXml() {
    let _xmlForSend = '<tasks>'
    let tasksCollection = $('.tasks-controls').children('#task');
    let timeCollection = $('.task-timer').children('#time-discussion');
    for (let element of tasksCollection) {
        tasks.add(`<task name="${element.value}" `);
    }

    for (let i = 0; i < timeCollection.length; i++) {
        _xmlForSend += `${Array.from(tasks)[i]}time-discussion="${timeCollection[i].value}"/>`;
    }
    _xmlForSend += "</tasks>";

    tasks.clear()
    return _xmlForSend;
}

$('#change-location').click(() => {
    location.href = $('#link-room').val();
});

$('#add-new-task').click(() => {
    $('.tasks-controls').append(`<label for="task">Task and time for discussion:</label><input type="text" class="form-control" id="task" name="task">`);
    $('.task-timer').append(`<label for="time-discussion">Minutes:</label><input type="number" min="1" max="10" class="form-control time-discuss" id="time-discussion" name="time-task">`);
});

$('#push-tasks').click(() => {
    xmlTasks = null
    nameRoom = null
    $(".error-data-new-room").css({ "display": "none" })
    nameRoom = $('#name-meeting').val();
    xmlTasks = createXml();

    let parserXml = new DOMParser();
    let currentXmlTasks = parserXml.parseFromString(xmlTasks, "text/xml");

    let tasksForSend = currentXmlTasks.getElementsByTagName('tasks')[0].getElementsByTagName('task')

    if (!validationNameRoom.validateNameRoom(nameRoom) ||
        !validationNameRoom.validateTasksNames(tasksForSend) ||
        !validationNameRoom.validateTimeDiscussions(tasksForSend)) {
        $(".error-data-new-room").css({ "display": "block" });
        return false
    }

    $.ajax({
        type: "POST",
        url: "../create-room",
        data: `nameRoom=${nameRoom}&xmlTasks=${xmlTasks}`,
        success: function (msg) {
            switch (msg) {
                case "error":
                    $(".error-data-new-room").css({ "display": "block" });
                    return;
            }
            $('#link-room').val(msg);
        }
    });
});

$('#signUp').click(() => {
    window.location.href = "/registrationform";
})