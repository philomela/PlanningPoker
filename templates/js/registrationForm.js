$('#submit-reg').click(() => {
    $('.notification-email, .notification-username, .notification-pass').css({ "display": "none" });
    $('#email, #userName, #password').css({ "border": "none" })

    email = $('#email').val();
    userName = $('#userName').val();
    password = $('#password').val();

    if (!validationEmail.validate(email)) {
      $('.notification-email').css({ "display": "block" });
      $('#email').css({ "border": "1px solid red" })
      return
    }
    if (!validationUsername.validate(userName)) {
      $('.notification-username').css({ "display": "block" });
      $('#userName').css({ "border": "1px solid red" })
      return
    }
    if (!validationPass.validate(password)) {
      $('.notification-pass').css({ "display": "block" })
      $('#password').css({ "border": "1px solid red" })
      return
    }

    $('#warning').css({ "display": "block" });

    $.ajax({
      type: "POST",
      url: "../registration",
      data: `email=${email}&password=${password}&userName=${userName}`,
      success: function (msg) {
        if (msg == "Succsess") {
          window.location.href = '../';
        }
        else if (msg == "Unsuccsess") {
          $('#warning').text("pass or email or username no valid");
        }
        else {
          $('#warning').text(msg);
        }
      }
    });
  });

  $('#signIn').click(() => {
    window.location.href = "/loginform";
  })