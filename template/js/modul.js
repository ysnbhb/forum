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

export { showPss, removespace };
