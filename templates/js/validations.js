const validationEmail = {
    patternRegexp: /^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$/, 
    validate: function(email) {
        return this.patternRegexp.test(email)
    }
}

const validationPass = {
    patternRegexp: /(?=.*[0-9!#$%&'()*+,-.\/:;<=>?@[\]^_{|}~])(?=.*[a-z])[0-9a-zA-Z!#$%&'()*+,-.\/:;<=>?@[\]^_{|}~]{6,24}/, 
    validate: function(password) {
        return this.patternRegexp.test(password)
    }
}

const validationUsername = {
    validate: function(username) {
        if (username.length >= 3 && username.length <= 24) 
            return true
        return false
    }
}