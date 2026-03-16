document.body.addEventListener("htmx:historyRestore", (ev) => {
  // tab title
  const splitPath = ev.detail.path.split("/");
  let title = splitPath[1];
  if (title === "") {
    title = "Dashboard";
  }
  title = title.charAt(0).toUpperCase() + title.slice(1).toLowerCase();
  document.title = title;

  // nav active link
  a.forEach((e) => e.classList.remove("active"));
  for (let i = 0; i < a.length; i++) {
    const text = a[i].text.trim().toLowerCase();
    a[i].classList.toggle("active", text === title.toLowerCase());
  }
});

// htmx custom events
document.body.addEventListener("incomeCategoryAdded", function (ev) {
  closeLastDialog();
  const li = document.querySelectorAll("#income-cat li");
  document.querySelector("#income-cat h4").innerText =
    `Income (${li.length + 1})`;
});

document.body.addEventListener("expenseCategoryAdded", function (ev) {
  closeLastDialog();
  const li = document.querySelectorAll("#expense-cat li");
  document.querySelector("#expense-cat h4").innerText =
    `Expenses (${li.length + 1})`;
});
