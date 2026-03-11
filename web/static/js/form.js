// disable child buttons
function disableBtns(e) {
  const btns = e.querySelectorAll(".btn");
  for (let i = 0; i < btns.length; i++) {
    btns[i].disabled = true;
  }
}
