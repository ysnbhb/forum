function Fetch() {
  const btn = document.getElementById("fetch");
  btn.addEventListener("click", async () => {
    const input = document.getElementById("username");
    if (input.value === "") {
      input.focus();
      return;
    }
    try {
      const res = await fetch("/auth/user/signup", {
        method: "POST",
        body: `name=${input.value}`,
      });
      if (res.ok) {
        window.location.href = "/";
      } else {
        const err = document.getElementById("error");
        const data = await res.json();
        err.innerText = data.error;
        err.style.display = "block";
      }
    } catch (error) {
      const err = document.getElementById("error");
      err.innerText = "Internet problem";
      err.style.display = "block";
    }
  });
}

Fetch();
