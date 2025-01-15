import { showPss, removespace, exists, headers, Sing } from "./modul.js";

showPss("changTYpe", "password");

headers("sign up", "signup");

exists().then((userExsit) => {
  if (userExsit) {
    window.location.href = "/";
  }
});

function singIn() {
  const bottum = document.getElementById("sing-in");
  const passwordd = document.getElementById("password");
  const userinf = document.getElementById("userInf");
  bottum.addEventListener("click", (event) => {
    event.preventDefault();
    if (userinf.value === "") {
      userinf.focus();
      userinf.style.border = "1.8px solid red";
      return;
    } else if (userinf.value) {
      userinf.style.border = "1px solid #ccc";
    }
    if (passwordd.value === "") {
      passwordd.focus();
      passwordd.style.border = "1.8px solid red";
      return;
    } else {
      passwordd.style.border = "1px solid #ccc";
    }
    const error = document.getElementById("error_message");
    error.style.color = "red";
    fetch("/user/signin", {
      method: "POST",
      body: new URLSearchParams({
        userInf: userinf.value,
        passwd: passwordd.value,
      }),
    })
      .then((res) => {
        if (res.ok) {
          window.location.href = "/";
        } else {
          res.text().then((data) => {
            error.style.width = "100%";
            error.style.height = "40px";
            error.style.textAlign = "center";
            error.style.color = "red";
            error.innerHTML = data;
            return;
          });
        }
      })
      .catch((erro) => {
        error.innerHTML = "check your internet please";
        error.style.width = "100%";
        error.style.height = "40px";
        error.style.textAlign = "center";
        error.style.color = "red";
      });
  });
}

removespace("userInf");

singIn();
Sing("github");
Sing("google");
