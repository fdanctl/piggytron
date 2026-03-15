// === sidebar === //
// let colapsed set in sidebar-initial
const sidebar = document.getElementById("sidebar");

const handleColapseSidebar = () => {
  document.documentElement.classList.toggle("sidebar-colapsed");
  colapsed = !colapsed;
  localStorage.setItem("colapsed", colapsed);
  sidebar
    .querySelectorAll(".sublinks")
    .forEach((e) => e.classList.remove("open"));
};

const sidebarTogglePopover = (ele) => {
  if (colapsed) {
    document.getElementById(ele.dataset.name + "-popover").togglePopover();
  } else {
    document.getElementById(ele.dataset.name + "-popover").hidePopover();
  }
};

const sidebarShowPopover = (ele) => {
  if (colapsed) {
    document.getElementById(ele.dataset.name + "-popover").showPopover();
  }
};

const sidebarHidePopover = (ele) => {
  const sub = document.getElementById(ele.dataset.name + "-popover");
  if (colapsed || sub.matches(":popover-open")) {
    sub.hidePopover();
  }
};

// === dialog x nav sheet === //
const dialogRoot = document.querySelector("#dialog-root");
const closeLastDialog = () => {
  dialogRoot.lastChild.classList.add("closing");

  setTimeout(() => {
    dialogRoot.lastChild.close();
    dialogRoot.lastChild.classList.remove("closing");
    // if not nav sheet remove from dom
    if (!dialogRoot.lastChild.matches("#nav-sheet")) {
      dialogRoot.lastChild.remove();
    }
  }, 200);
};

dialogRoot?.lastChild.addEventListener("toggle", (e) => {
  if (e.target === dialogRoot.lastChild) {
    dialogRoot.lastChild.focus();
  }
});

document.addEventListener("click", (e) => {
  if (e.target === dialogRoot.lastChild) {
    closeLastDialog();
  }
});

let noStartYPosition = 0;
document.addEventListener("touchstart", (e) => {
  if (e.target === dialogRoot) {
    noStartYPosition = e.touches[0].clientY;
  }
});

document.addEventListener("touchmove", (e) => {
  if (e.target === dialogRoot.lastChild) {
    e.preventDefault();
    const deltaY = e.touches[0].clientY - noStartYPosition;
    if (deltaY > 0) {
      dialogRoot.lastChild.style.transform = `translateY(${deltaY}px)`;
    }
  }
});

document.addEventListener("touchend", (e) => {
  if (e.target === dialogRoot.lastChild) {
    const deltaY = e.changedTouches[0].clientY - noStartYPosition;
    if (deltaY / e.view.outerHeight > 0.2) {
      closeLastDialog();
      setTimeout(() => {
        dialogRoot.lastChild.style.transform = "translateY(0)";
      }, 10);
    } else {
      dialogRoot.lastChild.style.transform = "translateY(0)";
    }
  }
});

// === nav === //
const path = window.location.pathname;
const a = document.querySelectorAll("nav a");
for (let i = 0; i < a.length; i++) {
  const text = a[i].text.trim().toLowerCase();
  a[i].classList.toggle(
    "active",
    path.includes(text) || path === a[i].pathname,
  );
}

const handleActiveLink = (ele) => {
  a.forEach((e) => e.classList.remove("active"));
  const pathname = ele.pathname;
  for (let i = 0; i < a.length; i++) {
    const text = a[i].text.trim().toLowerCase();
    a[i].classList.toggle("active", pathname.includes(text));
  }
  ele.classList.add("active");
};

const handleCloseAllSublinks = (ev) => {
  const allSublinks = ev.target.closest("nav").querySelectorAll(".sublinks");

  for (let i = 0; i < allSublinks.length; i++) {
    allSublinks[i].classList.remove("open");
  }
};
const handleExpandToggle = (ev) => {
  ev.stopPropagation();

  const closestSublinks = ev.target.closest("li").querySelector(".sublinks");

  const isOpen = closestSublinks.matches(".open");

  handleCloseAllSublinks(ev);

  if (!isOpen) {
    if (document.documentElement.classList.contains("sidebar-colapsed")) {
      return;
    }
    closestSublinks.classList.add("open");
  }
};

const handleExpandOpen = (ele) => {
  if (document.documentElement.classList.contains("sidebar-colapsed")) {
    return;
  }
  ele.classList.add("open");
};

const handleExpandRemove = (ele) => {
  ele.classList.remove("open");
};

const handleSublinkFocus = (ele) => {
  const a = document.querySelectorAll(`#${ele.dataset.name}-popover a`);
  for (let i = 0; i < a.length; i++) {
    a[i].classList.remove("focus");
    if (a[i].href === ele.href) {
      a[i].classList.add("focus");
    }
  }
};

function newPage(ev) {
  const splitPath = ev.detail.pathInfo.responsePath.split("/");
  console.log(splitPath);
  let title = splitPath[1];
  if (title === "") {
    title = "Dashboard";
  }
  console.log("title:", title);
  title = title.charAt(0).toUpperCase() + title.slice(1).toLowerCase();
  document.title = title;
  document.getElementsByTagName("main")[0].scrollTo(0, 0);
}
