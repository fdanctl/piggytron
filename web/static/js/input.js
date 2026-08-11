import { getTarget } from "./actions/ui";
import { buildCalendar, clickOption } from "./calendar";

/**
 * Verifies, on blur, if the input has error and adds error state to it
 * and its input-group.
 *
 * @param {HTMLInputElement} ele - The input being blurred.
 */
export function handleInputOnBlur(ele) {
  ele.classList.toggle("input--error", !ele.validity.valid);
  const parent = ele.parentElement.parentElement;
  if (parent.classList.contains("input-group")) {
    parent.classList.toggle("input-group--error", !ele.validity.valid);
  }
}

/**
 * Show/hide password. Changes between password input and text
 * input and icons accordingly.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The toggle button.
 */
export function passwordToggle({ ele }) {
  const pwdInput = ele.parentElement.parentElement.children[0];
  if (pwdInput.type === "password") {
    pwdInput.type = "text";
  } else if (pwdInput.type === "text") {
    pwdInput.type = "password";
  }
  ele.children[0].classList.toggle("hidden");
  ele.children[1].classList.toggle("hidden");
}

/**
 * Toggles the checkbox of a clickable pill and fires its input event.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The pill element.
 */
export function checkboxPillToggle({ ele }) {
  let cb = ele.querySelector("input");
  cb.checked = !cb.checked;
  cb.dispatchEvent(new Event("input", { bubbles: true })); // triggers change event
}

/**
 * Selects a custom dropdown option: writes the value into the hidden input,
 * marks the option and syncs the trigger button label.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The clicked option.
 * @param {DOMStringMap} param0.data - Dataset of the option; `data.value` is
 *   the selected value.
 */
export function selectSelect({ ele, data }) {
  const input = ele.parentElement.nextElementSibling;
  input.value = data.value;
  input.dispatchEvent(new Event("change", { bubbles: true })); // triggers change event

  const opts = ele.parentElement.querySelectorAll("li");
  for (let i = 0; i < opts.length; i++) {
    opts[i].classList.remove("options__item--selected");
  }
  ele.classList.add("options__item--selected");
  ele.closest(".popover").hidePopover();
  let drop = ele.closest(".dropdown");
  if (drop) {
    drop.querySelector("button > span").innerText = ele.firstChild.innerText;
    drop.querySelector("button").classList.remove("input--error");
    ele.closest(".input-group")?.classList.remove("input-group--error");
  }
}

/**
 * Selects a pill-style option, mirroring selectSelect for the
 * .select-pill variant.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The clicked option.
 * @param {DOMStringMap} param0.data - Dataset of the option; `data.value` is
 *   the selected value.
 */
export function selectPillSelect({ ele, data }) {
  ele.closest(".popover").hidePopover();
  const input = ele.parentElement.nextElementSibling;
  input.value = data.value;
  input.dispatchEvent(new Event("change", { bubbles: true })); // triggers change event

  const opts = ele.parentElement.querySelectorAll("li");
  for (let i = 0; i < opts.length; i++) {
    opts[i].classList.remove("select-pill__option--selected");
  }
  ele.classList.add("select-pill__option--selected");
  let drop = ele.closest(".select.select-pill");
  if (drop) {
    drop.querySelector("button").innerText = ele.innerText;
  }
}

/**
 * Scrolls the selected option into view inside an open popover, on the next
 * frame.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The triggering element.
 * @param {DOMStringMap} param0.data - Dataset of the trigger; `data.target`
 *   is the id of the popover element.
 */
export function centerSelected({ ele, data }) {
  const popover = getTarget(ele, data.target);

  requestAnimationFrame(() => {
    if (popover.matches(":popover-open")) {
      const selected = popover.querySelector(".options__item--selected");

      selected?.scrollIntoView({
        block: "center",
      });
    }
  });
}

/**
 * Strips non-numeric characters from a cash input while typing and clamps
 * the value to two decimal places.
 *
 * @param {Object} param0 - Action payload.
 * @param {HTMLInputElement} param0.ele - The cash input.
 */
export function sanitizeCashInput({ ele }) {
  let value = ele.value.replace(/[^0-9.]/g, "");
  const parts = value.split(".");
  if (parts.length > 2) {
    value = parts[0] + "." + parts.slice(1).join("");
  } else if (parts.length === 2 && parts[1].length > 2) {
    value = parts[0] + "." + parts[1].slice(0, 2);
  }

  ele.value = value;
}

/**
 * Formats the cash input on blur: pads the decimals to two digits and adds
 * thousands separators to the integer part.
 *
 * @param {Object} param0 - Action payload.
 * @param {HTMLInputElement} param0.ele - The cash input.
 */
