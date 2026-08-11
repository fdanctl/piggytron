import { makeSVG } from "./icons";

/**
 * Programmatically clicks the dropdown option whose data-value matches,
 * used to sync calendar month/year <select>.
 *
 * @param {Element} ele - Any element inside the dropdown.
 * @param {string} value - The data-value to click.
 */
export function clickOption(ele, value) {
  const options = ele.closest(".dropdown").querySelectorAll("li");
  for (let i = 0; i < options.length; i++) {
    if (options[i].dataset.value === value) {
      options[i].click();
      break;
    }
  }
}

/**
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The month-back button.
 */
export function prevMonth({ ele }) {
  const calendar = ele.closest(".calendar");
  const month = calendar.querySelector("input[name='month']");
  if (month.value === "0") {
    const year = calendar.querySelector("input[name='year']");
    clickOption(year, String(Number(year.value) - 1));
    clickOption(month, "11");
  } else {
    clickOption(month, String(Number(month.value) - 1));
  }
  month.dispatchEvent(new Event("change", { bubbles: true })); // triggers change event
}

/**
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The month-forward button.
 */
export function nextMonth({ ele }) {
  const calendar = ele.closest(".calendar");
  const month = calendar.querySelector("input[name='month']");
  if (month.value === "11") {
    const year = calendar.querySelector("input[name='year']");
    clickOption(year, String(Number(year.value) + 1));
    clickOption(month, "0");
  } else {
    clickOption(month, String(Number(month.value) + 1));
  }
  month.dispatchEvent(new Event("change", { bubbles: true })); // triggers change event
}

/**
 * Rebuilds the day grid of the calendar for the selected month, rendering
 * up to 42 day cells starting from the Sunday before the 1st.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - Any element inside the calendar.
 */
export function buildCalendar({ ele }) {
  const calendar = ele.closest(".calendar");
  const year = calendar.querySelector("input[name='year']").value;
  const month = calendar.querySelector("input[name='month']").value;

  const date = new Date(year, month, 1);
  const monthFirstWeekDay = date.getDay();
  const firstSunday = date - new Date(monthFirstWeekDay * 24 * 60 * 60 * 1000); // hours * minutes * seconds * miliseconds

  const daysContainer = calendar.querySelector(".calendar__days");
  daysContainer.innerHTML = "";

  for (let i = 0; i < 42; i++) {
    const d = new Date(firstSunday + i * 24 * 60 * 60 * 1000).getDate();
    if (i === 28 && d < 7) {
      break;
    }
    if (i === 35 && d < 14) {
      break;
    }

    let ddiv;
    if (!((i < 7 && d > 7) || (i > 28 && d < 7))) {
      ddiv = document.createElement("div");
      ddiv.innerText = d;
    } else {
      ddiv = document.createElement("span");
    }
    daysContainer.appendChild(ddiv);
  }
}

/**
 * Appends or prepends the next ten year options to the year dropdown,
 * depending on the scroll direction of the sentinel option.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The sentinel option that triggered the action.
 * @param {DOMStringMap} param0.data - Dataset of the sentinel; `data.direction`
 *   is "down" (older years) or "up" (newer years).
 */
export function generateYearLI({ ele, data }) {
  const range = 10;
  let dir = data.direction;

  let opts = [];
  let y =
    dir === "down"
      ? Number(ele.previousElementSibling.dataset.value)
      : Number(ele.nextElementSibling.dataset.value);
  dir = dir === "down" ? -1 : 1;

  let i = 1;
  while (i <= range) {
    y = y + 1 * dir;
    opts.push(makeYearOption(y));
    i++;
  }

  if (dir < 0) {
    ele.parentElement.append(...opts);
    ele.parentElement.appendChild(ele);
  } else {
    ele.parentElement.prepend(...opts.reverse());
    ele.parentElement.prepend(ele);
  }
}

/**
 * Builds a year option (selectable via the ui.select.option action) with a
 * check icon.
 *
 * @param {number} year - The year to render.
 * @returns {HTMLLIElement} The new <li> element.
 */
function makeYearOption(year) {
  const li = document.createElement("li");
  li.dataset.action = "ui.select.option";
  li.dataset.value = year;
  const span = document.createElement("span");
  span.innerText = year;
  const svg = makeSVG("success", 16);
  li.innerHTML = svg;
  li.prepend(span);
  return li;
}
