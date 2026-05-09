function handleInputOnBlur(ele) {
  ele.classList.toggle("error", !ele.validity.valid);
  const parent = ele.parentElement.parentElement;
  if (parent.classList.contains("input-group")) {
    parent.classList.toggle("error", !ele.validity.valid);
  }
}

function handlePasswordToggle(ele) {
  const pwdInput = ele.parentElement.parentElement.children[0];
  if (pwdInput.type === "password") {
    pwdInput.type = "text";
  } else if (pwdInput.type === "text") {
    pwdInput.type = "password";
  }
  ele.children[0].classList.toggle("hidden");
  ele.children[1].classList.toggle("hidden");
}

function handleCheckPillToggle(ele) {
  let cb = ele.querySelector("input");
  cb.checked = !cb.checked;
  cb.dispatchEvent(new Event("change")); // triggers change event
}

// select
function select(ele) {
  const input = ele.parentElement.nextElementSibling;
  input.value = ele.dataset.value;
  input.dispatchEvent(new Event("change")); // triggers change event

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

function handleSelectToggle(evt) {
  const popover = evt.currentTarget;

  if (popover.matches(":popover-open")) {
    const selected = popover.querySelector(".selected");
    selected?.scrollIntoView({ block: "center" });
  }
}

// date
function handleDateOnChange(evt) {
  let raw = evt.target.value.replace(/\D/g, "");

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

  evt.target.value = formatted;
}

// time
function handleTimeOnChange(evt) {
  let value = evt.target.value.replace(/\D/g, "");

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

  evt.target.value = formatted;
}

function handleClickTimePopover(ele) {
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

function openTimePopover(evt) {
  const inputValue = evt.target
    .closest(".popover")
    .previousSibling.querySelector("input").value;

  const popover = evt.target;
  const opts = popover.querySelectorAll("li");
  for (let i = 0; i < opts.length; i++) {
    opts[i].classList.remove("selected");
  }
  if (inputValue.length < 5) {
    return;
  }

  const hhmm = inputValue.split(":");

  let h = hhmm[0];
  let m = hhmm[1];

  const harr = popover.querySelectorAll('[data-type="hour"]');
  for (let i = 0; i < harr.length; i++) {
    if (harr[i].innerText === h) {
      harr[i].classList.add("selected");
      break;
    }
  }

  const marr = popover.querySelectorAll('[data-type="minutes"]');
  for (let i = 0; i < marr.length; i++) {
    if (marr[i].innerText === m) {
      marr[i].classList.add("selected");
      break;
    }
  }

  if (popover.matches(":popover-open")) {
    const selected = popover.querySelectorAll(".selected");
    for (let i = 0; i < selected.length; i++) {
      selected[i].scrollIntoView({ block: "center" });
    }
  }
}
