$('#submit-log').click(() => {
    $('#warning').css({"display": "none"});
    
    if (!validationEmail.validate($('#login').val())) {
      $('#warning').css({"display": "block"});
      return
    }
    if (!validationPass.validate($('#password').val())) {
      $('#warning').css({"display": "block"});
      return
    }

    $.ajax({
      type: "POST",
      url: "../login",
      data: `email=${$('#login').val()}&password=${$('#password').val()}`,
      success: function (msg) {
        if (msg == "Unsuccsess") {
          $('#warning').css({"display": "block"});
        }
        else {
          window.location.href = msg;
        }
      }
    });
  });

  $('#signUp').click(() => {
    window.location.href = "/registrationform";
  })