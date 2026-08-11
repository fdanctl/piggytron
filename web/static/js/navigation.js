// === sidebar === //
/** Whether the sidebar is collapsed; persisted in localStorage. */
let collapsed = localStorage.getItem("collapsed") === "true";
const sidebar = document.getElementById("sidebar");

/**
 * Toggles the collapsed sidebar state and persists it.
 */
export const collapseSidebar = () => {
  document.documentElement.classList.toggle("is-sidebar-collapsed");
  collapsed = !collapsed;
  localStorage.setItem("collapsed", collapsed);
  sidebar
    .querySelectorAll(".sublinks")
    .forEach((e) => e.classList.remove("sublinks--open"));
};

/**
 * Shows the sidebar item popover, only when the sidebar is collapsed.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The nav link.
 */
export function sidebarShowPopover({ ele }) {
  if (collapsed) {
    document.getElementById(ele.dataset.name + "-popover").showPopover();
  }
}

/**
 * Hides the sidebar item popover unless the pointer moved into the link or
 * the popover itself.
 *
 * @param {Object} param0 - Action payload.
 * @param {Event} param0.evt - The mouseout event.
 * @param {Element} param0.ele - The nav link.
 * @param {DOMStringMap} param0.data - Dataset of the link; `data.name`
 *   identifies the popover.
 */
export function sidebarHidePopover({ evt, ele, data }) {
  const sub = document.getElementById(data.name + "-popover");

  if (
    !(ele.contains(evt.toElement) || sub.contains(evt.toElement)) &&
    (collapsed || sub.matches(":popover-open"))
  ) {
    sub.hidePopover();
  }
}

// === dialog x nav sheet === //
const dialogRoot = document.querySelector("#dialog-root");

/**
 * Opens the mobile navigation sheet inside the dialog root.
 */
export function openNavSheet() {
  dialogRoot.classList.add("is-dialog-open");
  dialogRoot
    .querySelector("#nav-sheet")
    .parentElement.classList.remove("hidden");
}

/**
 * Closes the top-most dialog.
 */
export const closeLastDialog = () => {
  const lc = dialogRoot.lastElementChild.firstElementChild;
  closeDialog(lc);
};

/**
 * Closes every open dialog, top to bottom.
 */
export const closeAllDialog = () => {
  const children = [...dialogRoot.children];
  for (let i = children.length - 1; i >= 0; i--) {
    closeDialog(children[i].firstElementChild);
  }
};

/**
 * Closes a single dialog with a closing animation; removes it from the DOM
 * afterwards unless it is the persistent nav sheet.
 *
 * @param {Element} ele - The dialog element to close.
 */
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

/**
 * Closes the top dialog when its overlay is clicked or Escape is pressed.
 *
 * @param {Object} param0 - Action payload.
 * @param {Event} param0.evt - The click or keydown event.
 */
export function overlayClose({ evt }) {
  if (
    (evt.type === "click" || evt.key === "Escape") &&
    evt.target === dialogRoot?.lastChild
  ) {
    closeLastDialog();
  }
}

/** State used for the pull-down-to-close dialog drag */
let ddragState = {
  startY: 0,
  delta: 0,
  currentDialog: null,
  dragging: false,
};

/**
 * Starts tracking a touch drag on the top dialog (pull-down-to-close).
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The dialog overlay.
 * @param {PointerEvent} param0.evt - The pointerdown event.
 */
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

/**
 * Moves the dragged dialog down with the pointer, only on pull-down.
 *
 * @param {PointerEvent} evt - The pointermove event.
 */
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

/**
 * Ends the drag: closes the dialog when pulled past 20% of the viewport
 * height, otherwise snaps it back.
 *
 * @param {PointerEvent} evt - The pointerup event.
 */
function endDialogDrag(evt) {
  if (!ddragState.dragging || !ddragState.currentDialog) return;
  if (ddragState.currentDialog.scrollTop > 0) return;

  const clientY = evt.touches ? evt.touches[0].clientY : evt.clientY;

  const deltaY = clientY - ddragState.startY;
  const shouldClose = deltaY > window.innerHeight * 0.2; // 20% of the viewport

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

  document.removeEventListener("pointermove", moveDialogDrag);
  document.removeEventListener("pointerup", endDialogDrag);
  document.removeEventListener("pointercancel", cancelDialogDrag);
}

/**
 * Handle drag cancels of a in-progress dialog drag and snaps the dialog back.
 */
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

/**
 * Marks the nav link matching the current pathname as active, defaulting
 * to the dashboard link on the root path.
 *
 * @param {Element} [ele] - Optional element to force active.
 */
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

/**
 * Collapses every open sublink group in the nav.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - Any nav element used to locate the nav.
 */
export function collapseAllSublinks({ ele }) {
  const allSublinks = ele.closest("nav").querySelectorAll(".sublinks");

  for (let i = 0; i < allSublinks.length; i++) {
    allSublinks[i].classList.remove("sublinks--open");
  }
}

/**
 * Toggles the sublink group of the clicked nav item (sidebar or nav sheet).
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The clicked nav link.
 * @param {DOMStringMap} param0.data - Dataset of the link; `data.name`
 *   identifies the sublinks group.
 */
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

/**
 * Expands the sublink group of the focused nav item.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The focused nav link.
 * @param {DOMStringMap} param0.data - Dataset of the link; `data.name`
 *   identifies the sublinks group.
 */
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

/**
 * Collapses the sublink group of the blurred nav item.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The blurred nav link.
 * @param {DOMStringMap} param0.data - Dataset of the link; `data.name`
 *   identifies the sublinks group.
 */
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

/**
 * Marks the sublink matching the focused nav item's href inside the item's
 * popover.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The focused nav link.
 */
export function handleSublinkFocus({ ele }) {
  const a = document.querySelectorAll(`#${ele.dataset.name}-popover a`);
  for (let i = 0; i < a.length; i++) {
    a[i].classList.remove("focus");
    if (a[i].href === ele.href) {
      a[i].classList.add("focus");
    }
  }
}

/**
 * Post-navigation hook: closes open dialogs, re-marks the active nav link
 * and moves focus to <main>.
 */
export function navigate() {
  closeLastDialog();
  handleActiveLink();
  document.getElementsByTagName("main")[0].focus();
}
