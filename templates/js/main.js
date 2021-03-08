let socketInst, currentXmlString, currentXml;
const coffeeIcon = 999, questionIcon = 777, coffeeIconMedian = "999.00", questionIconMedian = "777.00"

$.extend({
    getUrlVars: function () {
        var vars = [], hash;
        var hashes = window.location.href.slice(window.location.href.indexOf('?') + 1).split('&');
        for (var i = 0; i < hashes.length; i++) {
            hash = hashes[i].split('=');
            vars.push(hash[0]);
            vars[hash[0]] = hash[1];
        }
        return vars;
    },
    getUrlVar: function (name) {
        return $.getUrlVars()[name];
    }
});

function parseXmlResponse(currentXml) {
    
    $('.right-menu-room ul, .name-meeting h1, .tasks-left-menu-room, #tb-results').empty();
    var parserXml = new DOMParser();
    countTimerStarted++;
    currentXml = parserXml.parseFromString(currentXmlString, "text/xml");
    
    var persons = currentXml.getElementsByTagName('Persons')[0].getElementsByTagName('Person')
    if (persons != null && persons != undefined)
        for (var i = 0; i < persons.length; i++) {
            $(`<li class="person${persons[i].getAttribute('Id')}"><i class="fa fa-user-o" aria-hidden="true"></i><span class="user-name-text-right-menu">${persons[i].getAttribute('UserName')}<span></li>`).appendTo($('.right-menu-room ul'));
        }
    var nameMeeting = currentXml.getElementsByTagName('Room')[0].getAttribute('NameRoom');
    if (nameMeeting != null & persons != undefined) {
        console.log(nameMeeting)
        $('.name-meeting h1').text(nameMeeting);
    }
    var tasks = currentXml.getElementsByTagName('Tasks')[0].getElementsByTagName('Task')
    if (tasks != null && tasks != undefined) {
        console.log(tasks)
             
        for (let i = 0; i < tasks.length; i++) {
                
            if (tasks[i].getAttribute('IsCurrentActive') == 1) {
                console.log(tasks[i].getAttribute('IsCurrentActive'));
                
                if (countTimerStarted <= 1){
                    StartTimerTask(tasks[i].getAttribute('TimeDiscussion'));
                }
                
                if (!timerStarted){
                    StartTimerTask(tasks[i].getAttribute('TimeDiscussion'));
                    
                    
                }
                
                var currentPersonTasks = tasks[i].getElementsByTagName('PersonTask')
                for (k = 0; k < currentPersonTasks.length; k++) {
                    if (currentPersonTasks[k].getAttribute('Vote') != 0) {
                        $(`.person${currentPersonTasks[k].getAttribute('PersonId')}`).css({ "background": "#4d387e", "color": "white" });
                    }
                }
            }
            $(`<div id="task-${tasks[i].getAttribute('Id')}" class="task-room is-current-active-${tasks[i].getAttribute('IsCurrentActive')}">${tasks[i].getAttribute('NameTask')}</div><div class="time-task-discussion">
                        ${tasks[i].getAttribute('TimeDiscussion')} :min</div>`).appendTo($('.tasks-left-menu-room'));
                        if (tasks[i].getAttribute('Completed') == 1){
                            $(`#task-${tasks[i].getAttribute('Id')}`).css({ "text-decoration": "line-through", "pointer-events": "none"});
                        }
                        var hasCurrentActive = false;
                        for (let i = 0; i < tasks.length; i++){
                            if (tasks[i].getAttribute('IsCurrentActive') == 1){
                                hasCurrentActive = true;
                            }
                        }

                        if (tasks[i].getAttribute('Completed') == 0 && tasks[i].getAttribute('IsCurrentActive') == 0 && hasCurrentActive){
                            $(`#task-${tasks[i].getAttribute('Id')}`).css({ "color": "#CCCCCC", "pointer-events": "none"});
                            $(`#task-${tasks[i].getAttribute('Id')}`).next().css({"color": "#CCCCCC", "pointer-events": "none"})
                        }

                        if (i == tasks.length - 1 && !hasCurrentActive) {
                            clearInterval(timerTask);
                            timerStarted = false;
                            $('.is-current-active-1').next().text("Completed");
                        }

            $(`<tr id="tsk-tb-${tasks[i].getAttribute('Id')}" class="tsk-tb"><td id="tb-name-task">${tasks[i].getAttribute('NameTask')}</tr>`).appendTo($('#tb-results'));

            var currPersonTasks = tasks[i].getElementsByTagName('PersonTask');

            for (let i = 0; i < currPersonTasks.length; i++) {
                var currentPersonId = currPersonTasks[i].getAttribute('PersonId');
                for (let k = 0; k < persons.length; k++) {
                    if (persons[k].getAttribute('Id') == currentPersonId) {
                        $(`<td id="tsk-tb-person-${currentPersonId}">${persons[k].getAttribute('UserName')}</td>`).appendTo($('#tb-results').children().last());

                    }
                }
            }
            if (!timerStarted) {
                $(`<td id="tsk-tb-median">Median:</td>`).appendTo($('#tb-results').children().last());
                $(`<tr id="tsk-tb-${tasks[i].getAttribute('Id')}" class="tsk-tb-nested"><td id="tb-name-task"></tr>`).appendTo($('#tb-results'));
                for (let j = 0; j < currPersonTasks.length; j++) {
                    var currentPersonId = currPersonTasks[j].getAttribute('PersonId');
                    for (let l = 0; l < persons.length; l++) {
                        if (persons[l].getAttribute('Id') == currentPersonId) {
                            var currentScore = currPersonTasks[j].getAttribute('Score');
                            if(currentScore == questionIcon){
                                currentScore = '<i class="fa fa-question" aria-hidden="true"></i>'
                            }
                            else if (currentScore == coffeeIcon)
                                currentScore = '<i class="fa fa-coffee" aria-hidden="true"></i>'
                            $(`<td id="tsk-tb-person-${currentPersonId}-score">${currentScore}</td>`).appendTo($('#tb-results').children().last());

                        }

                    }
                }
                var currentMedian = tasks[i].getAttribute('Median');
                console.log(typeof(currentMedian))
                console.log(typeof(coffeeIconMedian))
                if (!currentMedian){
                    currentMedian = "Unknown";
                }
                else if (currentMedian == questionIconMedian){
                    currentMedian = '<i class="fa fa-question" aria-hidden="true"></i>';
                }
                else if (currentMedian == coffeeIconMedian) {
                    currentMedian = '<i class="fa fa-coffee" aria-hidden="true"></i>';
                }
                else {
                    currentMedian = tasks[i].getAttribute('Median');
                } 
                
                $(`<td id="tsk-tb-median-tsk-${tasks[i].getAttribute('Id')}-score">${currentMedian}</td>`).appendTo($('#tb-results').children().last());
            }
        }

    }

    console.log(persons);
}

