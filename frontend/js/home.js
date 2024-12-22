import { exists } from "./modul.js";
import { ClosePop, HandelHearder, ShowPop } from "./homemudul.js";
async function Hande() {
  let islogin = false;
  islogin = await exists();

  HandelHearder(islogin);
}
Hande();
ClosePop();
