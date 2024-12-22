import { exists } from "./modul.js";
import { ClosePop, HandelHearder, ShowPop, addLastPost } from "./homemudul.js";
async function Hande() {
  let islogin = false;
  islogin = await exists();

  HandelHearder(islogin);
  addLastPost();
}
Hande();
ClosePop();
