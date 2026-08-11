import { confirmModal } from "./confirmModal";
import { closeAllDialog, closeLastDialog } from "./navigation";
import { getpreferredTheme } from "./theme";
import { showToast } from "./toast";

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

const DEFAULT_TRANSITION = "navigate-forward";
document.body.addEventListener("htmx:beforeTransition", (evt) => {
  evt.detail.target.style.viewTransitionName =
    evt.target.dataset.transition ?? DEFAULT_TRANSITION;
});

document.body.addEventListener("htmx:configRequest", function (evt) {
  // for the charts
  evt.detail.headers["theme"] = getpreferredTheme();
});

// htmx custom events - set by the server with HX-Trigger header

document.body.addEventListener("show-toast", function (evt) {
  showToast(evt.detail.level, evt.detail.message);
});

document.body.addEventListener("incomeCategoryAdded", function () {
  closeLastDialog();
  const li = document.querySelectorAll("#income-cat li");
  document.querySelector("#income-cat h4").innerText =
    `Income (${li.length + 1})`;
});

document.body.addEventListener("expenseCategoryAdded", function () {
  closeLastDialog();
  const li = document.querySelectorAll("#expense-cat li");
  document.querySelector("#expense-cat h4").innerText =
    `Expenses (${li.length + 1})`;
});

document.body.addEventListener("closeModal", function () {
  closeLastDialog();
});

document.body.addEventListener("closeAllModal", function () {
  closeAllDialog();
});

document.body.addEventListener("contentPush", function (evt) {
  htmx.ajax("GET", evt.detail.url, {
    target: "#content",
    swap: `innerHTML transition:${evt.detail.transition ?? "false"}`,
    push: "true",
  });
});

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
