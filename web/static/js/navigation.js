// === sidebar === //
let colapsed = localStorage.getItem("colapsed") === "true";
const sidebar = document.getElementById("sidebar");

export const colapseSidebar = () => {
  document.documentElement.classList.toggle("is-sidebar-collapsed");
  colapsed = !colapsed;
  localStorage.setItem("colapsed", colapsed);
  sidebar
    .querySelectorAll(".sublinks")
    .forEach((e) => e.classList.remove("sublinks--open"));
};

export function sidebarShowPopover({ ele }) {
  if (colapsed) {
    document.getElementById(ele.dataset.name + "-popover").showPopover();
  }
}

export function sidebarHidePopover({ evt, ele, data }) {
  const sub = document.getElementById(data.name + "-popover");

  if (
    !(ele.contains(evt.toElement) || sub.contains(evt.toElement)) &&
    (colapsed || sub.matches(":popover-open"))
  ) {
    sub.hidePopover();
  }
}

// === dialog x nav sheet === //
const dialogRoot = document.querySelector("#dialog-root");

export function openNavSheet() {
  dialogRoot.classList.add("is-dialog-open");
  dialogRoot
    .querySelector("#nav-sheet")
    .parentElement.classList.remove("hidden");
}

export const closeLastDialog = () => {
  const lc = dialogRoot.lastElementChild.firstElementChild;
  closeDialog(lc);
};

export const closeAllDialog = () => {
  const children = [...dialogRoot.children];
  for (let i = children.length - 1; i >= 0; i--) {
    closeDialog(children[i].firstElementChild);
  }
};

export const closeDialog = (ele) => {
  if (!ele) return;
  dialogRoot.classList.remove("is-dialog-open");
  ele.classList.add("dialog--closing");
  ele.addEventListener(
    "transitionend",
    () => {
      ele.parentElement.classList.add("hidden");
      ele.classList.remove("dialog--closing");
      // if not nav sheet remove from dom
      if (!ele.matches("#nav-sheet")) {
        ele.parentElement.remove();
      }
    },
    { once: true },
  );
};

export function overlayClose({ evt }) {
  if (
    (evt.type === "click" || evt.key === "Escape") &&
    evt.target === dialogRoot?.lastChild
  ) {
    closeLastDialog();
  }
}

// let startY = 0;
// let currentDialog = null;
// let dragging = false;

let ddragState = {
  startY: 0,
  delta: 0,
  currentDialog: null,
  dragging: false,
};

export function startDialogDrag({ ele, evt }) {
  if (ele !== dialogRoot?.lastElementChild) return;
  if (ele.firstElementChild.scrollTop > 0) return;
  if (ddragState.dragging || ddragState.currentDialog) return;
  if (evt.pointerType !== "touch") return;

  ddragState.startY = evt.touches ? evt.touches[0].clientY : evt.clientY;
  ddragState.delta = 0;
  ddragState.currentDialog = ele.firstElementChild;
  ddragState.dragging = true;

  document.addEventListener("pointermove", moveDialogDrag);
  document.addEventListener("pointerup", endDialogDrag);
  document.addEventListener("pointercancel", cancelDialogDrag);
}

function moveDialogDrag(evt) {
  if (!ddragState.currentDialog || !ddragState.dragging) return;

  const clientY = evt.touches ? evt.touches[0].clientY : evt.clientY;

  const deltaY = clientY - ddragState.startY;

  // allow normal scrolling
  if (ddragState.currentDialog.scrollTop > 0) return;

  // only handle pull-down
  if (deltaY <= 0) return;

  evt.preventDefault();
  ddragState.delta = deltaY;
  ddragState.currentDialog.style.transform = `translateY(${deltaY}px)`;
}