export function cashInputBlur({ ele }) {
  let value = ele.value.replace(/[^0-9.]/g, "");

  let parts = value.split(".");
  let intPart = parts[0] || "0";
  let decimalPart = parts[1] || "";

  if (decimalPart !== "") {
    while (decimalPart.length < 2) {
      decimalPart += "0";
    }
    decimalPart = decimalPart.slice(0, 2);
  }

  // TODO: locale config
  intPart = parseInt(intPart || "0", 10).toLocaleString("en-US");

  ele.value = `${intPart}${decimalPart != "" ? "." + decimalPart : ""}`;
}

/**
 * Strips thousands separators on focus and places the caret at the end, so
 * the user can keep typing.
 *
 * @param {Object} param0 - Action payload.
 * @param {HTMLInputElement} param0.ele - The cash input.
 */
export function cashInputFocus({ ele }) {
  ele.value = ele.value.replaceAll(",", "");
  const length = ele.value.length;
  ele.setSelectionRange(length, length);
}

/**
 * Filters budget-amount keystrokes: allows digits, dot and control keys,
 * and moves between amount inputs with j/k or the arrow keys.
 *
 * @param {Object} param0 - Action payload.
 * @param {KeyboardEvent} param0.evt - The keydown event.
 */
export function budgetInput({ evt }) {
  // allow text shortcuts and reload page
  if (
    (evt.ctrlKey || evt.metaKey) &&
    ["c", "v", "x", "a", "r"].includes(evt.key.toLowerCase())
  ) {
    return;
  }

  const allowedKeys = [
    "Backspace",
    "Tab",
    "ArrowLeft",
    "ArrowRight",
    "Delete",
    "Enter",
  ];

  if (allowedKeys.includes(evt.key) || /^[0-9.]$/.test(evt.key)) return;

  const key = evt.code;

  if (key === "KeyJ" || key === "ArrowDown") {
    nextInput();
  } else if (key === "KeyK" || key === "ArrowUp") {
    prevInput();
  }

  // avoid making a request if not a number or dot
  evt.preventDefault();
}

/**
 * Plays the "scan" highlight animation on the input.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The amount input.
 */
export function budgetInputScan({ ele }) {
  ele.classList.add("scan");
  ele.addEventListener(
    "animationend",
    () => {
      ele.classList.remove("scan");
    },
    { once: true },
  );
}

/**
 * @returns {HTMLInputElement[]} The amount inputs.
 */
function getInputs() {
  return [...document.querySelectorAll("input[name='amount']")];
}

function nextInput() {
  const items = getInputs();
  const current = document.activeElement;

  const idx = items.indexOf(current);

  if (idx !== -1 && idx < items.length - 1) {
    items[idx + 1].focus();
  }
}

function prevInput() {
  const items = getInputs();
  const current = document.activeElement;

  const idx = items.indexOf(current);

  if (idx > 0) {
    items[idx - 1].focus();
  }
}

/**
 * Masks a date input as DD/MM/YYYY while typing, clamping day to 01–31 and
 * month to 01–12.
 *
 * @param {Object} param0 - Action payload.
 * @param {HTMLInputElement} param0.ele - The date input.
 */
export function dateOnChange({ ele }) {
  let raw = ele.value.replace(/\D/g, "");

  // Limit to 8 digits (DDMMYYYY)
  raw = raw.slice(0, 8);

  let day = raw.slice(0, 2);
  let month = raw.slice(2, 4);
  let year = raw.slice(4, 8);

  // Clamp day (01–31)
  if (day.length === 2) {
    let d = Math.min(Math.max(parseInt(day, 10), 1), 31);
    day = d.toString().padStart(2, "0");
  }

  // Clamp month (01–12)
  if (month.length === 2) {
    let m = Math.min(Math.max(parseInt(month, 10), 1), 12);
    month = m.toString().padStart(2, "0");
  }

  let formatted = day;

  if (raw.length > 2) {
    formatted += "/" + month;
  }

  if (raw.length > 4) {
    formatted += "/" + year;
  }

  ele.value = formatted;
}

/**
 * Opens the calendar popover, rebuilding its grid and preselecting the
 * value (or today) of the linked date input.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The triggering element.
 * @param {DOMStringMap} param0.data - Dataset of the trigger; `data.target`
 *   is the id of the calendar popover.
 */
