import { exists } from "./modul.js";
import {
  ClosePop,
  HandelHearder,
  ShowPop,
  addLastPost,
  getCgt,
} from "./homemudul.js";
async function Hande() {
  let islogin = false;
  islogin = await exists();

  HandelHearder(islogin);
  addLastPost();
  // console.log(await getCgt());
}
Hande();
ClosePop();
