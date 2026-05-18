import { buildCalendar, clickOption } from "./calendar";

export function handleInputOnBlur(ele) {
  ele.classList.toggle("error", !ele.validity.valid);
  const parent = ele.parentElement.parentElement;
  if (parent.classList.contains("input-group")) {
    parent.classList.toggle("error", !ele.validity.valid);
  }
}

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

export function checkboxPillToggle({ ele }) {
  let cb = ele.querySelector("input");
  cb.checked = !cb.checked;
  cb.dispatchEvent(new Event("input", { bubbles: true })); // triggers change event
}

// select
export function selectSelect({ ele, data }) {
  const input = ele.parentElement.nextElementSibling;
  input.value = data.value;
  input.dispatchEvent(new Event("change", { bubbles: true })); // triggers change event

  const opts = ele.parentElement.querySelectorAll("li");
  for (let i = 0; i < opts.length; i++) {
    opts[i].classList.remove("selected");
  }
  ele.classList.add("selected");
  ele.closest(".popover").hidePopover();
  let drop = ele.closest(".dropdown");
  if (drop) {
    drop.querySelector("button > span").innerText = ele.firstChild.innerText;
    drop.querySelector("button").classList.remove("error");
    ele.closest(".input-group")?.classList.remove("error");
  }
}

export function selectToggle({ data }) {
  const popover = document.getElementById(data.target);

  requestAnimationFrame(() => {
    if (popover.matches(":popover-open")) {
      const selected = popover.querySelector(".selected");

      selected?.scrollIntoView({
        block: "center",
      });
    }
  });
}

// cash input
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

  intPart = parseInt(intPart || "0", 10).toLocaleString("en-US");

  ele.value = `${intPart}${decimalPart != "" ? "." + decimalPart : ""}`;
}

export function cashInputFocus({ ele }) {
  ele.value = ele.value.replaceAll(",", "");
  const length = ele.value.length;
  ele.setSelectionRange(length, length);
}

export function cashInputNav({ evt }) {
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

// date
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

  // Optional: limit year range
  // if (year.length === 4) {
  //   let y = parseInt(year, 10);
  //   y = Math.min(Math.max(y, 1900), 2100); // adjust if needed
  //   year = y.toString();
  // }

  let formatted = day;

  if (raw.length > 2) {
    formatted += "/" + month;
  }

  if (raw.length > 4) {
    formatted += "/" + year;
  }

  ele.value = formatted;
}

export function openCalendar({ data }) {
  const target = document.getElementById(data.target);
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
    day = ddmmyyyy[0];
    month = ddmmyyyy[1];
    year = ddmmyyyy[2];
  }

  const yearInput = target.querySelector("input[name='year']");
  const monthInput = target.querySelector("input[name='month']");
  clickOption(yearInput, String(year));
  clickOption(monthInput, String(month));
  monthInput.dispatchEvent(new Event("change"), { bubbles: true }); // triggers change event

  const days = target.querySelectorAll(".days > div");
  for (let i = 0; i < days.length; i++) {
    if (day == Number(days[i].innerText)) {
      days[i].classList.add("selected");
    }
  }
}

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

// time
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

export function selectTime({ ele }) {
  const popover = ele.closest(".popover");

  const opts = ele.parentElement.querySelectorAll("li");
  for (let i = 0; i < opts.length; i++) {
    opts[i].classList.remove("selected");
  }
  ele.classList.add("selected");

  let h = "00";
  let m = "00";

  const harr = popover.querySelectorAll('[data-type="hour"]');
  for (let i = 0; i < harr.length; i++) {
    if (harr[i].matches(".selected")) {
      h = harr[i].innerText;
      break;
    }
  }

  const marr = popover.querySelectorAll('[data-type="minutes"]');
  for (let i = 0; i < marr.length; i++) {
    if (marr[i].matches(".selected")) {
      m = marr[i].innerText;
      break;
    }
  }

  const input = popover.previousSibling.querySelector("input");
  input.value = `${h}:${m}`;
  input.dispatchEvent(new Event("input"));
}

export function openTimePopover({ data }) {
  const target = document.getElementById(data.target);
  const inputValue = target
    .closest(".popover")
    .previousSibling.querySelector("input").value;

  const opts = target.querySelectorAll("li");
  for (let i = 0; i < opts.length; i++) {
    opts[i].classList.remove("selected");
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
      harr[i].classList.add("selected");
      break;
    }
  }

  const marr = target.querySelectorAll('[data-type="minutes"]');
  for (let i = 0; i < marr.length; i++) {
    if (marr[i].innerText === m) {
      marr[i].classList.add("selected");
      break;
    }
  }

  if (target.matches(":popover-open")) {
    const selected = target.querySelectorAll(".selected");
    for (let i = 0; i < selected.length; i++) {
      selected[i].scrollIntoView({ block: "center" });
    }
  }
}
