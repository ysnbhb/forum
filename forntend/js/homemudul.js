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

function ClosePop() {
  let close = true;
  const click = document.getElementById("click");
  const btnclose = document.getElementById("close");
  const popup = document.getElementById("popup");
  click.addEventListener("click", (event) => {
    event.stopPropagation();
    popup.classList.remove("closepop");
    close = false;
  });
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

function HandelHearder(islogin) {
  const userName = localStorage.getItem("userName");
  const header = document.createElement("header");
  console.log(islogin);
  header.innerHTML = `
      <div class="img-div">
        <a href="/">
          <img src="../forntend/image/logo.png" alt="" />
        </a>
      </div>
    `;
  const div = document.createElement("div");
  if (islogin) {
    const popup = document.createElement("div");
    popup.className = "PrLog";
    popup.id = "PrLog";
    popup.innerHTML = `
  
      <div class="classic">
        <a href="/profile"> <span class="material-symbols-outlined"> account_circle </span>
            <p>profile</p>
        </a>
      </div>
      <div class="classic">
        <a href="/loguot"><span class="material-symbols-outlined">logout</span><p>loguot</p></a>
        
      </div>
       <div class="classic">
        <a href="/postAdd"><span class="material-symbols-outlined">add_circle</span><p>Add Post</p></a>
      </div>
    `;
    document.body.append(popup);
    div.className = "link";
    div.id = "link";
    div.innerHTML = `
    <div class="userName"> <p style="margin-top: 9px;" class="profile">${userName}</p>
          <span class="material-symbols-outlined icon" style="font-size: xx-large;"> person </span></a
        >`;

    header.append(div);
    ShowPop(div);
  } else {
    console.log("no log");
  }
  document.body.prepend(header);
}

export { ShowPop, ClosePop, HandelHearder };
