// === sidebar === //
// let colapsed = localStorage.getItem("colapsed") === "true";
// const sidebar = document.getElementById("sidebar");
// sidebar.classList.toggle("colapsed", colapsed);

const handleColapseSidebar = () => {
  sidebar.classList.toggle("colapsed");
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

// === nav sheet === //
const navSheet = document.querySelector("#nav-sheet");
const closeNavSheet = () => {
  navSheet.classList.add("closing");

  setTimeout(() => {
    navSheet.close();
    navSheet.classList.remove("closing");
  }, 200);
};

navSheet?.addEventListener("toggle", (e) => {
  if (e.target === navSheet) {
    navSheet.focus();
  }
});

document.addEventListener("click", (e) => {
  if (e.target === navSheet) {
    closeNavSheet();
  }
});

let noStartYPosition = 0;
document.addEventListener("touchstart", (e) => {
  if (e.target === navSheet) {
    noStartYPosition = e.touches[0].clientY;
  }
});

document.addEventListener("touchmove", (e) => {
  if (e.target === navSheet) {
    e.preventDefault();
    const deltaY = e.touches[0].clientY - noStartYPosition;
    if (deltaY > 0) {
      navSheet.style.transform = `translateY(${deltaY}px)`;
    }
  }
});

document.addEventListener("touchend", (e) => {
  if (e.target === navSheet) {
    const deltaY = e.changedTouches[0].clientY - noStartYPosition;
    if (deltaY / e.view.outerHeight > 0.2) {
      closeNavSheet();
      setTimeout(() => {
        navSheet.style.transform = "translateY(0)";
      }, 10);
    } else {
      navSheet.style.transform = "translateY(0)";
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

  handleCloseAllSublinks(event);

  if (!isOpen) {
    if (closestSublinks.closest("#sidebar")?.matches(".colapsed")) {
      return;
    }
    closestSublinks.classList.add("open");
  }
};

const handleExpandOpen = (ele) => {
  if (ele.closest("#sidebar")?.matches(".colapsed")) {
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

function newTitle(ev) {
  const splitPath = ev.detail.pathInfo.responsePath.split("/");
  let title = splitPath[splitPath.length - 1];
  if (title === "") {
    title = "Dashboard";
  }
  title = title.charAt(0).toUpperCase() + title.slice(1).toLowerCase();
  document.title = title;
}
