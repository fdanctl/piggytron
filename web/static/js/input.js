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