export function openCalendar({ ele, data }) {
  const target = getTarget(ele, data.target);
  buildCalendar({ ele: target.querySelector(".calendar") });
  const inputValue = target.previousSibling.querySelector("input").value;
  let day;
  let month;
  let year;
  if (inputValue.length < 10) {
    const presentDay = new Date(Date.now());
    day = presentDay.getDate();
    month = presentDay.getMonth();
    year = presentDay.getFullYear();
  } else {
    const ddmmyyyy = inputValue.split("/");
    day = Number(ddmmyyyy[0]);
    month = Number(ddmmyyyy[1]) - 1;
    year = Number(ddmmyyyy[2]);
  }

  const yearInput = target.querySelector("input[name='year']");
  const monthInput = target.querySelector("input[name='month']");
  clickOption(yearInput, String(year));
  clickOption(monthInput, String(month));
  monthInput.dispatchEvent(new Event("change"), { bubbles: true }); // triggers change event

  const days = target.querySelectorAll(".calendar__days > div");
  for (let i = 0; i < days.length; i++) {
    if (day == Number(days[i].innerText)) {
      days[i].classList.add("calendar__day--selected");
    }
  }
}

/**
 * Writes the clicked calendar day into the linked date input (DD/MM/YYYY)
 * and closes the popover.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The calendar popover.
 * @param {Event} param0.evt - The click event; the target must be a day cell.
 */
export function selectDay({ ele, evt }) {
  if (!Number.isNaN(Number(evt.target.innerHTML))) {
    const calendar = evt.target.closest(".calendar");
    let year = calendar.querySelector("input[name='year']").value;
    let month = calendar.querySelector("input[name='month']").value;

    const date = new Date(year, month, evt.target.innerText);

    const input = ele.previousElementSibling.querySelector("input");
    input.value = date.toLocaleDateString("en-GB");
    input.dispatchEvent(new Event("input"));

    ele.hidePopover();
  }
}

/**
 * Masks a time input as HH:MM while typing, clamping hours to 00–23 and
 * minutes to 00–59.
 *
 * @param {Object} param0 - Action payload.
 * @param {HTMLInputElement} param0.ele - The time input.
 */
export function timeOnChange({ ele }) {
  let value = ele.value.replace(/\D/g, "");

  value = value.slice(0, 4);

  let hours = value.slice(0, 2);
  let minutes = value.slice(2, 4);

  // Clamp hours (00–23)
  if (hours.length === 2) {
    let h = Math.min(parseInt(hours, 10), 23);
    hours = h.toString().padStart(2, "0");
  }

  // Clamp minutes (00–59)
  if (minutes.length === 2) {
    let m = Math.min(parseInt(minutes, 10), 59);
    minutes = m.toString().padStart(2, "0");
  }

  let formatted = hours;
  if (value.length > 2) {
    formatted += ":" + minutes;
  }

  ele.value = formatted;
}

/**
 * Combines the selected hour/minute options into the linked time input.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The clicked time option.
 */
export function selectTime({ ele }) {
  const popover = ele.closest(".date-input__popover");

  const opts = ele.parentElement.querySelectorAll("li");
  for (let i = 0; i < opts.length; i++) {
    opts[i].classList.remove("options__item--selected");
  }
  ele.classList.add("options__item--selected");

  let h = "00";
  let m = "00";

  const harr = popover.querySelectorAll('[data-type="hour"]');
  for (let i = 0; i < harr.length; i++) {
    if (harr[i].matches(".options__item--selected")) {
      h = harr[i].innerText;
      break;
    }
  }

  const marr = popover.querySelectorAll('[data-type="minutes"]');
  for (let i = 0; i < marr.length; i++) {
    if (marr[i].matches(".options__item--selected")) {
      m = marr[i].innerText;
      break;
    }
  }

  const input = popover.previousSibling.querySelector("input");
  input.value = `${h}:${m}`;
  input.dispatchEvent(new Event("input"));
}

/**
 * Opens the time popover, preselecting the hour/minute matching the linked
 * time input value.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The triggering element.
 * @param {DOMStringMap} param0.data - Dataset of the trigger; `data.target`
 *   is the id of the time options element.
 */
export function openTimePopover({ ele, data }) {
  const target = getTarget(ele, data.target);
  const inputValue = target
    .closest(".date-input__popover")
    .previousSibling.querySelector("input").value;

  const opts = target.querySelectorAll("li");
  for (let i = 0; i < opts.length; i++) {
    opts[i].classList.remove("options__item--selected");
  }
  if (inputValue.length < 5) {
    return;
  }

  const hhmm = inputValue.split(":");

  let h = hhmm[0];
  let m = hhmm[1];

  const harr = target.querySelectorAll('[data-type="hour"]');
  for (let i = 0; i < harr.length; i++) {
    if (harr[i].innerText === h) {
      harr[i].classList.add("options__item--selected");
      break;
    }
  }

  const marr = target.querySelectorAll('[data-type="minutes"]');
  for (let i = 0; i < marr.length; i++) {
    if (marr[i].innerText === m) {
      marr[i].classList.add("options__item--selected");
      break;
    }
  }

  if (target.matches(":popover-open")) {
    const selected = target.querySelectorAll(".options__item--selected");
    for (let i = 0; i < selected.length; i++) {
      selected[i].scrollIntoView({ block: "center" });
    }
  }
}
