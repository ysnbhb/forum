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

function runButtom(idButtom) {
  const buttom = document.getElementById(idButtom);
  buttom.addEventListener("click", (env) => {
    env.preventDefault();
    if (Checkvalid()) {
      return;
    } else {
    }
  });
}

function CheckExicest(idInput) {
  const input = document.getElementById(idInput);
  const debounc = debounce(() => {
    if (validator.isEmail(input.value)) {
      console.log("valid");
    } else {
      console.log("not valid");
    }
    console.log(input.value);
  }, 1000);
  input.addEventListener("input", debounc);
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

CheckExicest("user_name");
CheckExicest("email");

function Checkvalid() {
  const password = document.getElementById("password");
  const confpassword = document.getElementById("confpassword");
  const buttom = document.getElementById("sing-up");
  if (password.value === "" || confpassword.value !== password.value) {
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

runButtom("sing-up");
