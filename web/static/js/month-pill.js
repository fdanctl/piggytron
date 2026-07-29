const range = 12;
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

function makeMonthLI(date) {
  const li = document.createElement("li");
  const a = document.createElement("a");
  a.href = `?month=${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, "0")}`;
  a.innerText = `${monthAbvMap[date.getMonth()]} ${date.getFullYear()}`;
  li.append(a);
  return li;
}

export function prevBudgetMonth() {
  const params = new URLSearchParams(window.location.search);

  let prev = new Date();
  prev.setMonth(prev.getMonth() - 1);
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

export function nextBudgetMonth() {
  const params = new URLSearchParams(window.location.search);

  let next = new Date();
  next.setMonth(next.getMonth() + 1);
  const qmonth = params.get("month");
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
