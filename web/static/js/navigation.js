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
  const lc = dialogRoot.lastElementChild;
  closeDialog(lc);
};

export const closeAllDialog = () => {
  const children = [...dialogRoot.children];
  for (let i = children.length - 1; i >= 0; i--) {
    closeDialog(children[i]);
  }
};

export const closeDialog = (ele) => {
  ele.classList.add("closing");
  ele.addEventListener(
    "transitionend",
    () => {
      ele.close();
      ele.classList.remove("closing");
      // if not nav sheet remove from dom
      if (!ele.matches("#nav-sheet")) {
        ele.remove();
      }
    },
    { once: true },
  );
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

let startY = 0;
let currentDialog = null;
let dragging = false;

document.addEventListener("touchstart", (e) => {
  const dialog = e.target.closest(".dialog");

  if (dialog !== dialogRoot?.lastElementChild) return;
  if (dialog.scrollTop > 0) return;

  currentDialog = dialog;
  startY = e.touches[0].clientY;
  dragging = true;
});

document.addEventListener(
  "touchmove",
  (e) => {
    if (!dragging || !currentDialog) return;

    const deltaY = e.touches[0].clientY - startY;

    // allow normal scrolling
    if (currentDialog.scrollTop > 0) return;

    // only handle pull-down
    if (deltaY <= 0) return;

    e.preventDefault();

    currentDialog.style.transform = `translateY(${deltaY}px)`;
  },
  { passive: false },
);

document.addEventListener("touchend", (e) => {
  if (!dragging || !currentDialog) return;
  if (currentDialog.scrollTop > 0) return;

  const deltaY = e.changedTouches[0].clientY - startY;
  const shouldClose = deltaY > window.innerHeight * 0.2;

  currentDialog.style.transition = "transform 200ms ease";

  if (shouldClose) {
    currentDialog.classList.add("closing");

    let toClose = currentDialog;

    currentDialog.addEventListener(
      "transitionend",
      () => {
        toClose.style.transform = "translateY(0)";
        toClose.close();
        toClose.classList.remove("closing");
        // if not nav sheet remove from dom
        if (!toClose.matches("#nav-sheet")) {
          toClose.remove();
        }
      },
      { once: true },
    );
  } else {
    currentDialog.style.transform = "translateY(0)";
  }

  dragging = false;
  currentDialog = null;
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
      handleExpandOpen(nextSibling);
    } else if (link.closest(".sublinks")) {
      handleSublinkFocus(link);
      handleExpandOpen(link.closest(".sublinks"));
    }
    sidebarShowPopover(link);
  });

  link.addEventListener("click", (evt) => {
    const nextSibling = link.nextElementSibling;
    if (nextSibling?.classList.contains("sublinks")) {
      handleExpandToggle(nextSibling);
      sidebarShowPopover(link);
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
