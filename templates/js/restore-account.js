var second = 10;
var timer = window.setTimeout(redirectFunc, 10000);
var interval = window.setInterval(() => {
    second--
    $('.second-show').text(second)
}, 1000)

function redirectFunc() {
    window.location.href = "/"
}

$('#signUp').click(() => {
    window.location.href = "/registrationform";
})