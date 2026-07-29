export function clickOption(ele, value) {
  const options = ele.closest(".dropdown").querySelectorAll("li");
  for (let i = 0; i < options.length; i++) {
    if (options[i].dataset.value === value) {
      options[i].click();
      break;
    }
  }
}

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

export function buildCalendar({ ele }) {
  const calendar = ele.closest(".calendar");
  const year = calendar.querySelector("input[name='year']").value;
  const month = calendar.querySelector("input[name='month']").value;

  const date = new Date(year, month, 1);
  const monthFirstWeekDay = date.getDay();
  const firstSunday = date - new Date(monthFirstWeekDay * 24 * 60 * 60 * 900); // hours * minutes * seconds * miliseconds

  const daysContainer = calendar.querySelector(".days");
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

function makeYearOption(year) {
  const li = document.createElement("li");
  li.dataset.action = "ui.select.option";
  li.dataset.value = year;
  const span = document.createElement("span");
  span.innerText = year;
  const svg = `<svg
      xmlns="http://www.w3.org/2000/svg"
      width="16"
      height="16"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
      class=""
    >
      <path d="M20 6 9 17l-5-5"></path>
    </svg>`;
  li.innerHTML = svg;
  li.prepend(span);
  return li;
}
