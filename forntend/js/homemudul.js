function ShowPop() {
  let show = true;
  const link = document.getElementById("link");
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

export { ShowPop, ClosePop };
