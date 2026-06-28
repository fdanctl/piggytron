import { resetSlider } from "./slider";
import { closeLastDialog } from "./navigation";

export function filterAccordionToggle({ ele }) {
  const div = ele.parentElement.parentElement.children[1];
  div.classList.toggle("flex-wrap");
  ele.children[0].classList.toggle("hidden");
  ele.children[1].classList.toggle("hidden");
}

export function resetTransactionFiltersForm() {
  history.replaceState({}, "", window.location.pathname);
  htmx.ajax("GET", "/partials/transactions", {
    target: "#itransactions",
  });
  const filterBtn = document.getElementById("filter-btn");
  filterBtn.setAttribute("hx-get", "/partials/transaction-filters?");
  filterBtn.querySelector(".notification")?.remove();
  htmx.process(filterBtn);

  closeLastDialog();
}

/**
 * @param {HTMLInputElement} input
 */
export function toggleFilterPill({ ele }) {
  const pillBox = document.getElementById("curr-filters");
  if (ele.checked) {
    const newPill = document.createElement("div");
    newPill.classList.add("pill");
    newPill.dataset.id = ele.value;

    const span = document.createElement("span");
    span.innerText = ele.dataset.label;

    const btn = document.createElement("button");
    btn.classList.add("reset-btn", "flex", "justify-center", "items-center");
    btn.type = "button";
    btn.innerHTML =
      '<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class=""><path d="M18 6 6 18"></path> <path d="m6 6 12 12"></path></svg>';

    newPill.appendChild(span);
    newPill.appendChild(btn);
    pillBox.appendChild(newPill);
    newPill.dataset.action = "ui.filters.remove";
  } else {
    pillBox.querySelector(`[data-id="${ele.value}"]`)?.remove();
  }
}

export function removeFilterPill({ ele }) {
  const form = document.getElementById("transactions-filters");
  const inputs = form.querySelectorAll("input");

  // handle amount-range pill separately
  if (ele.dataset.id === "amount-range") {
    ele.remove();

    // reset slider inputs to empty so they aren't sent as filters
    const minInput = form.querySelector("[name='minamount']");
    const maxInput = form.querySelector("[name='maxamount']");
    if (minInput) minInput.value = "";
    if (maxInput) maxInput.value = "";

    document.getElementById("minamount-chip").innerText =
      minInput.dataset.default;
    document.getElementById("maxamount-chip").innerText =
      maxInput.dataset.default;

    resetSlider(form.querySelector(".slider"));

    maxInput.dispatchEvent(new Event("input", { bubbles: true }));
    return;
  }

  // handle regular filter pills (types, accounts, categories)
  for (let i = 0; i < inputs.length; i++) {
    if (inputs[i].value === ele.dataset.id) {
      inputs[i].checked = false;
      inputs[i]?.dispatchEvent(new Event("input", { bubbles: true })); // triggers change event
    }
  }
}

export function changeMinAmountChip({ ele }) {
  let str = ele.value;
  if (str === "") {
    const slider = ele.closest(".slider");
    const minInput = slider.querySelector("[name='minamount']");
    str = minInput.dataset.default;
  }
  document.getElementById("minamount-chip").innerText = str;
  changeAmountRangeChip(ele);
}

export function changeMaxAmountChip({ ele }) {
  let str = ele.value;
  if (str === "") {
    const slider = ele.closest(".slider");
    const maxInput = slider.querySelector("[name='maxamount']");
    str = maxInput.dataset.default;
  }
  document.getElementById("maxamount-chip").innerText = str;
  changeAmountRangeChip(ele);
}

function changeAmountRangeChip(ele) {
  const slider = ele.closest(".slider");
  const minInput = slider.querySelector("[name='minamount']");
  const maxInput = slider.querySelector("[name='maxamount']");
  const minVal =
    minInput.value !== "" ? minInput.value : minInput.dataset.default;
  const maxVal =
    maxInput.value !== "" ? maxInput.value : maxInput.dataset.default;
  // remove pill if both are the defaults
  if (
    (minVal === "" || minVal === minInput.dataset.default) &&
    (maxVal === "" || maxVal === maxInput.dataset.default)
  ) {
    document.querySelector(".pill[data-id='amount-range']")?.remove();
    return;
  }

  const pillBox = document.getElementById("curr-filters");
  const newPill = document.createElement("div");
  let pill = document.querySelector("div[data-id='amount-range']");
  if (!pill) {
    newPill.classList.add("pill");
    newPill.dataset.id = "amount-range";

    const span = document.createElement("span");
    span.innerText = document
      .getElementById("minamount-chip")
      .closest("div").innerText;

    const btn = document.createElement("button");
    btn.classList.add("reset-btn", "flex", "justify-center", "items-center");
    btn.type = "button";
    btn.innerHTML =
      '<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class=""><path d="M18 6 6 18"></path> <path d="m6 6 12 12"></path></svg>';

    newPill.appendChild(span);
    newPill.appendChild(btn);

    pillBox.appendChild(newPill);
    newPill.dataset.action = "ui.filters.remove";
  } else {
    const span = pill.firstElementChild;

    span.innerText = document
      .getElementById("minamount-chip")
      .closest("div").innerText;
  }
}
