import { exists } from "./modul.js";
function HandelHearder() {
  exists();
  const userName = localStorage.getItem("userName");
  console.log(userName);
}


HandelHearder()