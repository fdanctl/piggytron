import { resetSlider } from "./slider";
import { closeLastDialog } from "./navigation";
import { formatDate } from "./utils";
import { makeSVG } from "./icons";

function addPill(id, label) {
  const pillBox = document.getElementById("curr-filters");
  const newPill = document.createElement("div");
  newPill.classList.add("pill");
  newPill.dataset.id = id;

  const span = document.createElement("span");
  span.innerText = label;

  const btn = document.createElement("button");
  btn.classList.add("reset-btn", "flex", "justify-center", "items-center");
  btn.type = "button";
  btn.innerHTML = makeSVG("x", 14);

  newPill.appendChild(span);
  newPill.appendChild(btn);
  pillBox.appendChild(newPill);
  newPill.dataset.action = "ui.filters.remove";
}

export function filterAccordionToggle({ ele }) {
  const div = ele.parentElement.parentElement.children[1];
  div.classList.toggle("flex-wrap");
  ele.children[0].classList.toggle("hidden");
  ele.children[1].classList.toggle("hidden");
}

export function resetTransactionFiltersForm() {
  history.replaceState({}, "", window.location.pathname);
  htmx.ajax("GET", "/partials/ledger", {
    target: "#itransactions",
  });
  const filterBtn = document.getElementById("filter-btn");
  filterBtn.setAttribute("hx-get", "/partials/ledger-filters?");
  filterBtn.querySelector(".notification")?.remove();
  htmx.process(filterBtn);

  closeLastDialog();
}

/**
 * @param {HTMLInputElement} input
 */
export function toggleFilterPill({ ele }) {
  if (ele.checked) {
    addPill(ele.value, ele.dataset.label);
  } else {
    const pillBox = document.getElementById("curr-filters");
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

    resetSlider(minInput.closest(".slider"));

    const urlParams = new URLSearchParams(window.location.search);

    if (urlParams.has("minamount")) {
      minInput.dispatchEvent(new Event("input", { bubbles: true }));
    } else {
      maxInput.dispatchEvent(new Event("input", { bubbles: true }));
    }

    return;
  }

  // handle date-range pill separately
  if (ele.dataset.id === "date-range") {
    ele.remove();

    // reset slider inputs to empty so they aren't sent as filters
    const minInput = form.querySelector("[name='mindate']");
    const maxInput = form.querySelector("[name='maxdate']");
    if (minInput) minInput.value = "";
    if (maxInput) maxInput.value = "";

    document.getElementById("mindate-chip").innerText = formatDate(
      new Date(Number(minInput.dataset.default) * 1000),
    );
    document.getElementById("maxdate-chip").innerText = formatDate(
      new Date(Number(maxInput.dataset.default) * 1000),
    );

    resetSlider(minInput.closest(".slider"));

    const urlParams = new URLSearchParams(window.location.search);

    if (urlParams.has("mindate")) {
      minInput.dispatchEvent(new Event("input", { bubbles: true }));
    } else {
      maxInput.dispatchEvent(new Event("input", { bubbles: true }));
    }

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
  changeRangePill(ele, "amount");
}

export function changeMaxAmountChip({ ele }) {
  let str = ele.value;
  if (str === "") {
    const slider = ele.closest(".slider");
    const maxInput = slider.querySelector("[name='maxamount']");
    str = maxInput.dataset.default;
  }
  document.getElementById("maxamount-chip").innerText = str;
  changeRangePill(ele, "amount");
}

export function changeMinDateChip({ ele }) {
  let str = ele.value;
  if (str === "") {
    const slider = ele.closest(".slider");
    const minInput = slider.querySelector("[name='mindate']");
    str = minInput.dataset.default;
  }
  document.getElementById("mindate-chip").innerText = formatDate(
    new Date(Number(str) * 1000),
  );
  changeRangePill(ele, "date");
}

export function changeMaxDateChip({ ele }) {
  let str = ele.value;
  if (str === "") {
    const slider = ele.closest(".slider");
    const maxInput = slider.querySelector("[name='maxdate']");
    str = maxInput.dataset.default;
  }
  document.getElementById("maxdate-chip").innerText = formatDate(
    new Date(Number(str) * 1000),
  );
  changeRangePill(ele);
}

function changeRangePill(ele, attr) {
  const slider = ele.closest(".slider");
  const minInput = slider.querySelector(`[name='min${attr}']`);
  const maxInput = slider.querySelector(`[name='max${attr}']`);
  const minVal = minInput.value || minInput.dataset.default;
  const maxVal = maxInput.value || maxInput.dataset.default;
  // remove pill if both are the defaults
  if (
    minVal === minInput.dataset.default &&
    maxVal === maxInput.dataset.default
  ) {
    document.querySelector(`.pill[data-id='${attr}-range']`)?.remove();
    return;
  }

  const pill = document.querySelector(`div[data-id='${attr}-range']`);
  if (!pill) {
    addPill(
      `${attr}-range`,
      document.getElementById(`min${attr}-chip`).closest("div").innerText,
    );
  } else {
    const span = pill.firstElementChild;

    span.innerText = document
      .getElementById(`min${attr}-chip`)
      .closest("div").innerText;
  }
}