function endDialogDrag(evt) {
  if (!ddragState.dragging || !ddragState.currentDialog) return;
  if (ddragState.currentDialog.scrollTop > 0) return;

  const clientY = evt.touches ? evt.touches[0].clientY : evt.clientY;

  const deltaY = clientY - ddragState.startY;
  const shouldClose = deltaY > window.innerHeight * 0.2;

  ddragState.currentDialog.style.transition = "transform 200ms ease";

  if (shouldClose) {
    ddragState.currentDialog
      .closest(".is-dialog-open")
      ?.classList.remove("is-dialog-open");
    ddragState.currentDialog.classList.add("dialog--closing");

    let toClose = ddragState.currentDialog;

    ddragState.currentDialog.addEventListener(
      "transitionend",
      () => {
        toClose.style.transform = "translateY(0)";
        toClose.parentElement.classList.add("hidden");
        toClose.classList.remove("dialog--closing");
        // if not nav sheet remove from dom
        if (!toClose.matches("#nav-sheet")) {
          toClose.parentElement.remove();
        }
      },
      { once: true },
    );
  } else {
    ddragState.currentDialog.style.transform = "";
  }

  ddragState.dragging = false;
  ddragState.currentDialog = null;
}

function cancelDialogDrag() {
  if (!ddragState.dragging || !ddragState.currentDialog) return;
  ddragState.currentDialog.style.transform = "";
  ddragState.dragging = false;
  ddragState.currentDialog = null;
  document.removeEventListener("pointermove", moveDialogDrag);
  document.removeEventListener("pointerup", endDialogDrag);
  document.removeEventListener("pointercancel", cancelDialogDrag);
}

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

export function collapseAllSublinks({ ele }) {
  const allSublinks = ele.closest("nav").querySelectorAll(".sublinks");

  for (let i = 0; i < allSublinks.length; i++) {
    allSublinks[i].classList.remove("sublinks--open");
  }
}

export function toggleSublinks({ ele, data }) {
  const sb = ele.closest("#sidebar");
  const ns = ele.closest("#nav-sheet");

  if (!sb && !ns) {
    return;
  }

  let sub;
  if (
    sb &&
    !document.documentElement.classList.contains("is-sidebar-collapsed")
  ) {
    sub = sb.querySelector(`#${data.name}-sublinks`);
  }
  if (ns) {
    sub = ns.querySelector(`#${data.name}-sublinks`);
  }

  if (!sub) {
    return;
  }

  const isOpen = sub.classList.contains("sublinks--open");
  collapseAllSublinks({ ele });
  sub.classList.toggle("sublinks--open", !isOpen);
  return;
}

export function expandSublinks({ ele, data }) {
  collapseAllSublinks({ ele });
  const sb = ele.closest("#sidebar");
  const ns = ele.closest("#nav-sheet");

  if (!sb && !ns) {
    return;
  }
  if (
    sb &&
    !document.documentElement.classList.contains("is-sidebar-collapsed")
  ) {
    sb.querySelector(`#${data.name}-sublinks`)?.classList.add("sublinks--open");
    return;
  }
  if (ns) {
    ns.querySelector(`#${data.name}-sublinks`)?.classList.add("sublinks--open");
    return;
  }
}

export function collapseSublinks({ ele, data }) {
  const sb = ele.closest("#sidebar");
  const ns = ele.closest("#nav-sheet");

  if (!sb && !ns) {
    return;
  }
  if (
    sb &&
    !document.documentElement.classList.contains("is-sidebar-collapsed")
  ) {
    sb.querySelector(`#${data.name}-sublinks`)?.classList.remove(
      "sublinks--open",
    );
    return;
  }
  if (ns) {
    ns.querySelector(`#${data.name}-sublinks`)?.classList.remove(
      "sublinks--open",
    );
    return;
  }
}

export function handleSublinkFocus({ ele }) {
  const a = document.querySelectorAll(`#${ele.dataset.name}-popover a`);
  for (let i = 0; i < a.length; i++) {
    a[i].classList.remove("focus");
    if (a[i].href === ele.href) {
      a[i].classList.add("focus");
    }
  }
}

export function navigate() {
  closeLastDialog();
  handleActiveLink();
  document.getElementsByTagName("main")[0].focus();
}
