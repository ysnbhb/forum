function showPss(idButtom, idPassword) {
  let hiden = true;
  const button = document.getElementById(idButtom);
  button.addEventListener("click", (event) => {
    event.preventDefault();
    const password = document.getElementById(idPassword);
    if (hiden) {
      password.setAttribute("type", "text");
      button.innerHTML = "visibility_off";
    } else {
      password.setAttribute("type", "password");
      button.innerHTML = "visibility";
    }
    hiden = !hiden;
  });
}

function confermetPas() {
  const password = document.getElementById("password");
  const confpassword = document.getElementById("confpassword");
  const passwdlabel = document.getElementById("passwd_lable");
  const confpasswordlabel = document.getElementById("confermet_lable");
  if (password.value === "") {
    passwdlabel.innerHTML = "invalid password";
    passwdlabel.style.color = "red";
    passwdlabel.style.fontSize = "15px";
    passwdlabel.style.marginLeft = "8px";
    password.style.border = "1.8px solid red";
    return true;
  } else {
    passwdlabel.innerHTML = "";
    password.style.border = "1px solid #ccc";
  }
  if (password.value !== confpassword.value) {
    confpasswordlabel.innerHTML = "please confirm Password";
    confpasswordlabel.style.color = "red";
    confpasswordlabel.style.fontSize = "15px";
    confpasswordlabel.style.marginLeft = "8px";
    confpassword.style.border = "1.8px solid red";
    confpassword.focus();
    return true;
  } else {
    confpasswordlabel.innerHTML = "";
    confpassword.style.border = "1px solid green";
    return false;
  }
}

function runButton(idButton) {
  const button = document.getElementById(idButton);
  const userName = document.getElementById("user_name");
  const email = document.getElementById("email");
  const password = document.getElementById("password");
  button.addEventListener("click", (event) => {
    event.preventDefault();
    if (userName.value === "") {
      userName.focus();
      userName.style.border = "1.8px solid red";
      return;
    } else {
      userName.style.border = "1px solid #ccc";
    }
    if (!validator.isEmail(email.value)) {
      email.focus();
      email.style.border = "1.8px solid red";
      return;
    } else {
      email.style.border = "1px solid #ccc";
    }
    if (confermetPas()) {
      return;
    }
    const div = document.getElementById("error_message");
    fetch("../user/singup", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        user_name: userName.value,
        passwd: password.value,
        email: email.value,
      }),
    })
      .then((res) => {
        if (res.ok) {
          window.location.href = "/login";
        } else {
          res.json().then((data) => {
            if (data.error) {
              div.innerHTML = data.error;
              div.style.width = "100%";
              div.style.height = "40px";
              div.style.textAlign = "center";
              div.style.color = "red";
            }
          });
        }
      })
      .catch((error) => {
        div.innerHTML = "check your Network please";
        div.style.width = "100%";
        div.style.height = "40px";
        div.style.textAlign = "center";
        div.style.color = "red";
      });
  });
}
showPss("changTYpe", "password");
showPss("changContype", "confpassword");

runButton("sing-up");
