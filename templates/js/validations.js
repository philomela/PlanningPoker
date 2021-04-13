const validationEmail = {
    patternRegexp: /^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$/,
    validate: function (email) {
        return this.patternRegexp.test(email)
    }
}

const validationPass = {
    patternRegexp: /(?=.*[0-9!#$%&'()*+,-.\/:;<=>?@[\]^_{|}~])(?=.*[a-z])[0-9a-zA-Z!#$%&'()*+,-.\/:;<=>?@[\]^_{|}~]{6,24}/,
    validate: function (password) {
        return this.patternRegexp.test(password)
    }
}

const validationUsername = {
    validate: function (username) {
        if (username.length >= 3 && username.length <= 24)
            return true
        return false
    }
}

const validationNameRoom = {
    validateNameRoom: function (nameRoom) {
        if (nameRoom.length >= 1 && nameRoom.trim().length != 0)
            return true
        return false
    },
    validateTasksNames: function (tasks) {
        for(let task of tasks){
            if (task.getAttribute('name').length > 1)
                continue;
            else return false;
        }
        return true;
    },
    validateTimeDiscussions: function (tasks) {
        for(let task of tasks){
            console.log(task.getAttribute('time-discussion'))
            if (task.getAttribute('time-discussion') >= 1 && task.getAttribute('time-discussion') <= 10)
                continue;
            else return false;
        }
        return true;
    }
}
