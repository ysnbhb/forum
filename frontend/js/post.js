import { ClosePop } from "./homemudul.js";

export function CreateDiv(post, islogin, ispopup) {
  const articles = document.createElement("article");
  articles.className = "articles";

  const postInf = document.createElement("div");
  postInf.className = "postInfo";

  const h2 = document.createElement("h2");
  h2.className = "uerName";
  h2.innerText = post.userName;

  const span = document.createElement("span");
  span.className = "date";
  span.innerText = post.date;

  postInf.append(h2, span);
  articles.append(postInf);

  const h4 = document.createElement("h4");
  h4.className = "title";
  h4.innerText = post.title;

  const div = document.createElement("div");
  div.className = "contant";
  div.innerText = post.contant;

  articles.append(h4, div);
  if (post.img) {
    const imgdiv = document.createElement("div");
    imgdiv.className = "uerImg";
    const src = document.createElement("img");
    src.src = post.img;
    imgdiv.append(src);
    articles.append(imgdiv);
  }

  articles.append(CreateCtg(post.categories));
  const btndiv = document.createElement("div");
  btndiv.className = "bottom";

  const reaction = document.createElement("div");
  reaction.className = "reaction";
  const [like, dislike] = HandulLike(
    post.reaction.type,
    post.reaction.numLike,
    post.reaction.numDisLike,
    islogin,
    post.id
  );
  reaction.append(like, dislike);
  btndiv.append(reaction);
  if (!ispopup) {
    const commant = document.createElement("button");
    commant.className = "commant";
    commant.innerHTML = `
        <span class="material-symbols-outlined"> comment </span>
  `;
    commant.addEventListener("click", () => {
      HandulPostCommant(islogin, post.id);
    });
    btndiv.append(commant);
  }
  articles.append(btndiv);
  const commants = document.createElement("div");
  commants.id = "commant";
  articles.append(commants);
  return articles;
}

export async function FetchPost(offset, islogin) {
  const allpost = document.getElementById("allPost");
  try {
    const res = await fetch(`/api/posts?offset=${offset}`);
    const posts = await res.json();
    if (posts.length === 0 && offset == 0) {
      document.getElementById("content").style.display = "none";
      return;
    }
    for (let i = 0; i < posts.length; i++) {
      const post = CreateDiv(posts[i], islogin);
      // console.log(posts[i]);
      for (let commante of posts[i].commant) {
        const divcom = CreateCommate(commante, islogin);

        post.querySelector("#commant").append(divcom);
      }
      allpost.append(post);
    }
    return posts.length;
  } catch (err) {}
}

function CreateCtg(categories) {
  const cgt = document.createElement("div");
  for (let categorie of categories) {
    const span = document.createElement("span");
    span.innerText = categorie;
    span.className = "type";
    cgt.append(span);
  }
  cgt.className = "cgt";
  return cgt;
}

function HandulLike(type, numlike, numdislike, islogin, postid) {
  const likebtn = document.createElement("button");
  const dislikebtn = document.createElement("button");
  likebtn.className = "like";
  dislikebtn.className = "dislike";

  if (type === "likes") {
    likebtn.classList.add(type);
  } else if (type === "dislikes") {
    dislikebtn.classList.add(type);
  }

  const icondlike = document.createElement("span");
  const icondislike = document.createElement("span");
  icondlike.className = "material-symbols-outlined";
  icondislike.className = "material-symbols-outlined";
  icondislike.style.fontSize = "14px";
  icondlike.style.fontSize = "14px";
  icondislike.innerHTML = "thumb_down";
  icondlike.innerHTML = "thumb_up";
  likebtn.append(icondlike);
  dislikebtn.append(icondislike);
  const likespan = document.createElement("span");
  const dislikespan = document.createElement("span");
  likespan.innerHTML = numlike;
  dislikespan.innerHTML = numdislike;
  likebtn.append(likespan);
  dislikebtn.append(dislikespan);
  let close = true;
  likebtn.addEventListener("click", () => {
    if (!islogin) {
      ClosePop(close);
      return;
    }
    if (likebtn.classList.length == 2) {
      likebtn.classList.remove("likes");
      numlike--;
    } else {
      if (dislikebtn.classList.length == 2) {
        dislikebtn.classList.remove("dislikes");
        numdislike--;
      }
      numlike++;
      likebtn.classList.add("likes");
    }
    likespan.innerHTML = numlike;
    dislikespan.innerHTML = numdislike;
    fetch(`/api/post/like?postid=${postid}&type=likes`, { method: "POST" });
  });
  dislikebtn.addEventListener("click", () => {
    if (!islogin) {
      ClosePop(close);
      return;
    }
    if (dislikebtn.classList.length == 2) {
      dislikebtn.classList.remove("dislikes");
      numdislike--;
    } else {
      if (likebtn.classList.length == 2) {
        likebtn.classList.remove("likes");
        numlike--;
      }
      numdislike++;
      dislikebtn.classList.add("dislikes");
    }
    likespan.innerHTML = numlike;
    dislikespan.innerHTML = numdislike;
    fetch(`/api/post/like?postid=${postid}&type=dislikes`, {
      method: "POST",
    }).then((res) => {
      // console.log(res);
    });
  });
  return [likebtn, dislikebtn];
}

// FetchPost(0);

