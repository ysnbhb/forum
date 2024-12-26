export function CreateDiv(post) {
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
  // console.log(post.id);
  const [like, dislike] = HandulLike(
    post.reaction.type,
    post.reaction.numLike,
    post.reaction.numDisLike,
    null,
    post.id
  );
  reaction.append(like, dislike);

  const commant = document.createElement("button");
  commant.className = "commant";
  commant.innerHTML = `
      <span class="material-symbols-outlined"> comment </span>
`;
  btndiv.append(reaction, commant);
  articles.append(btndiv);
  return articles;
}

async function FetchPost(offset) {
  const allpost = document.getElementById("allPost");
  try {
    const res = await fetch(`/api/posts?offset=${offset}`);
    const posts = await res.json();
    if (posts.length === 0 && offset == 0) {
      document.getElementById("content").style.display = "none";
      return;
    }
    for (let i = 0; i < posts.length; i++) {
      const post = CreateDiv(posts[i]);
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
  // console.log(postid);
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
  likebtn.addEventListener("click", () => {
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

FetchPost(0);

function Inf() {
  let offset = 20;
  let lenghtpost;
  window.addEventListener("scrollend", async () => {
    let windowHight = window.innerHeight;
    let scrol = window.scrollY;
    if (scrol + windowHight > document.body.scrollHeight - 1000) {
      if (lenghtpost === 0) {
        window.removeEventListener("scroll", Inf);
        return;
      }
      lenghtpost = await FetchPost(offset);
      offset += 20;
      console.log(offset);
    }
  });
}

Inf();
