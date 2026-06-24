import { closeAllDialog, closeDialog, closeLastDialog } from "./navigation";
import { showToast } from "./toast";

export function confirmModal({
  title = "Confirm",
  message = "Are you sure?",
  acceptText = "Yes",
  refuseText = "No",
} = {}) {
  return new Promise((resolve) => {
    const modal = document.createElement("dialog");
    modal.tabIndex = "-1";
    modal.classList.add("dialog", "float");

    modal.innerHTML = `
		<div class="dialog__bar">
			<div></div>
		</div>
		<button class="reset-btn dialog__x" data-action="ui.dialog.close-last">
      <svg xmlns="http://www.w3.org/2000/svg" width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class=""><path d="M18 6 6 18"></path> <path d="m6 6 12 12"></path></svg>
		</button>
    <div>
      <h4 class="mb-md">${title}</h4>
      <p class="text-subtitle">${message}</p>
      <div class="flex gap-xs justify-end items-center mt-sm">
        <button class="btn btn-outline refuse">${refuseText}</button>
        <button class="btn btn-outline accept">${acceptText}</button>
      </div>
`;

    document.getElementById("dialog-root").appendChild(modal);
    modal.showModal();
    modal.focus();

    modal.querySelector(".accept").onclick = () => {
      closeDialog(modal);
      resolve(true);
    };

    modal.querySelector(".refuse").onclick = () => {
      closeDialog(modal);
      resolve(false);
    };
  });
}

document.addEventListener("htmx:confirm", (evt) => {
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

document.body.addEventListener("htmx:historyRestore", (ev) => {
  // nav active link
  let title = document.title;
  const a = document.querySelectorAll("nav a");
  a.forEach((e) => e.classList.remove("active"));
  for (let i = 0; i < a.length; i++) {
    const text = a[i].text.trim().toLowerCase();
    a[i].classList.toggle("active", text === title.toLowerCase());
  }
});

document.body.addEventListener("htmx:sendError", function (ev) {
  showToast("error", "Network error");
});

document.body.addEventListener("htmx:timeout", function (ev) {
  showToast("error", "Request timed out");
});

// htmx custom events
document.body.addEventListener("show-toast", function (ev) {
  showToast(ev.detail.level, ev.detail.message);
});

document.body.addEventListener("incomeCategoryAdded", function (ev) {
  closeLastDialog();
  const li = document.querySelectorAll("#income-cat li");
  document.querySelector("#income-cat h4").innerText =
    `Income (${li.length + 1})`;
});

document.body.addEventListener("expenseCategoryAdded", function (ev) {
  closeLastDialog();
  const li = document.querySelectorAll("#expense-cat li");
  document.querySelector("#expense-cat h4").innerText =
    `Expenses (${li.length + 1})`;
});

document.body.addEventListener("closeModal", function (ev) {
  closeLastDialog();
});

document.body.addEventListener("closeAllModal", function (ev) {
  closeAllDialog();
});

document.body.addEventListener("contentPush", function (ev) {
  htmx.ajax("GET", ev.detail.url, {
    target: "#content",
    swap: `innerHTML transition:${ev.detail.transition ?? "false"}`,
    push: "true",
  });
});

document.body.addEventListener("refetch-transactions", function (ev) {
  const isGoalPage = document.getElementsByClassName("goal-actions").length > 0;
  if (isGoalPage) {
    htmx.ajax("GET", window.location.pathname, {
      target: "#content",
      swap: "innerHTML",
      push: "true",
    });
  }
});

document.body.addEventListener("transaction-deleted", function (ev) {
  closeAllDialog();
  const countEle = document.getElementById("filter-result-count");
  const count = countEle.innerText.match(/^\d*/);

  if (count) {
    countEle.innerText = `${Number(count[0]) - 1} results`;
  }
});
