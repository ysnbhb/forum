import { showPss, removespace, exists, headers, Sing } from "./modul.js";
headers("sign in", "signin");

exists()
  .then((userExists) => {
    if (userExists) {
      window.location.href = "/";
    }
  })
  .catch((error) => {
    console.error("An error occurred:", error);
  });

function confermetPas() {
  const password = document.getElementById("password");
  const confpassword = document.getElementById("confpassword");
  const passwdlabel = document.getElementById("passwd_lable");
  const confpasswordlabel = document.getElementById("confermet_lable");
  if (password.value === "") {
    passwdlabel.innerHTML = "invalid password";
    passwdlabel.style.color = "red";
    passwdlabel.style.fontSize = "12px";
    passwdlabel.style.marginLeft = "8px";
    password.style.border = "1.8px solid red";
    password.focus();
    return true;
  } else {
    passwdlabel.innerHTML = "";
    password.style.border = "1px solid #ccc";
  }
  if (password.value !== confpassword.value) {
    confpasswordlabel.style.color = "red";
    confpasswordlabel.style.fontSize = "12px";
    confpasswordlabel.style.marginLeft = "8px";
    confpasswordlabel.innerHTML = "please confirm Password";
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
  const labe_user = document.getElementById("user_label");
  const passwd_lable = document.getElementById("email_label");
  button.addEventListener("click", (event) => {
    event.preventDefault();
    labe_user.style.color = "red";
    labe_user.style.marginLeft = "8px";
    labe_user.style.fontSize = "12px";
    passwd_lable.style.color = "red";
    passwd_lable.style.marginLeft = "8px";
    passwd_lable.style.fontSize = "12px";
    if (userName.value === "") {
      labe_user.innerHTML = "Invalid input";
      userName.focus();
      userName.style.border = "1.8px solid red";
      return;
    } else if (userName.value.length > 10) {
      labe_user.innerHTML = "user name is too long max 10 letter";
      userName.focus();
      userName.style.border = "1.8px solid red";
      return;
    } else {
      labe_user.innerHTML = "";
      userName.style.border = "1px solid #ccc";
    }
    if (!validator.isEmail(email.value)) {
      passwd_lable.innerHTML = "Invalid email form";
      email.focus();
      email.style.border = "1.8px solid red";
      return;
    } else {
      passwd_lable.innerHTML = "";
      email.style.border = "1px solid #ccc";
    }
    if (confermetPas()) {
      return;
    }
    const div = document.getElementById("error_message");
    fetch("../user/signup", {
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
          window.location.href = "/";
        } else {
          res.json().then((data) => {
            console.log(data);
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

Sing("google");
Sing("github");
removespace("email");
removespace("user_name");
runButton("sing-up");
showPss("changTYpe", "password");
showPss("changContype", "confpassword");
