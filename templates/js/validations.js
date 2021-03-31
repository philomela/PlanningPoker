const validationEmail = {
    patternRegexp:/^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$/, 
    validate: function(email) {
        return this.patternRegexp.test(email)
    }
}
