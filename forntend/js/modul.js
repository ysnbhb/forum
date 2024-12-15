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

function removespace(idInput) {
  const input = document.getElementById(idInput);
  input.addEventListener("input", () => {
    input.value = input.value.replace(/\s/g, ""); // Remove all spaces
  });
}

async function exists() {
  fetch("/user/exist", { method: "POST" }).then((res) => {
    if (res.status === 302 || res.ok) {
      window.location.href = "/";
    }
  });
}

function Format(sing1, sing2, link) {
  return `
        <div class="img-div">
              <a href="/">
                <img src="../forntend/image/login.jpg" alt="" />
              </a>
              </div>
                    <nav class="Navbar">welcome to page ${sing1}</nav>
              <div class="link">
                   <a href="/${link}"><button class="btn-link">${sing2}</button></a>
                       
        </div> 
  `;
}

function headers() {
  const singup = `sing up`;
  const singin = `sing in`;
  const headers = document.getElementById("header");
  if (window.location.href === "http://localhost:8081/singin") {
    headers.innerHTML = Format(singin, singup, "singup");
  } else {
    headers.innerHTML = Format(singup, singin, "singin");
  }
}

export { showPss, removespace, exists, headers };
