import { exists } from "./modul.js";
import {
  HandelHearder,
  notLogpopo,
} from "./homemudul.js";
import { FetchPost, Inf } from "./post.js";
async function Hande() {
  let islogin = false;
  islogin = await exists();
  if (!islogin) {
    notLogpopo();
  }
 await HandelHearder(islogin);
  FetchPost(20, islogin);
  Inf(islogin);
  // addLastPost();
}
Hande();
// ClosePop();
