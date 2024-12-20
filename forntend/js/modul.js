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
    input.value = input.value.replace(/\s/g, "");
  });
}

async function exists() {
  try {
    const res = await fetch("/user/exist", { method: "POST" });
    if (res.ok || res.status === 302) {
      const data = await res.json();
      localStorage.setItem("userName", data.userName);
    }

    return res.ok || res.status === 302;
  } catch (error) {
    console.error("Error checking user existence:", error);
    return false;
  }
}

function Format(sing, link) {
  return `
        <div class="img-div">
              <a href="/">
                <img src="../forntend/image/forum.jpeg" alt="" />
              </a>
              </div>
              <div class="link">
                   <a href="/${link}"><button class="btn-link">${sing}</button></a>
                       
        </div> 
  `;
}

function headers(sign1 , sign2) {
  const headers = document.createElement("header");
  headers.innerHTML = Format(sign1, sign2);
  document.body.prepend(headers);
}


export { showPss, removespace, exists, headers };
