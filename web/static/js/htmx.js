import { confirmModal } from "./confirmModal";
import { closeAllDialog, closeLastDialog } from "./navigation";
import { getPreferredTheme } from "./theme";
import { showToast } from "./toast";

// Replaces the built-in window.confirm with the custom dialog. The question
// may be a JSON string configuring the dialog; otherwise it is used verbatim.
document.body.addEventListener("htmx:confirm", (evt) => {
  if (!evt.detail.question) return;

  // This will prevent the request from being issued to later manually issue it
  evt.preventDefault();

  let config;

  try {
    config = JSON.parse(evt.detail.question);
  } catch {
    config = {
      title: "Confirm",
      message: evt.detail.question,
      acceptText: "Yes",
      refuseText: "No",
    };
  }

  confirmModal(config).then(function (result) {
    if (result) {
      evt.detail.issueRequest(true); // true to skip the built-in window.confirm()
    }
  });
});

// Re-marks the nav link matching the restored page title.
// NOTE: probably will not be needed in HTMX 4, because it handles
// history diferently
document.body.addEventListener("htmx:historyRestore", () => {
  // nav active link
  let title = document.title;
  const a = document.querySelectorAll("nav a");
  a.forEach((e) => e.classList.remove("active"));
  for (let i = 0; i < a.length; i++) {
    const text = a[i].text.trim().toLowerCase();
    a[i].classList.toggle("active", text === title.toLowerCase());
  }
});

// Toast a generic error unless the server already sent its own trigger.
document.body.addEventListener("htmx:responseError", function (evt) {
  if (!evt.detail.xhr.getResponseHeader("HX-Trigger")) {
    showToast("error", "Something went wrong");
  }
});

document.body.addEventListener("htmx:sendError", function () {
  showToast("error", "Network error");
});

document.body.addEventListener("htmx:timeout", function () {
  showToast("error", "Request timed out");
});

// Names the view transition on the swapped target before navigation.
const DEFAULT_TRANSITION = "navigate-forward";
document.body.addEventListener("htmx:beforeTransition", (evt) => {
  evt.detail.target.style.viewTransitionName =
    evt.target.dataset.transition ?? DEFAULT_TRANSITION;
  // TODO: htmx 4 has a after transtion event. after transition set it to "none"
});

// Sends the effective theme with every request so charts can render dark.
document.body.addEventListener("htmx:configRequest", function (evt) {
  evt.detail.headers["theme"] = getPreferredTheme();
});

// HTMX custom events - set by the server with HX-Trigger header

// Server-emitted toast
// (HX-Trigger:
//   {"show-toast": {
//     "level": "success"|"warning"|"error"|"info",
//     "message": string
//   }}
// ).
document.body.addEventListener("show-toast", function (evt) {
  showToast(evt.detail.level, evt.detail.message);
});

// Bumps the income category counter after one is added.
document.body.addEventListener("incomeCategoryAdded", function () {
  closeLastDialog();
  const li = document.querySelectorAll("#income-cat li");
  document.querySelector("#income-cat h4").innerText =
    `Income (${li.length + 1})`;
});

// Bumps the expense category counter after one is added.
document.body.addEventListener("expenseCategoryAdded", function () {
  closeLastDialog();
  const li = document.querySelectorAll("#expense-cat li");
  document.querySelector("#expense-cat h4").innerText =
    `Expenses (${li.length + 1})`;
});

// Closes the top-most dialog.
document.body.addEventListener("closeModal", function () {
  closeLastDialog();
});

document.body.addEventListener("closeAllModal", function () {
  closeAllDialog();
});

// Navigates making an hx-get to swap #contetn with an optional transition.
// (HX-Trigger: {"contentPush": { "url": string, "transition": bool}})
// TODO: choose the transtion
document.body.addEventListener("contentPush", function (evt) {
  htmx.ajax("GET", evt.detail.url, {
    target: "#content",
    swap: `innerHTML transition:${evt.detail.transition ?? "false"}`,
    push: "true",
  });
});

// Refeshes the current page (or just the ledger list on the ledger page).
document.body.addEventListener("refetch-transactions", function () {
  const isLedgerPage = window.location.pathname.includes("ledger");
  if (!isLedgerPage) {
    htmx.ajax("GET", window.location.pathname, {
      target: "#content",
      swap: "innerHTML",
      push: "true",
    });
  } else {
    document.getElementById("itransactions").innerHTML = "";
    htmx.ajax("GET", "/partials/ledger" + window.location.search, {
      target: "#itransactions",
      swap: "innerHTML",
    });
  }
});

// After a deletion: close dialogs and, refresh on non-ledger pages,
// or decrements the result count.
document.body.addEventListener("transaction-deleted", function () {
  closeAllDialog();
  const isLedgerPage = window.location.pathname.includes("ledger");
  if (!isLedgerPage) {
    htmx.ajax("GET", window.location.pathname, {
      target: "#content",
      swap: "innerHTML",
      push: "true",
    });
  } else {
    const countEle = document.getElementById("filter-result-count");
    const count = countEle.innerText.match(/^\d*/);

    if (count) {
      countEle.innerText = `${Number(count[0]) - 1} results`;
    }
  }
});
