import { getCgt } from "./homemudul.js";

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

function headers(sign1, sign2) {
  const headers = document.createElement("header");
  headers.innerHTML = Format(sign1, sign2);
  document.body.prepend(headers);
}

function ctreatAddpost() {
  const add = document.getElementById("div-icon");
  let show = true;
  const div = document.createElement("div");
  div.addEventListener("click", (event) => {
    event.stopPropagation();
  });
  add.addEventListener("click", (event) => {
    event.stopPropagation();
    if (show) {
      document.body.append(div);
      div.className = "postpop";
      div.innerHTML = `
            <form method="post">
          <label for="img">Upload Image:</label>
          <label class="upload-icon" for="img" id="icon-img">
            <span class="material-symbols-outlined">add_a_photo</span>
          </label>
          <input type="file" id="img" name="img" accept="image/*" />
  
          <label for="title">Post Title:</label>
          <input
            type="text"
            id="title"
            name="title"
            placeholder="Enter post title"
            required
          />
  
          <label for="contant">Content:</label>
          <textarea
            id="contant"
            name="content"
            placeholder="Write your content here"
            required
          ></textarea>
            <div class="checkbox" id="checkbox"></div>
          <button type="submit" id="submit">Submit</button>
        </form>
        `;
      FormatCheckbok();
      addPost();
    } else {
      div.remove();
    }
    show = !show;
  });
  window.addEventListener("click", () => {
    console.log(show);
    if (!show) {
      div.remove();
      show = !show;
    }
  });
}

function addPost() {
  const btn = document.getElementById("submit");
  btn.addEventListener("click", (event) => {
    event.preventDefault();
    const img = document.getElementById("img");
    const title = document.getElementById("title");
    const contant = document.getElementById("contant");
    if (title.value === "") {
      title.focus();
      return;
    } else {
      title.style.border = "";
    }
    if (contant.value === "") {
      contant.focus();
      return;
    } else {
      contant.style.border = "";
    }
    if (!img.type.startsWith("image/")) {
      img.focus();
      return;
    }
    const form = new FormData();
    form.append("img", img.files[0]);
    form.append("title", title.value);
    form.append("contant", contant.value);
    fetch("/api/addPost", { method: "POST", body: form });
  });
}

async function FormatCheckbok() {
  const div = document.getElementById("checkbox");
  const ctg = await getCgt();
  console.log(ctg);
  for (let i = 0; i < ctg.length; i++) {
    console.log(ctg);
    
    const chek = document.createElement("div");
    chek.className = "checkbox-container";
    chek.innerHTML = `
        <input type="checkbox" id="${ctg[i]}" />
            <label for="${ctg[i]}">${ctg[i]}</label>
    `;
    div.append(chek);
  }
}

export { showPss, removespace, exists, headers, ctreatAddpost };
