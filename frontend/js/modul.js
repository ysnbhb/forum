import { getCgt, getCheckedCheckboxes } from "./homemudul.js";
import { CreateDiv } from "./post.js";

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

function Sing(id) {
  const div = document.getElementById(id);
  div.addEventListener("click", () => {
    window.location.href = `/auth/${id}`;
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
                <img src="../frontend/image/forum.jpeg" alt="" />
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
    // console.log(show);
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
            <div id="error"></div>
         <button type="submit" id="submit">Submit</button> 
        </form>
        `;
      FormatCheckbok();
      ReplacePhoto();
      addPost(div);
      show = false;
    } else {
      show = true;
      div.remove();
    }
  });
  window.addEventListener("click", () => {
    if (!show) {
      div.remove();
      show = true;
    }
  });
}

async function addPost(div) {
  const btn = document.getElementById("submit");
  const err = document.getElementById("error");
  btn.addEventListener("click", async (event) => {
    event.preventDefault();

    const img = document.getElementById("img");
    const title = document.getElementById("title");
    const content = document.getElementById("contant");

    if (title.value.trim() === "") {
      title.focus();
      title.style.border = "2px solid red";
      err.innerText = "Title is required";
      return false;
    } else {
      title.style.border = "";
    }

    if (content.value.trim() === "") {
      content.focus();
      title.style.border = "2px solid red";
      err.innerText = "Title is required";
      return false;
    } else {
      content.style.border = "";
    }
    const categories = getCheckedCheckboxes();
    if (categories.length === 0) {
      err.innerText = ` please aprove  categorie of this post`;
      err.style.color = "red";
      err.style.height = "30px";
      err.style.width = "90%";
      err.style.textAlign = "center";
      // err.classList.add("error-visible");

      return;
    } else {
      err.innerText = ``;
      err.style.height = "0";
    }
    if (img.files[0] && !img.files[0].type.startsWith("image/")) {
      img.focus();
      return false;
    }
    const allpost = document.getElementById("allPost");
    const form = new FormData();
    form.append("img", img.files[0]);
    form.append("title", title.value);
    form.append("content", content.value);
    form.append("categories", JSON.stringify(categories));

    try {
      const response = await fetch("/api/addPost", {
        method: "POST",
        body: form,
      });
      // console.log(response.json());
      const post = await response.json();
      if (response.ok) {
        const postdiv = CreateDiv(post, true, false);
        allpost.prepend(postdiv);
        div.remove();
      } else {
        err.innerHTML = post.error;
      }
      return true;
    } catch (error) {
      err.innerText = `Check your internet connection!`;
      err.style.color = "red";
      err.style.height = "30px";
      err.style.width = "90%";
      err.style.textAlign = "center";
      console.error("Error:", error);
      alert("An error occurred. Please try again.");
      return false; // Return false if there was an error
    }
  });
}

async function FormatCheckbok() {
  const div = document.getElementById("checkbox");
  const ctg = await getCgt();
  for (let i = 0; i < ctg.length; i++) {
    const chek = document.createElement("div");
    chek.className = "checkbox-container";
    chek.innerHTML = `
        <input type="checkbox" id="${ctg[i]}" value=${ctg[i]} />
            <label for="${ctg[i]}">${ctg[i]}</label>
    `;
    div.append(chek);
  }
}

function ReplacePhoto() {
  const span = document.getElementById("icon-img");
  const inputFile = document.getElementById("img");

  inputFile.addEventListener("change", function () {
    const file = inputFile.files[0];
    if (file) {
      const reader = new FileReader();

      // Read the file and create a preview
      reader.onload = function (e) {
        span.innerHTML = `
            <img src="${e.target.result}" alt="Preview" style="max-width: 60px; height : 100%;">
            <button class="delete-btn" id="deletePhoto">x</button>
         `;

        const deleteBtn = document.getElementById("deletePhoto");
        deleteBtn.addEventListener("click", function () {
          inputFile.value = "";
          span.innerHTML = `<span class="material-symbols-outlined">add_a_photo</span>`;
        });
      };

      reader.readAsDataURL(file);
    }
  });
}

export { showPss, removespace, exists, headers, ctreatAddpost, Sing };
