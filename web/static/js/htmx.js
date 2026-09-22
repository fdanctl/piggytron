import { confirmModal } from "./confirmModal";
import { closeAllDialog, closeLastDialog } from "./navigation";
import { getPreferredTheme } from "./theme";
import { showToast } from "./toast";

// Replaces the built-in window.confirm with the custom dialog. The question
// may be a JSON string configuring the dialog; otherwise it is used verbatim.
htmx.on("htmx:confirm", (evt) => {
  if (!evt.detail.ctx.confirm) return;

  // This will prevent the request from being issued to later manually issue it
  evt.preventDefault();

  let config;

  try {
    config = JSON.parse(evt.detail.ctx.confirm);
  } catch {
    config = {
      title: "Confirm",
      message: evt.detail.ctx.confirm,
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
htmx.on("htmx:before:history:restore", (evt) => {
  // nav active link
  let pathname = evt.detail.path;
  if (pathname === "/") {
    pathname = "dashboard";
  }
  const a = document.querySelectorAll("nav a");
  a.forEach((e) => e.classList.remove("active"));
  for (let i = 0; i < a.length; i++) {
    const text = a[i].text.trim().toLowerCase();
    a[i].classList.toggle("active", pathname.includes(text));
  }
});

// Toast a generic error unless the server already sent its own trigger.
htmx.on("htmx:response:error", (evt) => {
  if (
    !evt.detail.ctx.response.headers.get("HX-Trigger") &&
    evt.detail.ctx.response.status !== 422
  ) {
    showToast("error", "Something went wrong");
  }
});

htmx.on("htmx:error", () => {
  showToast("error", "Network error");
});

function waitForCharts() {
  return new Promise((resolve) => {
    const interval = setInterval(() => {
      if (window.chartsLoaded === true) {
        clearInterval(interval);
        resolve();
      }
    }, 100);
  });
}

htmx.on("htmx:before:swap", async (evt) => {
  if (
    (evt.detail.ctx.request.action === undefined &&
      !evt.detail.ctx.request.action.includes("/partials/charts")) ||
    window.chartsLoaded
  ) {
    return;
  }

  evt.preventDefault();

  await waitForCharts();
  htmx.swap(evt.detail.ctx);
});

// scroll to the top when changing page
htmx.on("htmx:after:swap", (evt) => {
  if (evt.detail.ctx.target.id !== "content") return;

  document.querySelector("main").scrollTo({
    top: 0,
    behavior: "instant",
  });

  window.scrollTo({
    top: 0,
    behavior: "instant",
  });
});

// Names the view transition on the swapped target before navigation.
const DEFAULT_TRANSITION = "navigate-forward";
htmx.on("htmx:before:viewTransition", (evt) => {
  evt.detail.ctx.target.style.viewTransitionName =
    evt.target.dataset.transition ?? DEFAULT_TRANSITION;
});

htmx.on("htmx:after:viewTransition", (evt) => {
  evt.detail.ctx.target.style.viewTransitionName = "none";
});

// Sends the effective theme with every request so charts can render dark.
htmx.on("htmx:config:request", (evt) => {
  evt.detail.ctx.request.headers["theme"] = getPreferredTheme();
});

// HTMX custom events - set by the server with HX-Trigger header

// Server-emitted toast
// (HX-Trigger:
//   {"show-toast": {
//     "level": "success"|"warning"|"error"|"info",
//     "message": string
//   }}
// ).
document.body.addEventListener("show-toast", (evt) => {
  showToast(evt.detail.level, evt.detail.message);
});

// Bumps the income category counter after one is added.
document.body.addEventListener("incomeCategoryAdded", () => {
  closeLastDialog();
  const li = document.querySelectorAll("#income-cat li");
  document.querySelector("#income-cat h4").innerText =
    `Income (${li.length + 1})`;
});

// Bumps the expense category counter after one is added.
document.body.addEventListener("expenseCategoryAdded", () => {
  closeLastDialog();
  const li = document.querySelectorAll("#expense-cat li");
  document.querySelector("#expense-cat h4").innerText =
    `Expenses (${li.length + 1})`;
});

// Closes the top-most dialog.
document.body.addEventListener("closeModal", () => {
  closeLastDialog();
});

document.body.addEventListener("closeAllModal", () => {
  closeAllDialog();
});

// Navigates making an hx-get to swap #contetn with an optional transition.
// (HX-Trigger: {"contentPush": { "url": string, "transition": bool}})
// TODO: choose the transtion
document.body.addEventListener("contentPush", (evt) => {
  htmx.ajax("GET", evt.detail.url, {
    target: "#content",
    swap: `innerHTML transition:${evt.detail.transition ?? "false"}`,
    push: "true",
  });
});

// Refeshes the current page (or just the ledger list on the ledger page).
document.body.addEventListener("refetch-transactions", () => {
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
document.body.addEventListener("transaction-deleted", () => {
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
