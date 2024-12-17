import { exists, showPss } from "./modul.js";
import { ClosePop, ShowPop } from "./homemudul.js";
function HandelHearder() {
  exists().then((userExict) => {
    if (userExict) {
    } else {
    }
  });
  const userName = localStorage.getItem("userName");
  console.log(userName);
}

ShowPop();
ClosePop();

HandelHearder();
