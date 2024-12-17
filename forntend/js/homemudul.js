function ShowPop() {
  let show = true;
  const link = document.getElementById("link");
  link.addEventListener("click", () => {
    if (show) {
      document.getElementById("PrLog").classList.add("showPop");
    } else {
      document.getElementById("PrLog").classList.remove("showPop");
    }
    show = !show;
  });
}

export { ShowPop };
