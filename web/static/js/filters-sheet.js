function filterAccordionToggle({ ele }) {
  const div = ele.parentElement.parentElement.children[1];
  div.classList.toggle("flex-wrap");
  ele.children[0].classList.toggle("hidden");
  ele.children[1].classList.toggle("hidden");
}

function resetTransactionFiltersForm() {
  const form = document.getElementById("transactions-filters");
  const inputs = form.querySelectorAll("input");
  for (let i = 0; i < inputs.length; i++) {
    inputs[i].checked = false;
  }
  inputs[0]?.dispatchEvent(new Event("input", { bubbles: true })); // triggers change event
  const pillBox = document.getElementById("curr-filters");
  pillBox.innerHTML = "";
}

/**
 * @param {HTMLInputElement} input
 */
function toggleFilterPill({ ele }) {
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

function removeFilterPill({ ele }) {
  const form = document.getElementById("transactions-filters");
  const inputs = form.querySelectorAll("input");
  for (let i = 0; i < inputs.length; i++) {
    if (inputs[i].value === ele.dataset.id) {
      inputs[i].checked = false;
      inputs[i]?.dispatchEvent(new Event("input", { bubbles: true })); // triggers change event
    }
  }
}

document.addEventListener("click", (evt) => {
  const ele = evt.target.closest("[data-action]");

  if (!ele) return;

  if (ele.dataset.action === "ui.filters.remove") {
    removeFilterPill({ ele });
  }

  if (ele.dataset.action === "ui.filters-accordion.toggle") {
    filterAccordionToggle({ ele });
  }

  if (ele.dataset.action === "ui.filters.reset") {
    resetTransactionFiltersForm();
  }
});

document.addEventListener("input", (evt) => {
  const ele = evt.target.closest("[data-input]");

  if (!ele) return;

  if (ele.dataset.input === "ui.filters.toggle") {
    toggleFilterPill({ ele });
  }
});
