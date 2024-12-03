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

  if (password && confpassword) {
    confpassword.addEventListener("input", () => {
      if (password.value !== confpassword.value) {
        confpassword.style.border = "1px solid red";
        return true;
      } else {
        confpassword.style.border = "1px solid green";
        return false;
      }
    });
  } else {
    console.error("Password or confirm password input not found!");
  }
}

function runButton(idButton) {
  const button = document.getElementById(idButton);
  const userName = document.getElementById("user_name");
  const email = document.getElementById("email");
  const password = document.getElementById("password");
  const label = document.createElement("label");

  button.addEventListener("click", (event) => {
    event.preventDefault();

    // Validation logic
    if (
      Checkvalid() &&
      !CheckExists("user_name", "user_lable", label, userName.value) &&
      !CheckExists("email", "email_lable", label, email.value)
    ) {
      console.error("Validation failed.");
      return;
    }

    // Make the POST request
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
          window.open("/login");
        } else {
          const err = document.createElement("div");
          res.json().then((data) => {
            err.innerHTML = data.error;
          });
        }
      })
      .catch((error) => {
        console.error("Network error:", error);
      });
  });
}

function useDebount(idInput, idlable) {
  const input = document.getElementById(idInput);
  const div = document.createElement("label");
  const debounc = debounce(() => {
    CheckExists(idInput, idlable, div, input.value);
  }, 1000);
  input.addEventListener("input", debounc);
}

function CheckExists(idInput, idlable, div, value) {
  div.innerHTML = "";
  if (value === "") {
    console.log(30);
    return false;
  }
  const lable = document.getElementById(idlable);
  div.style.color = "red";
  div.style.fontSize = "10px";
  div.style.marginLeft = "5px";
  lable.append(div);
  if (idInput === "email") {
    if (!validator.isEmail(value) && value !== "") {
      div.innerHTML = "invalid email";
      return false;
    } else {
      div.innerHTML = "";
    }
  }
  fetch(`/user/check?checker=${value}`, {
    method: "POST",
  }).then((resp) => {
    if (resp.ok) {
      div.innerHTML = "";
      console.log(0);
      return true;
    } else {
      console.log(100);
      if (idInput === "email") {
        div.innerHTML = "email is ready used try anther email";
      } else {
        div.innerHTML = "user name is ready used try anther email";
      }
      return false;
    }
  });
}

function debounce(func, wait) {
  let timer;
  return function (...arg) {
    clearTimeout(timer);
    timer = setTimeout(() => {
      func(...arg);
    }, wait);
  };
}

useDebount("user_name", "user_lable");
useDebount("email", "email_lable");

function Checkvalid() {
  const password = document.getElementById("password");
  const confpassword = document.getElementById("confpassword");
  const user_name = document.getElementById("user_name");
  const email = document.getElementById("email");
  const buttom = document.getElementById("sing-up");
  if (
    password.value === "" ||
    confpassword.value !== password.value ||
    user_name.value === "" ||
    !validator.isEmail(email.value)
  ) {
    buttom.style.cursor = "not-allowed";
    return true;
  } else {
    buttom.style.cursor = "pointer";
    return false;
  }
}

setInterval(Checkvalid, 100);
confermetPas();
showPss("changTYpe", "password");
showPss("changContype", "confpassword");

runButton("sing-up");
