/** How many month options to append per lazy batch. */
const range = 12;
/** Month abbreviations used for option labels. */
const monthAbvMap = [
  "Jan",
  "Feb",
  "Mar",
  "Apr",
  "May",
  "Jun",
  "Jul",
  "Aug",
  "Sep",
  "Oct",
  "Nov",
  "Dec",
];

/**
 * Appends the next batch of older-month options to the dropdown, stopping
 * at the budget start month (data.break) or when it hits the range.
 * If the budget start month is not reached it appends a new sentinel <li>.
 *
 * @param {Object} param0 - Action payload.
 * @param {Element} param0.ele - The sentinel <li> that triggered the action.
 * @param {DOMStringMap} param0.data - Dataset of the sentinel; `data.last` is
 *   the newest month already listed ("YYYY-MM") and `data.break` is the
 *   budget start month ("YYYY-MM").
 */
export function generateMoreMonthsLI({ ele, data }) {
  let last = data.last;

  const [bpYear, bpMonth] = data.break.split("-").map(Number);
  const bp = new Date(bpYear, bpMonth - 1);

  const [year, month] = last.split("-").map(Number);
  let curr = new Date(year, month - 1);

  let opts = [];
  let i = range;
  let isAll = false;
  while (i > 0) {
    curr.setMonth(curr.getMonth() - 1);
    opts.push(makeMonthLI(curr));
    if (curr.toDateString() === bp.toDateString()) {
      isAll = true;
      break;
    }
    i--;
  }
  ele.parentElement.append(...opts);
  if (!isAll) {
    ele.parentElement.appendChild(ele);
  } else {
    ele.remove();
  }
}

/**
 * Builds a single month option linking to `?month=YYYY-MM`.
 *
 * @param {Date} date - The month to render.
 * @returns {HTMLLIElement} The new <li> element.
 */
function makeMonthLI(date) {
  const li = document.createElement("li");
  const a = document.createElement("a");
  a.href = `?month=${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, "0")}`;
  a.innerText = `${monthAbvMap[date.getMonth()]} ${date.getFullYear()}`;
  li.append(a);
  return li;
}

/**
 * Navigates to the previous budget month.
 */
export function prevBudgetMonth() {
  const params = new URLSearchParams(window.location.search);

  let prev = new Date();
  prev = new Date(prev.getFullYear(), prev.getMonth() - 1, 1);
  // if there is already a querystring
  const qmonth = params.get("month");
  if (qmonth) {
    const [year, month] = qmonth.split("-").map(Number);
    prev = new Date(year, month - 2);
  }

  htmx.ajax(
    "GET",
    `${window.location.pathname}?month=${prev.getFullYear()}-${String(prev.getMonth() + 1).padStart(2, "0")}`,
    {
      target: "#content",
      swap: "innerHTML",
      push: "true",
    },
  );
}

/**
 * Navigates to the next budget month.
 */
export function nextBudgetMonth() {
  const params = new URLSearchParams(window.location.search);

  let next = new Date();
  next = new Date(next.getFullYear(), next.getMonth() + 1, 1);
  const qmonth = params.get("month");
  // if there is already a querystring
  if (qmonth) {
    const [year, month] = qmonth.split("-").map(Number);
    next = new Date(year, month);
  }

  htmx.ajax(
    "GET",
    `${window.location.pathname}?month=${next.getFullYear()}-${String(next.getMonth() + 1).padStart(2, "0")}`,
    {
      target: "#content",
      swap: "innerHTML",
      push: "true",
    },
  );
}
