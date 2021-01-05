//Общие скрипты веб-сервиса
        let socketInst, currentXmlString, currentXml;

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

        function createWebSocket() {
            socketInst = new WebSocket('ws://localhost:8080/echo?roomId=' + $.getUrlVar('roomId'));
            socketInst.onopen = function (event) {
                socketInst.send("ChangeGetVM")

            }
            socketInst.onmessage = function (event) {
                currentXmlString = event.data;
                switch (currentXmlString) {
                    case "StartTask": 
                    setInterval(startTimeTask($('.is-current-active-1').next()), 100);
                }
                
                parseXmlResponse(currentXml);
                //alert(currentXmlString);
            }
        }

        function startTimeTask(timerElement){
            timerElement.text = timerElement.text - 1;
        }

        function parseXmlResponse(currentXml) {
            $('.right-menu-room ul, .name-meeting h1, .tasks-left-menu-room').empty();
            var parserXml = new DOMParser();
            currentXml = parserXml.parseFromString(currentXmlString, "text/xml");
            var persons = currentXml.getElementsByTagName('Persons')[0].getElementsByTagName('Person')
            if (persons != null && persons != undefined)
                for (var i = 0; i < persons.length; i++) {
                    $(`<li class="person${persons[i].getAttribute('Id')}">${persons[i].getAttribute('UserName')}</li>`).appendTo($('.right-menu-room ul'));
                }
            var nameMeeting = currentXml.getElementsByTagName('Room')[0].getAttribute('NameRoom');
            if (nameMeeting != null & persons != undefined) {
                console.log(nameMeeting)
                $('.name-meeting h1').text(nameMeeting);
            }
            var tasks = currentXml.getElementsByTagName('Tasks')[0].getElementsByTagName('Task')
            if (tasks != null && tasks != undefined) {
                console.log(tasks)
                for (var i = 0; i < tasks.length; i++) {
                    $(`<div id="task-${tasks[i].getAttribute('Id')}" class="task-room is-current-active-${tasks[i].getAttribute('IsCurrentActive')}">${tasks[i].getAttribute('NameTask')}</div><div class="time-task-discussion">
                        ${tasks[i].getAttribute('TimeDiscussion')} :min</div>`).appendTo($('.tasks-left-menu-room'));
                }
            } 
            var votes = currentXml.getElementsByTagName('Tasks')[0].getElementsByTagName('PersonTask');
            if (votes != null && votes != undefined) {
                console.log(votes)
                for (var i = 0; i < votes.length; i++) {
                    if (votes[i].getAttribute('Vote') != 0) {
                        $(`.person${votes[i].getAttribute('PersonId')}`).css({"background": "red"})
                        //alert('hello')
                    }
                    
                }
            } 

            console.log(persons);
        }

        function GetChosenTask() {
            
        }

        $('.button-room').click(() => {
            socketInst.send("ChangeVote==<Change><AddVote vote='1' /></Change>")
        });
        
        $('#yellow').click(() => {
            SendColor('yellow')
        });
        $('#green').click(() => {
            SendColor('green')
        });

       
        createWebSocket();


        parseXmlResponse(currentXml); //test