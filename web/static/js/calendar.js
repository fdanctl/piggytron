function clickOption(ele, value) {
  const options = ele.closest(".dropdown").querySelectorAll("li");
  for (let i = 0; i < options.length; i++) {
    if (options[i].dataset.value === value) {
      options[i].click();
      break;
    }
  }
}

function prevMonth(ele) {
  const calendar = ele.closest(".calendar");
  const month = calendar.querySelector("input[name='month']");
  if (month.value === "0") {
    const year = calendar.querySelector("input[name='year']");
    clickOption(year, String(Number(year.value) - 1));
    clickOption(month, "11");
  } else {
    clickOption(month, String(Number(month.value) - 1));
  }
  month.dispatchEvent(new Event("change")); // triggers change event
}

function nextMonth(ele) {
  const calendar = ele.closest(".calendar");
  const month = calendar.querySelector("input[name='month']");
  if (month.value === "11") {
    const year = calendar.querySelector("input[name='year']");
    clickOption(year, String(Number(year.value) + 1));
    clickOption(month, "0");
  } else {
    clickOption(month, String(Number(month.value) + 1));
  }
  month.dispatchEvent(new Event("change")); // triggers change event
}

function buildMap(ele) {
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

    let ddiv = document.createElement("div");
    if (!((i < 7 && d > 7) || (i > 28 && d < 7))) {
      ddiv.innerText = d;
    } else {
      <span></span>;
    }
    daysContainer.appendChild(ddiv);
  }
}
