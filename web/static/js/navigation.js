// === sidebar === //
let colapsed = localStorage.getItem("colapsed") === "true";
const sidebar = document.getElementById("sidebar");

export const colapseSidebar = () => {
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

export function openNavSheet() {
  dialogRoot.querySelector("#nav-sheet").showModal();
}

export const closeLastDialog = () => {
  const lc = dialogRoot.lastChild;
  lc.classList.add("closing");

  setTimeout(() => {
    lc.close();
    lc.classList.remove("closing");
    // if not nav sheet remove from dom
    if (!lc.matches("#nav-sheet")) {
      lc.remove();
    }
  }, 200);
};

dialogRoot?.lastChild.addEventListener("toggle", (e) => {
  if (e.target === dialogRoot.lastChild) {
    dialogRoot.lastChild.focus();
  }
});

document.addEventListener("click", (e) => {
  if (e.target === dialogRoot?.lastChild) {
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
  if (e.target === dialogRoot?.lastChild) {
    e.preventDefault();
    const deltaY = e.touches[0].clientY - noStartYPosition;
    if (deltaY > 0) {
      dialogRoot.lastChild.style.transform = `translateY(${deltaY}px)`;
    }
  }
});

document.addEventListener("touchend", (e) => {
  if (e.target === dialogRoot?.lastChild) {
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
  let pathname = window.location.pathname;
  if (pathname === "/") {
    pathname = "dashboard";
  }
  for (let i = 0; i < a.length; i++) {
    const text = a[i].text.trim().toLowerCase();
    a[i].classList.toggle("active", pathname.includes(text));
  }
  ele?.classList.add("active");
};

const handleCloseAllSublinks = (ev) => {
  const allSublinks = ev.target.closest("nav").querySelectorAll(".sublinks");

  for (let i = 0; i < allSublinks.length; i++) {
    allSublinks[i].classList.remove("open");
  }
};

const handleExpandToggle = (ele) => {
  if (
    document.documentElement.classList.contains("sidebar-colapsed") &&
    ele.closest("#sidebar")
  ) {
    return;
  }
  ele.classList.toggle("open");
};

const handleExpandOpen = (ele) => {
  if (
    document.documentElement.classList.contains("sidebar-colapsed") &&
    ele.closest("#sidebar")
  ) {
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

const sidebarLinks = document.querySelectorAll("#sidebar [data-name]");

for (let i = 0; i < sidebarLinks.length; i++) {
  sidebarLinks[i];
}

for (let i = 0; i < sidebarLinks.length; i++) {
  const link = sidebarLinks[i];
  link.addEventListener("mouseenter", (evt) => {
    sidebarShowPopover(link);
  });

  link.addEventListener("mouseleave", (evt) => {
    sidebarHidePopover(link);
  });

  link.addEventListener("blur", (evt) => {
    sidebarHidePopover(link);
    if (link.closest(".sublinks") && !link.nextElementSibling) {
      handleExpandRemove(link.closest(".sublinks"));
    }
  });

  link.addEventListener("mousedown", (evt) => {
    const nextSibling = link.nextElementSibling;
    if (nextSibling?.classList.contains("sublinks")) {
      evt.preventDefault();
      handleExpandOpen(link);
    }
  });

  link.addEventListener("focus", (evt) => {
    const nextSibling = link.nextElementSibling;
    if (nextSibling?.classList.contains("sublinks")) {
      // sidebarShowPopover(link);
      handleExpandOpen(nextSibling);
    } else if (link.closest(".sublinks")) {
      handleExpandOpen(link.closest(".sublinks"));
    }
  });

  link.addEventListener("click", (evt) => {
    const nextSibling = link.nextElementSibling;
    if (nextSibling?.classList.contains("sublinks")) {
      handleExpandToggle(nextSibling);
    } else if (link.closest(".sublinks")) {
      handleExpandRemove(link.closest(".sublinks"));
    } else {
      evt.stopPropagation();
    }
    const popover = link.closest(".popover");
    if (popover) {
      sidebarHidePopover(popover);
    }
  });
}

export function navigate() {
  closeLastDialog();
  handleActiveLink();
  document.getElementsByTagName("main")[0].focus();
}

const navSheetLinks = document.querySelectorAll("#nav-sheet [data-name]");

for (let i = 0; i < navSheetLinks.length; i++) {
  const link = navSheetLinks[i];

  link.addEventListener("click", (evt) => {
    const nextSibling = link.nextElementSibling;
    if (nextSibling?.classList.contains("sublinks")) {
      handleExpandToggle(nextSibling);
    } else if (link.closest(".sublinks")) {
      handleExpandRemove(link.closest(".sublinks"));
      handleActiveLink(link);
    } else {
      evt.stopPropagation();
      handleActiveLink(link);
      handleCloseAllSublinks(evt);
    }
    const popover = link.closest(".popover");
    if (popover) {
      sidebarHidePopover(popover);
    }
  });
}
