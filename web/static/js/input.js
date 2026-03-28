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
  ele.parentElement.nextElementSibling.value = ele.dataset.value;
  const opts = ele.parentElement.querySelectorAll("li");
  for (let i = 0; i < opts.length; i++) {
    opts[i].classList.remove("selected");
  }
  ele.classList.add("selected");
  ele.closest(".popover").hidePopover();
  let drop = ele.closest(".dropdown");
  drop.querySelector("button > span").innerText = ele.firstChild.innerText;
  drop.querySelector("button").classList.remove("error");
  ele.closest(".input-group").classList.remove("error");
}
