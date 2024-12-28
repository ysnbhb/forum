import { ctreatAddpost } from "./modul.js";
function ShowPop(link) {
  let show = true;
  const popup = document.getElementById("PrLog");
  link.addEventListener("click", (event) => {
    event.stopPropagation();
    if (show) {
      popup.classList.add("showPop");
    } else {
      popup.classList.remove("showPop");
    }
    show = !show;
  });

  document.body.addEventListener("click", () => {
    if (!show) {
      popup.classList.remove("showPop");
      show = true;
    }
  });
}

function ClosePop(close) {
  const btnclose = document.getElementById("close");
  const popup = document.getElementById("popup");
  if (close) {
    popup.classList.remove("closepop");
  }
  // click.addEventListener("click", (event) => {
  //   event.stopPropagation();
  //   close = false;
  // });
  btnclose.addEventListener("click", (event) => {
    event.stopPropagation();
    popup.classList.add("closepop");
    close = true;
  });
  popup.addEventListener("click", (event) => {
    event.stopPropagation();
  });
  document.body.addEventListener("click", () => {
    if (!close) {
      popup.classList.add("closepop");
      close = true;
    }
  });
}

async function HandelHearder(islogin) {
  const userName = localStorage.getItem("userName");
  const header = document.createElement("header");
  header.innerHTML = `
      <div class="img-div">
        <a href="/">
          <img src="../frontend/image/logo.png" alt="" />
        </a>
      </div>
    `;
  const div = document.createElement("div");
  if (islogin) {
    const popup = document.createElement("div");
    const divIcon = document.createElement("div");
    divIcon.className = "div-icon";
    divIcon.id = "div-icon";

    divIcon.innerHTML = `
          <span class="material-symbols-outlined" style="font-size: 35px ; color: #000;">
            add
          </span>
    `;
    popup.className = "PrLog";
    popup.id = "PrLog";
    popup.innerHTML = `
      <div class="classic">
        <a href="/logout"><span class="material-symbols-outlined">logout</span><p>logout</p></a>
      </div>

    `;
    document.body.append(popup, divIcon);
    div.className = "link";
    div.id = "link";
    div.innerHTML = `
    <div class="userName"> <p style="margin-top: 9px;" class="profile">${userName}</p>
          <span class="material-symbols-outlined icon" style="font-size: xx-large;"> person </span></a
        >`;

    ShowPop(div);
    ctreatAddpost();
  } else {
    div.className = "nolog_link";
    div.innerHTML = `
        <a href="/signup" class="sign"
        ><button class="nolog_bnt">Signup</button></a
        >
        <a href="/signin" class="sign"
        ><button class="nolog_bnt">Signin</button></a
        >
        `;
  }
  header.append(div);
  document.body.prepend(header);
  return;
}

async function addLastPost() {
  let lastId, newId;
  let errorSErve = 0;
  let timer;

  const res = await fetch("/post/lastId");
  lastId = await res.text();

  async function getLastID() {
    try {
      const res = await fetch("/post/lastId");
      newId = await res.text();

      if (lastId !== newId) {
        const div = document.createElement("div");
        div.innerHTML = "New post has been published!";
        div.className = "interErro";
        div.style.backgroundColor = "#71bb49c9";
        document.body.append(div);
        setTimeout(() => div.remove(), 3000);

        lastId = newId;
      }
    } catch (error) {
      console.error("Error fetching last ID:", error, errorSErve);
      errorSErve++;

      if (errorSErve > 1) {
        if (!document.getElementById("interErro")) {
          const div = document.createElement("div");
          div.innerHTML = "Check your internet connection!";
          div.className = "interErro";
          div.id = "interErro";
          div.style.backgroundColor = "#ff0000ab";
          document.body.append(div);
        }

        clearInterval(timer);
      }
    }
  }

  timer = setInterval(getLastID, 5000);
}

function getCheckedCheckboxes() {
  const checkboxes = document.querySelectorAll(
    'input[type="checkbox"]:checked'
  );
  const checkedValues = Array.from(checkboxes).map(
    (checkbox) => checkbox.value
  );
  return checkedValues;
}

async function getCgt() {
  return await fetch("/api/getCategorie").then(async (response) => {
    return await response.json().then((data) => {
      return data;
    });
  });
}

function notLogpopo() {
  const div = document.createElement("div");
  div.className = "singpopup";
  div.classList.add("closepop");
  div.id = "popup";
  div.innerHTML = `
        <button class="close" id="close">
          <span class="material-symbols-outlined"> close </span>
        </button>
        <div class="popup">
          <div class="sing">You are Not Login :</div>
          <div class="link-sing">
            <div class="sing-in">
              <a href="/signin"
                ><button class="btn-link">you have account?</button></a
              >
            </div>

            <div class="sing-up">
              <a href="/signup"
                ><button class="btn-link">you don't have account?</button></a
              >
            </div>
          </div>
        </div>
  `;
  document.body.append(div);
}

export {
  ShowPop,
  ClosePop,
  HandelHearder,
  addLastPost,
  getCgt,
  getCheckedCheckboxes,
  notLogpopo,
};