export function Inf(islogin) {
  let offset = 20;
  let lenghtpost;
  window.addEventListener("scroll", async () => {
    let windowHight = window.innerHeight;
    let scrol = window.scrollY;
    if (scrol + windowHight > document.body.scrollHeight - 1000) {
      if (lenghtpost === 0) {
        window.removeEventListener("scroll", Inf);
        return;
      }
      lenghtpost = await FetchPost(offset, islogin);
      offset += 20;
    }
  });
}

async function HandulPostCommant(islogin, postid) {
  const postopop = document.getElementById("postopop");
  postopop.innerHTML = `
     <button class="close" id="closepost">
        <span class="material-symbols-outlined"> close </span>
      </button>
  `;
  const res = await fetch(`/api/post?postid=${postid}`);
  const post = await res.json();
  const div = CreateDiv(post, islogin, true);
  postopop.append(div);
  const commant = document.createElement("div");
  commant.className = "postcommat";
  commant.id = "postcommat";

  postopop.classList.remove("closePose");
  ClosePost(postopop);
  const len = await Handulcomment(commant, postid, 0, islogin);
  if (len == 0) {
    commant.innerHTML = "no commant YET";
    commant.style.height = "auto";
  } else {
    const showmore = document.createElement("span");
    showmore.className = "showMore";
    showmore.innerHTML = "show more";
    let offset = 5;
    showmore.addEventListener("click", () => {
      Morecommate(showmore, offset, postid, islogin, commant);
      offset += 5;
    });
    commant.append(showmore);
  }
  postopop.append(commant);
}

function ClosePost(postpop) {
  const btn = document.getElementById("closepost");
  btn.addEventListener("click", () => {
    postpop.innerHTML = "";
    postpop.classList.add("closePose");
  });
}

async function Handulcomment(div, postId, offset, islogin) {
  const res = await fetch(`/api/commant?offset=${offset}&postid=${postId}`);
  const commants = await res.json();
  for (let commant of commants) {
    const divcom = CreateCommate(commant, islogin);
    div.prepend(divcom);
  }
  return commants.length;
}

function CreateCommate(commant, islogin) {
  const div = document.createElement("div");
  div.className = "commant_post";
  const commantinfo = document.createElement("div");
  commantinfo.className = "commantInfo";
  const h1 = document.createElement("h1");
  h1.className = "Username";
  h1.innerText = commant.userName;
  const date = document.createElement("div");
  date.className = "date";
  date.innerText = commant.date;
  commantinfo.append(h1, date);
  const contan = document.createElement("div");
  contan.innerText = commant.contant;
  div.append(commantinfo, contan);
  const reaction = document.createElement("div");
  reaction.className = "reaction";
  const [like, dislike] = HandulLikeCommat(
    commant.reaction.type,
    commant.reaction.numLike,
    commant.reaction.numDisLike,
    islogin,
    commant.id
  );
  reaction.append(like, dislike);
  div.append(reaction);
  return div;
}

function HandulLikeCommat(type, numlike, numdislike, islogin, commantid) {
  // console.log(postid);
  const likebtn = document.createElement("button");
  const dislikebtn = document.createElement("button");
  likebtn.className = "likecomate";
  dislikebtn.className = "likecomate";
  if (type === "likes") {
    likebtn.classList.add(type);
  } else if (type === "dislikes") {
    dislikebtn.classList.add(type);
  }

  const icondlike = document.createElement("span");
  const icondislike = document.createElement("span");
  icondlike.className = "material-symbols-outlined";
  icondislike.className = "material-symbols-outlined";
  icondislike.style.fontSize = "14px";
  icondlike.style.fontSize = "14px";
  icondislike.innerHTML = "thumb_down";
  icondlike.innerHTML = "thumb_up";
  likebtn.append(icondlike);
  dislikebtn.append(icondislike);
  const likespan = document.createElement("span");
  const dislikespan = document.createElement("span");
  likespan.innerHTML = numlike;
  dislikespan.innerHTML = numdislike;
  likebtn.append(likespan);
  dislikebtn.append(dislikespan);
  let close = true;
  likebtn.addEventListener("click", () => {
    if (!islogin) {
      ClosePop(close);
      return;
    }
    if (likebtn.classList.length == 2) {
      likebtn.classList.remove("likes");
      numlike--;
    } else {
      if (dislikebtn.classList.length == 2) {
        dislikebtn.classList.remove("dislikes");
        numdislike--;
      }
      numlike++;
      likebtn.classList.add("likes");
    }
    likespan.innerHTML = numlike;
    dislikespan.innerHTML = numdislike;
    fetch(`/api/commant/like?commateId=${commantid}&type=likes`, {
      method: "POST",
    });
  });
  dislikebtn.addEventListener("click", () => {
    if (!islogin) {
      ClosePop(close);
      return;
    }
    if (dislikebtn.classList.length == 2) {
      dislikebtn.classList.remove("dislikes");
      numdislike--;
    } else {
      if (likebtn.classList.length == 2) {
        likebtn.classList.remove("likes");
        numlike--;
      }
      numdislike++;
      dislikebtn.classList.add("dislikes");
    }
    likespan.innerHTML = numlike;
    dislikespan.innerHTML = numdislike;
    fetch(`/api/commant/like?commateId=${commantid}&type=dislikes`, {
      method: "POST",
    }).then((res) => {
      // console.log(res);
    });
  });
  return [likebtn, dislikebtn];
}

async function Morecommate(span, offset, postid, islogin, div) {
  const len = await Handulcomment(div, postid, offset, islogin);
  if (len == 0) {
    span.remove();
  }
}
