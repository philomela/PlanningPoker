$('#createNewRoomButton').click(() => {
    window.location.href = "/newroom";
   
});
$('#signUp').click(() => {
    window.location.href = "/registrationform";
   
});
$('#signIn').click(() => {
    window.location.href = "/loginform";
   
});

$(document).ready(function () {
    $("a.scrollto").click(function () {
        elementClick = $(this).attr("href")
        destination = $(elementClick).offset().top;
        $("html:not(:animated),body:not(:animated)").animate({ scrollTop: destination }, 1800);
        return false;
    });
});