var timerStarted = false;
var timerTask;
var countTimerStarted = 0;

function StartTimerTask(initTime) {
    timerStarted = true;
    if ($('.is-current-active-1') == null || $('.is-current-active-1') == undefined)
        return

    var timeTask = initTime * 60, stopTime = timeTask * 1000;
    console.log(timeTask)
        timerTask = setInterval(() => {
        timeTask -= 1;
        timeTaskOut = `${Math.trunc(timeTask / 60) }-m: ${Math.trunc(timeTask % 60)}-sec`
        $('.is-current-active-1').next().text(timeTaskOut);
        $('.is-current-active-1').next().css({"background": "#6a155d", "color": "white", "padding": "3px"})
        console.log(timeTask);
    }, 1000)
    setTimeout(() => { 
        clearInterval(timerTask); 
        timerStarted = false; //Пересмотреть логику когда никто не проголосовал за задачу, после истечения таймера пропадает complete
        $('.is-current-active-1').next().text("Completed");
    }, stopTime)
    console.log(stopTime)
}

$(document).on('click', '.button-room', (e) => {
    e.preventDefault()
    socketInst.send(`ChangeVote==<Change><AddVote vote="1" score="${e.target.id.split('-')[1]}"/></Change>`)
    console.log(e.target)
    console.log(e.target.id.split('-')[1])
});

$('#yellow').click(() => {
    SendColor('yellow')
});
$('#green').click(() => {
    SendColor('green')
});


createWebSocket();


parseXmlResponse(currentXml); //test